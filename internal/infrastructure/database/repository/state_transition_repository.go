package repository

import (
	gist "be/internal/domain/gist"
	"be/internal/infrastructure/database/postgres"
	"be/internal/shared/helper"
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
	db := helper.WithTx(ctx, r.db.GetGormDB())
	if err := db.Create(entity).Error; err != nil {
		return nil, err
	}
	return entity, nil
}
