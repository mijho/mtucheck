# mtucheck

A small Go utility to find the optimal (largest working) MTU for a network path.

It sends non-fragmentable pings (`DF` bit set) of decreasing size and reports the first MTU that gets through. Runs on Linux and macOS (Darwin).

## Install

```
go install github.com/mijho/mtucheck@latest
```

Or build from source:

```
git clone https://github.com/mijho/mtucheck.git
cd mtucheck
go build -o mtucheck .
```

## Usage

```
mtucheck [-host 1.1.1.1] [-max 1500] [-min 1200] [-step 8] [-iface IFNAME]
```

| Flag     | Default   | Meaning                            |
|----------|-----------|------------------------------------|
| `-host`  | `1.1.1.1` | Host to ping                       |
| `-max`   | `1500`    | Largest MTU to try                 |
| `-min`   | `1200`    | Smallest MTU to try                |
| `-step`  | `8`       | MTU decrement between probes       |
| `-iface` | `IFNAME`  | Interface name shown in the output |

Example output:

```
$ mtucheck -iface eth0
ERR: 1500
ERR: 1492
ERR: 1484
ERR: 1476
PAS: 1468
ip link set dev eth0 mtu 1468
```

Apply the suggested command as root, e.g. `sudo ip link set dev eth0 mtu 1468` on Linux or `sudo ifconfig en0 mtu 1468` on macOS.

## Notes

- `ping -s` sets the payload size, so each probe uses `MTU - 28` bytes (20-byte IPv4 header + 8-byte ICMP header).
- The probe checks the path to `-host`, not your whole network. Pick a host reachable through the interface you are tuning.
- The exit status is non-zero if no size in the range passed or `ping` is unavailable.

## Tests

```
go test ./...
```

## License

MIT. See [LICENSE](LICENSE).
