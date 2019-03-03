package ingest

import (
	"strconv"

	"github.com/richardcase/ingest-sample/pkg/api"
	"github.com/richardcase/ingest-sample/pkg/phone"
)

// FormatPhoneE164 is a processing step that will format a UK
// phone number as E.164.
// See https://en.wikipedia.org/wiki/E.164
func FormatPhoneE164() Processor {
	return func(in <-chan interface{}, out chan interface{}) {
		for m := range in {
			rec := m.(Record)
			mobile := rec.Get("mobile_number")

			formatter := phone.NewUKToE164Formatter()

			formatted, _ := formatter.Format(mobile.(string))

			rec.Put("mobile_number", formatted)
			out <- rec
		}
	}
}

// MapToPerson is a processing step that maps a record to Person
func MapToPerson() Processor {
	return func(in <-chan interface{}, out chan interface{}) {
		for m := range in {
			rec := m.(Record)

			person := &api.Person{}

			val, err := strconv.ParseInt(rec.Get("id").(string), 10, 64)
			if err != nil {
				//TODO: handle errors, error channel??
				panic(err.Error())
			}
			person.Id = val
			person.Name = rec.Get("name").(string)
			person.Email = rec.Get("email").(string)
			person.MobileNumber = rec.Get("mobile_number").(string)

			out <- person
		}
	}
}
