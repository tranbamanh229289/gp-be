package repository

import (
	"be/internal/domain/claim"
	"be/internal/infrastructure/database/postgres"
	"context"
)

type StateTransitionRepository struct {
	db *postgres.PostgresDB
}

func NewStateTransitionRepository(db *postgres.PostgresDB) claim.IStateTransition {
	return &StateTransitionRepository{
		db: db,
	}
}

func (r *StateTransitionRepository) CreateStateTransition(ctx context.Context, entity *claim.StateTransition) (*claim.StateTransition, error) {
	if err := r.db.GetGormDB().WithContext(ctx).Create(entity).Error; err != nil {
		return nil, err
	}
	return entity, nil
}
