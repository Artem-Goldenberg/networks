package main

import (
	"bufio"
	_ "embed"
	"io/ioutil"
	"log"
	"net"
	"strconv"
	"strings"
	"time"
)

// When program runs here will be bytes from file config.txt
//
//go:embed config.txt
var blacklistFile []byte

// blocked urls and servers
var blacklist map[string]struct{} = make(map[string]struct{})

func main() {
	for _, line := range strings.Split(string(blacklistFile), "\n") {
		blacklist[line] = struct{}{} // add line to blacklist set
	}

	l, err := net.Listen("tcp", ":8080")
	check(err)
	log.Println("Server is running")

	defer l.Close()
	for {
		conn, err := l.Accept()
		check(err)
		handle(conn)
	}
}

type entry = struct {
	url          string
	responseCode int
}

var history []entry = make([]entry, 0)

func handle(conn net.Conn) {
	defer conn.Close()

	reader := bufio.NewReader(conn)
	line, _ := reader.ReadString('\n')

	words := strings.Split(line, " ")
	requestType := words[0]
	if requestType != "GET" && requestType != "POST" {
		badRequest(conn, line, nil)
		return
	}

	fullPath := strings.Split(words[1][1:], "/")
	address := fullPath[0]

	var path string
	if len(fullPath) > 1 {
		path = strings.Join(fullPath[1:], "/")
	} else {
		path = ""
	}
	
	if _, ok := blacklist[address + path]; ok { 
		// url is blocked
		conn.Write([]byte("HTTP/1.1 403 Forbidden\r\n\r\n"))
		log.Printf("Blocked an attempt to access url: %q\n", address + path)
		return
	}

	proxyConn, err := net.Dial("tcp", address+":80")
	if err != nil {
		badRequest(conn, address, err)
		return
	}
	defer proxyConn.Close()

	log.Printf("Connected to %q\n", address)

	headLine := requestType + " /" + path + " HTTP/1.1\r\n"

	hostLine, err := reader.ReadString('\n')
	if strings.Split(hostLine, " ")[0] == "Host:" {
		// modify host in proxy request
		hostLine = "Host: " + address + ":80" + "\r\n"
	}

	conn.SetReadDeadline(time.Now().Add(time.Second))
	requestBytes, _ := ioutil.ReadAll(reader)

	log.Printf("Formed request head line: %q\n", headLine)
	log.Printf("bytes: %q\n", string(requestBytes))
	proxyConn.Write(append([]byte(headLine+hostLine), requestBytes...))

	responseReader := bufio.NewReader(proxyConn)
	line, err = responseReader.ReadString('\n')
	check(err)
	log.Printf("Got response: %q\n", line)

	word := strings.Split(string(line), " ")[1]
	responseCode, err := strconv.Atoi(word)
	check(err)
	log.Printf("Response code: %d\n", responseCode)

	history = append(history, entry{address, responseCode})
	printHistory()

	proxyConn.SetReadDeadline(time.Now().Add(time.Second))
	responseBytes, _ := ioutil.ReadAll(responseReader)

	readLine := []byte(line)
	conn.Write(append(readLine, responseBytes...))

	log.Println("Written response to client")
}

func printHistory() {
	log.Println("History by now:")
	for _, entry := range history {
		log.Printf("%q\t:\t%d\n", entry.url, entry.responseCode)
	}
}

func badRequest(conn net.Conn, req string, err error) {
	conn.Write([]byte("HTTP/1.1 400 Bad Request\r\n\r\n"))
	if err != nil {
		log.Printf("Bad request: %q, error: %v\n", req, err)
	} else {
		log.Printf("Bad request: %q\n", req)
	}
}

func check(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
