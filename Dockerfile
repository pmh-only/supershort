FROM golang:1.21.1-alpine AS build

WORKDIR /app

RUN apk add build-base

COPY go.mod go.sum ./

RUN go mod download

COPY main.go ./

RUN CGO_ENABLED=1 GOOS=linux go build -o /app/supershort

# ---

FROM alpine

WORKDIR /app

COPY --from=build /app/supershort /app/supershort

RUN chown 1000:1000 /app -R
RUN chmod u+rwx /app -R

USER 1000:1000

ENTRYPOINT ["/app/supershort"]
