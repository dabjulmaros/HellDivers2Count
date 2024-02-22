package main

import (
	"encoding/json"
	"html/template"
	"io"
	"log"
	"net/http"
	"os"
	"time"

	"golang.org/x/text/language"
	"golang.org/x/text/message"
)

type Response struct {
	Player_count int `json:"player_count"`
	Result       int `json:"result"`
}

type SteamResponse struct {
	Response Response
}

func main() {

	oldTime := time.Now()
	oldCount := 0

	getPlayers := func() int {
		count := oldCount
		currentTime := time.Now()
		if oldCount == 0 || currentTime.Sub(oldTime).Minutes() > 5 {
			// res, err := http.Get("https://steamcommunity.com/app/553850")
			// if err != nil {
			// 	log.Fatal(err)
			// }
			// content, err := io.ReadAll(res.Body)
			// res.Body.Close()
			// if err != nil {
			// 	log.Fatal(err)
			// }

			// rgx, _ := regexp.Compile("<span class=\"apphub_NumInApp\">(?P<count>.+) In-Game<\\/span>")

			// m := rgx.FindStringSubmatch(string(content))
			// c := rgx.SubexpIndex("count")

			// count,err = strconv.Atoi(m[c]);
			// if err != nil{
			// 	count = 0
			// }

			// Seems to be somewhat behind the steam community website
			res, err := http.Get("https://api.steampowered.com/ISteamUserStats/GetNumberOfCurrentPlayers/v1/?appid=553850")
			if err != nil {
				log.Fatal(err)
			}
			content, err := io.ReadAll(res.Body)
			res.Body.Close()
			if err != nil {
				log.Fatal(err)
			}

			steamResponse := SteamResponse{}

			jsonErr := json.Unmarshal(content, &steamResponse)

			if jsonErr != nil {
				log.Fatal(jsonErr)
			}

			count = steamResponse.Response.Player_count
			oldCount = count
			oldTime = time.Now()
		}

		return count
	}

	h1 := func(w http.ResponseWriter, r *http.Request) {

		tmpl := template.Must(template.ParseFiles("index.html"))
		tmpl.Execute(w, nil)

	}

	h2 := func(w http.ResponseWriter, r *http.Request) {
		// var count string = "init"

		count := getPlayers()

		p := message.NewPrinter(language.English)

		
		rStr := p.Sprintf("<kbd>%d In-Game</kbd><small><sub>Last Update: %s UTC</sub></small>", count, oldTime.UTC().Format("03:04:05PM"));

		if count > 450000{
			rStr += "<img src='public/shock.png'/></div>"; 
		}

		tmpl, _ := template.New("t").Parse(rStr)

		tmpl.Execute(w, nil)
	}



	http.HandleFunc("/", h1)
	http.HandleFunc("/divers/", h2)
	http.Handle("/public/", http.StripPrefix("/public/", http.FileServer(http.Dir("public"))))

	
	port := os.Getenv("PORT")
	if port == ""{
		port = "8000"
	}

	log.Println("Running 0.0.0.0:"+port)
	log.Fatal(http.ListenAndServe("0.0.0.0:"+port, nil))
}
