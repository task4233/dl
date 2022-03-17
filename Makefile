.PHONY: clean
build:
	go build -o ./bin/delog -trimpath -mod=readonly ./cmd/delog

test:
	go test -v -race -shuffle=on ./...

test-with-coverage:
	go test -v -race -cover -shuffle=on ./... -coverprofile=cover.out
	go tool cover -html=cover.out -o cover.html

clean:
	rm -rf ./bin/delog