include docker.env
test:
	go test ./... 

build: test
	go build

buildImage: test
	docker build -t videoscrapper .

runDocker: buildImage
	docker run -it --rm -v "$(DOCKER_OUTPUT_PATH):/output/" -v "$(DOCKER_APPLICATION_HOME):/home/" videoscrapper

coverage:
	go test -coverprofile=coverage.out.tmp ./...
	cat coverage.out.tmp | grep -v "_mock.go" > coverage.out
	go tool cover -func=coverage.out





