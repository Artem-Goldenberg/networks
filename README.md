# networks
Задания по курсу компьютерных сетей

# Задание 3

Чтобы запустить нужно установить golang https://go.dev/dl/

Потом перейти в нужную папку и выполнить нужную команду
* single-threaded - задание A  
```
go run main.go
```
также есть видео taskA.mov, показывает как работает сервер
* multi-threaded - многопоточный сервер из задания B  
```
go run main.go
```
 есть скриншоты сравнения, что по времени он работает лучше чем предыдущий
 (single-threaded.png ~1.5, multi-threaded ~0.5)
 проверено с помощью дополнительного threadsTest.go в папке client
* client - задание C
```
go run client.go -- host port filename
```
* limited-threads - задание D
```
go run main.go -- num
```