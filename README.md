# videoScrapper

VideoScrapper is a simple application to scrap and download series eposiodes from specific web pages.
To be used it's necessary to modify the app.yaml configuration file where you need to specify the nexts fields:
## How to use:
### app.yaml
The app.yaml configuration file is composed by the next fields
- OutputPath: This is the path where episodes will be stored  when the are downloaded. (only use this field if application is running outside docker).

- SerieConfigurations: This field specify the list of series you want to track and download. it`s composed by:
  - SerieLink: the root link to scrap the tv serie.
  - SerieName: the name you want to use to store the episodes.
  - Provider:  the provider related to the SerieLink, right now only animeshowtv is allowed.


See the app.yaml.example 

## Running with docker
- Create a new folder  to store the application data (application_home)
- Add an app.yaml in the application_home and  set the SerieConfigurations.
- Create a new folder for the application output (output_path). 
- Edit  the docker.env file setting the previously created folders as DOCKER_OUTPUT_PATH and DOCKER_APPLICATION_HOME
- execute the next command:
```make runDocker```


## command to compile for arm

env GOOS=linux GOARCH=arm GOARM=5 go build


## What can we do next:

The main goal of this project is to be able to download videos from different sources. So far animeshowtv and cuevana are the only available providers. But it is possible to add new providers by simply implementing the GeneralDownloadService interface , and registering them in the initializeDownloadServices function. 