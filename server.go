package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"strings"
)

type client chan<- string

var (
	entering = make(chan client)
	leaving  = make(chan client)
	messages = make(chan string)
)

func broadcaster() {
	clients := make(map[client]bool)

	for {
		select {
		case msg := <-messages:
			// Envia mensagem para todos os clientes
			for cli := range clients {
				cli <- msg
			}

		case cli := <-entering:
			clients[cli] = true

		case cli := <-leaving:
			delete(clients, cli)
			close(cli)
		}
	}
}

func clientWriter(conn net.Conn, ch <-chan string) {
	for msg := range ch {
		fmt.Fprintln(conn, msg)
	}
}

func handleConn(conn net.Conn) {
	ch := make(chan string)
	go clientWriter(conn, ch)

	apelido := conn.RemoteAddr().String()
	ch <- "Você é " + apelido
	messages <- apelido + " chegou!"
	entering <- ch

	input := bufio.NewScanner(conn)
	for input.Scan() {
		if strings.HasPrefix(input.Text(), "!nick ") {
			novoApelido := strings.TrimPrefix(input.Text(), "!nick ")
			messages <- apelido + " mudou o apelido para " + novoApelido
			apelido = novoApelido
		} else if input.Text() == "!sair" {
			break
		} else if strings.HasPrefix(input.Text(), "!msg ") {
			partes := strings.SplitN(strings.TrimPrefix(input.Text(), "!msg "), " ", 2)
			destinatario, mensagem := partes[0], partes[1]
			messages <- fmt.Sprintf("[Mensagem Privada de %s para %s]: %s", apelido, destinatario, mensagem)
		} else {
			messages <- apelido + ":" + input.Text()
		}
	}

	leaving <- ch
	messages <- apelido + " saiu"
	conn.Close()
}

func main() {
	fmt.Println("Iniciando servidor...")
	listener, err := net.Listen("tcp", "localhost:3000")
	if err != nil {
		log.Fatal(err)
	}
	go broadcaster()
	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Print(err)
			continue
		}
		go handleConn(conn)
	}
}
