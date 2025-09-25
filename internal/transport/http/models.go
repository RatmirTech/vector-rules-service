package http

import (
	"time"
)

// SwaggerRule represents a rule for Swagger documentation
type SwaggerRule struct {
	ID           int64     `json:"id" example:"1"`
	RuleTypeID   int64     `json:"rule_type_id" example:"1"`
	Content      string    `json:"content" example:"{\"description\":\"Sample rule content\"}"`
	CreatedAt    time.Time `json:"created_at" example:"2023-01-01T00:00:00Z"`
	UpdatedAt    time.Time `json:"updated_at" example:"2023-01-01T00:00:00Z"`
	RuleTypeName *string   `json:"rule_type_name,omitempty" example:"security"`
}

// SwaggerRuleType represents a rule type for Swagger documentation
type SwaggerRuleType struct {
	ID        int64     `json:"id" example:"1"`
	Name      string    `json:"name" example:"security"`
	CreatedAt time.Time `json:"created_at" example:"2023-01-01T00:00:00Z"`
	UpdatedAt time.Time `json:"updated_at" example:"2023-01-01T00:00:00Z"`
}

// SwaggerCreateRuleRequest represents a create rule request for Swagger documentation
type SwaggerCreateRuleRequest struct {
	Type    string `json:"type" example:"security" validate:"required"`
	Content string `json:"content" example:"{\"description\":\"Sample rule content\"}" validate:"required"`
}

// SwaggerUpdateRuleRequest represents an update rule request for Swagger documentation
type SwaggerUpdateRuleRequest struct {
	ID      int64  `json:"id" example:"1" validate:"required"`
	Type    string `json:"type" example:"security" validate:"required"`
	Content string `json:"content" example:"{\"description\":\"Updated rule content\"}" validate:"required"`
}

// SwaggerCreateRuleTypeRequest represents a create rule type request for Swagger documentation
type SwaggerCreateRuleTypeRequest struct {
	Name string `json:"name" example:"security" validate:"required"`
}

// SwaggerUpdateRuleTypeRequest represents an update rule type request for Swagger documentation
type SwaggerUpdateRuleTypeRequest struct {
	ID   int64  `json:"id" example:"1" validate:"required"`
	Name string `json:"name" example:"updated-security" validate:"required"`
}

// SwaggerRuleMatch represents a rule match for Swagger documentation
type SwaggerRuleMatch struct {
	SwaggerRule
	Score float64 `json:"score" example:"0.95"`
}

// SwaggerErrorResponse represents an error response for Swagger documentation
type SwaggerErrorResponse struct {
	Error string `json:"error" example:"Invalid request data"`
}

// SwaggerListResponse represents a list response for Swagger documentation
type SwaggerListResponse struct {
	Items []interface{} `json:"items"`
	Total int           `json:"total" example:"100"`
	Page  int           `json:"page" example:"1"`
	Limit int           `json:"limit" example:"10"`
}
