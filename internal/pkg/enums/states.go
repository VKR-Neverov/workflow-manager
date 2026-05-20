package enums

type State string

var (
	StateUnknown                  State = "UNKNOWN"
	StateNew                      State = "NEW"
	StateCompleted                State = "COMPLETED"
	StateFailed                   State = "FAILED"
	StateWaitingEmail             State = "WAITING_EMAIL"
	StateEmailSet                 State = "EMAIL_SET"
	StateWaitingEmailConfirmation State = "WAITING_EMAIL_CONFIRMATION"
	StateEmailConfirmed           State = "EMAIL_CONFIRMED"
	StateWaitingCredentials       State = "WAITING_CREDENTIALS"
	StateCredentialsSet           State = "CREDENTIALS_SET"
	StateWaitingPassword          State = "WAITING_PASSWORD"
	StatePasswordSet              State = "PASSWORD_SET"
)
