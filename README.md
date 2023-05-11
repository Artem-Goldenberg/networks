# networks
Задания по курсу компьютерных сетей

2. Programming  

Poke - my ping
--- 
Отсылает эхо запросы через `ICMP` раз в секунду, замеряя процент потери пакетов и "round trip time" (`RTT`)  
Запустить:
```
gcc -o poke poke.c
./poke <host>
```
`<host>` - IPv4 адрес хоста или имя, например `youtube.com`.  

Пример вывода: (Еще есть скрин `PokeAtWork.png`)
```
$ ./poke wikipedia.org
got 44 bytes from 91.198.174.192:
seq: 0  RTT: 42.743 ms
RTT stats so far: avg: 42.743 ms  min: 42.743 ms  max: 42.743 ms
lost 0% of packets

got 44 bytes from 91.198.174.192:
seq: 1  RTT: 49.894 ms
RTT stats so far: avg: 46.319 ms  min: 42.743 ms  max: 49.894 ms
lost 0% of packets

got 44 bytes from 91.198.174.192:
seq: 2  RTT: 49.558 ms
RTT stats so far: avg: 47.398 ms  min: 42.743 ms  max: 49.894 ms
lost 0% of packets

got 44 bytes from 91.198.174.192:
seq: 3  RTT: 49.544 ms
RTT stats so far: avg: 47.935 ms  min: 42.743 ms  max: 49.894 ms
lost 0% of packets

got 44 bytes from 91.198.174.192:
seq: 4  RTT: 49.753 ms
RTT stats so far: avg: 48.298 ms  min: 42.743 ms  max: 49.894 ms
lost 0% of packets

^C
```

Есть обработка всех icmp типов, которые обозначают ошибку (3 - Destination Unreachable, 11, ...)  

Работает только на MacOS. Даже на Linux скорее всего работать не будет, названия макросов для `icmp` разные. 
