package main

import (
	"context"
	"net"
	"os"
	"strings"

	homedir "github.com/mitchellh/go-homedir"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"golang.org/x/sync/errgroup"
	"google.golang.org/grpc"

	"github.com/richardcase/ingest-sample/pkg/api"
	controller "github.com/richardcase/ingest-sample/pkg/controller/person"
	"github.com/richardcase/ingest-sample/pkg/repository/mongodb"
	"github.com/richardcase/ingest-sample/pkg/signal"
)

const (
	DefaultListenAddress = "localhost:7777"
	DefaultLogLevel      = "Info"
	DefaultDbName        = "ingest"
	DefaultCollName      = "people"
)

var (
	options *Options

	rootCmd = &cobra.Command{
		Use:   "person-svc",
		Short: "A sample person service",
		Long:  "",
		Run: func(_ *cobra.Command, args []string) {
			if err := doRun(); err != nil {
				logrus.Fatalf("error running person service: %s", err.Error())
			}
		},
	}
)

func main() {
	cobra.OnInitialize(initConfig)

	if err := rootCmd.Execute(); err != nil {
		logrus.Fatal("error executing too command", err)
	}
}

func doRun() error {
	logger, err := configureLogging(options.LogLevel)
	if err != nil {
		logrus.Fatalf("failed to configure logging: %s", err.Error())
	}

	if options.DbURL == "" {
		logger.Fatal("database url must be supplied")
	}
	logger.Debugf("database URL is: %s\n", options.DbURL)

	repo, err := mongodb.NewRepository(options.DbName, options.CollName, options.DbURL)
	if err != nil {
		logger.WithError(err).Fatal("error creating Mongo repository")
	}

	controller := controller.New(repo, logger)

	stopChan := signal.SetupSignalHandler()

	lis, err := net.Listen("tcp", options.ListenAddress)
	if err != nil {
		logger.WithError(err).Fatalf("failed to listen on address: %s", options.ListenAddress)
	}
	//TODO: add TLS
	server := grpc.NewServer()
	api.RegisterPersonServiceServer(server, controller)

	ctx, cancel := context.WithCancel(context.Background())
	wg, ctx := errgroup.WithContext(ctx)
	wg.Go(func() error { return server.Serve(lis) })

	logger.Info("started person service")

	select {
	case <-stopChan:
		logger.Info("shutdown signal received, shutdown....")
	case <-ctx.Done():
	}
	server.GracefulStop()
	cancel()

	if err := wg.Wait(); err != nil {
		return errors.Wrap(err, "unhandled error, existing")
	}

	return nil

}

func init() {
	options = &Options{}

	rootCmd.PersistentFlags().StringVar(&options.ConfigFile, "config", "", "Confile file (default is $HOME/.person-svc.yaml")
	rootCmd.PersistentFlags().StringVarP(&options.LogLevel, "loglevel", "l", DefaultLogLevel, "Log level for the service")
	rootCmd.PersistentFlags().StringVarP(&options.ListenAddress, "listen-address", "", DefaultListenAddress, "the address the http server should listen on")
	rootCmd.PersistentFlags().StringVarP(&options.DbURL, "dburl", "d", "", "the database connection url")
	rootCmd.PersistentFlags().StringVarP(&options.DbName, "dbname", "n", DefaultDbName, "the database name")
	rootCmd.PersistentFlags().StringVarP(&options.CollName, "collname", "c", DefaultCollName, "the collection name")

	_ = viper.BindPFlag("config", rootCmd.PersistentFlags().Lookup("config"))
	_ = viper.BindPFlag("loglevel", rootCmd.PersistentFlags().Lookup("loglevel"))
	_ = viper.BindPFlag("listen-address", rootCmd.PersistentFlags().Lookup("listen-address"))
	_ = viper.BindPFlag("dburl", rootCmd.PersistentFlags().Lookup("dburl"))
	_ = viper.BindPFlag("dbname", rootCmd.PersistentFlags().Lookup("dbname"))
	_ = viper.BindPFlag("collname", rootCmd.PersistentFlags().Lookup("collname"))
}

func initConfig() {
	if options.ConfigFile != "" {
		viper.SetConfigFile(options.ConfigFile)
	} else {
		home, err := homedir.Dir()
		if err != nil {
			logrus.Fatalf("failed to get home directory: %v", err)
		}
		viper.AddConfigPath(home)
		viper.SetConfigName(".person-svc.yaml")
	}

	replacer := strings.NewReplacer(".", "-")
	viper.SetEnvKeyReplacer(replacer)
	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err == nil {
		logrus.Debugf("using config file: %s", viper.ConfigFileUsed())
	}
}

func configureLogging(logLevel string) (*logrus.Entry, error) {
	hostname, err := os.Hostname()
	if err != nil {
		return nil, errors.Wrap(err, "getting hostname")
	}

	if logLevel != "" {
		level, err := logrus.ParseLevel(strings.ToUpper(logLevel))
		if err != nil {
			return nil, errors.Wrapf(err, "parsing log level: %s", logLevel)
		}
		logrus.SetLevel(level)
	}

	logrus.SetFormatter(&logrus.JSONFormatter{})
	fields := logrus.Fields{
		"hostname":  hostname,
		"component": "person-svc",
	}

	return logrus.StandardLogger().WithFields(fields), nil
}
