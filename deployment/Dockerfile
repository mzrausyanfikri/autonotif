FROM alpine:3.11
WORKDIR /app
RUN chown root /app
USER root
COPY deployment/tmp/autonotif-scheduler-osmosis .
EXPOSE 8081
CMD ["./autonotif-scheduler-osmosis"]
RUN apk add --no-cache tzdata
ENV TZ="Asia/Jakarta"
