FROM node:16-alpine3.15 AS frontend-builder
WORKDIR /src
COPY . .
RUN yarn global add @vue/cli
RUN cd webui/frontend && \
    yarn && \
    yarn build

FROM golang:1.17 AS backend-builder
WORKDIR /src
COPY --from=frontend-builder /src .
RUN go mod download && \ 
    CGO_ENABLED=1 GOOS=linux go build -o /app -a -ldflags '-linkmode external -extldflags "-static"' .

############## Final stage ##############
FROM scratch

LABEL org.opencontainers.image.title="LinCast"
LABEL org.opencontainers.image.description="Your Open Source and 100% Free podcast player."
LABEL org.opencontainers.image.authors="Pegasus8" 
LABEL org.opencontainers.image.source="https://github.com/LinCast-Team/LinCast"
LABEL org.opencontainers.image.url="https://github.com/LinCast-Team/LinCast"
LABEL org.opencontainers.image.licenses="GPL-3.0"
LABEL build-date=""

COPY --from=backend-builder /app /app
EXPOSE 8080

ENTRYPOINT [ "/app" ]
