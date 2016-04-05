package main

import (
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
)

type Player struct {
	Name  string `json:"name"`
	Color string `json:"color"`
}

func (p *Player) randomize() {
	bigString := "akdfya87w36c4iungt2673tc5unyedc7g2j97sa6ged7f6cnrgnydgf7awgcj57g62cnybfubwhe897r6gc7c63k84r"
	pos := rand.Intn(65)
	p.Name = bigString[pos : pos+rand.Intn(10)+5]

	lilString := "f453ff78c8097e9807a098b7098c70ad7809ea"
	pos = rand.Intn(30)
	p.Color = "#" + lilString[pos:pos+6]
}

func playersJSON(w http.ResponseWriter, players []Player) {
	fmt.Fprintf(w, "{\"players\": [")
	separator := ""
	for i := 0; i < len(players); i++ {
		p := players[i]
		pJson, err := json.Marshal(p)
		if err != nil {
			log.Println(err)
			return
		}

		fmt.Fprintf(w, "%v%v", separator, string(pJson))
		separator = ",\n"
	}
	fmt.Fprintf(w, "]}")
}
