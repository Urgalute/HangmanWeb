package main

import (
	"bufio"
	"fmt"
	"html/template"
	"log"
	"math/rand"
	"net/http"
	"os"
)

var Cheat int         //Compteur de changement de page tant que la game n'est pas fini
var Gardefou int      // Evite de recommencer si l'on a pas fini
var HangmanImg string //
var Texte string      // Txt pour indiquer au joueur si son guess était bon ou non

func main() {

	temp, err := template.ParseGlob("./Templates/*.html")
	if err != nil {
		fmt.Println(fmt.Sprint("ERREUR => %s", err.Error()))
		return
	}

	http.HandleFunc("/Hangman", func(w http.ResponseWriter, r *http.Request) {
		if Gardefou != 0 {
			Cheat++
			http.Redirect(w, r, "/Cheater", http.StatusSeeOther)
		} else {
			temp.ExecuteTemplate(w, "Welcome", nil)
		}
	})

	http.HandleFunc("/ChoixMot", func(w http.ResponseWriter, r *http.Request) {
		if Gardefou != 0 {
			Cheat++
			http.Redirect(w, r, "/Cheater", http.StatusSeeOther)
		} else {
			ShowTextFromFile(r.FormValue("difficulté"))
			Gardefou = 1
			http.Redirect(w, r, "/Display", http.StatusSeeOther)
		}
	})
	type Data struct {
		TabUnder   []string
		HangmanImg string
		Texte      string
	}
	http.HandleFunc("/Display", func(w http.ResponseWriter, r *http.Request) {
		HangmanImg = "/static/GraphHangman/Hangman" + fmt.Sprint(Graph /*+48*/) + ".png"
		datapage := Data{
			TabUnder, HangmanImg, Texte,
		}
		if Gardefou == 0 {
			http.Redirect(w, r, "/Hangman", http.StatusSeeOther)
		} else {
			temp.ExecuteTemplate(w, "Display", datapage)
		}
		fmt.Println("Word: ", Word)
		fmt.Println("TabUnder: ", TabUnder)
		fmt.Println("Guessed: ", Guessed)
		fmt.Println("Graph: ", Graph)
		fmt.Println("GardeFou: ", Gardefou)
		fmt.Println("Cheat: ", Cheat)
		fmt.Println("Lien Img: ", HangmanImg)
	})

	http.HandleFunc("/Traitement", func(w http.ResponseWriter, r *http.Request) {
		guess := ToLower(r.FormValue("inputtxt"))
		if Gardefou == 0 {
			http.Redirect(w, r, "/Hangman", http.StatusSeeOther)
		} else {
			if IsInGuessed(guess, Guessed) {
				Texte = "Vous avez déjà éssayer ça"
				http.Redirect(w, r, "/Display", http.StatusSeeOther)
			} else {
				if len(guess) > 1 {
					IsTheWord(guess)
				} else {
					IsInWord(guess)
				}
			}
			if win {
				http.Redirect(w, r, "/Winner", http.StatusSeeOther)
			} else if Graph > 9 {
				http.Redirect(w, r, "/Loose", http.StatusSeeOther)
			} else {
				http.Redirect(w, r, "/Display", http.StatusSeeOther)
			}
		}
	})

	http.HandleFunc("/Winner", func(w http.ResponseWriter, r *http.Request) {
		if win {
			win = false
			TabUnder = nil
			Guessed = nil
			Graph = 0
			Gardefou = 0
			Cheat = 0
			temp.ExecuteTemplate(w, "Winner", nil)
		} else if Gardefou == 0 {
			http.Redirect(w, r, "/Hangman", http.StatusSeeOther)
		} else {
			Cheat++
			http.Redirect(w, r, "/Cheater", http.StatusSeeOther)
		}
	})

	http.HandleFunc("/Loose", func(w http.ResponseWriter, r *http.Request) {
		if Graph > 9 {
			win = false
			TabUnder = nil
			Guessed = nil
			Graph = 0
			Gardefou = 0
			Cheat = 0
			temp.ExecuteTemplate(w, "Loose", nil)
		} else if Gardefou == 0 {
			http.Redirect(w, r, "/Hangman", http.StatusSeeOther)
		} else {
			Cheat++
			http.Redirect(w, r, "/Cheater", http.StatusSeeOther)
		}
	})

	http.HandleFunc("/Cheater", func(w http.ResponseWriter, r *http.Request) {
		datapage := Cheat
		if Cheat > 3 {
			temp.ExecuteTemplate(w, "Cheater", datapage)
			Graph = Graph + 2
		} else {
			http.Redirect(w, r, "Display", http.StatusSeeOther)
		}
	})

	fileServer := http.FileServer(http.Dir("Asset"))
	http.Handle("/static/", http.StripPrefix("/static/", fileServer))

	http.ListenAndServe("localhost:8080", nil)
}

