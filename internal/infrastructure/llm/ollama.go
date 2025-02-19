package llm

import (
	"context"
	"log"
	"net"
	"net/http"
	"net/url"
	"time"

	"github.com/cnc-csku/task-nexus-task-management/config"
)

type OllamaClient struct {
	HTTPClient *http.Client
}

func NewOllamaClient(ctx context.Context, cfg *config.Config) *OllamaClient {
	var httpClient *http.Client

	if cfg.OllamaClient.UseProxy {
		proxyUrl := &url.URL{
			Scheme: "http",
			Host:   cfg.OllamaClient.HttpProxyHost + ":" + cfg.OllamaClient.HttpProxyPort,
		}

		transport := &http.Transport{
			Proxy:       http.ProxyURL(proxyUrl),
			DialContext: (&net.Dialer{Timeout: 30 * time.Second}).DialContext,
		}
		httpClient = &http.Client{
			Transport: transport,
		}
	} else {
		httpClient = &http.Client{
			Transport: http.DefaultTransport,
		}
	}

	_, err := httpClient.Get("http://" + cfg.OllamaClient.Endpoint)
	if err != nil {
		panic("🚫 Failed to connect to Ollama: " + err.Error())
	}

	log.Println("🦙 Connected to Ollama")
	return &OllamaClient{HTTPClient: httpClient}
}
