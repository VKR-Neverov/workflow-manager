package registration

import (
	"context"
	"github.com/fidesy-pay/workflow-manager/internal/pkg/dto"
	"github.com/fidesy-pay/workflow-manager/internal/pkg/enums"
	"github.com/fidesy-pay/workflow-manager/internal/pkg/model"
	auth_service "github.com/fidesy-pay/workflow-manager/pkg/auth-service"
	clients_service "github.com/fidesy-pay/workflow-manager/pkg/clients-service"
	crypto_service "github.com/fidesy-pay/workflow-manager/pkg/crypto-service"
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
		cryptoClient CryptoClient
	}

	Storage interface {
		ListFlows(ctx context.Context, params dto.ListFlowsParams) ([]*model.Flow, error)
		UpdateFlow(ctx context.Context, flow *model.Flow) (*model.Flow, error)
		CreateFlow(ctx context.Context, flow *model.Flow) (*model.Flow, error)
	}

	EmailClient interface {
		SendCode(ctx context.Context, in *email_service.SendCodeRequest, opts ...grpc.CallOption) (*email_service.SendCodeResponse, error)
		ConfirmCode(ctx context.Context, in *email_service.ConfirmCodeRequest, opts ...grpc.CallOption) (*email_service.ConfirmCodeResponse, error)
	}

	AuthClient interface {
		SignUp(ctx context.Context, in *auth_service.SignUpRequest, opts ...grpc.CallOption) (*auth_service.SignUpResponse, error)
	}

	ClientsClient interface {
		CreateClient(ctx context.Context, in *clients_service.CreateClientRequest, opts ...grpc.CallOption) (*clients_service.Client, error)
		ListClients(ctx context.Context, in *clients_service.ListClientsRequest, opts ...grpc.CallOption) (*clients_service.ListClientsResponse, error)
	}

	CryptoClient interface {
		CreateWallet(ctx context.Context, in *crypto_service.CreateWalletRequest, opts ...grpc.CallOption) (*crypto_service.Wallet, error)
	}
)

func NewRunner(
	storage Storage,
	emailClient EmailClient,
	authClient AuthClient,
	clientClient ClientsClient,
	cryptoClient CryptoClient,
) *Runner {
	stepper := model.NewStepper[Data](storage)

	r := &Runner{
		Stepper:      stepper,
		storage:      storage,
		emailClient:  emailClient,
		authClient:   authClient,
		clientClient: clientClient,
		cryptoClient: cryptoClient,
	}

	stepper.Add(enums.StateNew, r.next(enums.StateWaitingEmail))
	stepper.Add(enums.StateWaitingEmail, r.waitingEmailStep)
	stepper.Add(enums.StateEmailSet, r.next(enums.StateWaitingEmailConfirmation))
	stepper.Add(enums.StateWaitingEmailConfirmation, r.waitingEmailConfirmationStep)
	stepper.Add(enums.StateEmailConfirmed, r.next(enums.StateWaitingCredentials))
	stepper.Add(enums.StateWaitingCredentials, r.waitingCredentialsStep)
	stepper.Add(enums.StateCredentialsSet, r.next(enums.StateCompleted))

	return r
}

func (r *Runner) Create(_ context.Context, _ model.Options) (*model.Flow, error) {
	return &model.Flow{}, nil
}

func (r *Runner) Type() enums.FlowType {
	return enums.FlowTypeRegistration
}

func (r *Runner) next(state enums.State) model.StepFunc[Data] {
	return func(ctx context.Context, input model.Step[Data]) (model.Step[Data], error) {
		return input.Next(state), nil
	}
}
