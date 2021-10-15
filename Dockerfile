FROM golang:1.17-alpine AS builder
RUN apk add --no-cache ca-certificates git

# Set the working directory outside $GOPATH to enable the support for modules.
WORKDIR /builder

# Fetch dependencies first; they are less susceptible to change on every build
# and will therefore be cached for speeding up the next build
COPY ./go.mod ./go.sum ./
RUN go mod download

# Import the code from the context.
COPY . .
RUN go build -o /app cmd/api/*.go

FROM alpine
WORKDIR /app
COPY --from=builder /app ./
CMD ["./app"]