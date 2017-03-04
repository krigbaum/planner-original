package main

import (
	"crypto/tls"
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

const DEBUG = true

const qotdURL = "https://www.quotesdaddy.com/feed"
const wotdURL = "https://www.merriam-webster.com/word-of-the-day"
const weatherURL = "http://w1.weather.gov/xml/current_obs/KLAF.xml"
const forecastURL = "http://www.accuweather.com/en/us/west-lafayette-in/47906/weather-forecast/2135952"
const nytURL = "http://rss.nytimes.com/services/xml/rss/nyt/US.xml"

const wotdReloadInterval = 12
const qotdReloadInterval = 12
const weatherReloadInterval = 1
const forecastReloadInterval = 1
const nytReloadInterval = 1
const timeCheckInterval = 3

const HTMLFile = "/home/pi/devel/src/github.com/pi/planner/index.html"
//const HTMLFile = "c:/Users/lekrigbaum/Desktop/go/src/github.com/krigbaum/planner/index.html"
//const HTMLFile = "c:/wamp64/www/planner/index.html"
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
		log.Printf("\n***** Error: Exit on http.Get(%s) in function getWeather() *****\n\n", url)
		//os.Exit(1)
		return
	}
	b, err := ioutil.ReadAll(resp.Body)
	resp.Body.Close()
	if err != nil {
		fmt.Fprintf(os.Stderr, "fetch: reading %s: %v\n", url, err)
		log.Printf("\n***** Error: Exit on resp.Body,Close() in function getWeather() *****\n\n")
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
	winds := strings.Split(w.wind, "(")
	w.wind = winds[0]

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
		log.Printf("\n***** Error: Exit on ioutil.ReadFile(%s) with errcode %d in function getWeather() *****\n\n", file, err)
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
		log.Printf("\n***** Error: Exit on ioutil.WriteFile(%s) with errcode %d in function getWeather() *****\n\n", file, err)
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
	log.Printf("\n***** Error: Exit on range ticker in function Weather\n\n")
}

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
		log.Printf("\n***** Error: Exit on http.Get(%s) in function getForecast() *****\n\n", url)
		//os.Exit(1)
		return
	}
	b, err := ioutil.ReadAll(resp.Body)
	resp.Body.Close()
	if err != nil {
		fmt.Fprintf(os.Stderr, "fetch: reading %s: %v\n", url, err)
		log.Printf("\n***** Error: Exit on resp.Body,Close() in function getForecast() *****\n\n")
		os.Exit(1)
	}
	parts := strings.SplitN(string(b), "Tonight", 2)

	//fieldName := "day=1\">"
	period := "Tonight"

	fieldName := "<span class=\"large-temp\">"
	temp := extractForecast(parts[1], fieldName, 1)

	fieldName = "<span class=\"cond\">"
	cond := extractForecast(parts[1], fieldName, 1)

	tonight := period + ": " + cond + ", low of " + temp

	fieldName = "day=2\">"
	period = extractForecast(parts[1], fieldName, 2)

	fieldName = "<span class=\"large-temp\">"
	temp = extractForecast(parts[1], fieldName, 2)

	fieldName = "<span class=\"cond\">"
	cond = extractForecast(parts[1], fieldName, 2)

	tomorrow := period + ": " + cond + ", high of " + temp
	file := HTMLFile
	mutex.Lock()
	src, err := ioutil.ReadFile(file)
	if err != nil {
		fmt.Printf("Error reading file %s: %v\n", file, err)
		log.Printf("\n***** Error: Exit on ioutil.ReadFile(%s) with errcode %d in function getForecast() *****\n\n", file, err)
	}
	memoryFile := string(src)

	memoryFile = replaceByID(memoryFile, "id=\"tonight\">", tonight)
	memoryFile = replaceByID(memoryFile, "id=\"tomorrow\">", tomorrow)

	err = ioutil.WriteFile(HTMLFile, []byte(memoryFile), 644)
	mutex.Unlock()
	if err != nil {
		fmt.Printf("Error writing to file %s: %v\n", HTMLFile, err)
		log.Printf("\n***** Error: Exit on ioutil.WriteFile(%s) with errcode %d in function getForecast() *****\n\n", file, err)
	}
}

