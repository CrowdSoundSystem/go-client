# go-client

A simple cli tool for interacting with a crowdsound server. Usefully for testing and debugging.

# Getting

```
$ go get github.com/crowdsoundsystem/go-client
```


# Using

Since it's mostly used for hacking now, you can simply just hack away in main.go, and
run:

```
$ go run main.go
```

Additionally, you can specify the host and port for working with different environments

```
$ go run main.go -host [default=localhost] -port [default=50051]
```
