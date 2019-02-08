package main

import (
	"Mafia-Werewolf-Party-Game/src"
	"bufio"
	"bytes"
	"fmt"
	"log"
	"net"
	"regexp"
)

type Mesg struct {
	content string
	room    *src.Room
}

func main() {
	ln, err := net.Listen("tcp", "127.0.0.1:8080")
	if err != nil {
		log.Println(err.Error())
	}
	allRooms := make([]*src.Room, 0)
	currentRoles := make(map[src.Role]uint)
	allClients := make(map[net.Conn]*src.Player)
	newConnections := make(chan net.Conn)
	deadConnections := make(chan net.Conn)
	messages := make(chan Mesg)

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
				allClients[conn] = &src.Player{false, nil, src.ClientCount,
					string(nameByte[:readBytes-1]), src.RandomJob(&currentRoles),
					false, 0, false}
				newConnections <- conn
				messages <- Mesg{fmt.Sprintln(allClients[conn].Name, " joined the room!"),
					nil}
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
					mesg, err := rd.ReadString('\n')
					if err != nil {
						break
					}
					if allClients[conn].Room == nil {
						joinRoomReg := regexp.MustCompile(`#JOIN_ROOM (\w+)`)
						createRoomReg := regexp.MustCompile(`#CREATE_ROOM (?P<RoomName>\w+)`)
						allRoomsReg := regexp.MustCompile(`#ROOMS`)
						if createRoomReg.MatchString(mesg) == true {
							res := createRoomReg.FindStringSubmatch(mesg)
							newRoom := &src.Room{}
							newRoom.SetName(res[1])
							newRoom.AddPlayer(allClients[conn])
							allRooms = append(allRooms, newRoom)
						} else if joinRoomReg.MatchString(mesg) == true {
							res := joinRoomReg.FindStringSubmatch(mesg)
							exists := false
							for _, r := range allRooms {
								if r.GetName() == res[1] {
									exists = true
									r.AddPlayer(allClients[conn])
								}
							}
							if exists == false {
								conn.Write([]byte("Room " + res[1] + " doesn't exist!"))
							}

						} else if allRoomsReg.MatchString(mesg) == true {
							var buffer bytes.Buffer
							for _, r := range allRooms {
								buffer.WriteString(r.GetName() + " owned by " + r.GetOwner().Name + "\n")
							}
							conn.Write([]byte("The currently available rooms are: \n" + buffer.String()))
						}
					} else if allClients[conn].Room != nil && allClients[conn].RoomOwner {
						startGameReg := regexp.MustCompile(`#START_GAME`)
						if startGameReg.MatchString(mesg) == true {
							allClients[conn].StartGame()
							//TODO
						}
					}
					messages <- Mesg{fmt.Sprintln("\n", allClients[conn].Name, " : ", mesg),
						allClients[conn].Room}
				}
				deadConnections <- conn
			}(curCon)
		case msg := <-messages:
			for conn, client := range allClients {
				if msg.room == client.Room && msg.room != nil {
					conn.Write([]byte(msg.content))
				} else if msg.room == nil && client.Room == nil {
					conn.Write([]byte(msg.content))
				} else if msg.room.IsPlaying() {
					//TODO
				}
			}
		case lostClient := <-deadConnections:
			gonePlayer := allClients[lostClient]
			go func(playerName string) {
				messages <- Mesg{fmt.Sprintln("\n", playerName, " left"), gonePlayer.Room}
			}(gonePlayer.Name)
			delete(allClients, lostClient)
		}
	}
}
