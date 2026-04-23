package utils

import (
	"fmt"
	"os"
	"strconv"

	"gopkg.in/gomail.v2"
)

type EmailData struct {
	To      string
	Subject string
	Body    string
}

func SendEmail(data EmailData) error {
	host := os.Getenv("SMTP_HOST")
	portStr := os.Getenv("SMTP_PORT")
	user := os.Getenv("SMTP_USER")
	pass := os.Getenv("SMTP_PASS")

	port, err := strconv.Atoi(portStr)
	if err != nil {
		return fmt.Errorf("invalid SMTP port: %v", err)
	}

	m := gomail.NewMessage()
	m.SetHeader("From", user)
	m.SetHeader("To", data.To)
	m.SetHeader("Subject", data.Subject)
	m.SetBody("text/html", data.Body)

	d := gomail.NewDialer(host, port, user, pass)

	// Send the email
	if err := d.DialAndSend(m); err != nil {
		return err
	}

	return nil
}

func BuildInquiryEmailBody(customerName, customerEmail, customerPhone, message, productName string) string {
	return fmt.Sprintf(`
		<h2>New Product Inquiry</h2>
		<p><strong>Product:</strong> %s</p>
		<p><strong>Customer Name:</strong> %s</p>
		<p><strong>Customer Email:</strong> %s</p>
		<p><strong>Customer Phone:</strong> %s</p>
		<p><strong>Message:</strong></p>
		<p>%s</p>
		<hr>
		<p>This email was sent from Machine Katalog system.</p>
	`, productName, customerName, customerEmail, customerPhone, message)
}
