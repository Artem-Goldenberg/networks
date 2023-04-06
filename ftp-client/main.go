package main

import (
	"fmt"
	"os"

	"github.com/secsy/goftp"
)

func main() {
	// Create client object with default config
	client, err := goftp.DialConfig(goftp.Config{
		User: "admin",
		Password: "admin",
	}, "192.168.1.72:80")
	// client, err := goftp.Dial("ftp://192.168.1.72")
	if err != nil {
		panic(err)
	}

	arg := os.Args[1]
	if arg == "list" { 
		_ = client.Retrieve("some", os.Stdout)
	} else if arg == "upload" { 
		local := os.Args[2]
		remotePath := os.Args[3]
		data, err := os.Open(local)
		if err != nil { 
			panic(err)
		}
		client.Store(remotePath, data)
	} else if arg == "download" { 
		remotePath := os.Args[2]
		localPath := os.Args[3]
		file, err := os.Create(localPath)
		if err != nil { 
			panic(err)
		}
		err = client.Retrieve(remotePath, file)
		if err != nil { 
			panic(err)
		}
	} else { 
		fmt.Println("Usage: go run main.go serverAddress [list | upload | download] [path1 | path2]")
	}
}
