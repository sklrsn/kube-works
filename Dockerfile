FROM golang:1.22.3 as base

WORKDIR /go/src/github.com/kube-works

COPY main.go .
COPY go.mod .

RUN go build -o /opt/bin/welcome *.go

RUN chmod a+x /opt/bin/welcome

EXPOSE 8080

CMD [ "/opt/bin/welcome" ]