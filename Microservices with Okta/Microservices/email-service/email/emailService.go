package email

import (
	"fmt"
	"gopkg.in/gomail.v2"
	"log"
	"os"
)
  
type to struct{
	Name,Address string
}

func SendEmail(name,email,subject,message string){
	os.Setenv("serverpassrd","")
	reciever := new(to)
	reciever.Name=name
	reciever.Address=email
	SendMail(*reciever,subject,message,"calendar.ics")
}

func SendMail(to2 struct {Name string; Address string}, subject string, message string, filePath string) {

	d := gomail.NewDialer("smtp.gmail.com", 587, "sender@gmail.com", os.Getenv("serverpassrd"))
	s, err := d.Dial()
	if err != nil {
		panic(err)
	}

	m := gomail.NewMessage()
		m.SetHeader("From", "sender@gmail.com")
		m.SetAddressHeader("To", to2.Address, to2.Name)
		m.SetHeader("Subject", subject)
		m.SetBody("text/html", fmt.Sprintf("Hi %s!", to2.Name) + "\n " + message)

		m.Attach(filePath)

		if err := gomail.Send(s, m); err != nil {
			log.Printf("Could not send email to %q: %v", to2.Address, err)
		}
		m.Reset()
}
