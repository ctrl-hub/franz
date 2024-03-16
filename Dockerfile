FROM golang:1.22-alpine as build

WORKDIR /app

RUN apk add git

COPY go.mod go.mod
COPY go.sum go.sum

RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o bin .

FROM scratch

COPY --from=build /app/bin /app/bin

CMD ["/app/bin"]
