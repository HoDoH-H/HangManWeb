package main

import (
	"encoding/json"
	"html/template"
	"log"
	"net/http"
	"os"

	"github.com/HoDoH-H/SimpleHangman"
)

type PlayerScore struct {
	Username   string
	Score      int
	Difficulty string
}

type GameStruct struct {
	User        *PlayerScore
	Data        *SimpleHangman.Data
	Leaderboard *[]PlayerScore
}

type RouteStruct struct {
	Next string
}

func UptLead(leaderboard *[]PlayerScore, user *PlayerScore) {
	result := []PlayerScore{}
	userAdded := false
	if len(*leaderboard) == 0 {
		*leaderboard = append(*leaderboard, *user)
	} else {
		for _, e := range *leaderboard {
			if user.Score > e.Score && user.Username != e.Username && !userAdded {
				result = append(result, *user)
				result = append(result, e)
				userAdded = true
			} else if user.Score >= e.Score && user.Username == e.Username && !userAdded {
				result = append(result, *user)
				userAdded = true
			} else if user.Username != e.Username {
				result = append(result, e)
			} else {
				result = append(result, e)
				userAdded = true
			}
		}
		if !userAdded {
			result = append(result, *user)
			userAdded = true
		}

		if len(result) >= 5 {
			*leaderboard = result[0:5]
		} else {
			*leaderboard = result
		}
	}
}

func SaveLeaderBoard(players *[]PlayerScore) {
	data, _ := json.Marshal(players)

	err := os.WriteFile("Save/save.json", data, 0644)
	if err != nil {
		os.Create("Save/save.json")
		os.WriteFile("Save/save.json", data, 0644)
	}
}

func GetLeaderBoard(players *[]PlayerScore) {
	data, _ := os.ReadFile("Save/save.json")
	json.Unmarshal(data, players)
}

func ActionHandler(w http.ResponseWriter, r *http.Request, data *SimpleHangman.Data) {
	input := r.FormValue("playerInput")
	formatedInput := SimpleHangman.FormatAns(input)

	if !SimpleHangman.IsLetterAlreadyTried(formatedInput, data) {
		if len(data.LetterTriedFormatizedText) > 0 {
			data.LetterTriedFormatizedText += formatedInput + "|"
		} else {
			data.LetterTriedFormatizedText += "|" + formatedInput + "|"
		}

		SimpleHangman.SplitWordToFindLetter(formatedInput, data)

		SimpleHangman.UpdateLife(formatedInput, data)

		SimpleHangman.VisualWord(data)

		SimpleHangman.IsWordDiscovered(data)
	}

	if data.Life == 0 {
		http.Redirect(w, r, "/Lose", http.StatusSeeOther)
	} else if data.WordFound {
		http.Redirect(w, r, "/Su2-vL5KG*Xc@_^llM$-fv3qoha+d01XcZG", http.StatusSeeOther)
	} else {
		http.Redirect(w, r, "/Game", http.StatusSeeOther)
	}
}

func AddWin(w http.ResponseWriter, r *http.Request, user *PlayerScore) {
	user.Score++
	http.Redirect(w, r, "/Win", http.StatusSeeOther)
}

