package main

import (
"bufio"
"fmt"
"log"
"math/rand"
"net"
"time"
)

const ROLECOUNT = 4

type Role int

const (
	CITIZEN Role = iota
	MAFIA
	DOCTOR
	SHERIFF
)
var(
	citizenCount uint = 0
	mafiaCount uint = 0
	doctorCount uint = 0
	sheriffCount uint = 0
	clientCount uint = 0
)

type Player struct {
	number       uint
	name         string
	job          Role
	invulnerable bool
	votes        uint
}

func (pl *Player) save() {
	pl.invulnerable = true
}

func (pl *Player) kill() {

}

func randomJob(curRoles *map[Role]uint) Role {
	rand.Seed(time.Now().UnixNano())
	if clientCount < 4 {
		choice := rand.Intn(ROLECOUNT)
		return Role(choice)
	} else if doctorCount == 0 {
		doctorCount+=1
		return DOCTOR
	} else if sheriffCount == 0 {
		sheriffCount+=1
		return SHERIFF
	} else if citizenCount < mafiaCount {
		citizenCount+=1
		return CITIZEN
	} else {
		mafiaCount+=1
		return MAFIA
	}
}

func main() {
	ln, err := net.Listen("tcp", "127.0.0.1:8080")
	if err != nil {
		log.Println(err.Error())
	}
	currentRoles := make(map[Role]uint)
	allClients := make(map[net.Conn]Player)
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
			clientCount += 1
			conn.Write([]byte("Tell me your name, babe!"))
			nameByte := make([]byte, 1024)
			go func() {
				readBytes, err := conn.Read(nameByte)
				if err != nil {
					log.Println(err.Error())
				}
				allClients[conn] = Player{clientCount, string(nameByte[:readBytes-1]),
					randomJob(&currentRoles), false, 0}
				newConnections <- conn
				messages <- fmt.Sprintln(allClients[conn].name, " joined the room!")
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
			clientCount += 1
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
					messages <- fmt.Sprintln("\n", allClients[curCon].name, " : ", m)
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
			}(gonePlayer.name)
			delete(allClients, lostClient)
		}
	}
}
