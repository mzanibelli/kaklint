SRC = internal/config/config.go
SRC += internal/errfmt/errfmt.go

kaklint: kaklint.go $(SRC)
	go build -o kaklint cmd/kaklint/main.go

test:
	go vet ./...
	go test ./...
	docker-compose run --rm testing

install:
	go install ./cmd/kaklint

clean:
	rm -f kaklint
