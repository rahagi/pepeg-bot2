FROM golang:alpine as build

RUN mkdir /app

WORKDIR /app

COPY . .

RUN GOOS=linux go build pepeg_bot.go

FROM scratch as prod

COPY --from=build /app/pepeg_bot .

ENTRYPOINT ["./pepeg_bot"]
