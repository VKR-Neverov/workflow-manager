package storage

import (
	"context"
	"encoding/json"
	"fmt"
	sq "github.com/Masterminds/squirrel"
	"github.com/fidesy-pay/workflow-manager/internal/pkg/dto"
	"github.com/fidesy-pay/workflow-manager/internal/pkg/model"
	"github.com/fidesy/sdk/common/postgres"
)

func (s *Storage) ListFlows(ctx context.Context, params dto.ListFlowsParams) ([]*model.Flow, error) {
	query := postgres.Builder().
		Select(flowFields).
		From(flowsTable)

	if len(params.Filter.IDIn) > 0 {
		query = query.Where(sq.Eq{
			"id": params.Filter.IDIn,
		})
	}

	if len(params.Filter.StateNotIn) > 0 {
		query = query.Where(sq.NotEq{
			"state": params.Filter.StateNotIn,
		})
	}

	query = query.
		Limit(params.Pagination.Limit()).
		Offset(params.Pagination.Offset())

	return postgres.Select[model.Flow](ctx, s.pool, query)
}

func (s *Storage) CreateFlow(ctx context.Context, flow *model.Flow) (*model.Flow, error) {
	data := map[string]interface{}{
		"state": flow.State,
		"type":  flow.Type,
	}

	if flow.Data != nil {
		// Marshal the data field to JSON
		jsonData, err := json.Marshal(flow.Data)
		if err != nil {
			return nil, err
		}

		data["data"] = jsonData
	}

	query := postgres.Builder().
		Insert(flowsTable).
		SetMap(data).
		Suffix(fmt.Sprintf("RETURNING %s", flowFields))

	return postgres.Exec[model.Flow](ctx, s.pool, query)
}

func (s *Storage) UpdateFlow(ctx context.Context, flow *model.Flow) (*model.Flow, error) {
	// Marshal the data field to JSON
	jsonData, err := json.Marshal(flow.Data)
	if err != nil {
		return nil, err
	}

	query := postgres.Builder().
		Update(flowsTable).
		SetMap(map[string]interface{}{
			"state":        flow.State,
			"completed_at": flow.CompletedAt,
			"data":         jsonData,
		}).
		Where(sq.Eq{
			"id": flow.ID,
		}).
		Suffix(fmt.Sprintf("RETURNING %s", flowFields))

	return postgres.Exec[model.Flow](ctx, s.pool, query)
}
