SRC = internal/config/config.go
SRC += internal/errfmt/errfmt.go
SRC += internal/linter/linter.go

kaklint: kaklint.go $(SRC)
	go build -o kaklint cmd/kaklint/main.go

test:
	go vet ./...
	go test ./...

install:
	go install ./cmd/kaklint
