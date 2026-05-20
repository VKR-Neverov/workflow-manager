package dto

import (
	"github.com/fidesy-pay/workflow-manager/internal/pkg/enums"
	"github.com/fidesy/sdk/common/postgres"
)

type ListFlowsParams struct {
	Filter     ListFlowsFilter
	Pagination postgres.Pagination
}

type ListFlowsFilter struct {
	IDIn       []string
	StateNotIn []enums.State
}
