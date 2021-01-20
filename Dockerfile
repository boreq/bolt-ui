FROM golang:1.14-alpine

RUN apk add git 
WORKDIR /velo
COPY . /velo
RUN go install -v ./cmd/velo

CMD ["/bin/sh", "-c", "velo run --verbosity debug /data/config.json"]
