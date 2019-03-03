package person_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/ginkgo/extensions/table"
	. "github.com/onsi/gomega"

	"github.com/richardcase/ingest-sample/pkg/api"
	controller "github.com/richardcase/ingest-sample/pkg/controller/person"
)

var _ = Describe("Person Validation", func() {
	type resolveCase struct {
		Id          int64
		Email       string
		Name        string
		Phone       string
		ExpectError bool
	}

	DescribeTable("When validating a person",
		func(c resolveCase) {
			person := &api.Person{
				Id:           c.Id,
				Name:         c.Name,
				Email:        c.Email,
				MobileNumber: c.Phone,
			}
			validator := &controller.PersonValidator{}
			err := validator.Validate(person)

			if c.ExpectError {
				Expect(err).To(HaveOccurred())
			} else {
				Expect(err).NotTo(HaveOccurred())
			}
		},
		Entry("with a valid person", resolveCase{
			Id:          1,
			Name:        "Test",
			Email:       "test@test.com",
			Phone:       "07911287552",
			ExpectError: false,
		}),
		Entry("with an id that is less than 1", resolveCase{
			Id:          0,
			Name:        "Test",
			Email:       "test@test.com",
			Phone:       "07911287552",
			ExpectError: true,
		}),
		Entry("with an empty name", resolveCase{
			Id:          1,
			Name:        "",
			Email:       "test@test.com",
			Phone:       "07911287552",
			ExpectError: true,
		}),
		Entry("with an empty email", resolveCase{
			Id:          1,
			Name:        "Test",
			Email:       "",
			Phone:       "07911287552",
			ExpectError: true,
		}),
		Entry("with an email that isn't a valid address", resolveCase{
			Id:          1,
			Name:        "Test",
			Email:       "absfdfgfj",
			Phone:       "07911287552",
			ExpectError: true,
		}),
		Entry("with an empty phone", resolveCase{
			Id:          1,
			Name:        "Test",
			Email:       "test@test.com",
			Phone:       "",
			ExpectError: true,
		}),
	)
})
