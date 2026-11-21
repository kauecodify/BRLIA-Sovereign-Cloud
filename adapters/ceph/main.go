package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"
)

const (
	adapterName  = "ceph"
	endpointEnv  = "ceph_ENDPOINT"
	probeTimeout = 2 * time.Second
)

type adapterStatus struct {
	Adapter   string `json:"adapter"`
	Endpoint  string `json:"endpoint"`
	Status    string `json:"status"`
	Detail    string `json:"detail,omitempty"`
	CheckedAt string `json:"checked_at"`
}

func probe(ctx context.Context, endpoint string) adapterStatus {
	status := adapterStatus{
		Adapter:   adapterName,
		Endpoint:  endpoint,
		CheckedAt: time.Now().UTC().Format(time.RFC3339),
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, endpoint, nil)
	if err != nil {
		status.Status = "error"
		status.Detail = err.Error()
		return status
	}

	client := &http.Client{Timeout: probeTimeout}
	resp, err := client.Do(req)
	if err != nil {
		status.Status = "unreachable"
		status.Detail = err.Error()
		return status
	}
	defer resp.Body.Close()

	if resp.StatusCode >= http.StatusOK && resp.StatusCode < http.StatusInternalServerError {
		status.Status = "reachable"
		status.Detail = resp.Status
	} else {
		status.Status = "degraded"
		status.Detail = resp.Status
	}
	return status
}

func main() {
	endpoint := os.Getenv(endpointEnv)
	if endpoint == "" {
		status := adapterStatus{
			Adapter:   adapterName,
			Status:    "unconfigured",
			Detail:    fmt.Sprintf("set %s to enable probes", endpointEnv),
			CheckedAt: time.Now().UTC().Format(time.RFC3339),
		}
		print(status)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), probeTimeout*2)
	defer cancel()

	print(probe(ctx, endpoint))
}

func print(status adapterStatus) {
	enc := json.NewEncoder(os.Stdout)
	enc.SetIndent("", "  ")
	_ = enc.Encode(status)
}
