# use official Golang image
FROM golang:1.25

# set working directory
WORKDIR /app

# copy go.mod and go.sum first to leverage build cache
COPY go.mod go.sum ./

# download dependencies only
RUN go mod download

# copy the rest of the source
COPY . .

# build the server, which lives under cmd/server
WORKDIR /app/cmd/server
RUN go build -o /app/api .

# expose the port used by the service
EXPOSE 8000

# run the compiled binary
CMD ["/app/api"]
