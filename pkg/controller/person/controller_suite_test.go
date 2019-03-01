package person_test

import (
	"testing"
	"time"

	timestamp "github.com/golang/protobuf/ptypes/timestamp"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/sirupsen/logrus"

	"github.com/richardcase/ingest-sample/pkg/api"
)

func TestPersonController(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Person Controller Suite")
}

func getTestLogger() *logrus.Entry {
	logrus.SetOutput(GinkgoWriter)
	logrus.SetLevel(logrus.InfoLevel)
	logrus.SetFormatter(&logrus.TextFormatter{})

	return logrus.StandardLogger().WithField("test", true)
}

func createPerson(id int64) *api.Person {
	created := &timestamp.Timestamp{}
	created.Seconds = time.Now().Unix()

	return &api.Person{
		Pid:          id,
		Email:        "test@test.com",
		Name:         "Test Person",
		MobileNumber: "+44 (0)7833 567991",
		Created:      created,
		Updated:      created,
	}
}
