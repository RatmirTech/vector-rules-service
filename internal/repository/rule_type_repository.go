package repository

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/ratmirtech/vector-rules-service/internal/domain"
)

type ruleTypeRepository struct {
	db *pgxpool.Pool
}

// NewRuleTypeRepository creates a new rule type repository
func NewRuleTypeRepository(db *pgxpool.Pool) domain.RuleTypeRepository {
	return &ruleTypeRepository{db: db}
}

func (r *ruleTypeRepository) Create(ctx context.Context, ruleType *domain.RuleType) (*domain.RuleType, error) {
	const query = `
		INSERT INTO rule_types (name)
		VALUES ($1)
		RETURNING id, created_at, updated_at`

	var result domain.RuleType
	result.Name = ruleType.Name

	err := r.db.QueryRow(ctx, query, ruleType.Name).
		Scan(&result.ID, &result.CreatedAt, &result.UpdatedAt)
	if err != nil {
		return nil, fmt.Errorf("failed to create rule type: %w", err)
	}

	return &result, nil
}

func (r *ruleTypeRepository) GetByID(ctx context.Context, id int64) (*domain.RuleType, error) {
	const query = `
		SELECT id, name, created_at, updated_at
		FROM rule_types
		WHERE id = $1`

	var ruleType domain.RuleType
	err := r.db.QueryRow(ctx, query, id).Scan(
		&ruleType.ID,
		&ruleType.Name,
		&ruleType.CreatedAt,
		&ruleType.UpdatedAt,
	)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, domain.ErrRuleTypeNotFound
		}
		return nil, fmt.Errorf("failed to get rule type by id: %w", err)
	}

	return &ruleType, nil
}

func (r *ruleTypeRepository) GetByName(ctx context.Context, name string) (*domain.RuleType, error) {
	const query = `
		SELECT id, name, created_at, updated_at
		FROM rule_types
		WHERE name = $1`

	var ruleType domain.RuleType
	err := r.db.QueryRow(ctx, query, name).Scan(
		&ruleType.ID,
		&ruleType.Name,
		&ruleType.CreatedAt,
		&ruleType.UpdatedAt,
	)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, domain.ErrRuleTypeNotFound
		}
		return nil, fmt.Errorf("failed to get rule type by name: %w", err)
	}

	return &ruleType, nil
}

func (r *ruleTypeRepository) Update(ctx context.Context, ruleType *domain.RuleType) (*domain.RuleType, error) {
	const query = `
		UPDATE rule_types 
		SET name = $2, updated_at = NOW()
		WHERE id = $1
		RETURNING updated_at`

	err := r.db.QueryRow(ctx, query, ruleType.ID, ruleType.Name).
		Scan(&ruleType.UpdatedAt)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, domain.ErrRuleTypeNotFound
		}
		return nil, fmt.Errorf("failed to update rule type: %w", err)
	}

	return ruleType, nil
}

func (r *ruleTypeRepository) Delete(ctx context.Context, id int64) error {
	const query = `DELETE FROM rule_types WHERE id = $1`

	result, err := r.db.Exec(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete rule type: %w", err)
	}

	if result.RowsAffected() == 0 {
		return domain.ErrRuleTypeNotFound
	}

	return nil
}

func (r *ruleTypeRepository) List(ctx context.Context, limit, offset int) ([]*domain.RuleType, error) {
	const query = `
		SELECT id, name, created_at, updated_at
		FROM rule_types
		ORDER BY name ASC
		LIMIT $1 OFFSET $2`

	rows, err := r.db.Query(ctx, query, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to list rule types: %w", err)
	}
	defer rows.Close()

	var ruleTypes []*domain.RuleType
	for rows.Next() {
		var ruleType domain.RuleType
		err := rows.Scan(
			&ruleType.ID,
			&ruleType.Name,
			&ruleType.CreatedAt,
			&ruleType.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan rule type: %w", err)
		}
		ruleTypes = append(ruleTypes, &ruleType)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating rule types: %w", err)
	}

	return ruleTypes, nil
}