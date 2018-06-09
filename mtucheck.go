package main

import (
	"fmt"
	"log"
	"os/exec"
	"runtime"
	"strconv"
)

var err error

func mtuCheckDarwin(s string) error {
	cmd := exec.Command("ping", "-D", "-c2", "-s", s, "1.1.1.1")
	_, err := cmd.CombinedOutput()
	return err
}

func mtuCheckLinux(s string) error {
	cmd := exec.Command("ping", "-M", "Do", "-c2", "-s", s, "1.1.1.1")
	_, err := cmd.CombinedOutput()
	return err
}

func checkMtuResult(s string, err error) int {
	if err != nil {
		errorcode := fmt.Sprintf("%s", err)
		if errorcode == "exit status 2"{
			fmt.Printf("ERR: %s\n", s)
			return 2
		} else {
			log.Fatalf("cmd.Run failed with %s\n", err)
		}
	}
	fmt.Printf("PAS: %s\n", s)
	return 0
}

func main() {
	count := 1212
	fmt.Println("just to check gitlab CI :)")

	OuterLoop:
		for i := 1500; i > count; i -= 8 {
			s := strconv.Itoa(i)

			switch os := runtime.GOOS; os {
				case "darwin":
					err := mtuCheckDarwin(s)
					ec := checkMtuResult(s, err)
					if ec == 0 {
						fmt.Printf("ifconfig IFNAME mtu %s\n", s)
						break OuterLoop
					}
				case "linux":
					err := mtuCheckLinux(s)
					ec := checkMtuResult(s, err)
					if ec == 0 {
						fmt.Printf("ip link set dev IFNAME mtu %s\n", s)
						break OuterLoop
					}
				default:
					fmt.Println("This script only runs on Linux and Darwin")
			}
		}
}