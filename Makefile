#concept of environment for fizzle engine

name=phasarv

run:
	go build -o $(name) ./vendor/main && ./$(name)

client:
	go build -o $(name) ./vendor/main && ./$(name) network
	# go build -o phClient ./vendor/client && ./phClient

serv:
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

build:
	# mkdir -p build
	go build -o build/$(name) ./vendor/main
	#GOOS=windows go build -o build/$(name).exe ./vendor/main

.PHONY: build