# <img src="https://purepng.com/public/uploads/large/purepng.com-hired-gun-graves-skinsplashartchampionleague-of-legendsskingraves-3315199260108xiqj.png" height="120"><img src="https://cacophony.org.nz/sites/default/files/gopher.png" height="120">
# Mafia (Werewolf) Party GoLang Game

## Description
The Mafia party game presents a conflict between the Mafia – the informed minority – and the Innocents – the uninformed majority. Originated by Dmitry Davidoff of the USSR in 1986, this popular game has many variations and can be played by a group of seven or more people.

The game has two phases; night, when the Mafia might secretly “murder” an innocent, and “day” when Innocents vote to eliminate a Mafiosi suspect. The game ends when either all the Mafia members are eliminated or the Innocents.

There are also other roles such as:
* Doctor - has the ability to save 1 person on a random basis.
* Sheriff - apprehends the Mafioso suspect.

## Install
1) Make sure you have **GoLang** installed on your local machine
https://golang.org/doc/install
2) After you successfully install **go**, run the following command:
```
go get github.com/HakanSunay/Mafia-Werewolf-Party-Game
```

## Usage
In order to start the game, you will first have to run:
```
# start the server
go run server.go <port>
```
In order to join the game as a player, you need to run:
```
# run the client
go run client.go <hostname> <port>
```
### Gameplay
1. You will be prompted to either CREATE or JOIN a ROOM
2. Whatever your decision was, the game can be started by the room creator when a minimum of 6 players gather.
3. Random hidden roles are assigned to all of the players.
4. Everyone falls asleep.
5. The _Mafia_ wakes up, the chat is opened for Mafia members only and they vote to eliminate someone and fall asleep.
6. The _Doctor_ wakes up, the chat is opened only for him and he is prompted to select a player to save and then he falls asleep as well.
7. Everybody wakes up, the chat is opened for everyone and they vote to choose the mafia to send to prison.
8. The sheriff arrests the chosen suspect.
9. Go back to 4, until either one of _Mafia_/_Innocent_ **Count** becomes 0.

## Possible Future Upgrades
* Run the server on AWS
* Connect using the Client using the Server's public IP and Port
* Implement 2D Game Graphics
## Contributing
If you want to contribute to this repository and fullfil my future plans, you can simply do the following:
* Clone the repository:
```
$ git clone https://github.com/HakanSunay/Mafia-Werewolf-Party-Game.git
```
* Create a branch:
```
$ git checkout -b your_branch_name
```
* Bless the branch with your extraordinary genius.
* Create a pull(/merge) request.

## Bug Reporting
In case you want to report any bugs and I strongly advise you to do so if you happen to run across any, please do by using the **Issues** section.

To see how to create an issue, please check the following link:
https://help.github.com/articles/creating-an-issue/
## Contact

You can contact me at:

hakansunayhalil@gmail.com

## License

MIT.

