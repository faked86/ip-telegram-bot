FROM golang:1.21-alpine
RUN apk add git

WORKDIR /src/app

COPY go.mod ./
COPY go.sum ./
RUN go mod download

COPY . .

RUN go build -o server ./cmd/bot/main.go

EXPOSE 8080

ENTRYPOINT [ "./entrypoint.sh" ]

CMD [ "./server" ]