FROM golang:1.26 AS build

WORKDIR /app
# download go modules
COPY go.mod go.sum ./
RUN go mod download
# copy source code
COPY *.go ./
# build
RUN CGO_ENABLED=0 GOOS=linux go build -v -a -installsuffix cgo -o /app/ueberwachungskamera-server

FROM alpine:3.23

# install ffmpeg
RUN apk add --no-cache ffmpeg

WORKDIR /app
# copy binary from previous stage
COPY --from=build /app/ueberwachungskamera-server /app/ueberwachungskamera-server
# copy folders
COPY static ./static
COPY templates ./templates

EXPOSE 80
CMD ["/app/ueberwachungskamera-server"]
