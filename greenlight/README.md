# Greenlight

> Greenlight API

## Setup

```shell
go mod download
```

## Migrations

### Installation

```shell
brew install golang-migrate
```

### Verification

```shell
migrate -version
```

### Execution

```shell
make up
```

```shell
make down
```

## Help

```shell
make help
```

```shell
go run ./cmd/api -help
```

## Environment Variables

- I've already added `.envrc.example` in root of the project, rename it to `.envrc` and update the username & password

## Run

```shell
make run/api
```
