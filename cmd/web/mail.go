package main

import (
	"bytes"
	"fmt"
	"html/template"
	"sync"
	"time"

	"github.com/vanng822/go-premailer/premailer" //premailer designs the css that wll be displayed in the html
	mail "github.com/xhit/go-simple-mail/v2"
)

// mail server
type Mail struct {
	Domain      string
	Host        string
	Port        int
	Username    string
	Password    string
	Encryption  string
	FromAddress string
	FromName    string
	Wait        *sync.WaitGroup
	MailerChan  chan Message
	ErrorChan   chan error
	DoneChan    chan bool
}

// items to be sent to clients mail
type Message struct {
	From          string
	FromName      string
	To            string
	Subject       string
	Attachments   []string
	AttachmentMap map[string]string
	Data          any
	DataMap       map[string]any
	Template      string
}

// a  helper wrapper to send email easily
func (app *Config) sendemail(msg Message) {
	//add counter to waitgroup , increment wg by 1
	app.Wait.Add(1)
	/*
		send message to the Mail channel(speaks to the MailerChan object in the Mail
		struct which will be accessed by any of the handlers)
	*/
	app.Mailer.MailerChan <- msg
}

// a func that listens for different channels in the Channel objects
// defined in the Mail struct
func (app *Config) ListenForMail() {
	//listen to diff channels
	for {
		select {
		//receiving data from the Mail channel
		case msg := <-app.Mailer.MailerChan:
			//send email to the go routine
			go app.Mailer.SendMail(msg, app.Mailer.ErrorChan)

		case err := <-app.Mailer.ErrorChan:
			app.ErrorLog.Println("error in mail", err)
		case <-app.Mailer.DoneChan:
			return
		}
	}

}

// Displays the email message
// Displays the email sent by the goroutine from the ListenForMail()
func (m *Mail) SendMail(msg Message, errorChan chan error) {
	//decrement wait group
	defer m.Wait.Done()

	if msg.Template == "" {
		//gets both the plin text and html message
		msg.Template = "mail"
	}
	if msg.From == "" {
		msg.From = m.FromAddress
	}
	if msg.FromName == "" {
		msg.FromName = m.FromName
	}
	if msg.AttachmentMap == nil {
		msg.AttachmentMap = make(map[string]string)
	}
	//send info to the 2 templates
	if len(msg.DataMap) == 0 {
		msg.DataMap = make(map[string]any)
	}
	//.message is called from the template and
	//displays the message in the template
	msg.DataMap["message"] = msg.Data //calls the key field data and maps with a value interface

	// build html mail
	formattedMessage, err := m.buildHTMLMessage(msg)
	if err != nil {
		//sending error to the channel
		errorChan <- err
	}
	//build plain text mail
	plainMessage, err := m.buildPlainTextMessage(msg)
	if err != nil {
		errorChan <- err
	}
	server := mail.NewSMTPClient()
	server.Host = "smtp.gmail.com"
	server.Port = 587
	server.Username = m.Username
	server.Password = m.Password
	// server.Host = m.Host
	// server.Port = m.Port
	// server.Username = m.Username
	// server.Password = m.Password
	server.Encryption = m.getEncryption(m.Encryption)
	server.KeepAlive = false
	server.ConnectTimeout = 10 * time.Second
	server.SendTimeout = 10 * time.Second

	smtpClient, err := server.Connect()
	if err != nil {
		errorChan <- err
	}
	email := mail.NewMSG()
	email.SetFrom(msg.From).AddTo(msg.To).SetSubject(msg.Subject)

	email.SetBody(mail.TextPlain, plainMessage)
	email.AddAlternative(mail.TextHTML, formattedMessage)

	if len(msg.Attachments) > 0 {
		for _, x := range msg.Attachments {
			email.AddAttachment(x)
		}
	}
	//i.e if the map isn't empty
	if len(msg.AttachmentMap) > 0 {
		for key, value := range msg.AttachmentMap {
			email.AddAttachment(value, key)
		}
	}
	err = email.Send(smtpClient)
	if err != nil {
		errorChan <- err
	}

}

func (m *Mail) buildHTMLMessage(msg Message) (string, error) {
	templateToRender := fmt.Sprintf("./cmd/web/templates/%s.html.gohtml", msg.Template) //display the name of the template called from the handlers

	t, err := template.New("email-html").ParseFiles(templateToRender)

	if err != nil {
		return "", err
	}

	var tpl bytes.Buffer
	if err = t.ExecuteTemplate(&tpl, "body", msg.DataMap); err != nil {
		return "", err
	}
	formattedMessage := tpl.String()
	formattedMessage, err = m.inlineCSS(formattedMessage)
	if err != nil {
		return "", err
	}
	return formattedMessage, nil
}

func (m *Mail) buildPlainTextMessage(msg Message) (string, error) {
	templateToRender := fmt.Sprintf("./cmd/web/templates/%s.plain.gohtml", msg.Template)

	t, err := template.New("email-plain").ParseFiles(templateToRender)
	if err != nil {
		return "", err
	}
	var tpl bytes.Buffer
	if err = t.ExecuteTemplate(&tpl, "body", msg.DataMap); err != nil {
		return "", err
	}
	plainMessage := tpl.String()

	return plainMessage, nil
}

func (m *Mail) inlineCSS(s string) (string, error) {
	options := premailer.Options{
		RemoveClasses:     false,
		CssToAttributes:   false,
		KeepBangImportant: true,
	}

	prem, err := premailer.NewPremailerFromString(s, &options)
	if err != nil {
		return "", err
	}

	html, err := prem.Transform()
	if err != nil {
		return "", err
	}

	return html, nil
}

func (m *Mail) getEncryption(e string) mail.Encryption {
	switch e {
	case "tls":
		return mail.EncryptionSTARTTLS
	case "ssl":
		return mail.EncryptionSSLTLS
	case "none":
		return mail.EncryptionNone
	default:
		return mail.EncryptionSTARTTLS
	}
}
