FROM golang:alpine as build

ARG VERSION

RUN mkdir /app

WORKDIR /app

COPY go.mod .

COPY go.sum .

RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -ldflags "-X 'main.Version=${VERSION}'" -o pepeg-bot2 .

FROM scratch as prod

LABEL org.opencontainers.image.source=https://github.com/rahagi/pepeg-bot2

COPY --from=build /app/pepeg-bot2 /app/pepeg-bot2

ENTRYPOINT ["/app/pepeg-bot2"]
