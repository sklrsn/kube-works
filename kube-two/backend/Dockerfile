FROM golang:1.22.3 as base

WORKDIR /go/src/github.com/sklrsn/kube-works/kube-two/backend

COPY go.mod go.sum main.go ./

RUN go build -o /opt/bin/backend *.go

RUN chmod a+x /opt/bin/backend

EXPOSE 8080

CMD [ "/opt/bin/backend" ]