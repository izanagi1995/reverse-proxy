FROM golang:alpine AS build

WORKDIR /build
COPY . .
RUN go build -o reverseproxy .


FROM alpine:latest AS runtime

WORKDIR /app
COPY --from=build /build/reverseproxy .
ENTRYPOINT ./reverseproxy
