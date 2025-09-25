package domain

import (
	"encoding/json"
	"time"
)

// RuleType represents a category of rules
type RuleType struct {
	ID        int64     `json:"id"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// Rule represents a business rule with vector embedding
type Rule struct {
	ID         int64           `json:"id"`
	RuleTypeID int64           `json:"rule_type_id"`
	Content    json.RawMessage `json:"content"`
	Embedding  []float32       `json:"-"` // Vector embedding for similarity search
	CreatedAt  time.Time       `json:"created_at"`
	UpdatedAt  time.Time       `json:"updated_at"`
	
	// Populated from join
	RuleTypeName *string `json:"rule_type_name,omitempty"`
}

// RuleMatch represents a rule with similarity score
type RuleMatch struct {
	Rule
	Score float64 `json:"score"`
}

// CreateRuleRequest represents request to create a rule
type CreateRuleRequest struct {
	Type    string          `json:"type" validate:"required"`
	Content json.RawMessage `json:"content" validate:"required"`
}

// UpdateRuleRequest represents request to update a rule
type UpdateRuleRequest struct {
	ID      int64           `json:"id" validate:"required"`
	Type    string          `json:"type" validate:"required"`
	Content json.RawMessage `json:"content" validate:"required"`
}

// CreateRuleTypeRequest represents request to create a rule type
type CreateRuleTypeRequest struct {
	Name string `json:"name" validate:"required"`
}

// UpdateRuleTypeRequest represents request to update a rule type
type UpdateRuleTypeRequest struct {
	ID   int64  `json:"id" validate:"required"`
	Name string `json:"name" validate:"required"`
}

// RetrieveRulesQuery represents query parameters for rule retrieval
type RetrieveRulesQuery struct {
	N       int      `json:"n" validate:"required,min=1,max=100"`
	Type    *string  `json:"type,omitempty"`
	Queries []string `json:"queries" validate:"required,min=1"`
}