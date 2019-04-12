SHELL := /bin/bash

# The name of the executable (default is current directory name)
TARGET := $(shell echo $${PWD})
install:
	go get
	go install
run:
	go run main.go
fmt:
	gofmt -s -l $(TARGET)
test:
	echo "running storage benchmark test..."
	go test $(TARGET)/storage -bench=.
	echo "done"
	echo "running storage unit tests..."
	go test $(TARGET)/storage
	echo "done"