func UpdateLeaderboardHandler(w http.ResponseWriter, r *http.Request, leaderboard *[]PlayerScore, user *PlayerScore, data *SimpleHangman.Data) {
	UptLead(leaderboard, user)
	SaveLeaderBoard(leaderboard)
	data.GameOver = true
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func NewGameHandler(w http.ResponseWriter, r *http.Request, gameData *GameStruct) {
	NameIsGood := false
	if gameData.User.Username == "Username" || gameData.User.Username != r.FormValue("usernameInput") {
		if r.FormValue("usernameInput") == "Username" {
			http.Redirect(w, r, "/IncorrectUsername", http.StatusSeeOther)
		} else {
			gameData.User.Username = r.FormValue("usernameInput")
			NameIsGood = true
		}
	} else {
		NameIsGood = true
	}

	if NameIsGood {
		if gameData.Data.GameOver {
			gameData.User.Score = 0
			gameData.Data.GameOver = false
		}

		gameData.Data.Life = 10
		gameData.Data.AncientLetter, gameData.Data.LetterFound = []string{}, []string{}
		gameData.Data.WordFound = false
		gameData.User.Difficulty = r.FormValue("Difficulty")
		SimpleHangman.GetWord(gameData.Data, r.FormValue("Difficulty"))
		SimpleHangman.VisualWord(gameData.Data)
		gameData.Data.LetterTriedFormatizedText = ""

		http.Redirect(w, r, "/Game", http.StatusSeeOther)
	}
}

func RestartGameHandler(w http.ResponseWriter, r *http.Request, gameData *GameStruct) {
	if gameData.Data.GameOver {
		gameData.User.Score = 0
		gameData.Data.GameOver = false
	}

	gameData.Data.Life = 10
	gameData.Data.AncientLetter, gameData.Data.LetterFound = []string{}, []string{}
	gameData.Data.WordFound = false
	SimpleHangman.GetWord(gameData.Data, gameData.User.Difficulty)
	SimpleHangman.VisualWord(gameData.Data)
	gameData.Data.LetterTriedFormatizedText = ""

	http.Redirect(w, r, "/Game", http.StatusSeeOther)
}

func Home(w http.ResponseWriter, r *http.Request, data *GameStruct) {
	template, err := template.ParseFiles("./Html/menu.html")
	if err != nil {
		log.Fatal(err)
	}
	template.Execute(w, data)
}

func HomeBis(w http.ResponseWriter, r *http.Request, data *GameStruct) {
	template, err := template.ParseFiles("./Html/menubis.html")
	if err != nil {
		log.Fatal(err)
	}
	template.Execute(w, data)
}

func LoseScreen(w http.ResponseWriter, r *http.Request, gameData *GameStruct, leaderboard *[]PlayerScore) {
	UptLead(leaderboard, gameData.User)
	SaveLeaderBoard(leaderboard)
	gameData.Data.GameOver = true
	template, err := template.ParseFiles("./Html/loseScreen.html")
	if err != nil {
		log.Fatal(err)
	}
	template.Execute(w, gameData)
}

func WinScreen(w http.ResponseWriter, r *http.Request, gameData *GameStruct) {
	template, err := template.ParseFiles("./Html/winScreen.html")
	if err != nil {
		log.Fatal(err)
	}
	template.Execute(w, gameData)
}

func Game(w http.ResponseWriter, r *http.Request, gameData *GameStruct) {
	template, err := template.ParseFiles("./Html/hangman.html")
	if err != nil {
		log.Fatal(err)
	}
	template.Execute(w, gameData)
}

func main() {
	leaderboard := []PlayerScore{}
	GetLeaderBoard(&leaderboard)
	data := SimpleHangman.Data{GameOver: true}
	user := PlayerScore{Username: "Username", Score: 0}

	gameData := GameStruct{User: &user, Data: &data, Leaderboard: &leaderboard}

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		Home(w, r, &gameData)
	})
	http.HandleFunc("/IncorrectUsername", func(w http.ResponseWriter, r *http.Request) {
		HomeBis(w, r, &gameData)
	})
	http.HandleFunc("/Game", func(w http.ResponseWriter, r *http.Request) {
		Game(w, r, &gameData)
	})
	http.HandleFunc("/Action", func(w http.ResponseWriter, r *http.Request) {
		ActionHandler(w, r, &data)
	})
	http.HandleFunc("/NewGame", func(w http.ResponseWriter, r *http.Request) {
		NewGameHandler(w, r, &gameData)
	})
	http.HandleFunc("/RestartGame", func(w http.ResponseWriter, r *http.Request) {
		RestartGameHandler(w, r, &gameData)
	})
	http.HandleFunc("/Lose", func(w http.ResponseWriter, r *http.Request) {
		LoseScreen(w, r, &gameData, &leaderboard)
	})
	http.HandleFunc("/Win", func(w http.ResponseWriter, r *http.Request) {
		WinScreen(w, r, &gameData)
	})
	http.HandleFunc("/NeedUpdate", func(w http.ResponseWriter, r *http.Request) {
		UpdateLeaderboardHandler(w, r, &leaderboard, &user, &data)
	})

	// Tried to make the routes as hard as possible to evoid fake wins
	http.HandleFunc("/Su2-vL5KG*Xc@_^llM$-fv3qoha+d01XcZG", func(w http.ResponseWriter, r *http.Request) {
		AddWin(w, r, &user)
	})
	fs := http.FileServer(http.Dir("./server/"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))
	http.ListenAndServe(":8080", nil)
}
