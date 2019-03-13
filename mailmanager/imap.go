package mailmanager

import (
	"log"
	"time"

	"github.com/emersion/go-imap"
	"github.com/emersion/go-imap/client"
)

// FetchMail fetch email using imaps
func FetchMail(timeSince, timeBefore time.Time, mboxName, imapServerName, imapAuthUser, imapAuthPassword string) []imap.Message {
	if timeSince.IsZero() && timeBefore.IsZero() {
		return nil
	}

	c, err := client.DialTLS(imapServerName, nil)
	if err != nil {
		log.Panic(err)
	}

	defer c.Logout()

	if err := c.Login(imapAuthUser, imapAuthPassword); err != nil {
		log.Panic(err)
	}

	_, err = c.Select(mboxName, false)
	if err != nil {
		log.Panic(err)
	}

	// Set search criteria
	criteria := imap.NewSearchCriteria()
	if !timeSince.IsZero() {
		criteria.Since = timeSince
	}
	if !timeBefore.IsZero() {
		criteria.Before = timeBefore
	}
	ids, err := c.Search(criteria)
	if err != nil {
		log.Fatal(err)
	}
	//log.Println("IDs found:", ids)

	if len(ids) > 0 {
		seqset := new(imap.SeqSet)
		seqset.AddNum(ids...)

		messages := make(chan *imap.Message, 10)
		done := make(chan error, 1)
		go func() {
			done <- c.Fetch(seqset, []imap.FetchItem{imap.FetchEnvelope}, messages)
		}()

		//log.Println(len(ids), "messages:")
		var messageEntities []imap.Message
		for msg := range messages {
			//log.Println(msg.Envelope.Date.String() + ":" + msg.Envelope.Subject)
			messageEntities = append(messageEntities, *msg)
		}

		if err := <-done; err != nil {
			log.Print(err)
		}

		return messageEntities

	}

	return nil
}

// DeleteMail :delete mails since specified datetime
func DeleteMail(timeSince, timeBefore time.Time, mboxName, imapServerName, imapAuthUser, imapAuthPassword string) {
	if timeSince.IsZero() && timeBefore.IsZero() {
		return
	}

	c, err := client.DialTLS(imapServerName, nil)
	if err != nil {
		log.Panic(err)
	}

	defer c.Logout()

	if err := c.Login(imapAuthUser, imapAuthPassword); err != nil {
		log.Panic(err)
	}

	_, err = c.Select(mboxName, false)
	if err != nil {
		log.Panic(err)
	}

	// Set search criteria
	criteria := imap.NewSearchCriteria()
	if !timeSince.IsZero() {
		criteria.Since = timeSince
	}
	if !timeBefore.IsZero() {
		criteria.Before = timeBefore
	}
	ids, err := c.Search(criteria)
	if err != nil {
		log.Fatal(err)
	}
	//log.Println("IDs found:", ids)

	if len(ids) > 0 {
		seqset := new(imap.SeqSet)
		seqset.AddNum(ids...)

		//Mark as Deleted
		item := imap.FormatFlagsOp(imap.AddFlags, true)
		flags := []interface{}{imap.DeletedFlag}
		if err := c.Store(seqset, item, flags, nil); err != nil {
			log.Panic(err)
		}

		if err := c.Expunge(nil); err != nil {
			log.Panic(err)
		}
	}

}

// PopMail :fetch and delete mails
func PopMail(timeSince, timeBefore time.Time, mboxName, imapServerName, imapAuthUser, imapAuthPassword string) []imap.Message {
	if timeSince.IsZero() && timeBefore.IsZero() {
		return nil
	}

	c, err := client.DialTLS(imapServerName, nil)
	if err != nil {
		log.Panic(err)
	}

	defer c.Logout()

	if err := c.Login(imapAuthUser, imapAuthPassword); err != nil {
		log.Panic(err)
	}

	_, err = c.Select(mboxName, false)
	if err != nil {
		log.Panic(err)
	}

	// Set search criteria
	criteria := imap.NewSearchCriteria()
	if !timeSince.IsZero() {
		criteria.Since = timeSince
	}
	if !timeBefore.IsZero() {
		criteria.Before = timeBefore
	}
	ids, err := c.Search(criteria)
	if err != nil {
		log.Fatal(err)
	}

	if len(ids) < 1 {
		return nil
	}

	seqset := new(imap.SeqSet)
	seqset.AddNum(ids...)

	messages := make(chan *imap.Message, 10)
	done := make(chan error, 1)
	go func() {
		done <- c.Fetch(seqset, []imap.FetchItem{imap.FetchEnvelope}, messages)
	}()

	//log.Println(len(ids), "messages:")
	var messageEntities []imap.Message
	var deleteIds []uint32
	for msg := range messages {
		//log.Println(msg.Envelope.Date.String() + ":" + msg.Envelope.Subject)
		messageEntities = append(messageEntities, *msg)
		deleteIds = append(deleteIds, msg.Uid)
	}

	if err := <-done; err != nil {
		log.Print(err)
	}

	// Delete fetched mails

	//Mark as Deleted
	item := imap.FormatFlagsOp(imap.AddFlags, true)
	flags := []interface{}{imap.DeletedFlag}
	if err := c.Store(seqset, item, flags, nil); err != nil {
		log.Panic(err)
	}

	if err := c.Expunge(nil); err != nil {
		log.Panic(err)
	}

	return messageEntities

}

// FilterMessageByRecipientAddress ...
func FilterMessageByRecipientAddress(messages []imap.Message, targetAddresses []*imap.Address) []imap.Message {
	slicedMessages := make([]imap.Message, 0, len(messages))

	// OPTIMIZE: Is there any faster search algorithm??
	for _, msg := range messages {
		var addresses []string
		for _, addr := range msg.Envelope.To {
			addresses = append(addresses, addr.MailboxName+"@"+addr.HostName)
		}
		for _, addr := range msg.Envelope.Cc {
			addresses = append(addresses, addr.MailboxName+"@"+addr.HostName)
		}
		for _, addr := range msg.Envelope.Bcc {
			addresses = append(addresses, addr.MailboxName+"@"+addr.HostName)
		}
	FOR_LABEL:
		for _, address := range addresses {
			for _, taddr := range targetAddresses {
				taddress := taddr.MailboxName + "@" + taddr.HostName
				if address == taddress {
					slicedMessages = append(slicedMessages, msg)
					break FOR_LABEL
				}
			}
		}
	}

	return slicedMessages
}
