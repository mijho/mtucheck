A small Go utility to check for the optimal MTU value to set.

Runs on Linux and Darwin machines.

Grab with go get:
`go get gitlab.com/mijho/mtucheck`

Clone the repo and build the binary:
```
$ git clone https://gitlab.com/mijho/mtucheck.git
$ cd mtucheck
$ go build mtucheck.go
```

Example output:
```
$ mtucheck
ERR: 1500
ERR: 1492
ERR: 1484
ERR: 1476
PAS: 1468
ifconfig IFNAME mtu 1468
```
