# Telegram Feeder

[![Tests](https://github.com/ksysoev/tg-feeder/actions/workflows/tests.yml/badge.svg)](https://github.com/ksysoev/tg-feeder/actions/workflows/tests.yml)
[![Go Report Card](https://goreportcard.com/badge/github.com/ksysoev/tg-feeder)](https://goreportcard.com/report/github.com/ksysoev/tg-feeder)
[![Go Reference](https://pkg.go.dev/badge/github.com/ksysoev/tg-feeder.svg)](https://pkg.go.dev/github.com/ksysoev/tg-feeder)
[![License: MIT](https://img.shields.io/badge/License-MIT-blue.svg)](https://opensource.org/licenses/MIT)

Service for listening RSS feeds, pre-processing and publishing to telegram

## Installation

## Building from Source

```sh
RUN CGO_ENABLED=0 go build -o feeder -ldflags "-X main.version=dev -X main.name=feeder" ./cmd/feeder/main.go
```

### Using Go

If you have Go installed, you can install Telegram Feeder directly:

```sh
go install github.com/ksysoev/tg-feeder/cmd/feeder@latest
```


## Using

```sh
feeder --log-level=debug --log-text=true --config=runtime/config.yml
```

## License

Telegram Feeder is licensed under the MIT License. See the LICENSE file for more details.
