package person_test

import (
	"context"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"

	"github.com/richardcase/ingest-sample/pkg/api"
	controller "github.com/richardcase/ingest-sample/pkg/controller/person"
	"github.com/richardcase/ingest-sample/pkg/repository/mocks"
)

var _ = Describe("Get Person", func() {
	var (
		r        *mocks.Repository
		ctl      *controller.Controller
		err      error
		personId int64
		logger   *logrus.Entry
		req      *api.GetPersonRequest
		person   *api.Person
		valid    controller.Validator
	)

	BeforeEach(func() {
		personId = 43
		logger = getTestLogger()
		valid = &controller.PersonValidator{}
	})

	Describe("when calling GetByID", func() {

		Describe("with a negative person id", func() {
			BeforeEach(func() {
				r = &mocks.Repository{}
				r.On("GetByID", -1).Return(nil, nil)

				ctl = controller.New(r, logger, valid)
				req = &api.GetPersonRequest{
					Id: int64(-1),
				}
				person, err = ctl.GetByID(context.TODO(), req)
			})

			It("should have returned an error", func() {
				Expect(err).To(HaveOccurred())
			})
			It("should have not returned a person", func() {
				Expect(person).To(BeNil())
			})
			It("should not have called the repository", func() {
				Expect(r.AssertNotCalled(GinkgoT(), "GetByID")).To(BeTrue())
			})
		})

		Describe("with a id of a person that doesn't exist", func() {
			BeforeEach(func() {
				r = &mocks.Repository{}
				r.On("GetByID", personId).Return(nil, nil)

				ctl = controller.New(r, logger, valid)
				req = &api.GetPersonRequest{
					Id: int64(personId),
				}
				person, err = ctl.GetByID(context.TODO(), req)
			})
			It("should not have returned an error", func() {
				Expect(err).NotTo(HaveOccurred())
			})
			It("should have not returned a person", func() {
				Expect(person).To(BeNil())
			})
			It("should have called the repository", func() {
				Expect(r.AssertNumberOfCalls(GinkgoT(), "GetByID", 1)).To(BeTrue())
			})
		})

		Describe("with a id of a person that exists", func() {
			var (
				repoPerson *api.Person
			)
			BeforeEach(func() {
				r = &mocks.Repository{}
				repoPerson = createPerson(int64(personId))
				r.On("GetByID", personId).Return(repoPerson, nil)

				ctl = controller.New(r, logger, valid)
				req = &api.GetPersonRequest{
					Id: int64(personId),
				}
				person, err = ctl.GetByID(context.TODO(), req)
			})
			It("should not have returned an error", func() {
				Expect(err).NotTo(HaveOccurred())
			})
			It("should have called the repository", func() {
				Expect(r.AssertNumberOfCalls(GinkgoT(), "GetByID", 1)).To(BeTrue())
			})
			It("should have returned a person", func() {
				Expect(person).To(BeEquivalentTo(repoPerson))
			})
		})

		Describe("with a id of a person that exists but datastore errors", func() {
			BeforeEach(func() {
				r = &mocks.Repository{}
				r.On("GetByID", personId).Return(nil, errors.New("datastore had an error"))

				ctl = controller.New(r, logger, valid)
				req = &api.GetPersonRequest{
					Id: int64(personId),
				}
				person, err = ctl.GetByID(context.TODO(), req)
			})
			It("should have returned an error", func() {
				Expect(err).To(HaveOccurred())
			})
			It("should have called the repository", func() {
				Expect(r.AssertNumberOfCalls(GinkgoT(), "GetByID", 1)).To(BeTrue())
			})
			It("should have not returned a person", func() {
				Expect(person).To(BeNil())
			})
		})

	})
})
