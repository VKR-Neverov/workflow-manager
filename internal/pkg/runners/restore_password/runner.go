package restorepassword

import (
	"context"
	"errors"
	"fmt"
	"github.com/fidesy-pay/workflow-manager/internal/pkg/dto"
	"github.com/fidesy-pay/workflow-manager/internal/pkg/enums"
	"github.com/fidesy-pay/workflow-manager/internal/pkg/model"
	auth_service "github.com/fidesy-pay/workflow-manager/pkg/auth-service"
	clients_service "github.com/fidesy-pay/workflow-manager/pkg/clients-service"
	email_service "github.com/fidesy-pay/workflow-manager/pkg/email-service"
	"google.golang.org/grpc"
)

type (
	Runner struct {
		*model.Stepper[Data]

		storage      Storage
		emailClient  EmailClient
		authClient   AuthClient
		clientClient ClientsClient
	}

	Storage interface {
		ListFlows(ctx context.Context, params dto.ListFlowsParams) ([]*model.Flow, error)
		CreateFlow(ctx context.Context, flow *model.Flow) (*model.Flow, error)
		UpdateFlow(ctx context.Context, flow *model.Flow) (*model.Flow, error)
	}
	EmailClient interface {
		SendCode(ctx context.Context, in *email_service.SendCodeRequest, opts ...grpc.CallOption) (*email_service.SendCodeResponse, error)
		ConfirmCode(ctx context.Context, in *email_service.ConfirmCodeRequest, opts ...grpc.CallOption) (*email_service.ConfirmCodeResponse, error)
	}
	AuthClient interface {
		UpdatePassword(ctx context.Context, in *auth_service.UpdatePasswordRequest, opts ...grpc.CallOption) (*auth_service.UpdatePasswordResponse, error)
	}
	ClientsClient interface {
		ListClients(ctx context.Context, in *clients_service.ListClientsRequest, opts ...grpc.CallOption) (*clients_service.ListClientsResponse, error)
	}
)

func NewRunner(
	storage Storage,
	emailClient EmailClient,
	authClient AuthClient,
	clientClient ClientsClient,
) *Runner {
	stepper := model.NewStepper[Data](storage)

	r := &Runner{
		Stepper:      stepper,
		storage:      storage,
		emailClient:  emailClient,
		authClient:   authClient,
		clientClient: clientClient,
	}

	stepper.Add(enums.StateNew, r.next(enums.StateWaitingEmailConfirmation))
	stepper.Add(enums.StateWaitingEmailConfirmation, r.waitingEmailConfirmationStep)
	stepper.Add(enums.StateEmailConfirmed, r.next(enums.StateWaitingPassword))
	stepper.Add(enums.StateWaitingPassword, r.waitingPasswordStep)
	stepper.Add(enums.StatePasswordSet, r.next(enums.StateCompleted))

	return r
}

func (r *Runner) Create(ctx context.Context, options model.Options) (*model.Flow, error) {
	email := options.GetString("email")
	if email == "" {
		return nil, errors.New("email is required")
	}

	clients, err := r.clientClient.ListClients(ctx, &clients_service.ListClientsRequest{
		Filter: &clients_service.ListClientsRequest_Filter{
			EmailIn: []string{email},
		},
	})
	if err != nil {
		return nil, fmt.Errorf("clientClient.ListClients: %w", err)
	}

	if len(clients.Clients) == 0 {
		return nil, errors.New("client with such email not found")
	}

	resp, err := r.emailClient.SendCode(ctx, &email_service.SendCodeRequest{
		Email: email,
	})
	if err != nil {
		return nil, fmt.Errorf("emailClient.SendCode: %w", err)
	}

	flow := &model.Flow{
		Data: &Data{
			Email:       email,
			EmailCodeID: resp.Id,
		},
	}

	return flow, nil
}

func (r *Runner) Type() enums.FlowType {
	return enums.FlowTypeRestorePassword
}

func (r *Runner) next(state enums.State) model.StepFunc[Data] {
	return func(ctx context.Context, input model.Step[Data]) (model.Step[Data], error) {
		return input.Next(state), nil
	}
}
