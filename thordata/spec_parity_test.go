package tests

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
)

func loadSpec(t *testing.T) map[string]any {
	t.Helper()
	p := os.Getenv("THORDATA_SDK_SPEC_PATH")
	if p == "" {
		p = filepath.Join("sdk-spec", "v1.json")
	}
	b, err := os.ReadFile(p)
	if err != nil {
		t.Skip("spec file not found: " + p)
	}
	var obj map[string]any
	if err := json.Unmarshal(b, &obj); err != nil {
		t.Fatalf("invalid spec json: %v", err)
	}
	return obj
}

func TestSpecProxyPorts(t *testing.T) {
	spec := loadSpec(t)
	proxy := spec["proxy"].(map[string]any)
	products := proxy["products"].(map[string]any)

	res := products["residential"].(map[string]any)["port"].(float64)
	mob := products["mobile"].(map[string]any)["port"].(float64)
	dc := products["datacenter"].(map[string]any)["port"].(float64)
	isp := products["isp"].(map[string]any)["port"].(float64)

	if int(res) != 9999 || int(mob) != 5555 || int(dc) != 7777 || int(isp) != 6666 {
		t.Fatalf("unexpected ports from spec: %v", products)
	}
}

func TestSpecSerpMapping(t *testing.T) {
	spec := loadSpec(t)
	serp := spec["serp"].(map[string]any)
	mappings := serp["mappings"].(map[string]any)
	tbms := mappings["searchTypeToTbm"].(map[string]any)

	if tbms["news"].(string) != "nws" {
		t.Fatalf("expected news -> nws, got %v", tbms["news"])
	}
}
