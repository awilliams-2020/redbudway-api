package main

import (
    gomail "gopkg.in/gomail.v2"
)

func main() {

    msg := gomail.NewMessage()
    msg.SetHeader("From", "awilliams@redbudway.com")
    msg.SetHeader("To", "christinemhoo@gmail.com")
    msg.SetHeader("Subject", "This is a test subject")
    msg.SetBody("text/html", "<b>This is the body of the mail</b>")
    //msg.Attach("/home/User/cat.jpg")

    n := gomail.NewDialer("redbudway.com", 25, "awilliams", "MerCedEsAmgGt22$")

    // Send the email
    if err := n.DialAndSend(msg); err != nil {
        panic(err)
    }

}