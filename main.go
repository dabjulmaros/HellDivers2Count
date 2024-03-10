package main

import (
	"encoding/json"
	"html/template"
	"io"
	"log"
	"math"
	"net/http"
	"os"
	"time"

	"golang.org/x/text/language"
	"golang.org/x/text/message"
)

type SteamResponse struct {
	Response struct {
		Player_count int `json:"player_count"`
		Result       int `json:"result"`
	}
}

type DataCount struct {
	Count   int
	Updated time.Time
}

type DataStore struct {
	Data []DataCount
	Peak int
}

type TemplateResponse struct {
	Count   string
	Updated string
	Peak48  string
	Peak    string
	Shocked bool
}

func main() {
	var counts []DataCount
	peakCount := 0
	oldCount := 0
	oldTime := time.Now()
	oldTime = oldTime.AddDate(0, -1, 0)

	saveData := func() {
		file, _ := json.Marshal(DataStore{Data: counts, Peak: peakCount})

		err := os.WriteFile("StoredData.json", file, os.ModePerm)

		if err != nil {
			log.Fatal(err)
		}
	}

	loadData := func() {
		file, err := os.Open("StoredData.json")

		if err != nil {
			log.Println(err)
			return
		}

		defer file.Close()

		fileByte, _ := io.ReadAll(file)

		var StoredData DataStore
		jsonErr := json.Unmarshal(fileByte, &StoredData)

		if jsonErr != nil {
			log.Fatal(jsonErr)
		}

		peakCount = StoredData.Peak
		counts = StoredData.Data

	}
	loadData()

	fetchData := func() int {
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

		return steamResponse.Response.Player_count
	}

	storePeak := func(count int) {
		if count > peakCount {
			peakCount = count
		}
	}

	var intervalData func()

	intervalData = func() {

		now := time.Now()

		if len(counts) > 0 {
			latestStoredCount := counts[len(counts)-1]

			if math.Floor(now.Sub(latestStoredCount.Updated).Hours()) == 0 {
				return
			}
		}

		count := fetchData()
		storePeak(count)
		counts = append(counts, DataCount{Count: count, Updated: now.UTC()})

		if len(counts) > 48 {
			counts = counts[:48]
		}

		saveData()

		oldCount = count
		oldTime = now

		p := message.NewPrinter(language.English)

		log.Println(p.Sprintf("Count: %d, Time: %s, Count Length: %d", count, now.Format("15:04:05 MST"), len(counts)))

		time.AfterFunc(1*time.Hour, intervalData)
	}

	intervalData()

	getPlayers := func() DataCount {

		now := time.Now()

		latestStoredCount := counts[len(counts)-1]
		rData := latestStoredCount

		if now.Sub(latestStoredCount.Updated).Minutes() > 5 {
			if now.Sub(oldTime).Minutes() > 5 {
				oldCount = fetchData()
				storePeak(oldCount)
				oldTime = time.Now()
			}
			rData = DataCount{Count: oldCount, Updated: oldTime.UTC()}
		}
		return rData
	}

	getIndex := func(w http.ResponseWriter, r *http.Request) {

		tmpl := template.Must(template.ParseFiles("index.html"))
		tmpl.Execute(w, nil)

	}

	getDivers := func(w http.ResponseWriter, r *http.Request) {
		p := message.NewPrinter(language.English)

		count := getPlayers()

		p48 := count.Count
		for _, v := range counts {
			if v.Count > p48 {
				p48 = v.Count
			}
		}

		cStr := p.Sprintf("%d", count.Count)
		// uStr := p.Sprintf("%s", count.Updated.Format("15:04:05 MST"))
		uStr := p.Sprintf("%s", count.Updated)
		pStr := p.Sprintf("%d", peakCount)
		p48Str := p.Sprintf("%d", p48)

		tmpl := template.Must(template.ParseFiles("template.html"))

		tmpl.Execute(w, TemplateResponse{Count: cStr, Updated: uStr, Peak: pStr, Peak48: p48Str, Shocked: count.Count > 700000})
	}

	getData := func(w http.ResponseWriter, r *http.Request) {
		b, err := json.Marshal(counts)
		if err != nil {
			log.Println(err)
			return
		}
		io.WriteString(w, string(b))
	}

	http.HandleFunc("/", getIndex)
	http.HandleFunc("/divers/", getDivers)
	http.HandleFunc("/data/", getData)
	http.Handle("/public/", http.StripPrefix("/public/", http.FileServer(http.Dir("public"))))

	port := os.Getenv("PORT")
	if port == "" {
		port = "8000"
	}

	log.Println("Running 0.0.0.0:" + port)
	log.Fatal(http.ListenAndServe("0.0.0.0:"+port, nil))
}
