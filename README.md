# Auto Record

Automatically play music though Sonos automatically. 

## Pre-requisites 

- Git 
- Golang 
- OpenCV https://github.com/hybridgroup/gocv 
    Follow the steps on the README for your OS

## Installation 

Clone the repo to your Go /src directory. 

```bash 
export GO111MODULE=off
```

## Usage 

Download a release from the tags. 

Setup a project on https://console.cloud.google.com/ and export the API Key with "AR_API_KEY" to environment. 

Example 
```bash 
export AR_API_KEY=ABC123
```

Same for Spotify https://developer.spotify.com/

Example 
```bash 
export SPOTIFY_CLIENTID=ABC123
export SPOTIFY_CLIENTSECRET=ABC123
```

Setup the main autorecord program and button trigger program on boot. 

## Further Reading

https://developer.spotify.com/documentation/general/guides/authorization/code-flow/
https://godoc.org/gocv.io/x/gocv
https://developer.sonos.com/build/content-service-get-started/play-audio/account-matching/
https://cloud.google.com/vision/