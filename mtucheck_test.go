package main

import (
	"reflect"
	"testing"
)

func TestPingArgsPayloadIsMTUMinusHeaders(t *testing.T) {
	args, err := pingArgs("linux", 1500, "example.com")
	if err != nil {
		t.Fatal(err)
	}
	want := []string{"-M", "do", "-c", "2", "-w", "5", "-s", "1472", "example.com"}
	if !reflect.DeepEqual(args, want) {
		t.Fatalf("got %v, want %v", args, want)
	}
}

func TestPingArgsDarwin(t *testing.T) {
	args, err := pingArgs("darwin", 1468, "h")
	if err != nil {
		t.Fatal(err)
	}
	if args[len(args)-2] != "1440" {
		t.Fatalf("payload = %s, want 1440", args[len(args)-2])
	}
}

func TestPingArgsUnsupported(t *testing.T) {
	if _, err := pingArgs("windows", 1500, "h"); err == nil {
		t.Fatal("expected error for unsupported OS")
	}
}

func TestIfconfigHint(t *testing.T) {
	if got := ifconfigHint("linux", 1496, "eth0"); got != "ip link set dev eth0 mtu 1496" {
		t.Fatal(got)
	}
	if got := ifconfigHint("darwin", 1496, "en0"); got != "ifconfig en0 mtu 1496" {
		t.Fatal(got)
	}
}
