package main

import (
	"crypto/tls"
	"errors"

	"gopkg.in/gomail.v2"
)

type (
	GoMailerConfig struct {
		Provider string
		Host     string
		Port     int
		Username string // email address
		Password string
		From     string
	}

	GoMailerRepository struct {
		GoMailerConfig
	}

	GoMailerForm struct {
		To      []string
		CC      []string
		BCC     []string
		Subject string
		Body    string
	}
)

func NewGoMailer(conf GoMailerConfig) *GoMailerRepository {
	return &GoMailerRepository{
		conf,
	}
}

func (conf *GoMailerRepository) SendEmail(data GoMailerForm) error {
	switch conf.Provider {
	case "gmail":
		return conf.SendEmailWithGmail(data)
	default:
		return errors.New("provider not found")
	}
}

func (conf *GoMailerRepository) SendEmailWithGmail(data GoMailerForm) error {
	mail := gomail.NewMessage()
	mail.SetHeader("From", conf.From)
	mail.SetHeader("To", data.To...)
	if len(data.CC) > 0 {
		mail.SetHeader("Cc", data.CC...)
	}
	if len(data.BCC) > 0 {
		mail.SetHeader("Bcc", data.BCC...)
	}
	mail.SetHeader("Subject", data.Subject)
	mail.SetBody("text/html", data.Body)
	// mail.Attach("./sample.png")

	dialer := gomail.NewDialer(conf.Host, conf.Port, conf.Username, conf.Password)
	dialer.TLSConfig = &tls.Config{InsecureSkipVerify: true}

	err := dialer.DialAndSend(mail)
	if err != nil {
		// TODO : log
		return err
	}

	return nil
}
