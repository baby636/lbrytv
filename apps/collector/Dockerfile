FROM alpine:3.15.4
EXPOSE 8080

WORKDIR /app
COPY ./dist/collector_linux_amd64/collector /app
# COPY ./collector.yml ./collector.yml

CMD ["./collector"]
