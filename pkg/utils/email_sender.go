package utils

import "gopkg.in/gomail.v2"

func SendEmail(to, subject, body string) error {
	m := gomail.NewMessage()
	m.SetHeader("To", to)
	m.SetHeader("From", "django.blog.test.send@gmail.com")
	m.SetHeader("Subject", subject)
	m.SetBody("text/html", body)

	d := gomail.NewDialer("smtp.gmail.com", 587, "django.blog.test.send@gmail.com", "qunrvaxxahqzxzhn")

	// Send the email to Bob, Cora and Dan.
	if err := d.DialAndSend(m); err != nil {
		return err
	}
	return nil
}
