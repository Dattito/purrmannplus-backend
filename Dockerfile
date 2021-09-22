# ./Dockerfile

FROM golang:1.17-alpine AS build

RUN apk --no-cache add build-base=0.5-r2

WORKDIR /build

COPY src/go.mod src/go.sum ./
RUN go mod download

COPY ./src .

RUN CGO_ENABLED=1 go build -o /bin/app .

FROM alpine:3

COPY --from=build /bin/app /bin/app

ENTRYPOINT ["/bin/app"]