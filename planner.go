package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"
)

type wotd struct {
	word string
	def  string
}

type qotd struct {
	source string
	quote  string
}

type weather struct {
	location    string
	latitude    string
	longitude   string
	updated     string
	conditions  string
	temperature string
	humidity    string
	wind        string
	pressure    string
	visibility  string
}

type nyt struct {
	title1       string
	description1 string
	title2       string
	description2 string
	title3       string
	description3 string
}

//const qotdURL = "http://quotes.rest/qod.xml"
const qotdURL = "http://feeds.feedburner.com/quotationspage/qotd"

//const wotdURL = "http://www.macmillandictionary.com/us/wotd/wotdrss.xml"
const wotdURL = "https://wordsmith.org/awad/rss1.xml"
const weatherURL = "http://w1.weather.gov/xml/current_obs/KLAF.xml"
const forecastURL = "http://www.accuweather.com/en/us/west-lafayette-in/47906/weather-forecast/2135952"
const nytURL = "http://rss.nytimes.com/services/xml/rss/nyt/US.xml"

const wotdReloadInterval = 24
const qotdReloadInterval = 24
const weatherReloadInterval = 1
const forecastReloadInterval = 1
const nytReloadInterval = 1

//const HTMLFile = "/home/krigbaum/devel/go/src/github.com/krigbaum/planner/FamilyPlanner/index.html"
const HTMLFile = "c:/Users/lekrigbaum/Desktop/go/src/github.com/krigbaum/planner/FamilyPlanner/index.html"

var mutex = &sync.Mutex{}

func extractWeather(text string, str string, rep int) string {
	loc := 0
	start := 0
	for i := 1; i <= rep; i++ {
		loc = strings.Index(text[start:], str) + len(str)
		start = start + loc
	}

	end := strings.Index(text[start:], "</")
	end = start + end
	return text[start:end]
}

func getWeather() {
	url := weatherURL
	var w weather

	resp, err := http.Get(url)
	if err != nil {
		fmt.Fprintf(os.Stderr, "fetch: %v\n", err)
		os.Exit(1)
	}
	b, err := ioutil.ReadAll(resp.Body)
	resp.Body.Close()
	if err != nil {
		fmt.Fprintf(os.Stderr, "fetch: reading %s: %v\n", url, err)
		os.Exit(1)
	}

	fieldName := "<location>"
	w.location = extractWeather(string(b), fieldName, 1)

	fieldName = "<latitude>"
	w.latitude = extractWeather(string(b), fieldName, 1)
	lat, _ := strconv.ParseFloat(w.latitude, 64)
	latitude := strconv.FormatFloat(lat, 'f', 2, 64)

	fieldName = "<longitude>"
	w.longitude = extractWeather(string(b), fieldName, 1)
	long, _ := strconv.ParseFloat(w.longitude, 64)
	longitude := strconv.FormatFloat(long, 'f', 2, 64)

	location := fmt.Sprintf("%s (%s N / %s W)", w.location, latitude, longitude)

	fieldName = "<observation_time>"
	w.updated = extractWeather(string(b), fieldName, 1)

	fieldName = "<weather>"
	w.conditions = extractWeather(string(b), fieldName, 1)

	fieldName = "<temperature_string>"
	w.temperature = extractWeather(string(b), fieldName, 1)

	fieldName = "<relative_humidity>"
	w.humidity = extractWeather(string(b), fieldName, 1)
	humidity := fmt.Sprintf("%s&#37;", w.humidity)

	fieldName = "<wind_string>"
	w.wind = extractWeather(string(b), fieldName, 1)

	fieldName = "<pressure_in>"
	w.pressure = extractWeather(string(b), fieldName, 1)

	fieldName = "<visibility_mi>"
	w.visibility = extractWeather(string(b), fieldName, 1)
	visibility := fmt.Sprintf("%s mi", w.visibility)

	file := HTMLFile
	mutex.Lock()
	src, err := ioutil.ReadFile(file)
	if err != nil {
		fmt.Printf("Error reading file %s: %v\n", file, err)
	}
	memoryFile := string(src)

	memoryFile = replaceByID(memoryFile, "id=\"location\">", location)
	memoryFile = replaceByID(memoryFile, "id=\"updated\">", w.updated)
	memoryFile = replaceByID(memoryFile, "id=\"cond\">", w.conditions)
	memoryFile = replaceByID(memoryFile, "id=\"temp\">", w.temperature)
	memoryFile = replaceByID(memoryFile, "id=\"humid\">", humidity)
	memoryFile = replaceByID(memoryFile, "id=\"pressure\">", w.pressure)
	memoryFile = replaceByID(memoryFile, "id=\"visibility\">", visibility)
	memoryFile = replaceByID(memoryFile, "id=\"wind\">", w.wind)

	err = ioutil.WriteFile(HTMLFile, []byte(memoryFile), 0644)
	mutex.Unlock()
	if err != nil {
		fmt.Printf("Error writing to file %s: %v\n", HTMLFile, err)
	}
}

