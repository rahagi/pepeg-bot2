FROM golang:alpine as build

ARG VERSION

RUN mkdir /app

WORKDIR /app

COPY go.* ./

RUN go mod download

RUN go mod verify

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -ldflags "-X 'main.Version=${VERSION}'" -o pepeg-bot .

FROM alpine as prod

COPY --from=build /app/pepeg-bot /app/pepeg-bot

ENTRYPOINT ["/app/pepeg-bot"]
