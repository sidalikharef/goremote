.PHONY: build test vet race clean
build:
	go build ./...
test:
	go test ./... -count=1
vet:
	go vet ./...
race:
	go test -race ./...
clean:
	rm -f goremote
