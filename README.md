# networks
Задания по курсу компьютерных сетей

## Клиент и сервер с протоколом Stop and Wait

```
go run server/main.go
go run client/main.go filename <timeout>
```
`<timeout>` - таймаут в секундах  
`filename` - файл чтобы отправить серверу
например `client/some.txt`

Сервер примет файл и сохранит у себя в файле recieve.txt

Пропадание пакетов не успел сделать
