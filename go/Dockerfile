FROM golang:1.7

COPY gopath /go

WORKDIR /go/src/benchtest

RUN go-wrapper download
RUN go-wrapper install benchtest
ENV PARALLEL=10
ENV ITERATIONS=10

CMD ["go-wrapper", "run"]