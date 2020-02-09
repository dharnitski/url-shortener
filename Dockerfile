# 1) BUILD API
FROM golang:1.13 AS build-go
WORKDIR /url-shortener

# Copy the local files to the container's workspace.
ADD ./server .

# Build the service inside the container.
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o app .


# 2) BUILD UI
FROM node AS build-node 
WORKDIR /url-shortener

# Copy both projects to the container's workspace.
ADD . .

# go to client project
WORKDIR /url-shortener/client

RUN npm install
RUN npm run build:ci

# 3) BUILD FINAL IMAGE
FROM alpine
RUN apk --no-cache add ca-certificates
WORKDIR /app/server/
# copy Go binary
COPY --from=build-go /url-shortener/app /app/server
# copy DB migrations
COPY --from=build-go /url-shortener/migrations /app/server/migrations
# copy front end
COPY --from=build-node /url-shortener/server/ui /app/server/ui
EXPOSE 8080
CMD ["./app"]
