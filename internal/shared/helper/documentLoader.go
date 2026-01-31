package helper

import (
	"embed"
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/iden3/go-schema-processor/v2/loaders"
	"github.com/iden3/go-schema-processor/v2/verifiable"
	"github.com/piprate/json-gold/ld"
)

// Embed context files vào binary
//

//go:embed context/*.jsonld
var contextFS embed.FS

// cachedDocumentLoader implements caching với pre-populated contexts
type cachedDocumentLoader struct {
	cache      map[string]*ld.RemoteDocument
	cacheMutex sync.RWMutex
	baseLoader ld.DocumentLoader
}

func (c *cachedDocumentLoader) LoadDocument(url string) (*ld.RemoteDocument, error) {
	// Try get from cache first
	c.cacheMutex.RLock()
	if doc, exists := c.cache[url]; exists {
		c.cacheMutex.RUnlock()
		return doc, nil
	}
	c.cacheMutex.RUnlock()

	// Load from base loader
	remoteDoc, err := c.baseLoader.LoadDocument(url)
	if err != nil {
		return nil, fmt.Errorf("failed to load document from %s: %w", url, err)
	}

	// Cache the result
	c.cacheMutex.Lock()
	c.cache[url] = remoteDoc
	c.cacheMutex.Unlock()

	return remoteDoc, nil
}

// CacheLoaderOptions options for cache loader
type CacheLoaderOptions struct {
	HTTPClient  *http.Client
	IPFSClient  loaders.IPFSClient
	IPFSGateway string
}

// userAgentTransport adds User-Agent header to all requests
type userAgentTransport struct {
	Transport http.RoundTripper
	UserAgent string
}

func (t *userAgentTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	req.Header.Set("User-Agent", t.UserAgent)
	req.Header.Set("Accept", "application/ld+json, application/json")
	return t.Transport.RoundTrip(req)
}

// createDefaultHTTPClient creates HTTP client with User-Agent
func createDefaultHTTPClient() *http.Client {
	return &http.Client{
		Timeout: 30 * time.Second,
		Transport: &userAgentTransport{
			Transport: http.DefaultTransport,
			UserAgent: "Mozilla/5.0 (compatible; ZKCredentialApp/1.0)",
		},
	}
}

// loadJSONLDContext reads and parses JSON-LD context file
func loadJSONLDContext(filename string) (interface{}, error) {
	data, err := contextFS.ReadFile("context/" + filename)
	if err != nil {
		return nil, fmt.Errorf("failed to read context file %s: %w", filename, err)
	}

	var context interface{}
	if err := json.Unmarshal(data, &context); err != nil {
		return nil, fmt.Errorf("failed to parse context file %s: %w", filename, err)
	}

	return context, nil
}

// NewCacheLoader creates a document loader with pre-populated cache
func NewCacheLoader(opts *CacheLoaderOptions) ld.DocumentLoader {
	if opts == nil {
		opts = &CacheLoaderOptions{}
	}

	// Ensure we have HTTP client with User-Agent
	if opts.HTTPClient == nil {
		opts.HTTPClient = createDefaultHTTPClient()
	}

	// Create base loader with HTTP client
	loaderOpts := []loaders.DocumentLoaderOption{
		loaders.WithHTTPClient(opts.HTTPClient),
	}

	baseLoader := loaders.NewDocumentLoader(
		opts.IPFSClient,
		opts.IPFSGateway,
		loaderOpts...,
	)

	// Initialize cache with embedded contexts
	cache := make(map[string]*ld.RemoteDocument)

	// Load and cache W3C Credential 2018 context
	w3cContext, err := loadJSONLDContext("W3CCredential2018.jsonld")
	if err != nil {
		fmt.Printf("Warning: failed to load W3C context: %v\n", err)
	} else {
		cache[verifiable.JSONLDSchemaW3CCredential2018] = &ld.RemoteDocument{
			DocumentURL: verifiable.JSONLDSchemaW3CCredential2018,
			Document:    w3cContext, // Parsed JSON, not bytes!
		}
	}

	// Load and cache Iden3 Proofs context
	iden3ProofsContext, err := loadJSONLDContext("Iden3Proofs.jsonld")
	if err != nil {
		fmt.Printf("Warning: failed to load Iden3 Proofs context: %v\n", err)
	} else {
		cache[verifiable.JSONLDSchemaIden3Credential] = &ld.RemoteDocument{
			DocumentURL: verifiable.JSONLDSchemaIden3Credential,
			Document:    iden3ProofsContext,
		}
	}

	// Load and cache Iden3 Display Method context
	iden3DisplayContext, err := loadJSONLDContext("Iden3DisplayMethod.jsonld")
	if err != nil {
		fmt.Printf("Warning: failed to load Iden3 Display Method context: %v\n", err)
	} else {
		cache[verifiable.JSONLDSchemaIden3DisplayMethod] = &ld.RemoteDocument{
			DocumentURL: verifiable.JSONLDSchemaIden3DisplayMethod,
			Document:    iden3DisplayContext,
		}
	}

	return &cachedDocumentLoader{
		cache:      cache,
		baseLoader: baseLoader,
	}
}
