FROM alpine:latest

RUN apk update
RUN apk add git php go shellcheck rust npm curl

RUN curl --proto '=https' --tlsv1.2 -sSf https://sh.rustup.rs | sh -s -- -y

RUN . ~/.cargo/env && rustup target add x86_64-unknown-linux-musl

ENV GOPATH=/go
ENV GOBIN=/go/bin

ENV PATH=$PATH:/go/bin:/root/.cargo/bin

RUN go get golang.org/x/lint/golint
RUN npm -g install eslint

# Enable integration testing.
ENV KAKLINT_ENV=docker

COPY . /kaklint

WORKDIR /kaklint

RUN go get

ENTRYPOINT ["go", "test", "."]
