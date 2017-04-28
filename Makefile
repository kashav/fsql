.PHONY: fsql clean install lint

name = fsql

all: fsql

fsql:
	go build -o ./$(name) -v .

clean:
	rm -f ./$(name)

install:
	go install

lint:
	${GOPATH}/bin/golint . query

test: fsql
	go test
