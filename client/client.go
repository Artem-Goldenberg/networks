package main

import (
	"fmt"
	"os"
	"net"
	"io/ioutil"
)

func main() { 
	// [client.go, --, host, port, file]
	if len(os.Args) != 5 { 
		fmt.Print("Usage: go run client.go -- <server> <port> <filename>")
		return
	}

	host := os.Args[2]
	port := os.Args[3]
	filename := os.Args[4]

	conn, err := net.Dial("tcp", host + ":" + port)
	check(err)
	defer conn.Close()

	request := 
		"GET /" + filename + " HTTP/1.1\r\n" +
		"Host: " + host + "\r\n" +
		"User-Agent: Go client\r\n\r\n"
	
	_, err = conn.Write([]byte(request))
	check(err)

	responseData, err := ioutil.ReadAll(conn)
	check(err)

	fmt.Println(string(responseData))
}

func check(err error) { 
	if err != nil { 
		panic(err)
	}
}