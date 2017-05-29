name := fsql
sources := $(wildcard **/*.go)

build = GOOS=$(1) GOARCH=$(2) go build -o build/$(name)$(3)
tar = cd build && tar -cvzf $(1)_$(2).tar.gz $(name)$(3) && rm $(name)$(3)
zip = cd build && zip $(1)_$(2).zip $(name)$(3) && rm $(name)$(3)

.PHONY: clean install lint test

all: fsql

build: darwin linux windows

fsql: $(sources)
	go build -o ./$(name) -v .

clean:
	rm -rf ./$(name) build/

lint:
	${GOPATH}/bin/golint ./...

install:
	go get -u -v ./...

test:
	go test -v ./...

##### LINUX BUILDS #####
linux: build/linux_arm.tar.gz build/linux_arm64.tar.gz build/linux_386.tar.gz build/linux_amd64.tar.gz

build/linux_386.tar.gz: $(sources)
	$(call build,linux,386,)
	$(call tar,linux,386)

build/linux_amd64.tar.gz: $(sources)
	$(call build,linux,amd64,)
	$(call tar,linux,amd64)

build/linux_arm.tar.gz: $(sources)
	$(call build,linux,arm,)
	$(call tar,linux,arm)

build/linux_arm64.tar.gz: $(sources)
	$(call build,linux,arm64,)
	$(call tar,linux,arm64)

##### DARWIN BUILDS #####
darwin: build/darwin_amd64.tar.gz

build/darwin_amd64.tar.gz: $(sources)
	$(call build,darwin,amd64,)
	$(call tar,darwin,amd64)

##### WINDOWS BUILDS #####
windows: build/windows_386.zip build/windows_amd64.zip

build/windows_386.zip: $(sources)
	$(call build,windows,386,.exe)
	$(call zip,windows,386,.exe)

build/windows_amd64.zip: $(sources)
	$(call build,windows,amd64,.exe)
	$(call zip,windows,amd64,.exe)
