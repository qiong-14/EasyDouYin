init:
	cd ./ && source ./dy_secure_config.sh
build:
	go build
run:
	./EasyDouYin
all: init build run
clean:
	rm ./EasyDouYin