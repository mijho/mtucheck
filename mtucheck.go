// Command mtucheck finds the largest MTU that can reach a host without
// fragmentation, by sending non-fragmentable pings of decreasing size.
package main

import (
	"errors"
	"flag"
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strconv"
)

// headerOverhead is the IPv4 header (20 bytes) plus the ICMP header (8 bytes).
// ping's -s flag sets the ICMP payload size, so MTU = payload + headerOverhead.
const headerOverhead = 28

// pingArgs returns the ping arguments that send a single non-fragmentable
// probe with the given MTU, or an error if the OS is unsupported.
func pingArgs(goos string, mtu int, host string) ([]string, error) {
	payload := strconv.Itoa(mtu - headerOverhead)
	switch goos {
	case "darwin":
		return []string{"-D", "-c", "2", "-t", "5", "-s", payload, host}, nil
	case "linux":
		return []string{"-M", "do", "-c", "2", "-w", "5", "-s", payload, host}, nil
	default:
		return nil, fmt.Errorf("unsupported OS %q: only Linux and Darwin are supported", goos)
	}
}

// probe reports whether a packet of the given MTU gets through unfragmented.
// A non-nil error means ping could not be run at all; a failed probe is not an error.
func probe(goos string, mtu int, host string) (bool, error) {
	args, err := pingArgs(goos, mtu, host)
	if err != nil {
		return false, err
	}
	if err := exec.Command("ping", args...).Run(); err != nil {
		var exitErr *exec.ExitError
		if errors.As(err, &exitErr) {
			return false, nil // no reply, "frag needed", or "message too long"
		}
		return false, err // ping missing or not executable
	}
	return true, nil
}

// ifconfigHint returns the command that would apply the MTU on this OS.
func ifconfigHint(goos string, mtu int, iface string) string {
	if goos == "darwin" {
		return fmt.Sprintf("ifconfig %s mtu %d", iface, mtu)
	}
	return fmt.Sprintf("ip link set dev %s mtu %d", iface, mtu)
}

func main() {
	host := flag.String("host", "1.1.1.1", "host to probe")
	maxMTU := flag.Int("max", 1500, "largest MTU to try")
	minMTU := flag.Int("min", 1200, "smallest MTU to try")
	step := flag.Int("step", 8, "MTU decrement between probes")
	iface := flag.String("iface", "IFNAME", "interface name to show in the suggested command")
	flag.Parse()

	if *step <= 0 || *minMTU <= headerOverhead || *maxMTU < *minMTU {
		fmt.Fprintln(os.Stderr, "invalid -min/-max/-step values")
		os.Exit(2)
	}

	goos := runtime.GOOS
	for mtu := *maxMTU; mtu >= *minMTU; mtu -= *step {
		ok, err := probe(goos, mtu, *host)
		if err != nil {
			fmt.Fprintf(os.Stderr, "mtucheck: %v\n", err)
			os.Exit(1)
		}
		if !ok {
			fmt.Printf("ERR: %d\n", mtu)
			continue
		}
		fmt.Printf("PAS: %d\n", mtu)
		fmt.Println(ifconfigHint(goos, mtu, *iface))
		return
	}
	fmt.Fprintf(os.Stderr, "no MTU in [%d, %d] passed; check connectivity to %s\n", *minMTU, *maxMTU, *host)
	os.Exit(1)
}
