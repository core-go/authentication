package auth

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/common-go/mail"
)

type VerifiedCodeEmailSender struct {
	MailSender     mail.SimpleMailSender
	From           mail.Email
	TemplateLoader mail.TemplateLoader
}

func NewVerifiedCodeEmailSender(mailSender mail.SimpleMailSender, from mail.Email, templateLoader mail.TemplateLoader) *VerifiedCodeEmailSender {
	return &VerifiedCodeEmailSender{mailSender, from, templateLoader}
}

func truncatingSprintf(str string, args ...interface{}) string {
	n := strings.Count(str, "%s")
	if n > len(args) {
		n = len(args)
	}
	return fmt.Sprintf(str, args[0:n]...)
}

func (s *VerifiedCodeEmailSender) Send(ctx context.Context, to string, code string, expireAt time.Time, params interface{}) error {
	diff := expireAt.Sub(time.Now())
	strDiffMinutes := fmt.Sprint(diff.Minutes())
	subject, template, err := s.TemplateLoader.Load(ctx, to)
	if err != nil {
		return err
	}
	if strings.Index(subject, "%s") >= 0 {
		subject = fmt.Sprintf(subject, code)
	}
	content := truncatingSprintf(template,
		code, strDiffMinutes,
		code, strDiffMinutes,
		code, strDiffMinutes,
		code, strDiffMinutes,
		code, strDiffMinutes,
		code, strDiffMinutes,
		code, strDiffMinutes,
		code, strDiffMinutes,
		code, strDiffMinutes,
		code, strDiffMinutes)

	toMail := params.(string)
	mailTo := []mail.Email{{Address: toMail}}
	mailData := mail.NewSimpleHtmlMail(s.From, subject, mailTo, nil, content)
	return s.MailSender.Send(*mailData)
}