func Weather() {
	// Initial Weather load on startup
	log.Println("Initial Weather() Load")
	getWeather()

	//==================================
	// Repeat Weather load every weatherdReloadInterval
	ticker := time.NewTicker(time.Hour * weatherReloadInterval)
	for range ticker.C {
		log.Println("Periodic Weather() Load")
		getWeather()
	}
}

//=======================================================================

func extractForecast(text string, str string, rep int) string {
	loc := 0
	start := 0
	for i := 1; i <= rep; i++ {
		loc = strings.Index(text[start:], str) + len(str)
		start = start + loc
	}

	end := strings.Index(text[start:len(text)], "<")
	end = start + end
	return text[start:end]
}

func getForecast() {
	url := forecastURL
	//var w wotd

	resp, err := http.Get(url)
	if err != nil {
		fmt.Fprintf(os.Stderr, "fetch: %v\n", err)
		os.Exit(1)
	}
	b, err := ioutil.ReadAll(resp.Body)
	resp.Body.Close()
	if err != nil {
		fmt.Fprintf(os.Stderr, "fetch: reading %s: %v\n", url, err)
		os.Exit(1)
	}

	fieldName := "day=1\">"
	period := extractForecast(string(b), fieldName, 2)

	fieldName = "class=\"large-temp\">"
	temp := extractForecast(string(b), fieldName, 6)

	fieldName = "class=\"cond\">"
	cond := extractForecast(string(b), fieldName, 3)

	tonight := period + ": " + cond + ", low of " + temp

	fieldName = "day=2\">"
	period = extractForecast(string(b), fieldName, 2)

	fieldName = "class=\"large-temp\">"
	temp = extractForecast(string(b), fieldName, 7)

	fieldName = "class=\"cond\">"
	cond = extractForecast(string(b), fieldName, 3)

	tomorrow := period + ": " + cond + ", high of " + temp
	file := HTMLFile
	mutex.Lock()
	src, err := ioutil.ReadFile(file)
	if err != nil {
		fmt.Printf("Error reading file %s: %v\n", file, err)
	}
	memoryFile := string(src)

	memoryFile = replaceByID(memoryFile, "id=\"tonight\">", tonight)
	memoryFile = replaceByID(memoryFile, "id=\"tomorrow\">", tomorrow)

	err = ioutil.WriteFile(HTMLFile, []byte(memoryFile), 644)
	mutex.Unlock()
	if err != nil {
		fmt.Printf("Error writing to file %s: %v\n", HTMLFile, err)
	}
}

func Forecast() {
	// Initial WOTD load on startup
	log.Println("Initial Forecast() Load")
	getForecast()

	//==================================
	// Repeat WOTD load every forecastReloadInterval
	ticker := time.NewTicker(time.Hour * forecastReloadInterval)
	for range ticker.C {
		log.Println("Periodic Forecast() Load")
		getForecast()
	}
}

//=======================================================================
func extractWOTD(text string, str string, rep int) string {
	loc := 0
	start := 0
	for i := 1; i <= rep; i++ {
		loc = strings.Index(text[start:], str) + len(str)
		start = start + loc
	}

	end := strings.Index(text[start:len(text)], "</")
	end = start + end
	return text[start:end]
}

func getWOTD() {
	url := wotdURL
	var w wotd

	resp, err := http.Get(url)
	if err != nil {
		fmt.Fprintf(os.Stderr, "fetch: %v\n", err)
		os.Exit(1)
	}
	b, err := ioutil.ReadAll(resp.Body)
	resp.Body.Close()
	if err != nil {
		fmt.Fprintf(os.Stderr, "fetch: reading %s: %v\n", url, err)
		os.Exit(1)
	}

	fieldName := "<title>"
	w.word = extractWOTD(string(b), fieldName, 2)
	word := fmt.Sprintf("%s: ", w.word)

	fieldName = "<description>"
	w.def = extractWOTD(string(b), fieldName, 2)

	file := HTMLFile
	mutex.Lock()
	src, err := ioutil.ReadFile(file)
	if err != nil {
		fmt.Printf("Error reading file %s: %v\n", file, err)
	}
	memoryFile := string(src)

	memoryFile = replaceByID(memoryFile, "id=\"word\">", word)
	memoryFile = replaceByID(memoryFile, "id=\"definition\">", w.def)

	err = ioutil.WriteFile(HTMLFile, []byte(memoryFile), 644)
	mutex.Unlock()
	if err != nil {
		fmt.Printf("Error writing to file %s: %v\n", HTMLFile, err)
	}
}

func WOTD() {
	// Initial WOTD load on startup
	log.Println("Initial WOTD() Load")
	getWOTD()

	//==================================
	// Repeat WOTD load every wotdReloadInterval
	ticker := time.NewTicker(time.Hour * wotdReloadInterval)
	for range ticker.C {
		log.Println("Periodic WOTD() Load")
		getWOTD()
	}
}

func extractQOTD(text string, str string, rep int) string {
	loc := 0
	start := 0
	for i := 1; i <= rep; i++ {
		loc = strings.Index(text[start:], str) + len(str)
		start = start + loc
	}

	end := strings.Index(text[start:len(text)], "</")
	end = start + end
	return text[start:end]
}

