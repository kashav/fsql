.PHONY: fsql clean install lint

all: fsql

fsql:
	go build -o ./fsql -v .

clean:
	rm -f ./fsql

install:
	go get -u -v

lint:
	${GOPATH}/bin/golint . query compare

test: fsql
	go test
