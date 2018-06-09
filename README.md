A small Go utility to check for the optimal MTU value to set.

Runs on Linux and Darwin machines.

Grab the relevant binary or go get:
`go get https://gitlab.com/mijho/mtucheck"

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
