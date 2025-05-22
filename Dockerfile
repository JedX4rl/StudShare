FROM golang:1.23

WORKDIR /app

COPY go.mod ./
COPY go.sum ./
RUN go mod download

COPY . ./

CMD ["go", "run", "./cmd/stud_shar  e/main.go"]