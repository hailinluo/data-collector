FROM golang:1.13

COPY bin/data-collector /app/bin/

WORKDIR /app/bin

CMD ["./data-collector"]