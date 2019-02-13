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
	// currentRoles := make(map[src.Role]uint)
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
			// TODO: Make sure names are unique
			nameByte := make([]byte, 1024)
			go func() {
				readBytes, err := conn.Read(nameByte)
				if err != nil {
					log.Println(err.Error())
				}
				allClients[conn] = &src.Player{false, nil, src.ClientCount,
					string(nameByte[:readBytes-1]), src.CITIZEN,
					0, false, false, false}
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
				// TODO: IsEligible should block chat
				rd := bufio.NewReader(conn)
				for {
					curPlayer := allClients[conn]
					curRoom := allClients[conn].Room
					if !(curPlayer.IsEligibleToChat()) {
						continue
					}
					mesg, err := rd.ReadString('\n')
					if err != nil {
						break
					}
					if curRoom == nil {
						joinRoomReg := regexp.MustCompile(`#JOIN_ROOM (\w+)`)
						createRoomReg := regexp.MustCompile(`#CREATE_ROOM (?P<RoomName>\w+)`)
						allRoomsReg := regexp.MustCompile(`#ROOMS`)
						if createRoomReg.MatchString(mesg) == true {
							res := createRoomReg.FindStringSubmatch(mesg)
							newRoom := &src.Room{}
							newRoom.SetName(res[1])
							newRoom.AddPlayer(curPlayer)
							allRooms = append(allRooms, newRoom)
							conn.Write([]byte("You have successfully created room " + res[1]))
							continue
						} else if joinRoomReg.MatchString(mesg) == true {
							res := joinRoomReg.FindStringSubmatch(mesg)
							exists := false
							for _, r := range allRooms {
								if r.GetName() == res[1] {
									exists = true
									r.AddPlayer(curPlayer)
									conn.Write([]byte("You have successfully joined room " + r.GetName()))
									messages <- Mesg{fmt.Sprintln("\n", curPlayer.Name,
										" joined ", r.GetName()),
										r}
								}
							}
							if exists == false {
								conn.Write([]byte("Room " + res[1] + " doesn't exist!"))
							}
							continue

						} else if allRoomsReg.MatchString(mesg) == true {
							var buffer bytes.Buffer
							for _, r := range allRooms {
								buffer.WriteString(r.GetName() + " owned by " + r.GetOwner().Name + "\n")
							}
							if buffer.Len() == 0 {
								conn.Write([]byte("There are no rooms available!\n"))
							} else {
								conn.Write([]byte("The currently available rooms are: \n" + buffer.String()))
							}
							continue
						}
					} else if curRoom != nil &&
						!(curRoom.IsPlaying()) &&
						curPlayer.RoomOwner {
						startGameReg := regexp.MustCompile(`#START_GAME`)
						if startGameReg.MatchString(mesg) == true {
							outcome := curPlayer.StartGame()
							if outcome {
								conn.Write([]byte("Game started!"))
								messages <- Mesg{"Mafiozos, the game has begun! Show no mercy!",
									allClients[conn].Room}
							} else {
								conn.Write([]byte("Room can't be started!"))
							}
							continue
						}
					} else if curRoom.IsPlaying() {
						// swapped this
						// TODO: this part must be at the beginning!
						voteReg := regexp.MustCompile(`#VOTE (\w+)`)
						getPlayersReg := regexp.MustCompile(`#PLAYERS`)
						if voteReg.MatchString(mesg) == true {
							matchRes := voteReg.FindStringSubmatch(mesg)
							votedPlayerName := matchRes[1]
							if curPlayer.Job != src.DOCTOR {
								fmt.Println(curPlayer.Name, " wants to kick ", matchRes[1])
							} else {
								fmt.Println(curPlayer.Name, " wants to save ", matchRes[1])
							}
							// TODO : check if the name is correct
							curPlayer.CastVote(votedPlayerName)

							// Check if Next Stage is possible
							curStage := curRoom.GetStage()
							if curRoom.CanGoToNextStage() {
								hotSeatPlayer := curRoom.GetMostVotedPlayer()
								fmt.Println("HotSeatPlayer is ", hotSeatPlayer.Name)
								switch curStage {
									case src.MAFIASTAGE:
										hotSeatPlayer.AssignChosen()
									case src.DOCTORSTAGE:
										if curRoom.HasDoctor(){
											hotSeatPlayer.Save()
										}
									case src.ALLSTAGE:
										hotSeatPlayer.Die()
										messages <- Mesg{"Go town has decided to imprison " +
											hotSeatPlayer.Name + "\n",curRoom}
								}
								curRoom.NextStage()
								curRoom.Reset()
								if curRoom.GetStage() == src.MAFIASTAGE {
									messages <- Mesg{"MAFIA, IT IS TIME TO MURDER SOMEBODY!\n",
										curRoom}
									continue
								} else if curRoom.GetStage() == src.DOCTORSTAGE {
									if curRoom.HasDoctor(){
										messages <- Mesg{"DOC, SAVE A POOR OR A CORRUPT SOUL!\n",
											curRoom}
									} else {
										// added this else clause !!!!
										curRoom.NextStage()
										curRoom.Reset()
										deadWithoutDoctor := curRoom.FindChosenPlayerToDie()
										var announcement2 string
										if deadWithoutDoctor != nil {
											announcement2 = "GoTown! Our dear friend " + deadWithoutDoctor.Name +
												" has fallen into the hands of the mighty mafia last night!\n"
										}
										messages <- Mesg{announcement2,
											curRoom}
										continue
									}
								} else {
									deadPlayer := curRoom.FindChosenPlayerToDie()
									var announcement string
									if deadPlayer != nil {
										announcement = "GoTown! Our dear friend " + deadPlayer.Name +
											" has fallen into the hands of the mighty mafia last night!\n"
									} else {
										announcement = "GoTown! Our doctor did a great job last night!\n"
									}
									messages <- Mesg{announcement,
										curRoom}
									continue
								}
							}
							// no continue here, so that other players can see our VOTES!
						} else if getPlayersReg.MatchString(mesg) == true {
							var playersBuffer bytes.Buffer
							for _, pl := range curRoom.GetPlayers() {
								if pl.Dead == false{
									playersBuffer.WriteString(pl.Name + "\n")
								}
							}
							if playersBuffer.Len() == 0 {
								conn.Write([]byte("There are no players available!\n"))
							} else {
								conn.Write([]byte("The current players are: \n" + playersBuffer.String()))
							}
							continue
						}
					}
					messages <- Mesg{fmt.Sprintln("\n", curPlayer.Name, " : ", mesg),
						curRoom}
				}
				deadConnections <- conn
			}(curCon)
		case msg := <-messages:
			for conn, client := range allClients {
				if msg.room != nil && msg.room == client.Room && !(msg.room.IsPlaying()) {
					conn.Write([]byte(msg.content))
				} else if msg.room == nil && client.Room == nil {
					conn.Write([]byte(msg.content))
				} else if msg.room != nil && msg.room == client.Room && msg.room.IsPlaying() {
					currentStage := msg.room.GetStage()
					if client.Dead {
						conn.Write([]byte(msg.content))
					} else if currentStage == src.MAFIASTAGE {
						if client.Job == src.MAFIA {
							conn.Write([]byte(msg.content))
						}
					} else if currentStage == src.DOCTORSTAGE {
						if client.Job == src.DOCTOR {
							conn.Write([]byte(msg.content))
						}
					} else if currentStage == src.ALLSTAGE {
						conn.Write([]byte(msg.content))
					}
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