func Forecast() {
	// Initial WOTD load on startup
	log.Println("Initial Forecast() Load")
	getForecast()

	//==================================
	// Repeat forecast load every forecastReloadInterval
	ticker := time.NewTicker(time.Hour * forecastReloadInterval)
	for range ticker.C {
		log.Println("Periodic Forecast() Load")
		getForecast()
	}
	log.Printf("\n***** Error: Exit on range ticker in function Weather\n\n")
}

//=======================================================================

func extractWOTD(text string, str string, rep int) string {
	loc := 0
	start := 0
	for i := 1; i <= rep; i++ {
		loc = strings.Index(text[start:], str) + len(str)
		start = start + loc
	}
	end := -1
	if str == "<dt>" {
		end = strings.Index(text[start:len(text)], "</dt>")
	} else {
		end = strings.Index(text[start:len(text)], "</")
	}
	end = start + end
	return text[start:end]
}

func extractNumberOfDefs(text string, str string) int {
	//loc := 0
	//start := 0
	//reps = strconv.Atoi(strings.LastIndex(text[start:], str) + len(str))
	reps := strings.Count(text, str)

	return reps
}

func getQuery() string {
	url := wotdURL
	var w wotd

	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := &http.Client{Transport: tr}
	resp, err := client.Get(url)
	if err != nil {
		fmt.Fprintf(os.Stderr, "fetch: %v\n", err)
		//os.Exit(1)
		return "ERR"
	}
	b, err := ioutil.ReadAll(resp.Body)
	resp.Body.Close()
	if err != nil {
		fmt.Fprintf(os.Stderr, "fetch: reading %s: %v\n", url, err)
		os.Exit(1)
	}

	fieldName := "<title>"
	w.word = extractWOTD(string(b), fieldName, 1)
	wotd := fmt.Sprintf("%s: ", w.word)
	split1 := strings.Split(wotd, ": ")
	split2 := strings.Split(split1[1], " ")
	word := strings.ToLower(split2[0])

	queryPt1 := "http://www.dictionaryapi.com/api/v1/references/collegiate/xml/"
	queryPt2 := "?key=6becb26a-c8cc-4b72-813f-55849be7b7a5"
	query := queryPt1 + word + queryPt2

	return query
}

