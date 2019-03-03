package person

import (
	"context"
	"io"
	"time"

	"github.com/golang/protobuf/ptypes/empty"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"

	"github.com/richardcase/ingest-sample/pkg/api"
	"github.com/richardcase/ingest-sample/pkg/repository"
)

// Controller contains the implementation of the people service
type Controller struct {
	logger    *logrus.Entry
	repo      repository.Repository
	validator Validator
}

// New creates a new Controller
func New(repo repository.Repository, logger *logrus.Entry, validator Validator) *Controller {
	return &Controller{
		logger:    logger,
		repo:      repo,
		validator: validator,
	}
}

// GetByID returns a Person for a supplied ID
func (c *Controller) GetByID(ctx context.Context, req *api.GetPersonRequest) (*api.Person, error) {
	if req.Id < 0 {
		return nil, errors.New("person id must be 0 or higher")
	}

	person, err := c.repo.GetByID(req.Id)
	if err != nil {
		return nil, errors.Wrapf(err, "getting person with id %d", req.Id)
	}

	return person, nil
}

// Store will store (upsert) People from a stream into a datastore. When the
// stream is closed a summary will be returned
func (c *Controller) Store(stream api.PersonService_StoreServer) error {
	personCount := 0
	errorCount := 0

	start := time.Now()

	for {
		person, err := stream.Recv()
		if err == io.EOF {
			end := time.Now()
			return stream.SendAndClose(&api.PersonSummary{
				PersonCount: int32(personCount),
				ErrorCount:  int32(errorCount),
				ElapsedTime: int32(end.Sub(start).Seconds()),
			})
		}
		if err != nil {
			return errors.Wrapf(err, "receiving person from stream")
		}

		err = c.validator.Validate(person)
		if err != nil {
			c.logger.Debugf("error validating person %d: %s", person.Id, err)
			//TODO: store the error and return
			errorCount++
			continue
		}

		err = c.repo.Store(person)
		if err != nil {
			c.logger.WithError(err).Error("error storing person")
			return errors.Wrapf(err, "saving person with id %d", person.Id)
		}
		personCount++
	}
}

// Delete will delete a specific Person given an id
func (c *Controller) Delete(ctx context.Context, req *api.DeletePersonRequest) (*empty.Empty, error) {
	if req.Id < 0 {
		return nil, errors.New("person id must be 0 or higher")
	}

	err := c.repo.Delete(req.Id)
	if err != nil {
		return nil, errors.Wrapf(err, "deleting person with id %d", req.Id)
	}

	return &empty.Empty{}, nil
}
