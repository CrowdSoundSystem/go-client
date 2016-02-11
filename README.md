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
$ go run cmd.go -host [default=localhost] -port [default=50051] -cmd [default=queue|post|vote]
```

