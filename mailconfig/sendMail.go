package mailconfig

import (
	"bytes"
	"fmt"
	"net/smtp"
	"text/template"
)

func SendMail(m interface{}, index bool) {

	// Sender data.
	from := "devesh.tiwari12141@gmail.com"
	password := "8871789749"

	// Receiver email address.
	to := []string{
		"devesh.tiwari444@gmail.com",
		// "ashutosh2890@gmail.com",
	}

	// smtp server configuration.
	smtpHost := "smtp.gmail.com"
	smtpPort := "587"

	// Authentication.
	auth := smtp.PlainAuth("", from, password, smtpHost)

	t, err := template.ParseFiles("/home/devesh/go/src/gin-poc/mailconfig/template.html")
	if err != nil {
		fmt.Println(err.Error())
	}
	var body bytes.Buffer

	var b struct {
		Name    string
		Message interface{}
	}
	if index {
		b.Name = "Index"
	} else {
		b.Name = "Stock"
	}
	b.Message = m
	mimeHeaders := "MIME-version: 1.0;\nContent-Type: text/html; charset=\"UTF-8\";\n\n"
	body.Write([]byte(fmt.Sprintf("Subject: Found Arbitrage %v \n%s\n\n", b.Name, mimeHeaders)))
	err = t.Execute(&body, b)
	if err != nil {
		fmt.Println(err.Error())
	}
	// Sending email.
	err = smtp.SendMail(smtpHost+":"+smtpPort, auth, from, to, body.Bytes())
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("Email Sent!", b.Name)
}
