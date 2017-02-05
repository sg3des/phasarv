#concept of environment for fizzle engine

name=phasarv

run:
	go build -o $(name) ./vendor/main && ./$(name) ${ARGS}

client:
	go build -o $(name) ./vendor/main && ./$(name) client

server:
	go build -o phServer ./vendor/server && ./phServer

get:
	go get -u github.com/go-gl/glfw/v3.1/glfw
	go get -u github.com/go-gl/mathgl/mgl32
	go get -u github.com/golang/freetype
	go get -u github.com/tbogdala/groggy
	go get -u github.com/tbogdala/gombz
	
	go get -u github.com/go-gl/gl/v3.3-core/gl

	go get -u github.com/tbogdala/fizzle

test:
	go test -v vendor/main

network:
	cd vendor/network && go test -v ${ARGS} ./ 

network-bench:
	go run ./vendor/network/test-server/server.go &
	cd vendor/network && go test -run nil -bench . -benchmem



build:
	# mkdir -p build
	go build -o build/$(name) ./vendor/main
	#GOOS=windows go build -o build/$(name).exe ./vendor/main

.PHONY: build