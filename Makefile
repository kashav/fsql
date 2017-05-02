.PHONY: fsql clean install lint

all: fsql

fsql:
	go build -o ./fsql -v .

clean:
	rm -f ./fsql

install:
	go install

lint:
	${GOPATH}/bin/golint . query compare

test: fsql
	go test
