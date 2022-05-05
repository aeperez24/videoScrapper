# videoScrapper

VideoScrapper is a simple application to scrap and download series eposiodes from specific web pages.
To be used is necessary to modify the app.yaml configuration file where you need to specify the nexts fields:
## How to use:
- OutputPath: This is the path where episodes will be stored  when the are downloaded.

- SerieConfigurations: This field is meant to specify the list of series you want to track and download. it`s composed by:
  - SerieLink: the root link to scrap the tv serie.
  - SerieName: the name you want to use to store the episodes.
  - Provider:  the provider related to the SerieLink, right now only animeshowtv is allowed.



## What can we do next:

The goal of this  project is to be able of managing differents provider to download. right now  animeshowtv is the only provider avaliable. but it is possible to add new providers by  just implementing  the  GeneralDownloadService, and registering them on the main.go 


## command to compile for arm

env GOOS=linux GOARCH=arm GOARM=5 go build