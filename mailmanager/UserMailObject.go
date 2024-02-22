package mailmanager

import (
	"github.com/tasshi-me/mail-notice-linebot/mongodb"

	"github.com/emersion/go-imap"
)

// MailObject ..
type MailObject struct {
	TargetLineID        string
	MailFromName        string
	MailFromAddress     string
	MailReceivedAddress string
	MailSubject         string
}

// UserMailObject ..
type UserMailObject struct {
	TargetLineID string
	MailObjects  []MailObject
}

// ConvertMessagesToUserMailObject ..
func ConvertMessagesToUserMailObject(messages []imap.Message, lineUsers []mongodb.LineUser) []UserMailObject {
	var userMailObjects []UserMailObject
	for _, lineUser := range lineUsers {
		var mailObjects []MailObject
	MSG_LOOP:
		for _, msg := range messages {
			var addresses []*imap.Address
			for _, addr := range msg.Envelope.To {
				addresses = append(addresses, addr)
			}
			for _, addr := range msg.Envelope.Cc {
				addresses = append(addresses, addr)
			}
			for _, addr := range msg.Envelope.Bcc {
				addresses = append(addresses, addr)
			}

			for _, registeredAddress := range lineUser.RegisteredAddresses {
				for _, address := range addresses {
					fullAddress := address.MailboxName + "@" + address.HostName
					if registeredAddress == fullAddress {
						mailFromAddress := msg.Envelope.From[0]
						mailObject := MailObject{
							TargetLineID:        lineUser.LineID,
							MailFromName:        mailFromAddress.MailboxName + "@" + mailFromAddress.HostName,
							MailFromAddress:     mailFromAddress.PersonalName,
							MailReceivedAddress: fullAddress,
							MailSubject:         msg.Envelope.Subject,
						}
						mailObjects = append(mailObjects, mailObject)
						//log.Println(mailObject)
						continue MSG_LOOP
					}
				}
			}
		}
		if len(mailObjects) > 0 {
			userMailObject := UserMailObject{
				TargetLineID: lineUser.LineID,
				MailObjects:  mailObjects,
			}
			userMailObjects = append(userMailObjects, userMailObject)
		}
	}

	return userMailObjects
}