func getWOTD() {
	query := getQuery()
	if query == "ERR" {
		return
	}


	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := &http.Client{Transport: tr}
	resp, err := client.Get(query)
	if err != nil {
		fmt.Fprintf(os.Stderr, "fetch: %v\n", err)
		//os.Exit(1)
		return
	}
	b, err := ioutil.ReadAll(resp.Body)
	resp.Body.Close()
	if err != nil {
		fmt.Fprintf(os.Stderr, "fetch: reading %s: %v\n", query, err)
		os.Exit(1)
	}

    if DEBUG {
        fmt.Println("WOTD returned: ", string(b))
    }
    
	fieldName := "<ew>"
	word := extractWOTD(string(b), fieldName, 1)

	fieldName = "<pr>"
	pronounce := "(" + extractWOTD(string(b), fieldName, 1) + ")"

	fieldName = "<fl>"
	pos := extractWOTD(string(b), fieldName, 1)

	fieldName = "<sn>"
	//numDefs, _ := strconv.Atoi(extractNumberOfDefs(string(b), fieldName))
	numDefs := extractNumberOfDefs(string(b), fieldName)
    
    if DEBUG {
        fmt.Println("numDefs is: ", numDefs)
    }
    
	fieldName = "<dt>"
	defs := ""
	reps := 1
	for reps <= numDefs {
		//tmp := extractWOTD(string(b), fieldName, reps + 1)
		//fmt.Println("tmp = ", tmp)
		defs = defs + strconv.Itoa(reps) + stripHTML(extractWOTD(string(b), fieldName, reps)) + "<br>"
		
        if DEBUG {
		  fmt.Printf("line %d) %s\n", reps, defs)
        }
        reps++
	}
    
    if DEBUG {
        fmt.Println("definition: ", defs)
    }

	file := HTMLFile
	var mutex = &sync.Mutex{}
	mutex.Lock()
	src, err := ioutil.ReadFile(file)
	if err != nil {
		fmt.Printf("Error reading file %s: %v\n", file, err)
	}
	memoryFile := string(src)
	//fmt.Println(memoryFile)
	memoryFile = replaceByID(memoryFile, "id=\"word\">", word)
	memoryFile = replaceByID(memoryFile, "id=\"pronounce\">", pronounce)
	memoryFile = replaceByID(memoryFile, "id=\"pos\">", pos)
	memoryFile = replaceByID(memoryFile, "id=\"definitions\">", defs)

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
		//os.Exit(1)
		return
	}
	b, err := ioutil.ReadAll(resp.Body)
	resp.Body.Close()
	if err != nil {
		fmt.Fprintf(os.Stderr, "fetch: reading %s: %v\n", url, err)
		os.Exit(1)
	}

	fieldName := "<description>"
	q.quote = extractQOTD(string(b), fieldName, 2)

	//fieldName = "<description>"
	//q.quote = extractQOTD(string(b), fieldName, 3)

	file := HTMLFile
	mutex.Lock()
	src, err := ioutil.ReadFile(file)
	if err != nil {
		fmt.Printf("Error reading file %s: %v\n", file, err)
	}
	memoryFile := string(src)

	memoryFile = replaceByID(memoryFile, "id=\"quote\">", q.quote)
	//memoryFile = replaceByID(memoryFile, "id=\"qsource\">", q.source)

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
	// Repeat QOTD load every qotdReloadInterval
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
	log.Printf("Wrote \"%s\" to title1", n.title1)
	memoryFile = replaceByID(memoryFile, "id=\"desc1\">", n.description1)
	log.Printf("Wrote \"%s\" to description1", n.description1)
	memoryFile = replaceByID(memoryFile, "id=\"title2\">", n.title2)
	log.Printf("Wrote \"%s\" to title2", n.title2)
	memoryFile = replaceByID(memoryFile, "id=\"desc2\">", n.description2)
	log.Printf("Wrote \"%s\" to description2", n.description2)
	memoryFile = replaceByID(memoryFile, "id=\"title3\">", n.title3)
	log.Printf("Wrote \"%s\" to title3", n.title3)
	memoryFile = replaceByID(memoryFile, "id=\"desc3\">", n.description3)
	log.Printf("Wrote \"%s\" to description3", n.description3)

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

func stripHTML(orig string) string {
	//fmt.Println("Original string: ", orig)
	//fmt.Println()
	result := ""
	start := 0
	end := strings.Index(orig, "<")
	// Another HTML Tag found.
	for end != -1 {
		result = result + orig[start:end]
		//fmt.Println("1) result = ", result)
		orig = strings.SplitN(orig, ">", 2)[1]
		//fmt.Println("2) str = ", orig)
		end = strings.Index(orig, "<")
		//fmt.Println("start =", start)
		//fmt.Println("end = ", end)
		// No other HTML tags found
		if end == -1 {
			result = result + orig
			return result
		}
		//fmt.Println()
		//fmt.Println("start =", start)
		//fmt.Println("result = ", result)
		end = strings.Index(orig, "<")
		//fmt.Println("end = ", end)
		//fmt.Println()
		//end := strings.Index(str, ">")
	}
	result = orig
	return result
}

func TimeCheck() {
	// Initial Weather load on startup
	log.Println("***** Time Check *****")

	// Repeat Time Check every timeCheckInterval minutes
	ticker := time.NewTicker(time.Minute * timeCheckInterval)
	for range ticker.C {
		log.Println("***** Time Check *****")
	}
	log.Printf("\n***** Error: Exit on range ticker in function TimeCheck\n\n")
}

func main() {
	log.Printf("\n\n")
	log.Printf("****************************************************************\n")
	log.Printf("*                                                              *\n")
	log.Printf("*                  Starting Family Planner                     *\n")
	log.Printf("*                                                              *\n")
	log.Printf("****************************************************************\n\n")

	go Weather()
	time.Sleep(10 * time.Second)
	go Forecast()
	time.Sleep(10 * time.Second)
	go QOTD()
	time.Sleep(1 * time.Second)
	go WOTD()
	go TimeCheck()
	//time.Sleep(15 * time.Second)
	//go NYT()

	select {}
}
