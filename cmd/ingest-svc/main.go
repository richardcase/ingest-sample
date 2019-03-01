package main

import (
	"bufio"
	"context"
	"os"
	"runtime"
	"strings"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"

	"github.com/richardcase/ingest-sample/pkg/ingest"
	"github.com/richardcase/ingest-sample/pkg/signal"
	"golang.org/x/sync/errgroup"
)

const (
	DefaultServiceAddress = "127.0.0.1:7777"
	DefaultLogLevel       = "Info"
)

var (
	serviceAddress string
	logLevel       string
	numWorkers     int
	sampleFile     string

	rootCmd = &cobra.Command{
		Use:   "ingest-svc",
		Short: "A sample ingest service",
		Long:  "",
		Run: func(c *cobra.Command, args []string) {
			if err := doRun(); err != nil {
				logrus.Fatalf("error running ingest service: %s", err.Error())
			}
		},
	}
)

func main() {
	if err := rootCmd.Execute(); err != nil {
		logrus.Fatal("error executing command", err)
	}
}

func init() {
	rootCmd.PersistentFlags().StringVarP(&logLevel, "loglevel", "l", DefaultLogLevel, "Log level for the service")
	rootCmd.PersistentFlags().StringVarP(&serviceAddress, "service-address", "", DefaultServiceAddress, "the address the http server should listen on")
	rootCmd.PersistentFlags().IntVarP(&numWorkers, "workers", "w", defaultWorkersNum(), "the number of workers to start")
	rootCmd.PersistentFlags().StringVarP(&sampleFile, "sample", "s", "", "the path to the sample file")
}

func doRun() error {
	logger, err := configureLogging(logLevel)
	if err != nil {
		logrus.Fatalf("failed to configure logging: %s", err.Error())
	}
	if serviceAddress == "" {
		logger.Fatal("you must supply the address of the person service")
	}

	stopChan := signal.SetupSignalHandler()

	ctx, cancel := context.WithCancel(context.Background())
	wg, ctx := errgroup.WithContext(ctx)

	sourceStream := make(chan interface{}, numWorkers)

	logger.Info("starting source worker")
	wg.Go(func() error { return source(sampleFile, sourceStream) })

	logger.Info("starting workers")
	for i := 0; i < numWorkers; i++ {
		logger.Debugf("starting worker %d", i)
		wg.Go(func() error { return worker(logger, serviceAddress, sourceStream) })
	}

	go func() {
		select {
		case <-stopChan:
			logger.Info("shutdown signal received, shutting down")
		case <-ctx.Done():
		}

		cancel()
	}()

	if err := wg.Wait(); err != nil {
		return errors.Wrapf(err, "error occured waiting")
	}

	return nil
}

func defaultWorkersNum() int {
	numProcs := runtime.NumCPU() - 1
	if numProcs < 1 {
		return 1
	}
	return numProcs
}

func worker(logger *logrus.Entry, serverAddress string, in chan interface{}) error {
	client, err := ingest.NewPersonSvcClient(serverAddress, logger)
	if err != nil {
		return err
	}
	personPipeline := ingest.NewPipeline("ingest-person", client.PersonSvcDestination(), ingest.FormatEmail(), ingest.MapToPerson())
	personPipeline.Run(in)

	return nil
}

func source(dataPath string, out chan interface{}) error {
	file, err := os.Open(dataPath)
	if err != nil {
		return err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanLines)

	lineNo := 0
	var headers []string
	containsHeader := true
	for scanner.Scan() {
		line := scanner.Text()
		if containsHeader && lineNo == 0 {
			headers = strings.Split(line, ",")
			lineNo++
			continue
		}

		parts := strings.Split(line, ",")
		record := ingest.NewMapRecord()
		for index := range parts {
			record.Put(headers[index], parts[index])
		}
		out <- record
		lineNo++
	}
	close(out)
	return nil
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
