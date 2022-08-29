include docker.env
test: generateMocks
	go test ./... 

build: test
	go build

buildImage: test
	docker build -t videoscrapper .

runDocker: buildImage
	docker run -it --rm -v "$(DOCKER_OUTPUT_PATH):/output/" -v "$(DOCKER_APPLICATION_HOME):/home/" videoscrapper

coverage:
	go test -coverprofile=coverage.out  ./...
	go tool cover -func=coverage.out

coverageHtml: coverage
	go tool cover -html=coverage.out

generateMocks:
	mockery --dir=provider/animeshow  --all  --output=mock/animeshow/  --outpkg=animeshow
	mockery --dir=service  --all  --output=mock/service/  --outpkg=service
	mockery --dir=port  --all  --output=mock/port/  --outpkg=port




