
// package main

// import (
// 	"bufio"
// 	"fmt"
// 	"log"
// 	"net"
// 	"strconv"
// 	"strings"
// )

// const (
// 	PORT = ":8080"
// )






// func handleConnection(conn net.Conn) {
// 	defer conn.Close()
// 	msg, err := bufio.NewReader(conn).ReadString('\n')
// 	if err != nil {
// 		fmt.Println("Could not read message from connection: ", err)
// 		return
// 	}
// 	res, ok := parseCommand(msg)
// 	if !ok {
// 		conn.Write([]byte("Could not parse message...\n"))
// 		return
// 	}
// 	conn.Write([]byte(res))
// }


// func main() {
// 	list, err := net.Listen("tcp", PORT)
// 	if err != nil {
// 		log.Fatalf("Could not start server: %s\n", err)
// 	}
// 	defer list.Close()
// 	fmt.Printf("Now accepting connections on port %v...\n", PORT)
// 	for {
// 		conn, err := list.Accept()
// 		if err != nil {
// 			fmt.Println("Could not accept incoming request: ", err)
// 			continue
// 		}
// 		go handleConnection(conn)
// 	}
// }