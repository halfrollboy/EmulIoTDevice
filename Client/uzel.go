package main

import (
	"bufio"
	"fmt"
	"math/rand"
	"net"
	"runtime"
	"strconv"
	"sync"
	"time"
)

// Подключаемся к сокету
var conn, _ = net.Dial("tcp", "127.0.0.1:8081")

//Создание структуры узла
type node struct {
	id   int
	name string
}

//Возвращает id того IoT устройства к которому хочет подключиться
func (n node) getIDDevice() int {
	return rand.Intn(1000)
}

//Константы для работы с каналами
const (
	goroutineCount = 20
	iterationCount = 10000
	quotaCount     = 10 // количество горутин, которые должны работать, пока остальные будут ждать их завершения
)

//Главная функция обработчик
func worker(in int, wg *sync.WaitGroup, quotaChan chan struct{}) {
	/*
		Занимаем слот в канале. Если места не будет, то горутина
		будет ждать и не начнет работу, пока не освободиться место
	*/
	quotaChan <- struct{}{}
	defer wg.Done()

	for j := 0; j < iterationCount; j++ {
		/*
			Функция обмена информацией
			Здесь узел опрашивает устройство в диапозоне 1000
		*/
		formatWork(in)

		if j%2 == 0 {
			<-quotaChan             // делимся ресурсами с другими горутинами
			quotaChan <- struct{}{} // но при этом лимит работащих горутин по прежнему тот же
		}
		runtime.Gosched() // передает управление другой горутине
	}
	<-quotaChan // освобождает слот
}

func main() {
	var massNode [10]node

	//Создадим узлы
	for i := 0; i < quotaCount; i++ {
		massNode[i] = node{id: i, name: "Node " + strconv.Itoa(rand.Intn(1000))}
	}

	wg := &sync.WaitGroup{}                      //Для ожидания горутин
	quotaChan := make(chan struct{}, quotaCount) // буфферезированный канал(асинхронный), с пустыми структурами(они не занимают места в памяти)

	for i := 0; i < goroutineCount; i++ {
		wg.Add(1)
		go worker(i, wg, quotaChan)
	}
	time.Sleep(time.Millisecond)
	wg.Wait()
}

func formatWork(in int) {
	fmt.Fprintf(conn, strconv.Itoa(rand.Intn(1000))+"\n")
	message, _ := bufio.NewReader(conn).ReadString('\n')
	fmt.Print("Message from server: " + message)
}

// func formatWork(in, j int) string {
// 	return fmt.Sprintln(strings.Repeat("  ", in), "█",
// 		strings.Repeat("  ", goroutineCount-in),
// 		"th", in,
// 		"iter", j, strings.Repeat("■", j))
// }
