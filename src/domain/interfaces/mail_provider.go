package interfaces

type MailProvider interface {
	Send(to, subject, body string) error
}
