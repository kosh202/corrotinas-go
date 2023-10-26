package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"strings"
)

func mustCopy(dst io.Writer, src io.Reader) {
	if _, err := io.Copy(dst, src); err != nil {
		log.Fatal(err)
	}
}

func main() {
	conn, err := net.Dial("tcp", "localhost:3000")
	fmt.Println("Connected!")
	if err != nil {
		log.Fatal(err)
	}

	done := make(chan struct{})

	go func() {
		io.Copy(os.Stdout, conn)
		log.Println("done")
		done <- struct{}{}
	}()

	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Print("Enter message: ")
		text, err := reader.ReadString('\n')
		if err != nil {
			log.Fatal(err)
		}

		text = strings.TrimSpace(text)

		if strings.HasPrefix(text, "!nick ") {
			_, err := fmt.Fprintln(conn, text)
			if err != nil {
				log.Fatal(err)
			}
		} else if text == "!sair" {
			_, err := fmt.Fprintln(conn, text)
			if err != nil {
				log.Fatal(err)
			}
			conn.Close()
			break
		} else if strings.HasPrefix(text, "!msg ") {
			_, err := fmt.Fprintln(conn, text)
			if err != nil {
				log.Fatal(err)
			}
		} else {
			_, err := fmt.Fprintln(conn, text)
			if err != nil {
				log.Fatal(err)
			}
		}
	}

	<-done
}
