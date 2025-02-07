FROM golang:latest AS build-golang

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o /server ./cmd/server

FROM node:23-alpine AS build-scripts

WORKDIR /tmp/build

COPY ./web/package.json ./web/package-lock.json ./
RUN npm install

# Build scripts
COPY ./web ./
RUN mkdir dist
RUN npx browserify ./index.js > ./dist/bundle.js
RUN npx uglify-js --compress -o ./dist/bundle.min.js ./dist/bundle.js


FROM alpine:latest AS app

WORKDIR /app

COPY --from=build-golang /server ./

COPY ./assets ./assets
COPY ./templates ./templates
COPY --from=build-scripts /tmp/build/dist/bundle.min.js ./assets/bundle.js

EXPOSE ${PUBLIC_PORT:-8080}

CMD ["./server"]
