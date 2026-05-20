package app

import (
	"context"
	"github.com/fidesy-pay/workflow-manager/internal/pkg/builders"
	"github.com/fidesy-pay/workflow-manager/internal/pkg/enums"
	"github.com/fidesy-pay/workflow-manager/internal/pkg/model"
	desc "github.com/fidesy-pay/workflow-manager/pkg/workflow-manager"
	validation "github.com/go-ozzo/ozzo-validation"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (i *Implementation) CreateRestorePassword(ctx context.Context, req *desc.CreateRestorePasswordRequest) (*desc.FlowResponse, error) {
	if err := validateCreateRestorePasswordRequest(req); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	createOptions := model.NewOptions()
	createOptions.SetString("email", req.Email)

	flow, err := i.runner.Create(ctx, enums.FlowTypeRestorePassword, createOptions)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "runner.Create: %v", err)
	}

	flow, err = i.runner.Complete(ctx, flow.ID.String(), nil)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "runner.Complete: %v", err)
	}

	return builders.BuildFlowResponse(flow), nil
}

func validateCreateRestorePasswordRequest(req *desc.CreateRestorePasswordRequest) error {
	return validation.ValidateStruct(req,
		validation.Field(&req.Email, validation.Required),
	)
}