func getQOTD() {
	url := qotdURL
	var q qotd

	resp, err := http.Get(url)
	if err != nil {
		fmt.Fprintf(os.Stderr, "fetch: %v\n", err)
		os.Exit(1)
	}
	b, err := ioutil.ReadAll(resp.Body)
	resp.Body.Close()
	if err != nil {
		fmt.Fprintf(os.Stderr, "fetch: reading %s: %v\n", url, err)
		os.Exit(1)
	}

	fieldName := "<title>"
	q.source = extractQOTD(string(b), fieldName, 3)

	fieldName = "<description>"
	q.quote = extractQOTD(string(b), fieldName, 3)

	file := HTMLFile
	mutex.Lock()
	src, err := ioutil.ReadFile(file)
	if err != nil {
		fmt.Printf("Error reading file %s: %v\n", file, err)
	}
	memoryFile := string(src)

	memoryFile = replaceByID(memoryFile, "id=\"quote\">", q.quote)
	memoryFile = replaceByID(memoryFile, "id=\"qsource\">", q.source)

	err = ioutil.WriteFile(HTMLFile, []byte(memoryFile), 644)
	mutex.Unlock()
	if err != nil {
		fmt.Printf("Error writing to file %s: %v\n", HTMLFile, err)
	}
}

func QOTD() {
	// Initial WOTD load on startup
	log.Println("Initial QOTD() Load")
	getQOTD()

	//==================================
	// Repeat WOTD load every qotdReloadInterval
	ticker := time.NewTicker(time.Hour * qotdReloadInterval)
	for range ticker.C {
		log.Println("Periodic QOTD() Load")
		getQOTD()
	}
}

func extractNYT(text string, str string, rep int) string {
	loc := 0
	start := 0
	for i := 1; i <= rep; i++ {
		loc = strings.Index(text[start:], str) + len(str)
		start = start + loc
	}

	end := strings.Index(text[start:len(text)], "</")
	end = start + end
	return text[start:end]
}

func getNYT() {
	url := nytURL
	var n nyt

	resp, err := http.Get(url)
	if err != nil {
		fmt.Fprintf(os.Stderr, "fetch: %v\n", err)
		os.Exit(1)
	}
	b, err := ioutil.ReadAll(resp.Body)
	resp.Body.Close()
	if err != nil {
		fmt.Fprintf(os.Stderr, "fetch: reading %s: %v\n", url, err)
		os.Exit(1)
	}

	fieldName := "<title>"
	n.title1 = "\"" + extractNYT(string(b), fieldName, 3) + "\""

	fieldName = "<description>"
	n.description1 = extractQOTD(string(b), fieldName, 1)

	fieldName = "<title>"
	n.title2 = "\"" + extractNYT(string(b), fieldName, 4) + "\""

	fieldName = "<description>"
	n.description2 = extractQOTD(string(b), fieldName, 2)

	fieldName = "<title>"
	n.title3 = "\"" + extractNYT(string(b), fieldName, 5) + "\""

	fieldName = "<description>"
	n.description3 = extractQOTD(string(b), fieldName, 3)

	file := HTMLFile
	mutex.Lock()
	src, err := ioutil.ReadFile(file)
	if err != nil {
		fmt.Printf("Error reading file %s: %v\n", file, err)
	}
	memoryFile := string(src)

	memoryFile = replaceByID(memoryFile, "id=\"title1\">", n.title1)
	memoryFile = replaceByID(memoryFile, "id=\"desc1\">", n.description1)
	memoryFile = replaceByID(memoryFile, "id=\"title2\">", n.title2)
	memoryFile = replaceByID(memoryFile, "id=\"desc2\">", n.description2)
	memoryFile = replaceByID(memoryFile, "id=\"title3\">", n.title3)
	memoryFile = replaceByID(memoryFile, "id=\"desc3\">", n.description3)

	err = ioutil.WriteFile(HTMLFile, []byte(memoryFile), 644)
	mutex.Unlock()
	if err != nil {
		fmt.Printf("Error writing to file %s: %v\n", HTMLFile, err)
	}
}

func NYT() {
	// Initial NYT load on startup
	log.Println("Initial NYT() Load\n")
	getNYT()

	//==================================
	// Repeat NYT load every nytReloadInterval
	ticker := time.NewTicker(time.Hour * nytReloadInterval)
	for range ticker.C {
		log.Println("Periodic NYT() Load\n")
		getNYT()
	}
}

func replaceByID(src string, old string, new string) string {
	i := strings.Index(src, old)
	i = i + len(old)
	substr1 := src[:i]
	substr2 := src[i:]
	i = strings.Index(substr2, "</")
	substr2 = substr2[i:]
	src = substr1 + new + substr2
	return src
}

func main() {
	go Weather()
	time.Sleep(15 * time.Second)
	go getForecast()
	time.Sleep(15 * time.Second)
	go QOTD()
	time.Sleep(15 * time.Second)
	go WOTD()
	time.Sleep(15 * time.Second)
	go NYT()

	select {}
}
