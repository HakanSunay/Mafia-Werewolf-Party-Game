package main

import (
	"Mafia-Werewolf-Party-Game/src"
	"bufio"
	"fmt"
	"log"
	"net"
	"regexp"
)

func main() {
	ln, err := net.Listen("tcp", "127.0.0.1:8080")
	if err != nil {
		log.Println(err.Error())
	}
	allRooms := make([]src.Room, 0)
	currentRoles := make(map[src.Role]uint)
	allClients := make(map[net.Conn]src.Player)
	newConnections := make(chan net.Conn)
	deadConnections := make(chan net.Conn)
	messages := make(chan string)

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
				allClients[conn] = src.Player{false, nil, src.ClientCount,
					string(nameByte[:readBytes-1]), src.RandomJob(&currentRoles),
					false, 0, false}
				newConnections <- conn
				messages <- fmt.Sprintln(allClients[conn].Name, " joined the room!")
			}()
		}
	}()

	for {
		select {
		case curCon := <-newConnections:
			src.ClientCount += 1
			go func(conn net.Conn) {
				rd := bufio.NewReader(conn)
				for {
					m, err := rd.ReadString('\n')
					if err != nil {
						break
					}
					if allClients[curCon].Room == nil {
						jre := regexp.MustCompile(`#JOIN_ROOM (\w+)`)
						cre := regexp.MustCompile(`#CREATE_ROOM (?P<RoomName>\w+)`)
						if cre.MatchString(m) == true {
							fmt.Println("here1")
							res := cre.FindStringSubmatch(m)
							curPlayer := allClients[curCon]
							newRoom := curPlayer.CreateRoom(res[1])
							allRooms = append(allRooms, *newRoom)
							fmt.Println("here2")
						} else if jre.MatchString(m) == true {
							res := jre.FindStringSubmatch(m)
							existingRoom := src.FindRoom(allRooms, res[1])
							if existingRoom != nil {
								curPlayer := allClients[curCon]
								existingRoom.AddPlayer(&curPlayer)
							} else {
								curCon.Write([]byte("Room " + res[1] + " doesn't exist!"))
							}
						}
					}
					messages <- fmt.Sprintln("\n", allClients[curCon].Name, " : ", m)
					fmt.Println(allRooms)
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
