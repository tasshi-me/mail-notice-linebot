package mailmanager

import (
	"crypto/tls"
	"fmt"
	"log"
	"net"
	"net/mail"
	"net/smtp"
)

// SendMail with SMTPS
func SendMail(from, to mail.Address, subject, body, smptServerName, smtpAuthUser, smtpAuthPassword string) {

	// Parse header
	headers := make(map[string]string)
	headers["From"] = from.String()
	headers["To"] = to.String()
	headers["Subject"] = subject

	// Concat header and body
	var message string
	for key, value := range headers {
		message += fmt.Sprintf("%s: %s\r\n", key, value)
	}
	message += "\r\n"
	message += body

	// SMTP and SMTP-Auth Setting
	host, _, _ := net.SplitHostPort(smptServerName)
	auth := smtp.PlainAuth("", smtpAuthUser, smtpAuthPassword, host)

	// TLS setting
	tlsconfig := &tls.Config{ServerName: host}

	// Dial up SMTP Server
	c, err := smtp.Dial(smptServerName)
	if err != nil {
		log.Panic(err)
	}

	// Starting TLS
	c.StartTLS(tlsconfig)

	// SMTP-Auth over TLS
	if err = c.Auth(auth); err != nil {
		log.Panic(err)
	}

	// SMTP COMMAND: MAIL FROM
	if err = c.Mail(from.Address); err != nil {
		log.Panic(err)
	}

	// SMTP COMMAND: RCPT TO
	if err = c.Rcpt(to.Address); err != nil {
		log.Panic(err)
	}

	// SMTP COMMAND: DATA
	w, err := c.Data()
	if err != nil {
		log.Panic(err)
	}

	// Send DATA
	_, err = w.Write([]byte(message))
	if err != nil {
		log.Panic(err)
	}

	err = w.Close()
	if err != nil {
		log.Panic(err)
	}

	c.Quit()

}

// FetchMail fetch email using imaps
func FetchMail() {

}
