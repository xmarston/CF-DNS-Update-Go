package cloudflare

import (
	"testing"
)

func TestCloudflareInit(t *testing.T) {
	cloudflareInitialization := Init("../")
	var expectedResult error = nil
	if cloudflareInitialization != expectedResult {
		t.Fatalf("Expected %s but got %s", expectedResult, cloudflareInitialization)
	}
}
