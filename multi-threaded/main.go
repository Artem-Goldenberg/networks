package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"net"
	"net/http"
	"strings"
)

func main() { 
	l, err := net.Listen("tcp", ":8080")
	check(err)
	defer l.Close()
	for { 
		conn, err := l.Accept()
		check(err)
		go handle(conn)
	}
}

func handle(conn net.Conn) { 
	defer conn.Close()

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

func readLines(reader *bufio.Reader) ([]string, error) { 
	var err error
	str := ""
	lines := make([]string, 0)
	for ;err != io.EOF; str, err = reader.ReadString('\n') { 
		if err != nil { return nil, err }
		lines = append(lines, str)
	}
	return lines, nil
}

func check(err error) { 
	if err != nil { 
		panic(err)
	}
}