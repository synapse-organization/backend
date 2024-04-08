package utils

import "gopkg.in/gomail.v2"

func SendEmail(to, subject, body string) error {
	m := gomail.NewMessage()
	m.SetHeader("To", to)
	m.SetHeader("From", "synapseteam.proj@gmail.com")
	m.SetHeader("Subject", subject)
	m.SetBody("text/html", body)

	d := gomail.NewDialer("smtp.gmail.com", 587, "synapseteam.proj@gmail.com", "xoscrmcssxbfmgqe")

	// Send the email to Bob, Cora and Dan.
	if err := d.DialAndSend(m); err != nil {
		return err
	}
	return nil
}
