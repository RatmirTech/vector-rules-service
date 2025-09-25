package usecase

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/ratmirtech/vector-rules-service/internal/domain"
	"github.com/ratmirtech/vector-rules-service/internal/infra/embeddings"
)

type ruleService struct {
	ruleRepo       domain.RuleRepository
	ruleTypeRepo   domain.RuleTypeRepository
	embeddingProvider domain.EmbeddingProvider
}

// NewRuleService creates a new rule service
func NewRuleService(
	ruleRepo domain.RuleRepository,
	ruleTypeRepo domain.RuleTypeRepository,
	embeddingProvider domain.EmbeddingProvider,
) domain.RuleService {
	return &ruleService{
		ruleRepo:          ruleRepo,
		ruleTypeRepo:      ruleTypeRepo,
		embeddingProvider: embeddingProvider,
	}
}

func (s *ruleService) RetrieveSimilar(ctx context.Context, query *domain.RetrieveRulesQuery) ([]*domain.RuleMatch, error) {
	// Generate embeddings for all queries
	embeds, err := s.embeddingProvider.GenerateBatchEmbeddings(ctx, query.Queries)
	if err != nil {
		return nil, fmt.Errorf("failed to generate embeddings: %w", err)
	}

	// Average all query embeddings into a single embedding
	avgEmbedding, err := embeddings.AverageEmbeddings(embeds)
	if err != nil {
		return nil, fmt.Errorf("failed to average embeddings: %w", err)
	}

	// Find similar rules using the averaged embedding
	matches, err := s.ruleRepo.FindSimilar(ctx, avgEmbedding, query.Type, query.N)
	if err != nil {
		return nil, fmt.Errorf("failed to find similar rules: %w", err)
	}

	return matches, nil
}

func (s *ruleService) CreateRule(ctx context.Context, req *domain.CreateRuleRequest) (*domain.Rule, error) {
	// Validate and get rule type
	ruleType, err := s.ruleTypeRepo.GetByName(ctx, req.Type)
	if err != nil {
		return nil, fmt.Errorf("invalid rule type '%s': %w", req.Type, err)
	}

	// Generate embedding for the content
	contentStr := string(req.Content)
	embedding, err := s.embeddingProvider.GenerateEmbedding(ctx, contentStr)
	if err != nil {
		return nil, fmt.Errorf("failed to generate embedding: %w", err)
	}

	// Create rule
	rule := &domain.Rule{
		RuleTypeID: ruleType.ID,
		Content:    req.Content,
		Embedding:  embedding,
	}

	createdRule, err := s.ruleRepo.Create(ctx, rule)
	if err != nil {
		return nil, fmt.Errorf("failed to create rule: %w", err)
	}

	return createdRule, nil
}

func (s *ruleService) GetRule(ctx context.Context, id int64) (*domain.Rule, error) {
	rule, err := s.ruleRepo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get rule: %w", err)
	}
	return rule, nil
}

func (s *ruleService) UpdateRule(ctx context.Context, req *domain.UpdateRuleRequest) (*domain.Rule, error) {
	// Validate and get rule type
	ruleType, err := s.ruleTypeRepo.GetByName(ctx, req.Type)
	if err != nil {
		return nil, fmt.Errorf("invalid rule type '%s': %w", req.Type, err)
	}

	// Check if rule exists
	existingRule, err := s.ruleRepo.GetByID(ctx, req.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to get existing rule: %w", err)
	}

	// Generate new embedding for updated content
	contentStr := string(req.Content)
	embedding, err := s.embeddingProvider.GenerateEmbedding(ctx, contentStr)
	if err != nil {
		return nil, fmt.Errorf("failed to generate embedding: %w", err)
	}

	// Update rule
	existingRule.RuleTypeID = ruleType.ID
	existingRule.Content = req.Content
	existingRule.Embedding = embedding

	updatedRule, err := s.ruleRepo.Update(ctx, existingRule)
	if err != nil {
		return nil, fmt.Errorf("failed to update rule: %w", err)
	}

	return updatedRule, nil
}

func (s *ruleService) DeleteRule(ctx context.Context, id int64) error {
	err := s.ruleRepo.Delete(ctx, id)
	if err != nil {
		return fmt.Errorf("failed to delete rule: %w", err)
	}
	return nil
}

func (s *ruleService) ListRules(ctx context.Context, ruleType *string, limit, offset int) ([]*domain.Rule, error) {
	rules, err := s.ruleRepo.List(ctx, ruleType, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to list rules: %w", err)
	}
	return rules, nil
}