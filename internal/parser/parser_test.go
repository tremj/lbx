package parser

import (
	"strings"
	"testing"
)

func TestValidYAML(t *testing.T) {
	config := []byte(`
name: my-lb
description: Testing ts

listeners:
  - name: my-http
    protocol: http
    port: 80

backends:
  - name: backend1
    port: 8080
  - name: backend2
    port: 8080
`)
	err := ValidateConfig(config)
	if err != nil {
		t.Fatalf("Expected 0 errors, got %v", err)
	}
}

func TestUnknownYAMLField(t *testing.T) {
	config := []byte(`
name: my-lb
description: Testing ts

hello: hiiii

listeners:
  - name: my-http
    protocol: http
    port: 80

backends:
  - name: backend1
    port: 8080
  - name: backend2
    port: 8080
`)
	err := ValidateConfig(config)
	expectedSection := "field hello not found in type types.Config"
	if !strings.Contains(err.Error(), expectedSection) {
		t.Fatalf("Unexpected error: %v", err)
	}
}

func TestMissingNameDescription(t *testing.T) {
	config := []byte(`
listeners:
  - name: my-http
    protocol: http
    port: 80

backends:
  - name: backend1
    port: 8080
  - name: backend2
    port: 8080
`)
	err := ValidateConfig(config)
	expectedMsg := []string{"missing name of your LB configuration", "missing description of your LB configuration"}
	for _, msg := range expectedMsg {
		if !strings.Contains(err.Error(), msg) {
			t.Fatalf("Did not find %s in: \n%v", msg, err)
		}
	}
}

func TestMissingListener(t *testing.T) {
	config := []byte(`
name: my-lb
description: Testing ts

backends:
  - name: backend1
    port: 8080
  - name: backend2
    port: 8080
`)
	err := ValidateConfig(config)
	expectedSection := "missing or empty listener config"
	if !strings.Contains(err.Error(), expectedSection) {
		t.Fatalf("Unexpected error: %v", err)
	}
}

func TestEmptyListener(t *testing.T) {
	config := []byte(`
name: my-lb
description: Testing ts

listeners:

backends:
  - name: backend1
    port: 8080
  - name: backend2
    port: 8080
`)
	err := ValidateConfig(config)
	expectedMsg := []string{"missing or empty listener config"}
	for _, msg := range expectedMsg {
		if !strings.Contains(err.Error(), msg) {
			t.Fatalf("Did not find %s in: \n%v", msg, err)
		}
	}
}

func TestListenerMissingNamePort(t *testing.T) {
	config := []byte(`
name: my-lb
description: Testing ts

listeners:
  - protocol: http

backends:
  - name: backend1
    port: 8080
  - name: backend2
    port: 8080
`)
	err := ValidateConfig(config)
	expectedMsg := []string{"listener 1: Missing name", "listener 1: Missing port number"}
	for _, msg := range expectedMsg {
		if !strings.Contains(err.Error(), msg) {
			t.Fatalf("Did not find %s in: \n%v", msg, err)
		}
	}
}

func TestListenerMissingProtocol(t *testing.T) {
	config := []byte(`
name: my-lb
description: Testing ts

listeners:
  - name: my-http
    port: 80

backends:
  - name: backend1
    port: 8080
  - name: backend2
    port: 8080
`)
	err := ValidateConfig(config)
	expectedMsg := []string{"listener 1: Missing protocol"}
	for _, msg := range expectedMsg {
		if !strings.Contains(err.Error(), msg) {
			t.Fatalf("Did not find %s in: \n%v", msg, err)
		}
	}
}

func TestListenerWrongHTTPPort(t *testing.T) {
	config := []byte(`
name: my-lb
description: Testing ts

listeners:
  - name: my-http
    protocol: http
    port: 81

backends:
  - name: backend1
    port: 8080
  - name: backend2
    port: 8080
`)
	err := ValidateConfig(config)
	expectedMsg := []string{"listener 1: Invalid port 81 for protocol \"http\""}
	for _, msg := range expectedMsg {
		if !strings.Contains(err.Error(), msg) {
			t.Fatalf("Did not find %s in: \n%v", msg, err)
		}
	}
}

