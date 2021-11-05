FROM golang:alpine as build

RUN mkdir /app

WORKDIR /app

COPY go.* ./

RUN go mod download

RUN go mod verify

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o pepeg_bot .

FROM scratch as prod

COPY --from=build /app/pepeg_bot /app/pepeg_bot

ENTRYPOINT ["/app/pepeg_bot"]
