package client

import (
	"context"
	"net"
	"net/http"
	"time"
)

var (
	apiVersion   = "v1.24"
	dockerSocket = "/var/run/docker.sock"
	baseURL      = "dockerhost"
	unixType     = "unix"
)

// APIClient contains the necessary configuration for a client to interact
// with the Docker API
type APIClient struct {
	NetworkType string
	Version     string
	Location    string
	Host        string
	Client      *http.Client
}

// New generates a new instance of APIClient
func New(netType, host string) (c *APIClient, err error) {
	c = &APIClient{
		NetworkType: netType,
		Version:     apiVersion,
		Host:        host,
	}

	return
}

// NewDefault generates a new instance of APIClient with sensible defaults
func NewDefault() (c *APIClient, err error) {
	c = &APIClient{
		NetworkType: unixType,
		Version:     apiVersion,
		Location:    dockerSocket,
		Host:        baseURL,
		Client: &http.Client{
			Transport: &http.Transport{
				DialContext: func(_ context.Context, _, _ string) (net.Conn, error) {
					return net.Dial(unixType, dockerSocket)
				},
				TLSHandshakeTimeout: 10 * time.Second,
			},
		},
	}

	return
}
