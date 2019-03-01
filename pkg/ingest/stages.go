package ingest

import (
	"strconv"
	"strings"

	"github.com/richardcase/ingest-sample/pkg/api"
)

func FormatEmail() Processor {
	return func(in <-chan interface{}, out chan interface{}) {
		for m := range in {
			rec := m.(Record)
			email := rec.Get("email")
			rec.Put("email", strings.ToUpper(email.(string)))
			out <- rec
		}
	}
}

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
			person.Pid = val
			person.Name = rec.Get("name").(string)
			person.Email = rec.Get("email").(string)
			person.MobileNumber = rec.Get("mobile_number").(string)

			out <- person
		}
	}
}
