package mocks

import _ "github.com/vektra/mockery"

//go:generate ${GOPATH}/bin/mockery -tags netgo -dir=../ -name=PersonService_StoreServer -output=./
