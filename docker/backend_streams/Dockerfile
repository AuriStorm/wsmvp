FROM golang:1.22.4

ENV GOPATH "$HOME/go"
ENV PATH "$GOPATH/bin:$PATH"

WORKDIR /app

COPY ./backend_streams /app

CMD ["go", "run", "main.go"]
