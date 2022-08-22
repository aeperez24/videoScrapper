include docker.env
test:
	go test ./... 

build: test
	go build

buildImage: test
	docker build -t videoscrapper .

runDocker: buildImage
	docker run -it --rm -v "$(DOCKER_OUTPUT_PATH):/output/" -v "$(DOCKER_APPLICATION_HOME):/home/" videoscrapper


