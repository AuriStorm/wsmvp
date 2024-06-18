FROM golang:1.22.0

ENV GOPATH "$HOME/go"
ENV PATH "$GOPATH/bin:$PATH"

WORKDIR /app
COPY ./frontend ./
RUN go mod init wsmvp/frontend-server

RUN go build .
CMD ["./frontend-server"]
