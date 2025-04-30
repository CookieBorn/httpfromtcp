FROM golang:1.23.4

RUN apt-get update && apt-get install -y netcat-openbsd

RUN go install github.com/bootdotdev/bootdev@latest

ENV PATH=$PATH:/go/bin

WORKDIR /usr/src/app
COPY . .

CMD ["go", "run", "."]
