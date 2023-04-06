# networks
Задания по курсу компьютерных сетей

Установить го   
Перейти в ftp-client  
```
go mod tidy
go run main.go download remotePath localPath
```

---
20 Minute After Deadline

- Доделал правильный вывод файлов на консоль
- Теперь клиент читает имя сервера, порт, логин и пароль из конфигурационного файла, к которому получает путь.
 Стандартный файл это `config.yaml` 

 Теперь варианты запуска такие
 ```
 go run main.go config.yaml list
 go run main.go config.yaml download remotePath localPath
 go run main.go config.yaml upload localPath remotePath
 ```