package statistic

import "context"

type IStatisticRepository interface {
	FindHolderStatisticByDID(ctx context.Context, did string) (*HolderStatistic, error)
	CreateHolderStatistic(ctx context.Context) (*HolderStatistic, error)
	UpdateHolderStatisticByDID(ctx context.Context, entity *HolderStatistic, changes map[string]interface{}) error

	FindIssuerStatisticByDID(ctx context.Context, did string) (*IssuerStatistic, error)
	CreateIssuerStatistic(ctx context.Context) (*IssuerStatistic, error)
	UpdateIssuerStatisticByDID(ctx context.Context, entity *IssuerStatistic, changes map[string]interface{}) error

	FindVerifierStatisticByDID(ctx context.Context, did string) (*VerifierStatistic, error)
	CreateVerifierStatistic(ctx context.Context) (*VerifierStatistic, error)
	UpdateVerifierStatisticByDID(ctx context.Context, entity *VerifierStatistic, changes map[string]interface{}) error
}
