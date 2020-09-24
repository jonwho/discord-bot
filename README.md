# discord-bot

## Requirements
* Go 1.13
* Redis
* Discord account
* IEX Cloud account
* Alpaca account

## Optional
* Docker
* Minikube (or some other Kubernetes cluster host)
* Kubernetes

## ENV
* `cp .env.example .env`
* Fill in `.env` with your credentials
* Get your bot token from discord from [here](https://discordapp.com/developers/applications/me).
* Enable developer mode for Discord then right click channel or user to get ID.
* Grab your test/real tokens from [https://iexcloud.io/console/](https://iexcloud.io/console/)
* Grab your Alpaca ID and SECRET KEY from [https://app.alpaca.markets](https://app.alpaca.markets)

## Kubernetes Secrets
* Create secrets with `kubectl create secrets generic <uri>`
> Create from file is easier `kubectl create secret generic botsecrets --from-env-file=.env`
* Edit secrets with `kubectl edit secrets <uri>`
* View secrets with `kubectl get secret <uri> -o jsonpath='{.data}'`

## Run tests
Assuming you have filled in the `.env` file you can now run tests with:
```
make test
```

## Get it running one of three ways
### Binary
* Export ENV vars in .env to your shell
* Build the binary `make build-discord-bot`
* Run the binary `make run-discord-bot`
* Done

### Docker Compose
* Create a `.env` file see [example](#ENV).
* Build docker image and run compose `make up`
* Done

### Kubernetes
* Create a `.env` file see [example](#ENV).
* Create secrets `kubectl create secret generic botsecrets --from-env-file=.env`
* Create your cluster `minikube start`
* Apply k8s deployment to cluster `kubectl apply -f k8s/k8s.yml`
* Done

#### Docker Hub
For every bot update the docker image needs to be bumped so that Kubernetes can get the latest image.

* Build and tag image `docker build -t jonwho/discord-bot:runbot-v{n} -f Dockerfile.runbot .`
> Where n is the bump number
* Push the image to Docker Hub `docker push jonwho/discord-bot:runbot-v{n}`
