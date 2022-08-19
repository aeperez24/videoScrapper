FROM golang:1.19.0-alpine3.15 as build
RUN mkdir /videoscrapper
WORKDIR /videoscrapper
COPY .  /videoscrapper/
RUN go build

FROM  alpine:3.14
COPY --from=build /videoscrapper/animewatcher ./
WORKDIR home
RUN mkdir /output/
RUN mkdir traking_files/

ENTRYPOINT [ "/animewatcher" ]

