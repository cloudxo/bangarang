package email

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log"
	"net/mail"
	"net/smtp"
	"strconv"
	"strings"

	"github.com/Sirupsen/logrus"
	"github.com/eliothedeman/bangarang/escalation"
	"github.com/eliothedeman/bangarang/event"
)

func init() {
	escalation.LoadFactory("email", NewEmail)
}

type Email struct {
	conf *EmailConfig
	Auth *smtp.Auth
}

func NewEmail() escalation.Escalation {
	e := &Email{
		conf: &EmailConfig{},
		Auth: nil,
	}
	return e
}

func encodeRFC2047(String string) string {
	//Make sure mail is rfc2047 compliant
	addr := mail.Address{String, ""}
	return strings.Trim(addr.String(), " <>")
}

func writeEmailBuffer(headers map[string]string, body string) string {
	//Could pre-allocate an email buffer, but since an alert is a rather exceptional
	//	event the gains are negligable
	buf := ""
	for k, v := range headers {
		buf += fmt.Sprintf("%s: %s\r\n", k, v)
	}
	buf += "\r\n" + base64.StdEncoding.EncodeToString([]byte(body))
	return buf
}

// Send an email via smtp
func (e *Email) Send(i *event.Incident) error {

	//For now set the description as both the subject and body
	headers := make(map[string]string)
	headers["From"] = e.conf.Sender
	headers["To"] = strings.Join(e.conf.Recipients, ",")
	headers["Subject"] = i.FormatDescription()
	headers["MIME-Version"] = "1.0"
	headers["Content-Type"] = "text/plain; charset=\"utf-8\""
	headers["Content-Transfer-Encoding"] = "base64"

	// make the body a json encoded representation of the incidnent
	body, err := json.MarshalIndent(i, "", "    ")
	if err != nil {
		logrus.Errorf("Unable to encode incidnet for email %s", err.Error())
	}

	log.Println("sending email")
	err = smtp.SendMail(e.conf.Host+":"+strconv.Itoa(e.conf.Port), *e.Auth,
		e.conf.Sender, e.conf.Recipients, []byte(writeEmailBuffer(headers, string(body))))
	if err != nil {
		logrus.Errorf("Unable to send email via smtp %s", err)
	}
	log.Println("done sending mail")
	return err
}

func (e *Email) ConfigStruct() interface{} {
	return e.conf
}

func (e *Email) Init(conf interface{}) error {
	auth := smtp.PlainAuth("", e.conf.User, e.conf.Password, e.conf.Host)
	e.Auth = &auth
	return nil
}

type EmailConfig struct {
	Sender     string   `json:"source_email"`
	Recipients []string `json:"dest_emails"`
	Host       string   `json:"host"`
	User       string   `json:"user"`
	Password   string   `json:"password"`
	Port       int      `json:"port"`
}
