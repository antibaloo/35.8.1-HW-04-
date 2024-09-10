package main

import (
	"bufio"
	"fmt"
	"io"
	"math/rand"
	"net"
	"os"
	"strings"
	"time"
)

// Протокол службы
const proto = "tcp4"

// Порт службы, поставил по умолчанию порт telnet
const addr = ":23"

var proverbs []string

func main() {
	f, err := os.Open("go-proverbs.txt")
	if err != nil {
		fmt.Printf("При при чтении файла с поговорками произошла ошибка: %v", err)
		return
	}
	reader := bufio.NewReader(f)
	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			if err == io.EOF {
				break
			} else {
				fmt.Printf("При при чтении поговорок из файла произошла ошибка: %v", err)
				return
			}
		}

		// Убираем лишнее
		line = strings.Trim(line, "\r")
		line = strings.Trim(line, "\n")
		proverbs = append(proverbs, line)
	}

	// Запускаем сетевую службу
	service, err := net.Listen(proto, addr)
	if err != nil {
		fmt.Printf("При запуске сетевой службы произошла ошибка: %v", err)
		return
	}

	// Не забываем освобождать ресурсы
	defer service.Close()

	// Ожидаем подключения
	for {
		// Принимаем подключение.
		c, err := service.Accept()
		if err != nil {
			fmt.Printf("При подключении клиента к сетевой службе произошла ошибка: %v", err)
			return
		}
		// Вызов обработчика подключения.
		go randomProverb(c)
	}
}

func randomProverb(c net.Conn) {
	// Бесконечный цикл
	for {
		// Спим 3 секунды
		time.Sleep(3 * time.Second)
		// Получаем случайную поговорку
		proverb := proverbs[rand.Intn(len(proverbs))]
		// Отправляем клиенту
		c.Write([]byte(proverb + "\n"))
	}
}
