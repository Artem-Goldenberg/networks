package main

import (
	"crypto/tls"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/go-gomail/gomail"
)

const (
	server = "smtp.gmail.com"
	port   = 25
)

func main() {
	if len(os.Args) != 5 {
		fmt.Println("Usage: go run main.go <sender> <sender password> <receiver> <filename>")
		return
	}

	from := os.Args[1]
	password := os.Args[2]
	to := os.Args[3]
	filename := os.Args[4]
	extension := filepath.Ext(filename)

	var bodyType string
	if extension == ".txt" {
		bodyType = "text/plain"
	} else if extension == ".html" {
		bodyType = "text/html"
	} else {
		fmt.Printf("Invalid file format: %q\n", extension)
		return
	}

	content, err := ioutil.ReadFile(filename)
	if err != nil { 
		panic(err)
	}

	m := gomail.NewMessage()
	m.SetHeader("From", from)
	m.SetHeader("To", to)
	m.SetHeader("Subject", "testing")
	m.SetBody(bodyType, string(content))
	// m.SetBody("text/plain", "just a message")

	d := gomail.NewDialer("smtp.gmail.com", port, from, password)
	d.TLSConfig = &tls.Config{InsecureSkipVerify: true}

	if err := d.DialAndSend(m); err != nil {
		panic(err)
	}

	fmt.Println("Mail sent")

	return
}
