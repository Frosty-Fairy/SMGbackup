package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"
	"time"

	scp "github.com/bramvdbogaerde/go-scp"
	"github.com/bramvdbogaerde/go-scp/auth"
	"golang.org/x/crypto/ssh"
)

func main() {
	var sliceIP []string
	var chanAnswer = make(chan string, len(sliceIP))

	ipList, err := ioutil.ReadFile("IP/ip_list.txt")
	if err != nil {
		log.Fatal(err)

	}

	a := strings.Split(string(ipList), "\n")
	sliceIP = append(sliceIP, a...)

	for i := 0; i < len(sliceIP); i++ {
		a := sliceIP[i]
		go backup(a, chanAnswer)

	}

	for i := 0; i < len(sliceIP); i++ {
		fmt.Println(<-chanAnswer)

	}

	fmt.Println("Успешное завершение работы")

}

func backup(a string, chanAnswer chan string) {

	clientConfig, _ := auth.PasswordKey("login", "password", ssh.InsecureIgnoreHostKey())

	port := ":22"

	client := scp.NewClient(a+port, &clientConfig)

	fmt.Println("Подключаюсь к хосту: ", a)

	err := client.Connect()
	if err != nil {

		log.Printf("Couldn't establish a connection to the remote server: %s", err)
		chanAnswer <- "Ошибка подключения к " + a
		return

	}

	defer client.Close()

	dt := time.Now()

	date_string := dt.Format("2006-Jan-02")

	file_path := "./backup/" + a + "(" + date_string + ")" + ".yaml"

	f, err := os.OpenFile(file_path, os.O_RDWR|os.O_CREATE, 0777)
	if err != nil {
		log.Fatalf("Couldn't open the output file")
	}

	defer f.Close()

	fmt.Println("Скачиваю конфиг с ", a)

	err = client.CopyFromRemote(f, "/etc/config/cfg.yaml")
	if err != nil {
		log.Fatalf("Copy failed from remote")
	}

	chanAnswer <- "Конфиг " + a + " сохранён"

}
