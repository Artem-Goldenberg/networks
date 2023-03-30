package main

import (
	"bytes"
	"crypto/tls"
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"net/smtp"
	"os"
)

const (
	server   = "smtp.gmail.com"
	port     = 25
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
	bytes1, err := ioutil.ReadFile(filename)
	checkError(err)

	con, err := net.Dial("tcp", "smtp.gmail.com:25")
	checkError(err)

	config := &tls.Config{InsecureSkipVerify: true}
	c, err := smtp.NewClient(con, "smtp.gmail.com")

	if ok, _ := c.Extension("STARTTLS"); ok {
		if err := c.StartTLS(config); err != nil {
			c.Close()
			panic(err)
		}
	}

	checkError(err)

	var auth smtp.Auth
	if ok, _ := c.Extension("AUTH"); ok {
		auth = smtp.PlainAuth("", from, password, server)
	}
	c.Auth(auth)

	c.Mail(from)
	c.Rcpt(to)
	w, _ := c.Data()
	bytes.NewBufferString(fmt.Sprintf("Content-Type: %s\n", http.DetectContentType(bytes1))).WriteTo(w)
	bytes.NewBufferString("Content-Transfer-Encoding: base64\n").WriteTo(w)
	b := make([]byte, base64.StdEncoding.EncodedLen(len(bytes1)))
	base64.StdEncoding.Encode(b, bytes1)
	bytes.NewBuffer(b).WriteTo(w)

	w.Close()

	fmt.Println("mail sent")
}

func checkError(err error) {
	if err != nil {
		panic(err)
	}
}
