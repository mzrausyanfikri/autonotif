APP_NAME=autonotif-scheduler-osmosis
OUTPUT_DIR=../tmp
SAVE_DIR=../tmp
SERVER_HOST_PATH?=user@127.0.0.1:/home/user
SERVER_PORT?=80

.PHONY: clean
clean:
	rm -rf $(OUTPUT_DIR)
	rm -rf $(SAVE_DIR)
	docker image prune --force --filter='label=$(APP_NAME)'

.PHONY: compile
compile:
	mkdir -p $(OUTPUT_DIR)
	env GO111MODULE=on CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -ldflags '-extldflags "-static"' -o $(OUTPUT_DIR)/$(APP_NAME) ../../cmd/scheduler/main.go

.PHONY: docker-build
docker-build:
	docker build -t $(APP_NAME) --label $(APP_NAME) -f ../Dockerfile ../..

.PHONE: docker-save
docker-save:
	mkdir -p $(SAVE_DIR)
	docker save -o $(SAVE_DIR)/$(APP_NAME)-latest.tar $(APP_NAME):latest

.PHONY: scp
scp:
	scp -P $(SERVER_PORT) ../server/Makefile $(SERVER_HOST_PATH)
	scp -P $(SERVER_PORT) -r ../../db $(SERVER_HOST_PATH)
	scp -P $(SERVER_PORT) $(SAVE_DIR)/$(APP_NAME)-latest.tar $(SERVER_HOST_PATH)/app-images

.PHONE: upload-latest
upload-latest: clean compile docker-build docker-save scp
