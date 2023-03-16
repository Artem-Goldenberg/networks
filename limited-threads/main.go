package main

import (
	"bufio"
	"fmt"
	"os"
	"net"
	"net/http"
	"strings"
	"strconv"
)

// dummy object
type coin struct{}

func main() { 
	if len(os.Args) != 3 { 
		fmt.Println("Usage: go run main.go -- <concurrency-level: int>")
		return
	}

	limit, err := strconv.Atoi(os.Args[2])
	check(err)

	l, err := net.Listen("tcp", ":8080")
	check(err)
	defer l.Close()

	// only `limit` threads at the same time
	lock := make(chan coin, limit)
	for { 
		// block if channel is filled
		lock <- coin{}

		conn, err := l.Accept()
		check(err)
		go func() { 
			handle(conn)
			// free space when finished
			<-lock
			conn.Close()
		}()
	}
}

func handle(conn net.Conn) { 
	// defer conn.Close()

	reader := bufio.NewReader(conn)
	line, err := reader.ReadString('\n')
	check(err)

	words := strings.Split(line, " ")
	fmt.Print(words)
	if words[0] == "GET" {
		path := words[1][1:]
		data, err := os.ReadFile(path)
		if err != nil { 
			conn.Write([]byte("HTTP/1.1 404 Not Found\r\n\r\n"))
			fmt.Print("Written error")
			return
		}
		responseString := 
			"HTTP/1.1 200 OK\r\n" + 
			"Content-Type: " + http.DetectContentType(data) + "\r\n" + 
			"Content-Length: " + fmt.Sprintf("%d", len(data)) + "\r\n\r\n"

		conn.Write(append([]byte(responseString), data...))
		fmt.Printf("Written response: %v\n", responseString)
	} else { 
		conn.Write([]byte("HTTP/1.1 400 Bad Request\r\n\r\n"))
	}
}

func check(err error) { 
	if err != nil { 
		panic(err)
	}
}