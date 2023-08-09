FROM golang:1.21

WORKDIR /go/src/app
COPY . .

RUN go get -d -v ./...
RUN go install -v ./...
RUN go build .

CMD ["iam-role-policy-changes-check"]
