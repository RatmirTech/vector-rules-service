package repository

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/pgvector/pgvector-go"
	"github.com/ratmirtech/vector-rules-service/internal/domain"
)

type ruleRepository struct {
	db *pgxpool.Pool
}

// NewRuleRepository creates a new rule repository
func NewRuleRepository(db *pgxpool.Pool) domain.RuleRepository {
	return &ruleRepository{db: db}
}

func (r *ruleRepository) Create(ctx context.Context, rule *domain.Rule) (*domain.Rule, error) {
	const query = `
		INSERT INTO rules (rule_type_id, content, embedding)
		VALUES ($1, $2, $3)
		RETURNING id, created_at, updated_at`

	var embedding interface{}
	if rule.Embedding != nil {
		embedding = pgvector.NewVector(rule.Embedding)
	}

	var result domain.Rule
	result = *rule

	err := r.db.QueryRow(ctx, query, rule.RuleTypeID, rule.Content, embedding).
		Scan(&result.ID, &result.CreatedAt, &result.UpdatedAt)
	if err != nil {
		return nil, fmt.Errorf("failed to create rule: %w", err)
	}

	return &result, nil
}

func (r *ruleRepository) GetByID(ctx context.Context, id int64) (*domain.Rule, error) {
	const query = `
		SELECT r.id, r.rule_type_id, r.content, r.embedding, r.created_at, r.updated_at,
		       rt.name as rule_type_name
		FROM rules r
		JOIN rule_types rt ON r.rule_type_id = rt.id
		WHERE r.id = $1`

	var rule domain.Rule
	var embedding pgvector.Vector
	var embeddingNull sql.NullString

	err := r.db.QueryRow(ctx, query, id).Scan(
		&rule.ID,
		&rule.RuleTypeID,
		&rule.Content,
		&embeddingNull,
		&rule.CreatedAt,
		&rule.UpdatedAt,
		&rule.RuleTypeName,
	)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, domain.ErrRuleNotFound
		}
		return nil, fmt.Errorf("failed to get rule by id: %w", err)
	}

	if embeddingNull.Valid {
		if err := embedding.Scan(embeddingNull.String); err == nil {
			rule.Embedding = embedding.Slice()
		}
	}

	return &rule, nil
}

func (r *ruleRepository) Update(ctx context.Context, rule *domain.Rule) (*domain.Rule, error) {
	const query = `
		UPDATE rules 
		SET rule_type_id = $2, content = $3, embedding = $4, updated_at = NOW()
		WHERE id = $1
		RETURNING updated_at`

	var embedding interface{}
	if rule.Embedding != nil {
		embedding = pgvector.NewVector(rule.Embedding)
	}

	err := r.db.QueryRow(ctx, query, rule.ID, rule.RuleTypeID, rule.Content, embedding).
		Scan(&rule.UpdatedAt)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, domain.ErrRuleNotFound
		}
		return nil, fmt.Errorf("failed to update rule: %w", err)
	}

	return rule, nil
}

func (r *ruleRepository) Delete(ctx context.Context, id int64) error {
	const query = `DELETE FROM rules WHERE id = $1`

	result, err := r.db.Exec(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete rule: %w", err)
	}

	if result.RowsAffected() == 0 {
		return domain.ErrRuleNotFound
	}

	return nil
}

func (r *ruleRepository) List(ctx context.Context, ruleType *string, limit, offset int) ([]*domain.Rule, error) {
	query := `
		SELECT r.id, r.rule_type_id, r.content, r.created_at, r.updated_at,
		       rt.name as rule_type_name
		FROM rules r
		JOIN rule_types rt ON r.rule_type_id = rt.id`

	var args []interface{}
	argIndex := 1

	if ruleType != nil {
		query += " WHERE rt.name = $" + fmt.Sprintf("%d", argIndex)
		args = append(args, *ruleType)
		argIndex++
	}

	query += fmt.Sprintf(" ORDER BY r.created_at DESC LIMIT $%d OFFSET $%d", argIndex, argIndex+1)
	args = append(args, limit, offset)

	rows, err := r.db.Query(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to list rules: %w", err)
	}
	defer rows.Close()

	var rules []*domain.Rule
	for rows.Next() {
		var rule domain.Rule
		err := rows.Scan(
			&rule.ID,
			&rule.RuleTypeID,
			&rule.Content,
			&rule.CreatedAt,
			&rule.UpdatedAt,
			&rule.RuleTypeName,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan rule: %w", err)
		}
		rules = append(rules, &rule)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating rules: %w", err)
	}

	return rules, nil
}

func (r *ruleRepository) FindSimilar(ctx context.Context, embedding []float32, ruleType *string, limit int) ([]*domain.RuleMatch, error) {
	query := `
		SELECT r.id, r.rule_type_id, r.content, r.created_at, r.updated_at,
		       rt.name as rule_type_name,
		       1 - (r.embedding <=> $1) as similarity_score
		FROM rules r
		JOIN rule_types rt ON r.rule_type_id = rt.id
		WHERE r.embedding IS NOT NULL`

	var args []interface{}
	args = append(args, pgvector.NewVector(embedding))
	argIndex := 2

	if ruleType != nil {
		query += " AND rt.name = $" + fmt.Sprintf("%d", argIndex)
		args = append(args, *ruleType)
		argIndex++
	}

	query += fmt.Sprintf(" ORDER BY r.embedding <=> $1 LIMIT $%d", argIndex)
	args = append(args, limit)

	rows, err := r.db.Query(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to find similar rules: %w", err)
	}
	defer rows.Close()

	var matches []*domain.RuleMatch
	for rows.Next() {
		var match domain.RuleMatch
		err := rows.Scan(
			&match.ID,
			&match.RuleTypeID,
			&match.Content,
			&match.CreatedAt,
			&match.UpdatedAt,
			&match.RuleTypeName,
			&match.Score,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan rule match: %w", err)
		}
		matches = append(matches, &match)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating rule matches: %w", err)
	}

	return matches, nil
}

func (r *ruleRepository) UpdateEmbedding(ctx context.Context, id int64, embedding []float32) error {
	const query = `
		UPDATE rules 
		SET embedding = $2, updated_at = NOW()
		WHERE id = $1`

	result, err := r.db.Exec(ctx, query, id, pgvector.NewVector(embedding))
	if err != nil {
		return fmt.Errorf("failed to update rule embedding: %w", err)
	}

	if result.RowsAffected() == 0 {
		return domain.ErrRuleNotFound
	}

	return nil
}