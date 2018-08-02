start:betconstructProxy.bin
	./start.sh

betconstructProxy.bin:main.go betconstruct.go helperFuncs.go
	go build -o $@ $^

proxy.linux:main.go betconstruct.go helperFuncs.go
	GOOS=linux go build -o $@ $^
.PHONY:deploy
deploy:proxy.linux
	rsync -vaurz --progress --remove-source-files ./proxy.linux bproxy:/home/ubuntu/
