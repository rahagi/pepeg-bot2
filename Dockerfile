FROM golang:alpine as build

ARG VERSION

RUN mkdir /app

WORKDIR /app

COPY go.mod .

COPY go.sum .

RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -ldflags "-X 'main.Version=${VERSION}'" -o pepeg-bot .

FROM scratch as prod

COPY --from=build /app/pepeg-bot /app/pepeg-bot

ENTRYPOINT ["/app/pepeg-bot"]
