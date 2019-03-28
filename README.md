# discord-bot

## Requirements
* Go 1.11
    * Download from here: https://golang.org/dl/
    * Update Go using instructions here: https://gist.github.com/nikhita/432436d570b89cab172dcf2894465753
    * Verify with `go version`
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
```
BOT_TOKEN=aaldhflj23roi0v8aaj1j13b.DLKLHAlkasjf9__1lk12hvaha-1-2987Q0
BOT_LOG_CHANNEL_ID=819231023981718081
MAINTAINERS=164098129888809171,917410876781231900
```
