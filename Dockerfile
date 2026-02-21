FROM golang:1.26-alpine AS build

RUN adduser --uid 1000 --disabled-password serverip

WORKDIR /app

COPY go.mod /app/
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build .

FROM scratch
COPY --from=build /etc/passwd /etc/passwd
COPY --from=build /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=build /app/serverip /serverip
USER serverip
CMD ["/serverip"]
