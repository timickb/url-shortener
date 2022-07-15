FROM golang:latest

RUN go version
ENV GOPATH=/

COPY ./ ./

RUN apt update
RUN apt -y install postgresql-client

RUN chmod +x wait-for-postgres.sh

RUN go mod download
RUN make

CMD ["./artifacts/bin/urlapi"]