//
//
//
// 			FONCTION
//
//
//

var Word string       //mot a deviner
var TabUnder []string //tableau d'underscore
var win bool          //verif si on a win (si y'a plus d'underscore)
var Guessed []string  //liste des lettres déjà rentrer
var Graph int         //compteur de point pour le graph

func ShowTextFromFile(path string) {
	file, err := os.Open(path)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()
	sc := bufio.NewScanner(file)
	lines := make([]string, 0)
	for sc.Scan() {
		lines = append(lines, sc.Text())
	}
	if err := sc.Err(); err != nil {
		log.Fatal(err)
	}
	randomIndex := rand.Intn(len(lines))
	Word = ToLower(lines[randomIndex])
	Underscore(Word)
	RandomLetter()
}

func ToLower(s string) string {
	var n string
	for _, c := range s {
		if c > 64 && c < 91 {
			n += string(c + 32)
		} else {
			n += string(c)
		}
	}
	return n
}

func Underscore(Word string) {
	for _, i := range Word {
		if i == '-' {
			TabUnder = append(TabUnder, ("-"))
		} else {
			TabUnder = append(TabUnder, ("_"))
		}
	}
}

func RandomLetter() {
	if len(Word) > 5 {
		id := rand.Intn(len(Word))
		ChangeTableau(string(Word[id]))
		Guessed = append(Guessed, string(Word[id]))
	}
	if len(Word) > 7 {
		ind := rand.Intn(len(Word))
		ChangeTableau(string(Word[ind]))
		Guessed = append(Guessed, string(Word[ind]))
	}
}

func ChangeTableau(c string) {
	for id, i := range Word {
		if string(i) == c {
			TabUnder[id] = c
		}
	}
}

func IsTheWord(w string) {
	if w == Word {
		win = true
		Gardefou = 0
		return
	} else {
		Texte = "Raté, c'est pas le bon mot"
		Guessed = append(Guessed, w)
	}
	win = false
	Graph += 2
}

func IsInGuessed(g string, Tab []string) bool {
	for _, i := range Tab {
		if i == g {
			Texte = "Déjà éssayé"
			return true
		}
	}
	return false
}

func IsComplete() {
	for _, i := range TabUnder {
		if i == "_" {
			win = false
			return
		}
	}
	win = true
	Gardefou = 0
}

func IsInWord(l string) {
	Guessed = append(Guessed, l)
	for _, i := range Word {
		if string(i) == l {
			ChangeTableau(l)
			IsComplete()
			Texte = "Bien joué !!"
			return
		}
	}
	Texte = "la lettre n'est pas dans le mot..."
	Graph++
}

/*

func Menu() {
	if win {
		DisplayHangman()
		fmt.Printf("Bien jouer, vous avez deviner le mot : %s", Word)
		return
	} else if Graph > 9 {
		fmt.Printf("Vous avez perdu, le mot à deviner était : %s", Word)
		return
	} else {
		Display()
		WordOrLetter()
		Menu()
	}
}

func Display() {
	fmt.Print("Voici le mot à deviner :")
	for _, i := range TabUnder {
		fmt.Print(i, " ")
	}
	fmt.Println("")
}

func WordOrLetter() {
	var Input string
	fmt.Println("1 : Entrez le mot entier")
	fmt.Println("2 : Entrez une seule lettre")
	fmt.Scan(&Input)
	fmt.Scan()
	switch Input {
	case "1":
		fmt.Println("Entrez un mot :")
		fmt.Scan(&Input)
		if IsInGuessed(Input, Guessed) {
			fmt.Println("Vous avez déjà essayé ce mot")
			return
		}
		IsTheWord(Input)
	case "2":
		fmt.Println("Entrez une lettre :")
		fmt.Scan(&Input)
		if IsInGuessed(Input, Guessed) {
			fmt.Println("Vous avez déjà essayé cette lettre")
			return
		}
		IsInWord(Input)
		return
	default:
		fmt.Println("Choix invalide")
		WordOrLetter()
	}
}

func DisplayHangman() {
	if Graph == 0 {
		Graph = 1
	}
	file, err := os.Open("Asset/GraphHangman/hangman.txt")
	if err != nil {
		log.Fatalf("Error: %s", err)
	}
	defer file.Close()
	var lines []string
	for i := 0; i < 7; i++ {
		scanner := bufio.NewScanner(file)
		for scanner.Scan() {
			lines = append(lines, scanner.Text())
		}
	}
	set := Graph * 10
	for i := 10 * (Graph - 1); i < set; i++ {
		fmt.Println(lines[i])
	}
}



*/
