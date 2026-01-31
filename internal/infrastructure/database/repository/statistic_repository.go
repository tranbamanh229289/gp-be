package repository

import (
	"be/config"
	"be/internal/domain/statistic.go"
	"be/internal/infrastructure/database/postgres"
	"be/internal/shared/helper"
	"context"
)

type StatisticRepository struct {
	db     *postgres.PostgresDB
	config *config.Config
}

func NewStatisticRepository(db *postgres.PostgresDB, config *config.Config) statistic.IStatisticRepository {
	return &StatisticRepository{
		db:     db,
		config: config,
	}
}

func (r *StatisticRepository) FindHolderStatisticByDID(ctx context.Context, did string) (*statistic.HolderStatistic, error) {
	var entity statistic.HolderStatistic
	if err := r.db.GetGormDB().WithContext(ctx).Where("holder_did = ?", did).First(&entity).Error; err != nil {
		return nil, err
	}
	return &entity, nil
}

func (r *StatisticRepository) CreateHolderStatistic(ctx context.Context) (*statistic.HolderStatistic, error) {
	db := helper.WithTx(ctx, r.db.GetGormDB())
	var entity *statistic.HolderStatistic = &statistic.HolderStatistic{}
	if err := db.Create(entity).Error; err != nil {
		return nil, err
	}
	return entity, nil
}

func (r *StatisticRepository) UpdateHolderStatisticByDID(ctx context.Context, entity *statistic.HolderStatistic, changes map[string]interface{}) error {
	db := helper.WithTx(ctx, r.db.GetGormDB())
	if err := db.Model(entity).Updates(changes).Error; err != nil {
		return err
	}
	return nil
}

func (r *StatisticRepository) FindIssuerStatisticByDID(ctx context.Context, did string) (*statistic.IssuerStatistic, error) {
	var entity statistic.IssuerStatistic
	if err := r.db.GetGormDB().WithContext(ctx).Where("issuer_did = ?", did).First(&entity).Error; err != nil {
		return nil, err
	}
	return &entity, nil
}

func (r *StatisticRepository) CreateIssuerStatistic(ctx context.Context) (*statistic.IssuerStatistic, error) {
	db := helper.WithTx(ctx, r.db.GetGormDB())
	var entity *statistic.IssuerStatistic = &statistic.IssuerStatistic{}
	if err := db.Create(entity).Error; err != nil {
		return nil, err
	}
	return entity, nil
}

func (r *StatisticRepository) UpdateIssuerStatisticByDID(ctx context.Context, entity *statistic.IssuerStatistic, changes map[string]interface{}) error {
	db := helper.WithTx(ctx, r.db.GetGormDB())
	if err := db.Model(entity).Updates(changes).Error; err != nil {
		return err
	}
	return nil
}

func (r *StatisticRepository) FindVerifierStatisticByDID(ctx context.Context, did string) (*statistic.VerifierStatistic, error) {
	var entity statistic.VerifierStatistic
	if err := r.db.GetGormDB().WithContext(ctx).Where("verifier_did = ?", did).First(&entity).Error; err != nil {
		return nil, err
	}
	return &entity, nil
}

func (r *StatisticRepository) CreateVerifierStatistic(ctx context.Context) (*statistic.VerifierStatistic, error) {
	db := helper.WithTx(ctx, r.db.GetGormDB())
	var entity *statistic.VerifierStatistic = &statistic.VerifierStatistic{}
	if err := db.Create(entity).Error; err != nil {
		return nil, err
	}
	return entity, nil
}

func (r *StatisticRepository) UpdateVerifierStatisticByDID(ctx context.Context, entity *statistic.VerifierStatistic, changes map[string]interface{}) error {
	db := helper.WithTx(ctx, r.db.GetGormDB())
	if err := db.Model(entity).Updates(changes).Error; err != nil {
		return err
	}
	return nil
}
