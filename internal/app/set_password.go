package app

import (
	"context"

	"github.com/fidesy-pay/workflow-manager/internal/pkg/builders"
	"github.com/fidesy-pay/workflow-manager/internal/pkg/model"
	desc "github.com/fidesy-pay/workflow-manager/pkg/workflow-manager"
	validation "github.com/go-ozzo/ozzo-validation"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (i *Implementation) SetPassword(ctx context.Context, req *desc.SetPasswordRequest) (*desc.FlowResponse, error) {
	if err := validateSetPasswordRequest(req); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	options := model.NewOptions()

	options.SetString("password", req.GetPassword())

	flow, err := i.runner.Complete(ctx, req.GetId(), options)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "runner.Complete: %v", err)
	}

	return builders.BuildFlowResponse(flow), nil
}

func validateSetPasswordRequest(req *desc.SetPasswordRequest) error {
	return validation.ValidateStruct(req,
		validation.Field(&req.Password, validation.Required),
	)
}
