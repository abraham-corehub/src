## How To Test Go HTTPS services

### Usage

```shell
➜ mkdir -p $GOPATH/src/github.com/jodosha && cd $GOPATH/src/github.com/jodosha
➜ git clone https://gist.github.com/jodosha/885dd981c657f599952b9c5df8f6b812 microservice && cd microservice
➜ chmod +x certificate.sh && ./certificate.sh
➜ go test -v
```