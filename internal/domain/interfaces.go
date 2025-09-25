package domain

import (
	"context"
	"errors"
)

var (
	ErrRuleNotFound     = errors.New("rule not found")
	ErrRuleTypeNotFound = errors.New("rule type not found")
	ErrInvalidInput     = errors.New("invalid input")
	ErrDuplicateEntry   = errors.New("duplicate entry")
)

// RuleRepository defines the interface for rule data access
type RuleRepository interface {
	// Create creates a new rule
	Create(ctx context.Context, rule *Rule) (*Rule, error)
	
	// GetByID retrieves a rule by ID
	GetByID(ctx context.Context, id int64) (*Rule, error)
	
	// Update updates an existing rule
	Update(ctx context.Context, rule *Rule) (*Rule, error)
	
	// Delete deletes a rule by ID
	Delete(ctx context.Context, id int64) error
	
	// List retrieves rules with optional type filter
	List(ctx context.Context, ruleType *string, limit, offset int) ([]*Rule, error)
	
	// FindSimilar finds rules similar to the given embedding
	FindSimilar(ctx context.Context, embedding []float32, ruleType *string, limit int) ([]*RuleMatch, error)
	
	// UpdateEmbedding updates the embedding of a rule
	UpdateEmbedding(ctx context.Context, id int64, embedding []float32) error
}

// RuleTypeRepository defines the interface for rule type data access
type RuleTypeRepository interface {
	// Create creates a new rule type
	Create(ctx context.Context, ruleType *RuleType) (*RuleType, error)
	
	// GetByID retrieves a rule type by ID
	GetByID(ctx context.Context, id int64) (*RuleType, error)
	
	// GetByName retrieves a rule type by name
	GetByName(ctx context.Context, name string) (*RuleType, error)
	
	// Update updates an existing rule type
	Update(ctx context.Context, ruleType *RuleType) (*RuleType, error)
	
	// Delete deletes a rule type by ID
	Delete(ctx context.Context, id int64) error
	
	// List retrieves all rule types
	List(ctx context.Context, limit, offset int) ([]*RuleType, error)
}

// EmbeddingProvider defines the interface for generating embeddings
type EmbeddingProvider interface {
	// GenerateEmbedding generates an embedding for the given text
	GenerateEmbedding(ctx context.Context, text string) ([]float32, error)
	
	// GenerateBatchEmbeddings generates embeddings for multiple texts
	GenerateBatchEmbeddings(ctx context.Context, texts []string) ([][]float32, error)
}

// RuleService defines business logic operations for rules
type RuleService interface {
	// RetrieveSimilar retrieves rules similar to the given queries
	RetrieveSimilar(ctx context.Context, query *RetrieveRulesQuery) ([]*RuleMatch, error)
	
	// CreateRule creates a new rule
	CreateRule(ctx context.Context, req *CreateRuleRequest) (*Rule, error)
	
	// GetRule retrieves a rule by ID
	GetRule(ctx context.Context, id int64) (*Rule, error)
	
	// UpdateRule updates an existing rule
	UpdateRule(ctx context.Context, req *UpdateRuleRequest) (*Rule, error)
	
	// DeleteRule deletes a rule by ID
	DeleteRule(ctx context.Context, id int64) error
	
	// ListRules retrieves rules with optional filters
	ListRules(ctx context.Context, ruleType *string, limit, offset int) ([]*Rule, error)
}

// RuleTypeService defines business logic operations for rule types
type RuleTypeService interface {
	// CreateRuleType creates a new rule type
	CreateRuleType(ctx context.Context, req *CreateRuleTypeRequest) (*RuleType, error)
	
	// GetRuleType retrieves a rule type by ID
	GetRuleType(ctx context.Context, id int64) (*RuleType, error)
	
	// UpdateRuleType updates an existing rule type
	UpdateRuleType(ctx context.Context, req *UpdateRuleTypeRequest) (*RuleType, error)
	
	// DeleteRuleType deletes a rule type by ID
	DeleteRuleType(ctx context.Context, id int64) error
	
	// ListRuleTypes retrieves all rule types
	ListRuleTypes(ctx context.Context, limit, offset int) ([]*RuleType, error)
}