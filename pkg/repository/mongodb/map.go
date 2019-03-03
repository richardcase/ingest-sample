package mongodb

import (
	"github.com/mongodb/mongo-go-driver/bson"

	"github.com/richardcase/ingest-sample/pkg/api"
)

func mapPersonToBson(person *api.Person) bson.D {
	return bson.D{
		{Key: "_id", Value: person.Id},
		{Key: "name", Value: person.Name},
		{Key: "email", Value: person.Email},
		{Key: "mobile_number", Value: person.MobileNumber},
	}
}

func mapRawBsonToPerson(raw bson.Raw) *api.Person {
	person := &api.Person{}
	val := raw.Lookup("_id")
	person.Id = val.Int64()
	val = raw.Lookup("name")
	person.Name = val.String()
	val = raw.Lookup("email")
	person.Email = val.String()
	val = raw.Lookup("mobile_number")
	person.Email = val.String()

	return person
}
