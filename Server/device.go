package main

import (
	"bufio"
	"fmt"
	"math/rand"
	"net"
	"strconv"
)

type iot struct {
	id     int
	result int
	name   string
}

func (c iot) getResult() string {
	return c.name + strconv.Itoa(c.result)
}

func main() {
	//Структура устройства
	var listIoT [1000]iot

	//создаём устройства
	for i := 0; i < 1000; i++ {
		listIoT[i] = iot{
			id:     i,
			result: rand.Intn(200),
			name:   "{ IoT device | " + strconv.Itoa(i) + " }",
		}
	}

	fmt.Println("Запускаю сервер...")
	// Устанавливаем прослушивание порта
	ln, _ := net.Listen("tcp", ":8081")
	// Открываем порт
	conn, _ := ln.Accept()
	//Запускаем цикл
	for {
		// Будем прослушивать все сообщения разделенные \n
		message, err := bufio.NewReader(conn).ReadString('\n')
		// Распечатываем полученое сообщение

		if err != nil {
			fmt.Println(" Error ")
		}

		var name string
		go func() {
			for _, elem := range listIoT {
				//fmt.Print(message[0 : len(message)-1])
				if strconv.Itoa(elem.id) == message[0:len(message)-1] {
					// request = strconv.Itoa(elem.result)
					name = elem.name
					break
				}
			}
			// Процесс выборки для полученной строки
			newmessage := "request " + name + message
			// Отправить новую строку обратно клиенту
			conn.Write([]byte(newmessage + "\n"))
		}()
	}
}
