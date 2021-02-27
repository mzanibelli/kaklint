SRC = internal/config/config.go
SRC += internal/errfmt/errfmt.go

kaklint: kaklint.go $(SRC)
	go build -o kaklint cmd/kaklint/main.go

test:
	go vet ./...
	go test ./...
	docker build -t kaklint-testing .
	docker run --rm --volume $(shell pwd):/kaklint kaklint-testing

install:
	go install ./cmd/kaklint

clean:
	rm -f kaklint
