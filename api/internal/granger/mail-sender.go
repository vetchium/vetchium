package granger

import (
	"context"
	"errors"
	"time"

	"github.com/wneessen/go-mail"

	"github.com/psankar/vetchi/api/internal/db"
	"github.com/psankar/vetchi/api/pkg/vetchi"
)

func (g *Granger) mailSender(quit <-chan struct{}) {
	defer g.wg.Done()

	for {
		select {
		case <-quit:
			g.log.Dbg("mailSender quitting")
			return
		case <-time.After(5 * time.Second):
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
					db.EmailStateChange{
						EmailDBKey: email.EmailKey,
						EmailState: db.EmailStateProcessed,
					},
				)
				if err != nil {
					continue
				}
			}

		}
	}
}

func (g *Granger) sendEmail(email db.Email) error {
	if len(email.EmailTo) == 0 {
		g.log.Err("Email has no recipients", "email", email.EmailKey)
		return errors.New("email has no recipients")
	}

	var err error

	m := mail.NewMsg()
	if err = m.From(email.EmailFrom); err != nil {
		g.log.Err("failed to set From address", "error", err)
		return err
	}

	if err = m.To(email.EmailTo...); err != nil {
		g.log.Err("failed to set To address", "error", err)
		return err
	}

	if err = m.Cc(email.EmailCC...); err != nil {
		g.log.Err("failed to set Cc address", "error", err)
		return err
	}

	if err = m.Bcc(email.EmailBCC...); err != nil {
		g.log.Err("failed to set Bcc address", "error", err)
		return err
	}

	m.Subject(email.EmailSubject)
	m.SetBodyString(mail.TypeTextHTML, email.EmailHTMLBody)
	m.AddAlternativeString(mail.TypeTextPlain, email.EmailTextBody)

	g.log.Inf("sending email", "email", email.EmailKey, "env", g.env)
	var c *mail.Client
	if g.env == vetchi.ProdEnv {
		c, err = mail.NewClient(
			g.smtp.host,
			mail.WithPort(g.smtp.port),
			mail.WithUsername(g.smtp.user),
			mail.WithPassword(g.smtp.password),
			mail.WithSMTPAuth(mail.SMTPAuthLogin),
		)
		if err != nil {
			g.log.Err("failed to create PROD mail client", "error", err)
			return err
		}
	} else {
		c, err = mail.NewClient(
			g.smtp.host,
			mail.WithPort(g.smtp.port),
			mail.WithTLSPortPolicy(mail.NoTLS),
			mail.WithSMTPAuth(mail.SMTPAuthCustom),
		)
		if err != nil {
			g.log.Err("failed to create DEV mail client", "error", err)
			return err
		}
	}

	if err := c.DialAndSend(m); err != nil {
		g.log.Err("failed to send mail", "error", err)
		return err
	}

	g.log.Inf("email sent", "email", email.EmailKey, "to", email.EmailTo)

	return nil
}
