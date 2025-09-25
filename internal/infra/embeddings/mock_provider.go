package embeddings

import (
	"context"
	"fmt"
	"math"
	"math/rand"
	"strings"

	"github.com/ratmirtech/vector-rules-service/internal/domain"
)

// mockEmbeddingProvider is a stub implementation of EmbeddingProvider for development/testing
type mockEmbeddingProvider struct {
	dimensions int
}

// NewMockEmbeddingProvider creates a new mock embedding provider
// This is a stub implementation - replace with real embedding service (OpenAI, Cohere, etc.)
func NewMockEmbeddingProvider(dimensions int) domain.EmbeddingProvider {
	if dimensions <= 0 {
		dimensions = 1536 // Default to OpenAI ada-002 dimensions
	}
	return &mockEmbeddingProvider{
		dimensions: dimensions,
	}
}

// GenerateEmbedding generates a mock embedding based on text content
// This is a simple hash-based approach for demonstration
// TODO: Replace with real embedding service (OpenAI API, Sentence Transformers, etc.)
func (m *mockEmbeddingProvider) GenerateEmbedding(ctx context.Context, text string) ([]float32, error) {
	if text == "" {
		return nil, fmt.Errorf("text cannot be empty")
	}

	// Create a deterministic "embedding" based on text content
	embedding := make([]float32, m.dimensions)
	
	// Use text hash to create consistent embeddings for same text
	textLower := strings.ToLower(text)
	seed := int64(0)
	for _, char := range textLower {
		seed += int64(char)
	}
	
	// Create a seeded random generator for consistency
	rng := rand.New(rand.NewSource(seed))
	
	// Generate normalized vector
	var norm float64
	for i := 0; i < m.dimensions; i++ {
		val := rng.Float64()*2.0 - 1.0 // Random value between -1 and 1
		embedding[i] = float32(val)
		norm += val * val
	}
	
	// Normalize the vector to unit length
	norm = math.Sqrt(norm)
	for i := 0; i < m.dimensions; i++ {
		embedding[i] = embedding[i] / float32(norm)
	}
	
	return embedding, nil
}

// GenerateBatchEmbeddings generates embeddings for multiple texts
func (m *mockEmbeddingProvider) GenerateBatchEmbeddings(ctx context.Context, texts []string) ([][]float32, error) {
	if len(texts) == 0 {
		return nil, fmt.Errorf("texts slice cannot be empty")
	}

	embeddings := make([][]float32, len(texts))
	for i, text := range texts {
		embedding, err := m.GenerateEmbedding(ctx, text)
		if err != nil {
			return nil, fmt.Errorf("failed to generate embedding for text %d: %w", i, err)
		}
		embeddings[i] = embedding
	}
	
	return embeddings, nil
}

// AverageEmbeddings computes the average of multiple embeddings
// This is a utility function for aggregating multiple query embeddings
func AverageEmbeddings(embeddings [][]float32) ([]float32, error) {
	if len(embeddings) == 0 {
		return nil, fmt.Errorf("embeddings slice cannot be empty")
	}

	dimensions := len(embeddings[0])
	if dimensions == 0 {
		return nil, fmt.Errorf("embeddings cannot be empty")
	}

	// Validate all embeddings have same dimensions
	for i, emb := range embeddings {
		if len(emb) != dimensions {
			return nil, fmt.Errorf("embedding %d has different dimensions: expected %d, got %d", i, dimensions, len(emb))
		}
	}

	// Compute average
	avgEmbedding := make([]float32, dimensions)
	for _, embedding := range embeddings {
		for i, val := range embedding {
			avgEmbedding[i] += val
		}
	}

	// Normalize by count
	count := float32(len(embeddings))
	for i := range avgEmbedding {
		avgEmbedding[i] /= count
	}

	return avgEmbedding, nil
}