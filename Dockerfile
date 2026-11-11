FROM node:22-alpine AS assets
WORKDIR /src
COPY package.json package-lock.json ./
RUN npm ci
COPY framework/tooling/assets.mjs framework/tooling/assets.mjs
COPY resources resources
RUN npm run build

FROM golang:1.26.0-alpine AS build
WORKDIR /src
COPY go.mod go.sum ./
RUN go mod download
COPY app app
COPY cmd cmd
COPY config config
COPY framework framework
COPY public public
COPY resources resources
COPY routes routes
COPY --from=assets /src/public/assets public/assets
RUN CGO_ENABLED=0 go build -trimpath -ldflags="-s -w" -o /ticket-sync ./cmd/server

FROM alpine:3.23
RUN apk add --no-cache ca-certificates tzdata && adduser -D -u 10001 app && mkdir /data && chown app:app /data
WORKDIR /app
COPY --from=build /ticket-sync /app/ticket-sync
USER app
ENV PORT=8001 DATABASE_PATH=/data/ticket-sync.db APP_ENV=production
EXPOSE 8001
ENTRYPOINT ["/app/ticket-sync"]
