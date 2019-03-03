package mongodb

import (
	"context"

	"github.com/mongodb/mongo-go-driver/bson"
	"github.com/mongodb/mongo-go-driver/mongo"
	"github.com/mongodb/mongo-go-driver/mongo/options"
	"github.com/pkg/errors"

	"github.com/richardcase/ingest-sample/pkg/api"
)

type Repository struct {
	client     *mongo.Client
	collection *mongo.Collection
}

func NewRepository(database string, collection string, uri string, opts ...*options.ClientOptions) (*Repository, error) {
	client, err := mongo.Connect(context.TODO(), uri, opts...)
	if err != nil {
		return nil, errors.Wrapf(err, "connecting to mongodb %s", uri)
	}

	err = client.Ping(context.TODO(), nil)
	if err != nil {
		return nil, errors.Wrapf(err, "pinging mongodb")
	}

	coll := client.Database(database).Collection(collection)

	return &Repository{
		client:     client,
		collection: coll,
	}, nil
}

func (r *Repository) GetByID(id int64) (*api.Person, error) {
	filter := bson.D{{Key: "_id", Value: id}}
	result := r.collection.FindOne(context.TODO(), filter)
	if result.Err() != nil {
		return nil, errors.Wrapf(result.Err(), "getting person with id %d", id)
	}

	raw, err := result.DecodeBytes()
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, errors.Wrapf(err, "decoding bson for id %d", id)
	}
	person := mapRawBsonToPerson(raw)

	return person, nil

}

func (r *Repository) Store(person *api.Person) error {
	mapped := mapPersonToBson(person)

	existing, err := r.GetByID(person.Id)
	if err != nil {
		return errors.Wrapf(err, "checking if person already exists %d", person.Id)
	}

	if existing != nil {
		filter := bson.D{{Key: "_id", Value: person.Id}}
		result := r.collection.FindOneAndReplace(context.TODO(), filter, mapped)
		if result.Err() != nil {
			return errors.Wrapf(result.Err(), "updating person with id %d", person.Id)
		}
	} else {
		_, err := r.collection.InsertOne(context.TODO(), mapped)
		if err != nil {
			return errors.Wrapf(err, "inserting new person %d", person.Id)
		}
	}

	return nil
}

func (r *Repository) Delete(id int64) error {
	filter := bson.D{{Key: "_id", Value: id}}

	result := r.collection.FindOneAndDelete(context.TODO(), filter)
	if result.Err() != nil {
		return errors.Wrapf(result.Err(), "finding & deleting person with id %d", id)
	}

	return nil
}

func (r *Repository) Check() bool {
	if r.client == nil {
		return false
	}
	err := r.client.Ping(context.TODO(), nil)
	return err == nil
}
