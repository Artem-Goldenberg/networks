package main

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"log"
	"net"
	"os"
	"strconv"
	"strings"
)

const MaxPackageSize = 10000 // 10 Kb
const headerLength = 6

type Package struct { 
	Data []byte
}

func main() {
	l, err := net.ListenPacket("udp", ":80")
	check(err)
	defer l.Close()

	for {
		buffer := make([]byte, 1000)
		n, address, err := l.ReadFrom(buffer)
		if err != nil {
			log.Printf("Error from address: %q when reading: %q\n", address, string(buffer[:n]))
			return
		}
		// log.Printf("Received: %q, from: %q\n", string(buffer[:n]), address)
		words := strings.Split(string(buffer[:n]), " ") 
		if words[0] == "SEND" { 
			n, err := strconv.Atoi(words[1])
			check(err)
			data := initiateProtocol(l, n)
			file, _ := os.Create("recieved.txt")
			file.Write(data)
			file.Close()
			// log.Printf("Received: %q, from: %q\n", string(data), address)
		}
	}
}

func initiateProtocol(l net.PacketConn, n int) []byte { 
	all_data := new(bytes.Buffer)
	buf := make([]byte, MaxPackageSize + headerLength)
	for i := 0; i < n; i++ {
		n, address, err := l.ReadFrom(buf)
		log.Printf("String: %q\n", string(buf[6:26]))
		check(err)
		l.WriteTo(makeAckPackage(getNum(buf)), address)
		length := getLength(buf)
		log.Printf("Received %d bytes, length = %d\n", n, length)
		if len(buf) >= headerLength {
			// fmt.Fprint(all_data, buf[headerLength:length])
			all_data.Write(buf[headerLength:length])
		}
	}
	return all_data.Bytes()
}

func getLength(data []byte) int { 
	return int(binary.LittleEndian.Uint32(data[:4]))
}

func getNum(data []byte) int {
	return int(data[4])
}

func makeAckPackage(num int) []byte { 
	l := int32(headerLength)
	n := int8(num)
	buf := new(bytes.Buffer)
	fmt.Printf("Package: length: %d, num: %d\n", l, n)
	binary.Write(buf, binary.LittleEndian, l)
	binary.Write(buf, binary.LittleEndian, n)
	binary.Write(buf, binary.LittleEndian, int8(1))
	return buf.Bytes()
}

func check(err error) {
	if err != nil {
		panic(err)
	}
}
