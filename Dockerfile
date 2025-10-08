FROM golang:1.25 AS build

WORKDIR /app

COPY . .

RUN go mod download
RUN go build -o /bambulabs-exporter .

FROM debian:bookworm-slim

COPY --from=build /bambulabs-exporter /bambulabs-exporter

EXPOSE 9101

CMD [ "/bambulabs-exporter" ]
