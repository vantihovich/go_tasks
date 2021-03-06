package email

import (
	"github.com/mailjet/mailjet-apiv3-go"
	log "github.com/sirupsen/logrus"

	cnfg "github.com/vantihovich/go_tasks/tree/master/swagger/configuration"
)

type Client struct {
	mjClient    *mailjet.Client
	senderEmail string
}

func New(cfg cnfg.MailJetParameters) Client {
	client := mailjet.NewMailjetClient(cfg.User, cfg.Password)
	return Client{
		mjClient:    client,
		senderEmail: cfg.SenderEmail,
	}
}

func (m Client) SendForgottenPasswordEmail(recipient, secret string) error {
	var messagesInfo = []mailjet.InfoMessagesV31{
		mailjet.InfoMessagesV31{
			From: &mailjet.RecipientV31{
				Email: m.senderEmail,
				Name:  "Messager from login services",
			},
			To: &mailjet.RecipientsV31{
				mailjet.RecipientV31{
					Email: recipient,
				},
			},
			TemplateLanguage: true,
			Subject:          "Your secret key for changing password!",
			TextPart:         "Dear logginer, here is your link for changing the password in our services.May the force be with you!",
			HTMLPart:         "<h3>Dear logginer, here is your <a href=\"http://localhost:3000/auth/forgot_password_reset_password?parameter={{var:secret}}\">link</a> for changing the password in our services.</h3><br />May the force be with you!",
			Variables: map[string]interface{}{
				"secret": secret},
			Headers: map[string]interface{}{
				"X-My-header": "X2332X-324-432-534"},
		},
	}

	messages := mailjet.MessagesV31{Info: messagesInfo}
	sendEmailResult, err := m.mjClient.SendMailV31(&messages)
	if err != nil {
		log.WithError(err).Info("Error sending the email")
		return err
	}
	log.Debug("A successfull attempt to send email:", sendEmailResult)
	return nil
}
