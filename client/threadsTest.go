package main

import (
	"fmt"
	// "os"
	"net"
	"io/ioutil"
)

func main() { 
	for i := 1; i <= 100; i++ { 
		conn, err := net.Dial("tcp", "localhost:8080")
		check(err)
		filename := "samples/" + fmt.Sprintf("%d", i % 11 + 1) + ".txt"
		request := 
			"GET /" + filename + " HTTP/1.1\r\n" +
			"Host: " + "localhost" + "\r\n" +
			"User-Agent: Go client " + filename + "\r\n\r\n"
		_, err = conn.Write([]byte(request))
		check(err)
		fmt.Println("Written")

		responseData, err := ioutil.ReadAll(conn)
		check(err)

		fmt.Println(string(responseData))
		// conn.Close()
	}
}

func check(err error) { 
	if err != nil { 
		panic(err)
	}
}
