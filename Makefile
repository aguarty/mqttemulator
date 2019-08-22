BIN_PATH?=.bin/mqtt-emulator

vendors:
	go mod download
	go mod vendor

docker: 
	docker build -t mqtt-emulator .	
	
docker-run:
	docker run -it --rm mqtt-emulator

build:
	GOOS=linux CGO_ENABLED=0 go build -mod=vendor -o ${BIN_PATH} cmd/*.go

run:
	.bin/mqtt-emulator