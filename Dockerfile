# syntax=docker/dockerfile:1
# BUILD-STAGE
FROM golang:1.18-alpine as build

WORKDIR /go/app

COPY . .

# Download dependencies
RUN go mod download

# Without CGO_ENABLED=0, the resulting binary is not found in the CMD: `exec ./app: no such file or directory`
RUN CGO_ENABLED=0 go build -o /go/bin/app ./cmd

# RUN STAGE
# Use distroless image to reduce image size
FROM gcr.io/distroless/static-debian11

COPY --from=build /go/bin/app .

# Make port 80 available to the host machine
EXPOSE 80

CMD ["./app"]
