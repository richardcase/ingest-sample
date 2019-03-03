package person_test

import (
	"io"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/mock"

	"github.com/richardcase/ingest-sample/pkg/api"
	apimocks "github.com/richardcase/ingest-sample/pkg/api/mocks"
	controller "github.com/richardcase/ingest-sample/pkg/controller/person"
	"github.com/richardcase/ingest-sample/pkg/repository/mocks"
)

var _ = Describe("Store Person", func() {
	var (
		r        *mocks.Repository
		ctl      *controller.Controller
		err      error
		personId int64
		logger   *logrus.Entry
		person   *api.Person
		stream   *apimocks.PersonService_StoreServer
		valid    controller.Validator
	)

	BeforeEach(func() {
		personId = 43
		logger = getTestLogger()
		valid = &controller.PersonValidator{}
	})

	Describe("when calling Store", func() {

		Describe("with an stream thats EOF", func() {
			var summary *api.PersonSummary
			BeforeEach(func() {
				stream = &apimocks.PersonService_StoreServer{}
				stream.On("Recv").Return(nil, io.EOF)
				stream.On("SendAndClose", mock.MatchedBy(func(input *api.PersonSummary) bool {
					summary = input
					return true
				})).Return(nil)

				r = &mocks.Repository{}
				r.On("Store", -1).Return(nil, nil)

				ctl = controller.New(r, logger, valid)
				err = ctl.Store(stream)
			})

			It("should NOT have returned an error", func() {
				Expect(err).NotTo(HaveOccurred())
			})
			It("should not have called the repository", func() {
				Expect(r.AssertNotCalled(GinkgoT(), "Store")).To(BeTrue())
			})
			It("should have returned a person summary", func() {
				Expect(summary).NotTo(BeNil())
				Expect(summary.PersonCount).To(Equal(int32(0)))
			})
		})

		Describe("with an stream with 1 person", func() {
			var summary *api.PersonSummary
			BeforeEach(func() {
				person = createPerson(personId)
				stream = &apimocks.PersonService_StoreServer{}
				recvCallsPerson := 0
				recvCallsError := 0
				mockPersonResultsFn := func() *api.Person {
					recvCallsPerson++
					if recvCallsPerson == 1 {
						return person
					} else {
						return nil
					}
				}
				mockErrorResultsFn := func() error {
					recvCallsError++
					if recvCallsError == 1 {
						return nil
					} else {
						return io.EOF
					}
				}

				stream.On("Recv").Return(mockPersonResultsFn, mockErrorResultsFn)
				stream.On("SendAndClose", mock.MatchedBy(func(input *api.PersonSummary) bool {
					summary = input
					return true
				})).Return(nil)

				r = &mocks.Repository{}
				r.On("Store", person).Return(nil)

				ctl = controller.New(r, logger, valid)
				err = ctl.Store(stream)
			})

			It("should NOT have returned an error", func() {
				Expect(err).NotTo(HaveOccurred())
			})
			It("should have called the repository", func() {
				Expect(r.AssertNumberOfCalls(GinkgoT(), "Store", 1)).To(BeTrue())
			})
			It("should have returned a person summary", func() {
				Expect(summary).NotTo(BeNil())
				Expect(summary.PersonCount).To(Equal(int32(1)))
			})
		})

		Describe("with an stream that returns an error", func() {
			var summary *api.PersonSummary
			BeforeEach(func() {
				person = createPerson(personId)
				stream = &apimocks.PersonService_StoreServer{}

				stream.On("Recv").Return(nil, errors.New("stream error"))
				stream.On("SendAndClose", mock.MatchedBy(func(input *api.PersonSummary) bool {
					summary = input
					return true
				})).Return(nil)

				r = &mocks.Repository{}
				r.On("Store", person).Return(nil)

				ctl = controller.New(r, logger, valid)
				err = ctl.Store(stream)
			})

			It("should have returned an error", func() {
				Expect(err).To(HaveOccurred())
			})
			It("should not have called the repository", func() {
				Expect(r.AssertNotCalled(GinkgoT(), "Store")).To(BeTrue())
			})
			It("should NOT have returned a person summary", func() {
				Expect(summary).To(BeNil())
			})
		})

		Describe("with a stream of 1 person but the datastore errors", func() {
			var summary *api.PersonSummary
			BeforeEach(func() {
				person = createPerson(personId)
				stream = &apimocks.PersonService_StoreServer{}

				stream.On("Recv").Return(person, nil)
				stream.On("SendAndClose", mock.MatchedBy(func(input *api.PersonSummary) bool {
					summary = input
					return true
				})).Return(nil)

				r = &mocks.Repository{}
				r.On("Store", person).Return(errors.New("datastore error"))

				ctl = controller.New(r, logger, valid)
				err = ctl.Store(stream)
			})

			It("should have returned an error", func() {
				Expect(err).To(HaveOccurred())
			})
			It("should have called the repository", func() {
				Expect(r.AssertNumberOfCalls(GinkgoT(), "Store", 1)).To(BeTrue())
			})
			It("should NOT have returned a person summary", func() {
				Expect(summary).To(BeNil())
			})
		})
	})
})
