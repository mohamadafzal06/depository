FROM golang:alpine as builder

RUN apk add git

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

WORKDIR /app/depository
RUN go build -o /depository

FROM alpine:latest as release


WORKDIR /app/

COPY --from=builder /depository .

EXPOSE 1378

ENTRYPOINT ["./depository"]

