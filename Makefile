#concept of environment for fizzle engine

name=phasarv

run: install
	go build -o $(name) ./vendor/ && ./$(name) ${a}


client: install
	go build -o $(name) ./vendor/ && ./$(name) client

clients: install
	go build -o $(name) ./vendor/
	./$(name) client & 
	./$(name) client & 
	sleep 1
	./bspwm_place.sh phasarv-client


server:
	go build -o phServer ./vendor/server && ./phServer



install:
	go install ./vendor/...

get:
	go get -u github.com/go-gl/glfw/v3.1/glfw
	go get -u github.com/go-gl/mathgl/mgl32
	go get -u github.com/golang/freetype
	go get -u github.com/tbogdala/groggy
	go get -u github.com/tbogdala/gombz
	
	go get -u github.com/go-gl/gl/v3.3-core/gl

	go get -u github.com/tbogdala/fizzle

test:
	go test -v vendor/

network:
	cd vendor/network && go test -v ${a} ./ 

network-bench:
	go run ./vendor/network/test-server/server.go &
	cd vendor/network && go test -run nil -bench . -benchmem




.PHONY: build