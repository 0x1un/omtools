package main

import (
	"crypto/tls"
	"gopkg.in/gomail.v2"
)

type SMTPConfig struct {
	server   string
	port     int
	username string
	password string
}

func NewSMTPSender(server string, port int, username string, password string) *SMTPConfig {
	return &SMTPConfig{
		server:   server,
		port:     port,
		username: username,
		password: password,
	}
}

func (s *SMTPConfig) SendToMail(msg *gomail.Message) error {
	d := gomail.NewDialer(s.server, s.port, s.username, s.password)
	d.TLSConfig = &tls.Config{InsecureSkipVerify: true}
	err := d.DialAndSend(msg)
	if err != nil {
		return err
	}
	return nil
}

func (s *SMTPConfig) writeMessage(body string, from string, filename string,subject string, to ...string) *gomail.Message {
	m := gomail.NewMessage()
	m.SetHeader("From", from)
	m.SetHeader("To", to...)
	m.SetBody("text/html", body)
	m.SetHeader("Subject", subject)
	m.Attach(filename)
	return m
}
