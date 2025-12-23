package tests

import (
	"encoding/json"
	"os"
	"path/filepath"
	"runtime"
	"testing"
)

func loadSpec(t *testing.T) map[string]any {
	t.Helper()

	// 1) Allow explicit override
	if p := os.Getenv("THORDATA_SDK_SPEC_PATH"); p != "" {
		return readSpecFile(t, p)
	}

	// 2) Resolve repo root based on this test file location
	_, thisFile, _, ok := runtime.Caller(0)
	if !ok {
		t.Skip("unable to resolve current file path")
	}
	repoRoot := filepath.Clean(filepath.Join(filepath.Dir(thisFile), ".."))
	specPath := filepath.Join(repoRoot, "sdk-spec", "v1.json")

	if _, err := os.Stat(specPath); err != nil {
		t.Skip("spec file not found: " + specPath + " (did you init submodules?)")
	}
	return readSpecFile(t, specPath)
}

func readSpecFile(t *testing.T, path string) map[string]any {
	t.Helper()

	b, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("failed to read spec: %v", err)
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

	res := int(products["residential"].(map[string]any)["port"].(float64))
	mob := int(products["mobile"].(map[string]any)["port"].(float64))
	dc := int(products["datacenter"].(map[string]any)["port"].(float64))
	isp := int(products["isp"].(map[string]any)["port"].(float64))

	if res != 9999 || mob != 5555 || dc != 7777 || isp != 6666 {
		t.Fatalf("unexpected ports from spec: %+v", products)
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
