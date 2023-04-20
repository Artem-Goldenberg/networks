package main

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"math"
	"net"
	"os"
	"time"
)

// const MaxPackageSize = 10000 // 10 Kb
const MaxPackageSize = 2 // 10 Kb

func main() { 
	// if len(os.Args) < 2 { 
	// 	fmt.Println("Usage: go run client/go filepath")
	// 	return
	// }
	data, err := os.ReadFile("some.txt")
	check(err)

	segments := split(data)
	
	conn, err := net.Dial("udp", ":80")
	check(err)
	defer conn.Close()

	// let server know that we are about to send some segments
	conn.Write([]byte(fmt.Sprintf("SEND %d", len(segments))))
	time.Sleep(time.Second)

	timeout := 3 * time.Second

	flip := 0
	for _, segment := range segments { 
		p := makePackage(flip, segment)
		conn.Write(p)

		conn.SetReadDeadline(time.Now().Add(timeout))
		packageBuf := make([]byte, headerLength)
		_, err := conn.Read(packageBuf)
		for !isAck(flip, packageBuf) {
			for errors.Is(err, os.ErrDeadlineExceeded) { 
				conn.SetReadDeadline(time.Now().Add(timeout))
				_, err = conn.Read(packageBuf)
			}
			if err != nil { 
				fmt.Println(err)
				return
			}
		}

		flip = (flip + 1) % 2
	}
	n, err := conn.Write(data)
	if n != len(data) { 
		fmt.Printf("Written: %d, have: %d\n", n, len(data))
	}
	check(err)
}

const headerLength = 6
// package :: 
// length (4 bytes)
// 0 / 1 (1 byte)
// isAck (1 byte)
// data
func makePackage(num int, data []byte) []byte { 
	l := int32(headerLength + len(data))
	n := int8(num)
	buf := new(bytes.Buffer)
	fmt.Printf("Package: length: %d, num: %d, conte: %q\n", l, n, string(data))
	binary.Write(buf, binary.LittleEndian, l) // length
	binary.Write(buf, binary.LittleEndian, n) // 0 / 1
	binary.Write(buf, binary.LittleEndian, int8(0)) // isAck
	return []byte(string(buf.Bytes()) + string(data))
} 

func makeAckPackage(num int) []byte { 
	l := int32(headerLength)
	n := int8(num)
	buf := new(bytes.Buffer)
	binary.Write(buf, binary.LittleEndian, l)
	binary.Write(buf, binary.LittleEndian, n)
	binary.Write(buf, binary.LittleEndian, int8(1))
	return buf.Bytes()
}

func isAck(num int, data []byte) bool { 
	n := int8(data[4])
	isAck := int8(data[5])
	return isAck == 1 && int(n) == num
}

func split(data []byte) [][]byte {
	n := int(math.Ceil(float64(len(data)) / MaxPackageSize))
	segments := make([][]byte, n)
	for i := 0; i < n - 1; i++ { 
		segments[i] = data[i * MaxPackageSize : (i + 1) * MaxPackageSize]
	}
	segments[n - 1] = data[(n - 1) * MaxPackageSize :]
	return segments
}

func check(err error) {
	if err != nil {
		panic(err)
	}
}