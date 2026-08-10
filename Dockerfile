FROM golang:1.26-alpine AS build
WORKDIR /src
COPY go.mod ./
COPY *.go ./
COPY static ./static
RUN CGO_ENABLED=0 GOOS=linux go build -trimpath -ldflags="-s -w" -o /khoan-chuyen .

FROM gcr.io/distroless/static-debian12:nonroot
COPY --from=build /khoan-chuyen /khoan-chuyen
ENV PORT=8080
EXPOSE 8080
ENTRYPOINT ["/khoan-chuyen"]
