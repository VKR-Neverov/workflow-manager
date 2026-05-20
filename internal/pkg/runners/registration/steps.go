package registration

import (
	"context"
	"errors"
	"fmt"

	"github.com/fidesy-pay/workflow-manager/internal/pkg/enums"
	"github.com/fidesy-pay/workflow-manager/internal/pkg/model"
	auth_service "github.com/fidesy-pay/workflow-manager/pkg/auth-service"
	clients_service "github.com/fidesy-pay/workflow-manager/pkg/clients-service"
	crypto_service "github.com/fidesy-pay/workflow-manager/pkg/crypto-service"
	email_service "github.com/fidesy-pay/workflow-manager/pkg/email-service"
)

func (r *Runner) waitingEmailStep(ctx context.Context, input model.Step[Data]) (model.Step[Data], error) {
	email := input.Options.GetString("email")
	if email == "" {
		return input.EmptyStep(), nil
	}

	clients, err := r.clientClient.ListClients(ctx, &clients_service.ListClientsRequest{
		Filter: &clients_service.ListClientsRequest_Filter{
			EmailIn: []string{email},
		},
	})
	if err != nil {
		return input.EmptyStep(), fmt.Errorf("clientClient.ListClients: %w", err)
	}

	if len(clients.Clients) > 0 {
		return input.EmptyStep(), errors.New("email already exists")
	}

	resp, err := r.emailClient.SendCode(ctx, &email_service.SendCodeRequest{
		Email: email,
	})
	if err != nil {
		return input.EmptyStep(), fmt.Errorf("emailClient.SendCode: %w", err)
	}

	input.Data.Email = email
	input.Data.EmailCodeID = resp.Id

	return input.Next(enums.StateEmailSet), nil
}

func (r *Runner) waitingEmailConfirmationStep(ctx context.Context, input model.Step[Data]) (model.Step[Data], error) {
	emailCode := input.Options.GetString("email_code")
	if emailCode == "" {
		return input.EmptyStep(), nil
	}

	_, err := r.emailClient.ConfirmCode(ctx, &email_service.ConfirmCodeRequest{
		Id:   input.Data.EmailCodeID,
		Code: emailCode,
	})
	if err != nil {
		return input.EmptyStep(), fmt.Errorf("emailClient.ConfirmCode: %w", err)
	}

	return input.Next(enums.StateEmailConfirmed), nil
}

func (r *Runner) waitingCredentialsStep(ctx context.Context, input model.Step[Data]) (model.Step[Data], error) {
	username := input.Options.GetString("username")
	if username == "" {
		return input.EmptyStep(), nil
	}

	password := input.Options.GetString("password")
	if password == "" {
		return input.EmptyStep(), nil
	}

	client, err := r.clientClient.CreateClient(ctx, &clients_service.CreateClientRequest{
		Username: username,
		Email:    input.Data.Email,
	})
	if err != nil {
		return input.EmptyStep(), fmt.Errorf("clientClient.CreateClient: %w", err)
	}

	_, err = r.cryptoClient.CreateWallet(ctx, &crypto_service.CreateWalletRequest{
		ClientId: client.Id,
	})
	if err != nil {
		return input.EmptyStep(), fmt.Errorf("cryptoClient.CreateWallet: %w", err)
	}

	_, err = r.authClient.SignUp(ctx, &auth_service.SignUpRequest{
		Username: username,
		Password: password,
		ClientId: client.Id,
	})
	if err != nil {
		return input.EmptyStep(), fmt.Errorf("authClient.SignUp: %w", err)
	}

	return input.Next(enums.StateCredentialsSet), nil
}
