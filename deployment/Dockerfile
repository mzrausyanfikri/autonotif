FROM alpine:3.11
WORKDIR /app
RUN chown root /app
USER root
COPY deployment/tmp/autonotif-scheduler .
EXPOSE 8080
CMD ["./autonotif-scheduler"]
RUN apk add --no-cache tzdata
ENV TZ="Asia/Jakarta"
