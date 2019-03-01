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

type Controller struct {
	logger *logrus.Entry
	repo   repository.Repository
}

func New(repo repository.Repository, logger *logrus.Entry) *Controller {

	return &Controller{
		logger: logger,
		repo:   repo,
	}
}

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

func (c *Controller) Store(stream api.PersonService_StoreServer) error {
	personCount := 0
	start := time.Now()

	for {
		person, err := stream.Recv()
		if err == io.EOF {
			end := time.Now()
			return stream.SendAndClose(&api.PersonSummary{
				PersonCount: int32(personCount),
				ElapsedTime: int32(end.Sub(start).Seconds()),
			})
		}
		if err != nil {
			return errors.Wrapf(err, "receiving person from stream")
		}

		err = c.repo.Store(person)
		if err != nil {
			return errors.Wrapf(err, "saving person with id %d", person.Pid)
		}
		personCount++
	}
}

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
