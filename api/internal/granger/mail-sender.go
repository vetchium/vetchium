package granger

import (
	"context"
	"errors"
	"time"

	"github.com/psankar/vetchi/api/internal/db"
	"github.com/psankar/vetchi/api/pkg/libvetchi"
	"github.com/wneessen/go-mail"
)

func (g *Granger) mailSender(quit <-chan struct{}) {
	defer g.wg.Done()

	for {
		select {
		case <-quit:
			return
		case <-time.After(30 * time.Second):
			ctx := context.Background()
			emails, err := g.db.GetOldestUnsentEmails(ctx)
			if err != nil {
				continue
			}

			for _, email := range emails {
				err = g.sendEmail(email)
				if err != nil {
					continue
				}

				ctx := context.Background()
				err = g.db.UpdateEmailState(
					ctx,
					email.ID,
					db.EmailStateProcessed,
				)
				if err != nil {
					g.log.Error(
						"Updating email state",
						"error",
						err,
						"email",
						email.ID,
					)
					return
				}
			}

		}
	}
}

func (g *Granger) sendEmail(email db.Email) error {
	if len(email.EmailTo) == 0 {
		g.log.Error("Email has no recipients", "email", email.ID)
		return errors.New("email has no recipients")
	}

	var err error

	m := mail.NewMsg()
	if err = m.From(email.EmailFrom); err != nil {
		g.log.Error("failed to set From address", "error", err)
		return err
	}

	if err = m.To(email.EmailTo...); err != nil {
		g.log.Error("failed to set To address", "error", err)
		return err
	}

	if err = m.Cc(email.EmailCC...); err != nil {
		g.log.Error("failed to set Cc address", "error", err)
		return err
	}

	if err = m.Bcc(email.EmailBCC...); err != nil {
		g.log.Error("failed to set Bcc address", "error", err)
		return err
	}

	m.Subject(email.EmailSubject)
	m.SetBodyString(mail.TypeTextHTML, email.EmailHTMLBody)
	m.AddAlternativeString(mail.TypeTextPlain, email.EmailTextBody)

	g.log.Info("sending email", "email", email.ID, "env", g.env)
	var c *mail.Client
	if g.env == libvetchi.ProdEnv {
		c, err = mail.NewClient(
			g.SMTPHost,
			mail.WithPort(g.SMTPPort),
			mail.WithUsername(g.SMTPUser),
			mail.WithPassword(g.SMTPPassword),
			mail.WithSMTPAuth(mail.SMTPAuthLogin),
		)
		if err != nil {
			g.log.Error("failed to create PROD mail client", "error", err)
			return err
		}
	} else {
		c, err = mail.NewClient(
			g.SMTPHost,
			mail.WithPort(g.SMTPPort),
			mail.WithTLSPortPolicy(mail.NoTLS),
			mail.WithSMTPAuth(mail.SMTPAuthCustom),
		)
		if err != nil {
			g.log.Error("failed to create DEV mail client", "error", err)
			return err
		}
	}

	if err := c.DialAndSend(m); err != nil {
		g.log.Error("failed to send mail", "error", err)
		return err
	}

	g.log.Info("email sent", "email", email.ID, "to", email.EmailTo)

	return nil
}
