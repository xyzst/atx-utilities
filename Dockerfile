FROM golang:1.21.3-bookworm
LABEL authors="darren"
ENV TARGET_CSV=input.csv

WORKDIR /usr/src/atx-utilities

COPY go.mod ./
COPY *.csv ./
RUN go mod download && go mod verify

COPY . .
RUN go build -v -o . ./...
RUN cp find-city-council-district /usr/local/bin

CMD find-city-council-district $TARGET_CSV