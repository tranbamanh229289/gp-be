package repository

import (
	gist "be/internal/domain/gist"
	"be/internal/infrastructure/database/postgres"
	"context"
)

type StateTransitionRepository struct {
	db *postgres.PostgresDB
}

func NewStateTransitionRepository(db *postgres.PostgresDB) gist.IStateTransition {
	return &StateTransitionRepository{
		db: db,
	}
}

func (r *StateTransitionRepository) CreateStateTransition(ctx context.Context, entity *gist.StateTransition) (*gist.StateTransition, error) {
	if err := r.db.GetGormDB().WithContext(ctx).Create(entity).Error; err != nil {
		return nil, err
	}
	return entity, nil
}
