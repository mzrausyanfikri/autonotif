FROM alpine:3.11

RUN adduser -S -D -H -h /app appuser

WORKDIR /app
RUN chown appuser /app
USER appuser

COPY deployment/tmp/autonotif-scheduler .

EXPOSE 8080

CMD ["./autonotif-scheduler"]

RUN apk add --no-cache tzdata

ENV TZ="Asia/Jakarta"
