#concept of environment for fizzle engine

name=phasarv

run:
	go build -o $(name) ./vendor/main && ./$(name)

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