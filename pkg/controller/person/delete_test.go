package person_test

import (
	"context"

	"github.com/golang/protobuf/ptypes/empty"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"

	"github.com/richardcase/ingest-sample/pkg/api"
	controller "github.com/richardcase/ingest-sample/pkg/controller/person"
	"github.com/richardcase/ingest-sample/pkg/repository/mocks"
)

var _ = Describe("Delete Person", func() {
	var (
		r        *mocks.Repository
		ctl      *controller.Controller
		err      error
		personId int64
		logger   *logrus.Entry
		req      *api.DeletePersonRequest
		emp      *empty.Empty
		valid    controller.Validator
	)

	BeforeEach(func() {
		personId = 43
		logger = getTestLogger()
		valid = &controller.PersonValidator{}
	})

	Describe("when calling Delete", func() {

		Describe("with a negative person id", func() {
			BeforeEach(func() {
				r = &mocks.Repository{}
				r.On("Delete", -1).Return(nil)

				ctl = controller.New(r, logger, valid)
				req = &api.DeletePersonRequest{
					Id: int64(-1),
				}
				emp, err = ctl.Delete(context.TODO(), req)
			})

			It("should have returned an error", func() {
				Expect(err).To(HaveOccurred())
			})
			It("should not have called the repository", func() {
				Expect(r.AssertNotCalled(GinkgoT(), "Delete")).To(BeTrue())
			})
		})

		Describe("with a id of a person that doesn't exist", func() {
			BeforeEach(func() {
				r = &mocks.Repository{}
				r.On("Delete", personId).Return(nil)

				ctl = controller.New(r, logger, valid)
				req = &api.DeletePersonRequest{
					Id: int64(personId),
				}
				emp, err = ctl.Delete(context.TODO(), req)
			})
			It("should not have returned an error", func() {
				Expect(err).NotTo(HaveOccurred())
			})
			It("should have not returned a person", func() {
				Expect(emp).NotTo(BeNil())
			})
			It("should have called the repository", func() {
				Expect(r.AssertNumberOfCalls(GinkgoT(), "Delete", 1)).To(BeTrue())
			})
		})

		Describe("with a id of a person that exists", func() {
			BeforeEach(func() {
				r = &mocks.Repository{}
				r.On("Delete", personId).Return(nil)

				ctl = controller.New(r, logger, valid)
				req = &api.DeletePersonRequest{
					Id: int64(personId),
				}
				emp, err = ctl.Delete(context.TODO(), req)
			})
			It("should not have returned an error", func() {
				Expect(err).NotTo(HaveOccurred())
			})
			It("should have not returned a person", func() {
				Expect(emp).NotTo(BeNil())
			})
			It("should have called the repository", func() {
				Expect(r.AssertNumberOfCalls(GinkgoT(), "Delete", 1)).To(BeTrue())
			})
		})

		Describe("with a id of a person that exists but datastore errors", func() {
			BeforeEach(func() {
				r = &mocks.Repository{}
				r.On("Delete", personId).Return(errors.New("datastore had an error"))

				ctl = controller.New(r, logger, valid)
				req = &api.DeletePersonRequest{
					Id: int64(personId),
				}
				emp, err = ctl.Delete(context.TODO(), req)
			})
			It("should have returned an error", func() {
				Expect(err).To(HaveOccurred())
			})
			It("should have called the repository", func() {
				Expect(r.AssertNumberOfCalls(GinkgoT(), "Delete", 1)).To(BeTrue())
			})
		})

	})
})
