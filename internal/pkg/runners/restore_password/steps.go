package restorepassword

import (
	"context"
	"errors"
	"fmt"

	"github.com/fidesy-pay/workflow-manager/internal/pkg/enums"
	"github.com/fidesy-pay/workflow-manager/internal/pkg/model"
	auth_service "github.com/fidesy-pay/workflow-manager/pkg/auth-service"
	clients_service "github.com/fidesy-pay/workflow-manager/pkg/clients-service"
	email_service "github.com/fidesy-pay/workflow-manager/pkg/email-service"
)

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

func (r *Runner) waitingPasswordStep(ctx context.Context, input model.Step[Data]) (model.Step[Data], error) {
	password := input.Options.GetString("password")
	if password == "" {
		return input.EmptyStep(), nil
	}

	clients, err := r.clientClient.ListClients(ctx, &clients_service.ListClientsRequest{
		Filter: &clients_service.ListClientsRequest_Filter{
			EmailIn: []string{input.Data.Email},
		},
	})
	if err != nil {
		return input.EmptyStep(), fmt.Errorf("clientClient.ListClients: %w", err)
	}

	if len(clients.Clients) == 0 {
		return input.EmptyStep(), errors.New("client with such email not found")
	}

	_, err = r.authClient.UpdatePassword(ctx, &auth_service.UpdatePasswordRequest{
		Username: clients.Clients[0].Username,
		Password: password,
	})
	if err != nil {
		return input.EmptyStep(), fmt.Errorf("authClient.UpdatePassword: %w", err)
	}

	return input.Next(enums.StatePasswordSet), nil
}
