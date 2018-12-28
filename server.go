package main

import (
	"Mafia-Werewolf-Party-Game/src"
	"bufio"
	"fmt"
	"log"
	"net"
)

func main() {
	ln, err := net.Listen("tcp", "127.0.0.1:8080")
	if err != nil {
		log.Println(err.Error())
	}
	currentRoles := make(map[src.Role]uint)
	allClients := make(map[net.Conn]src.Player)
	newConnections := make(chan net.Conn)
	deadConnections := make(chan net.Conn)
	messages := make(chan string)
	// clientCount := 0

	go func() {
		for {
			conn, err := ln.Accept()
			if err != nil {
				log.Println(err.Error())
			}
			// src.clientCount += 1
			conn.Write([]byte("Tell me your name, babe!"))
			nameByte := make([]byte, 1024)
			go func() {
				readBytes, err := conn.Read(nameByte)
				if err != nil {
					log.Println(err.Error())
				}
				allClients[conn] = src.Player{false,nil,src.ClientCount,
				string(nameByte[:readBytes-1]), src.RandomJob(&currentRoles),
				false, 0}
				newConnections <- conn
				messages <- fmt.Sprintln(allClients[conn].Name, " joined the room!")
				fmt.Println(allClients[conn].Name," connected ",allClients[conn].Job)
			}()
			//readBytes, err := conn.Read(nameByte)
			//if err != nil {
			//	log.Println(err.Error())
			//}
			//allClients[conn] = Player{clientCount,string(nameByte[:readBytes-1])}
			//newConnections <- conn
		}
	}()

	for {
		select {
		case curCon := <-newConnections:
			src.ClientCount += 1
			//curCon.Write([]byte("Tell me your name, babe!"))
			//nameByte := make([]byte, 1024)
			//readBytes, err := curCon.Read(nameByte)
			//if err != nil {
			//	log.Println(err.Error())
			//}
			//allClients[curCon] = Player{clientCount,string(nameByte[:readBytes-1])}
			go func(conn net.Conn) {
				rd := bufio.NewReader(conn)
				for {
					m, err := rd.ReadString('\n')
					if err != nil {
						break
					}
					messages <- fmt.Sprintln("\n", allClients[curCon].Name, " : ", m)
				}
				deadConnections <- conn
			}(curCon)
		case msg := <-messages:
			for client := range allClients {
				client.Write([]byte(msg))
			}
		case lostClient := <-deadConnections:
			gonePlayer := allClients[lostClient]
			go func(playerName string) {
				messages <- fmt.Sprintln("\n", playerName, " left")
			}(gonePlayer.Name)
			delete(allClients, lostClient)
		}
	}
}
