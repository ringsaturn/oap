GOBASEPATH=$(shell go env var GOPATH | xargs)

gen:
	mockgen -source=$(GOBASEPATH)/pkg/mod/github.com/philchia/agollo/v4@v4.1.4/client.go -destination=mock_client_test.go -package=oap_test
