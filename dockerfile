FROM golang:1.16-alpine AS build

## Erstelle Arbeits folder
WORKDIR /app

## Kopiere die mod datei rein
COPY go.mod .
COPY go.sum .
RUN go mod download

## Kopiere alle .go datein rein
COPY *.go ./

## Erstelle den Container mit name
RUN go build -o /docker-monitoring-sw

FROM gcr.io/distroless/base-debian10

WORKDIR /

COPY --from=build /docker-monitoring-sw /docker-monitoring-sw

CMD [ "/docker-monitoring-sw" ]