func TestListenerWrongHTTPSPortAndMissingTLS(t *testing.T) {
	config := []byte(`
name: my-lb
description: Testing ts

listeners:
  - name: my-https
    protocol: https
    port: 81

backends:
  - name: backend1
    port: 8080
  - name: backend2
    port: 8080
`)
	err := ValidateConfig(config)
	expectedMsg := []string{"listener 1: Invalid port 81 for protocol \"https\"", "listener 1: Missing certificate information"}
	for _, msg := range expectedMsg {
		if !strings.Contains(err.Error(), msg) {
			t.Fatalf("Did not find %s in: \n%v", msg, err)
		}
	}
}

func TestListenerInvalidProtocol(t *testing.T) {
	config := []byte(`
name: my-lb
description: Testing ts

listeners:
  - name: my-ssh
    protocol: ssh
    port: 81

backends:
  - name: backend1
    port: 8080
  - name: backend2
    port: 8080
`)
	err := ValidateConfig(config)
	expectedMsg := []string{"listener 1: Invalid protocol \"ssh\""}
	for _, msg := range expectedMsg {
		if !strings.Contains(err.Error(), msg) {
			t.Fatalf("Did not find %s in: \n%v", msg, err)
		}
	}
}

func TestListenerValidTLS(t *testing.T) {
	config := []byte(`
name: my-lb
description: Testing ts

listeners:
  - name: my-https
    protocol: https
    port: 443
    tls_cert: "/path/to/cert.pem"
    tls_key: "/path/to/key.pem"

backends:
  - name: backend1
    port: 8080
  - name: backend2
    port: 8080
`)
	err := ValidateConfig(config)
	if err != nil {
		t.Fatalf("Expected 0 errors, got %v", err)
	}
}

func TestMissingBackend(t *testing.T) {
	config := []byte(`
name: my-lb
description: Testing ts

listeners:
  - name: my-http
    protocol: http
    port: 80
`)
	err := ValidateConfig(config)
	expectedMsg := []string{"missing or empty backend config"}
	for _, msg := range expectedMsg {
		if !strings.Contains(err.Error(), msg) {
			t.Fatalf("Did not find %s in: \n%v", msg, err)
		}
	}
}

func TestEmptyBackend(t *testing.T) {
	config := []byte(`
name: my-lb
description: Testing ts

listeners:
  - name: my-http
    protocol: http
    port: 80

backends:
`)
	err := ValidateConfig(config)
	expectedMsg := []string{"missing or empty backend config"}
	for _, msg := range expectedMsg {
		if !strings.Contains(err.Error(), msg) {
			t.Fatalf("Did not find %s in: \n%v", msg, err)
		}
	}
}

func TestBackendMissingNamePort(t *testing.T) {
	config := []byte(`
name: my-lb
description: Testing ts

listeners:
  - name: my-http
    protocol: http
    port: 80

backends:
  - port: 8080
  - name: backend2
`)
	err := ValidateConfig(config)
	expectedMsg := []string{"backend 1: Missing name", "backend 2: Missing port"}
	for _, msg := range expectedMsg {
		if !strings.Contains(err.Error(), msg) {
			t.Fatalf("Did not find %s in: \n%v", msg, err)
		}
	}
}

func TestBackendInvalidPort(t *testing.T) {
	config := []byte(`
name: my-lb
description: Testing ts

listeners:
  - name: my-http
    protocol: http
    port: 80

backends:
  - name: backend9
    port: 100000
`)
	err := ValidateConfig(config)
	expectedMsg := []string{"backend 1: Invalid port 100000"}
	for _, msg := range expectedMsg {
		if !strings.Contains(err.Error(), msg) {
			t.Fatalf("Did not find %s in: \n%v", msg, err)
		}
	}
}
