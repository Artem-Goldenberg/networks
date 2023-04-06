package main

import (
	"fmt"
	"os"

	"github.com/secsy/goftp"
	"gopkg.in/yaml.v2"
)

func main() {
	configPath := os.Args[1]

	configs := make(map[string]interface{})
	bs, err := os.ReadFile(configPath)
	check(err)
	err = yaml.Unmarshal(bs, &configs)
	check(err)

	// Create client object with default config
	client, err := goftp.DialConfig(goftp.Config{
		User: configs["user"].(string),
		Password: configs["password"].(string),
	}, configs["host"].(string) + ":" + fmt.Sprint(configs["port"].(int)))
	if err != nil {
		panic(err)
	}

	arg := os.Args[2]
	if arg == "list" { 
		files, err := client.ReadDir("/")
		check(err)
		for _, file := range files { 
			fmt.Println(file.Name())
		}
	} else if arg == "upload" { 
		local := os.Args[3]
		remotePath := os.Args[4]
		data, err := os.Open(local)
		if err != nil { 
			panic(err)
		}
		client.Store(remotePath, data)
	} else if arg == "download" { 
		remotePath := os.Args[3]
		localPath := os.Args[4]
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

func check(err error) { 
	if err != nil { 
		panic(err)
	}
}