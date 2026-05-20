package app

import (
	"context"
	"github.com/fidesy-pay/workflow-manager/internal/pkg/builders"
	"github.com/fidesy-pay/workflow-manager/internal/pkg/enums"
	"github.com/fidesy-pay/workflow-manager/internal/pkg/model"
	desc "github.com/fidesy-pay/workflow-manager/pkg/workflow-manager"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (i *Implementation) CreateRegistration(ctx context.Context, req *desc.CreateRegistrationRequest) (*desc.FlowResponse, error) {
	createOptions := model.NewOptions()

	flow, err := i.runner.Create(ctx, enums.FlowTypeRegistration, createOptions)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "runner.Create: %v", err)
	}

	flow, err = i.runner.Complete(ctx, flow.ID.String(), nil)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "runner.Complete: %v", err)
	}

	return builders.BuildFlowResponse(flow), nil
}
