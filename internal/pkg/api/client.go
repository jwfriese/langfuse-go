package api

import (
	"context"
	"encoding/base64"
	"net/http"
	"os"

	"github.com/henomis/restclientgo"
)

const (
	langfuseDefaultEndpoint = "https://cloud.langfuse.com"
)

type Client struct {
	restClient *restclientgo.RestClient
}

// New creates a new LangFuse client with the default endpoint and credentials from the environment.
// It reads the LANGFUSE_HOST, LANGFUSE_PUBLIC_KEY, and LANGFUSE_SECRET_KEY environment variables
// to configure the client. If the LANGFUSE_HOST environment variable is not set, it falls back to
// a default endpoint.
func New() *Client {
	langfuseHost := os.Getenv("LANGFUSE_HOST")
	if langfuseHost == "" {
		langfuseHost = langfuseDefaultEndpoint
	}

	publicKey := os.Getenv("LANGFUSE_PUBLIC_KEY")
	secretKey := os.Getenv("LANGFUSE_SECRET_KEY")

	return NewFromConfig(Config{
		Host:      langfuseHost,
		PublicKey: publicKey,
		SecretKey: secretKey,
	})
}

type Config struct {
	Host      string
	PublicKey string
	SecretKey string
}

// NewFromConfig creates a new LangFuse client with the given configuration.
func NewFromConfig(cfg Config) *Client {
	restClient := restclientgo.New(cfg.Host)
	restClient.SetRequestModifier(func(req *http.Request) *http.Request {
		req.Header.Set("Authorization", basicAuth(cfg.PublicKey, cfg.SecretKey))
		return req
	})

	return &Client{
		restClient: restClient,
	}
}

func (c *Client) Ingestion(ctx context.Context, req *Ingestion, res *IngestionResponse) error {
	return c.restClient.Post(ctx, req, res)
}

func basicAuth(publicKey, secretKey string) string {
	auth := publicKey + ":" + secretKey
	return "Basic " + base64.StdEncoding.EncodeToString([]byte(auth))
}
