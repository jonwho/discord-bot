# discord-bot

## Requirements
* Go 1.12
### Method 1
  * Download from here: https://golang.org/dl/
  * Update Go using instructions here: https://gist.github.com/nikhita/432436d570b89cab172dcf2894465753
  * Verify with `go version`
### Method 2
  * Follow instructions at https://github.com/udhos/update-golang
* Discord Application with Bot support
* Docker (optional)
## Get it running
* Create a `.env` file see [example](#env-example).
* Install the dependencies
    * Verify you've enabled go modules by setting to your environment variables `GO111MODULE=on`
    * Run `go mod download` to install dependencies to your local cache.
    * Or run `go mod vendor` to install dependencies to a vendor folder in the project.
* Verify dependencies are installed with `go mod verify`
* Run `docker-compose build` then `docker-compose up`
* Done!

## ENV example
These values are made up you must supplement with your own credentials.
* Get your bot token from discord from [here](https://discordapp.com/developers/applications/me).
* Enable developer mode for Discord then right click channel or user to get ID.
* Grab your test/real tokens from [https://iexcloud.io/console/](https://iexcloud.io/console/)
* `cp .env.example .env`
* Fill in `.env` with your credentials
