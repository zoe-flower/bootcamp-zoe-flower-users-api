//go:build integration

package main

import (
	"flag"
	"net/http"
	"os"
	"testing"

	"github.com/sethgrid/pester"
)

var serviceAddr = flag.String("service", os.Getenv("TEST_SERVICE_ADDR"), "Service address against which to run tests")

func TestPingService(t *testing.T) {
	resp, err := pester.Get(*serviceAddr + "/ping")
	if err != nil {
		t.Fatalf("❌ ping service: %+v", err)
	}

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("❌ ping status: %s", resp.Status)
	}

	t.Log("✅ Service alive")
}
