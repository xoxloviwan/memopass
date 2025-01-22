# Memopass (Gophkeeper)
![coverage](https://raw.githubusercontent.com/xoxloviwan/memopass/badges/.badges/main/coverage.svg)

A secure password and sensitive data management system written in Go.
_This is my study project on Go._

## Features

  - Secure credential storage
  - Encrypted data synchronization
  - TUI client

## Installation (on Windows)

  - Install `task` from https://taskfile.dev/installation/
  
  ```bash
  go install github.com/go-task/task/v3/cmd/task@latest
  ```
  - Run for build project from source code:  
  ```bash
  task install
  task build
  ```

For troubleshouting see [taskfile.yml](./taskfile.yml) and install same for your OS.

## Run checks and get coverage:

  ```bash
  task lint
  task cover # count coverage
  task coverv # open html coverage report in browser
  ```


## Usage

### Server
```bash
task run
```
### Client

```bash
task run_cli
```

## Usage (with TLS)

For use with TLS certificates, you will need to generate a self-signed certificate. For Windows OS it can be done with the following command:
```bash
task run_tls
task run_cli_tls
```

## License
MIT License
