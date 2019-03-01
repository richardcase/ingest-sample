package mongodb

import (
	"github.com/mongodb/mongo-go-driver/bson"
	"github.com/richardcase/ingest-sample/pkg/api"
)

func mapPersonToBson(person *api.Person) bson.D {
	return bson.D{
		//nolint
		{"pid", person.Pid},
		//nolint
		{"name", person.Name},
		//nolint
		{"email", person.Email},
		//nolint
		{"mobile_number", person.MobileNumber},
		//nolint
		{"created", person.Created},
		//nolint
		{"updated", person.Updated},
	}
}

func mapRawBsonToPerson(raw bson.Raw) *api.Person {
	person := &api.Person{}
	val := raw.Lookup("pid")
	person.Pid = val.Int64()
	val = raw.Lookup("name")
	person.Pid = val.Int64()
	val = raw.Lookup("email")
	person.Email = val.String()
	val = raw.Lookup("mobile_number")
	person.Email = val.String()
	val = raw.Lookup("created")
	//person.Created = timestamp.Timestamp & updated

	return person
}
