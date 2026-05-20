package app

import (
	"context"
	"github.com/fidesy-pay/workflow-manager/internal/pkg/builders"
	"github.com/fidesy-pay/workflow-manager/internal/pkg/dto"
	desc "github.com/fidesy-pay/workflow-manager/pkg/workflow-manager"
	"github.com/fidesy/sdk/common/postgres"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (i *Implementation) ListFlows(ctx context.Context, req *desc.ListFlowsRequest) (*desc.ListFlowsResponse, error) {
	if req.Filter == nil {
		return nil, status.Error(codes.InvalidArgument, "filter is required")
	}

	flows, err := i.storage.ListFlows(ctx, dto.ListFlowsParams{
		Filter: dto.ListFlowsFilter{
			IDIn: req.Filter.GetIdIn(),
		},
		Pagination: postgres.NewPagination(req.GetPage(), req.GetPerPage()),
	})
	if err != nil {
		return nil, status.Errorf(codes.Internal, "storage.ListFlows: %v", err)
	}

	return &desc.ListFlowsResponse{
		Flows: builders.BuildFlowsResponse(flows),
	}, nil
}
