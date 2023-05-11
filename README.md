# networks
Задания по курсу компьютерных сетей


# Отслеживание пути (traceroute)

Программа pathfind рабоатет аналогично `traceroute`
```
make pathfind
./pathfind <host> <numRequests>
```
`<host>` - Хост назначения, либо имя либо ip4 адресс

`<numRequests>` - Опциональный количество посылок с одинм и тем же `ttl` (default: 3)

Пример работы 
```
$ ./pathfind stanford.edu 4     
Tracing web.stanford.edu: 171.67.215.200 ...
 1. 192.168.1.1 (192.168.1.1) -- 3.476 ms 3.893 ms 3.457 ms 3.114 ms 
 2. 10.213.35.254 (10.213.35.254) -- 3.6 ms 3.486 ms 3.745 ms 3.428 ms 
 3. 10.210.116.25 (10.210.116.25) -- 3.758 ms 3.71 ms 3.578 ms 3.685 ms 
 4. 10.210.116.26 (10.210.116.26) -- 4.167 ms 4.392 ms 4.521 ms 4.925 ms 
 5. 10.210.116.30 (10.210.116.30) -- 4.808 ms 4.312 ms 4.499 ms 4.099 ms 
 6. 10.210.116.22 (10.210.116.22) -- 4.505 ms 5.545 ms 4.51 ms 4.5 ms 
 7. * * * * 
 8. 100ge14-1.core1.osl1.he.net (184.105.64.230) -- 26.823 ms 23.922 ms 24.028 ms 22.832 ms 
 9. * * * * 
10. * * * * 
11. * * * * 
12. * * * * 
13. * * * * 
14. stanford-university.100gigabitethernet5-1.core1.pao1.he.net (184.105.177.238) -- 165.826 ms 167.735 ms 164.994 ms 172.061 ms 
15. woa-west-rtr-vl2.sunet (171.64.255.132) -- 165.872 ms 165.097 ms 164.941 ms 165.167 ms 
16. * * * * 
17. web.stanford.edu (171.67.215.200) -- 173.864 ms 166.649 ms 166.918 ms 166.733 ms 
```

# Клиент и сервер на IPv6
Базовый клиент и сервер. Собрать так же:
```
cd ipv6
make client
make server
./server <host> <port>
./client <host> <port> <message>
```
`<host>` - ipv6 адрес

`<port>` - порт (число)

`message` - сообщения для отправки на сервер

Cкорее всего получится запустить на Unix
Пример работы
## Сервер
```
$ ./server ::1 8080
Listening on ::1: 8080 ...
Accepted connection from: ::1
Recieved something
Sending back: SOMETHING
 
```
## Клиент
```
$ ./client ::1 8080 something
Sending something to localhost: ::1 ...
got back: SOMETHING
```
 Также скриншоты в папке `ipv6`