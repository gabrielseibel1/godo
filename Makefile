all: verify build

verify: vet check test

.PHONY: vet
vet:
	go vet ./...

.PHONY: check
check:
	staticcheck ./...

.PHONY: test
test:
	go test --count 1 --cover --coverprofile=./cover.out ./...

.PHONY: build
build:
	go build .

.PHONY: debug
debug:
	go build -gcflags=all="-N -l"

