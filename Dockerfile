FROM golang:1.19.0-alpine3.15 as build
RUN go install github.com/vektra/mockery/v2@latest
RUN apk add  make

RUN mkdir /videoscrapper
WORKDIR /videoscrapper
COPY .  /videoscrapper/
RUN make generateMocks
RUN go build

FROM  alpine:3.14
COPY --from=build /videoscrapper/videoScrapper ./
WORKDIR /home
RUN mkdir /output/
RUN mkdir traking_files/

ENTRYPOINT [ "/videoScrapper" ]

