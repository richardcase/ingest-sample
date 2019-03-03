package ingest

import (
	"context"
	"log"
	"time"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"

	"github.com/richardcase/ingest-sample/pkg/api"
)

// Destination represents a function that represents the last step in the pipeline
type Destination func(in <-chan interface{})

// Run will run the destination
func (d Destination) Run(in <-chan interface{}) {
	go func() {
		d(in)
	}()
}

// PrintDestination is a destination that prints the records it receives
func PrintDestination() Destination {
	return func(in <-chan interface{}) {
		for m := range in {
			log.Printf("%+v\n", m)
		}
	}
}

// PersonSvcClient is a client for the Person service
type PersonSvcClient struct {
	conn   *grpc.ClientConn
	client api.PersonServiceClient
	logger *logrus.Entry
}

// NewPersonSvcClient creates a new PersonSvcClient
func NewPersonSvcClient(serverAddress string, logger *logrus.Entry) (*PersonSvcClient, error) {
	//TODO: add TLS
	conn, err := grpc.Dial(serverAddress, grpc.WithInsecure())
	if err != nil {
		return nil, errors.Wrapf(err, "connecting to person svc: %s", serverAddress)
	}

	client := api.NewPersonServiceClient(conn)

	return &PersonSvcClient{
		conn:   conn,
		client: client,
		logger: logger,
	}, nil
}

// PersonSvcDestination is a destination that will call the Person service to store
// a Person.
func (c *PersonSvcClient) PersonSvcDestination() Destination {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Minute)
	stream, err := c.client.Store(ctx)
	if err != nil {
		//TODO: Handle errors propely
		panic(err.Error())
	}

	return func(in <-chan interface{}) {
		for m := range in {
			person := m.(*api.Person)
			if err := stream.Send(person); err != nil {
				c.logger.WithError(err).Fatalf("error saving person %v", person)
			}
		}
		reply, err := stream.CloseAndRecv()
		if err != nil {
			c.logger.WithError(err).Fatal("error closing person stream")
		}
		c.logger.Infof("%d people save in %d seconds", reply.PersonCount, reply.ElapsedTime)
		cancel()
	}
}
