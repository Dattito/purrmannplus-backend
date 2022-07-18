# ./Dockerfile

FROM golang:1.18.4-alpine AS build

RUN apk --no-cache add build-base=0.5-r2

WORKDIR /build

COPY src/go.mod src/go.sum ./
RUN go mod download

COPY ./src .

RUN CGO_ENABLED=1 go build -o /app/app .

FROM alpine:3

RUN apk add --no-cache tzdata=2021e-r0 

ENV TZ Europe/Berlin

WORKDIR /app

COPY --from=build /app/app /app/app

ENV PATH_TO_API_VIEWS="/app/api/providers/rest/views" \
    PATH_TO_API_STATIC="/app/api/providers/rest/static"
COPY --from=build /build/api/providers/rest/views ${PATH_TO_API_VIEWS}
COPY --from=build /build/api/providers/rest/static ${PATH_TO_API_STATIC}

RUN mkdir /data

ENV DATABASE_TYPE SQLITE
ENV DATABASE_URI /data/db.sqlite

ENV LISTENING_PORT=3000

EXPOSE ${LISTENING_PORT}

ENTRYPOINT ["/app/app"]