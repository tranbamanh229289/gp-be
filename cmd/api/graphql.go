package main

import (
	"context"
	"graduate-project/internal/transport/graph"
	"log"
	"net/http"
	"os"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/handler/extension"
	"github.com/99designs/gqlgen/graphql/handler/transport"
	"github.com/99designs/gqlgen/graphql/playground"
	lru "github.com/hashicorp/golang-lru/v2"
	"github.com/vektah/gqlparser/v2/ast"
)

const defaultPort = "8080"

// LRUCache wraps golang-lru cache to implement graphql.Cache interface
type LRUCache[K comparable, V any] struct {
	cache *lru.Cache[K, V]
}

func NewLRUCache[K comparable, V any](size int) (*LRUCache[K, V], error) {
	cache, err := lru.New[K, V](size)
	if err != nil {
		return nil, err
	}
	return &LRUCache[K, V]{cache: cache}, nil
}

func (c *LRUCache[K, V]) Get(ctx context.Context, key K) (V, bool) {
	return c.cache.Get(key)
}

func (c *LRUCache[K, V]) Add(ctx context.Context, key K, value V) {
	c.cache.Add(key, value)
}

func gql() {
	port := os.Getenv("PORT")
	if port == "" {
		port = defaultPort
	}

	// create graphql server
	srv := handler.New(graph.NewExecutableSchema(graph.Config{Resolvers: &graph.Resolver{}}))

	// add transport http
	srv.AddTransport(transport.Options{})
	srv.AddTransport(transport.GET{})
	srv.AddTransport(transport.POST{})

	// Create LRU cache for query documents
	queryCache, err := NewLRUCache[string, *ast.QueryDocument](1000)
	if err != nil {
		log.Fatal("Failed to create query cache:", err)
	}
	srv.SetQueryCache(queryCache)

	srv.Use(extension.Introspection{})

	// Create LRU cache for APQ
	apqCache, err := NewLRUCache[string, string](100)
	if err != nil {
		log.Fatal("Failed to create APQ cache:", err)
	}
	srv.Use(extension.AutomaticPersistedQuery{
		Cache: apqCache,
	})

	// Playground UI
	http.Handle("/", playground.Handler("GraphQL playground", "/query"))
	http.Handle("/query", srv)

	log.Printf("connect to http://localhost:%s/ for GraphQL playground", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}