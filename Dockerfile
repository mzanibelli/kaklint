FROM archlinux:latest

RUN pacman -Sy --noconfirm git php go shellcheck rust eslint

ENV GOPATH=/go
ENV GOBIN=/go/bin
ENV PATH=$PATH:/go/bin:/root/.cargo/bin

RUN go get golang.org/x/lint/golint

# Enable integration testing.
ENV KAKLINT_ENV=docker

COPY . /kaklint

WORKDIR /kaklint

RUN go get -u

ENTRYPOINT ["go", "test", "-v", "."]
