FROM alpine:latest

WORKDIR /app

COPY pocketbase /app/

RUN chmod +x /app/pocketbase && \
    adduser -D pocketbase && \
    mkdir -p /pb_data && \
    chown -R pocketbase:pocketbase /pb_data

USER pocketbase

EXPOSE 8090

CMD ["/app/pocketbase", "serve", "--http=0.0.0.0:8090", "--dir=/pb_data"]
