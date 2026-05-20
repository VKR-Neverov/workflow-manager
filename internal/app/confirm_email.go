package app

import (
	"context"
	"github.com/fidesy-pay/workflow-manager/internal/pkg/builders"
	"github.com/fidesy-pay/workflow-manager/internal/pkg/model"
	desc "github.com/fidesy-pay/workflow-manager/pkg/workflow-manager"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (i *Implementation) ConfirmEmail(ctx context.Context, req *desc.ConfirmEmailRequest) (*desc.FlowResponse, error) {
	options := model.NewOptions()

	options.SetString("email_code", req.GetCode())

	flow, err := i.runner.Complete(ctx, req.GetId(), options)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "runner.Complete: %v", err)
	}

	return builders.BuildFlowResponse(flow), nil
}
