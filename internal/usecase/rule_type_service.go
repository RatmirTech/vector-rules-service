package usecase

import (
	"context"
	"fmt"

	"github.com/ratmirtech/vector-rules-service/internal/domain"
)

type ruleTypeService struct {
	ruleTypeRepo domain.RuleTypeRepository
}

// NewRuleTypeService creates a new rule type service
func NewRuleTypeService(ruleTypeRepo domain.RuleTypeRepository) domain.RuleTypeService {
	return &ruleTypeService{
		ruleTypeRepo: ruleTypeRepo,
	}
}

func (s *ruleTypeService) CreateRuleType(ctx context.Context, req *domain.CreateRuleTypeRequest) (*domain.RuleType, error) {
	ruleType := &domain.RuleType{
		Name: req.Name,
	}

	createdRuleType, err := s.ruleTypeRepo.Create(ctx, ruleType)
	if err != nil {
		return nil, fmt.Errorf("failed to create rule type: %w", err)
	}

	return createdRuleType, nil
}

func (s *ruleTypeService) GetRuleType(ctx context.Context, id int64) (*domain.RuleType, error) {
	ruleType, err := s.ruleTypeRepo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get rule type: %w", err)
	}
	return ruleType, nil
}

func (s *ruleTypeService) UpdateRuleType(ctx context.Context, req *domain.UpdateRuleTypeRequest) (*domain.RuleType, error) {
	// Check if rule type exists
	existingRuleType, err := s.ruleTypeRepo.GetByID(ctx, req.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to get existing rule type: %w", err)
	}

	// Update rule type
	existingRuleType.Name = req.Name

	updatedRuleType, err := s.ruleTypeRepo.Update(ctx, existingRuleType)
	if err != nil {
		return nil, fmt.Errorf("failed to update rule type: %w", err)
	}

	return updatedRuleType, nil
}

func (s *ruleTypeService) DeleteRuleType(ctx context.Context, id int64) error {
	err := s.ruleTypeRepo.Delete(ctx, id)
	if err != nil {
		return fmt.Errorf("failed to delete rule type: %w", err)
	}
	return nil
}

func (s *ruleTypeService) ListRuleTypes(ctx context.Context, limit, offset int) ([]*domain.RuleType, error) {
	ruleTypes, err := s.ruleTypeRepo.List(ctx, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to list rule types: %w", err)
	}
	return ruleTypes, nil
}