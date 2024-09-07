FROM alpine:3.20.3
EXPOSE 8080

WORKDIR /app
COPY ./dist/collector_linux_amd64/collector /app
# COPY ./collector.yml ./collector.yml

CMD ["./collector"]
