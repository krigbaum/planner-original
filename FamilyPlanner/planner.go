package main

import (
	"fmt"
    "io/ioutil"
	"net/http"
	"os"
	"strings"
    "time"
)

type wotd struct {
	word       string
	def        string
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

const qotdURL = "http://quotes.rest/qod.xml"
const wotdURL = "http://www.macmillandictionary.com/us/wotd/wotdrss.xml"
const weatherURL = "http://w1.weather.gov/xml/current_obs/KLAF.xml"
const wotdReloadInterval = 5
const qotdReloadInterval = 5
const weatherReloadInterval = 5
const HTMLFile = "/home/krigbaum/devel/go/src/github.com/krigbaum/planner/FamilyPlanner/index.html"

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

func getWeather(memoryFile string) {
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

	fieldName = "<longitude>"
	w.longitude = extractWeather(string(b), fieldName, 1)

	fieldName = "<observation_time>"
	w.updated = extractWeather(string(b), fieldName, 1)

	fieldName = "<weather>"
	w.conditions = extractWeather(string(b), fieldName, 1)
	
	fieldName = "<temperature_string>"
	w.temperature = extractWeather(string(b), fieldName, 1)

	fieldName = "<relative_humidity>"
	w.humidity = extractWeather(string(b), fieldName, 1)

	fieldName = "<wind_string>"
	w.wind = extractWeather(string(b), fieldName, 1)

	fieldName = "<pressure_in>"
	w.pressure = extractWeather(string(b), fieldName, 1)

	fieldName = "<visibility_mi>"
	w.visibility = extractWeather(string(b), fieldName, 1)
    
    fmt.Println("Updated: ", w.updated)

    memoryFile = replaceByID(memoryFile, "id=\"location\">", w.location)
	memoryFile = replaceByID(memoryFile, "id=\"updated\">", w.updated)
	memoryFile = replaceByID(memoryFile, "id=\"cond\">", w.conditions)
	memoryFile = replaceByID(memoryFile, "id=\"temp\">", w.temperature)
	memoryFile = replaceByID(memoryFile, "id=\"humid\">", w.humidity)
	memoryFile = replaceByID(memoryFile, "id=\"pressure\">", w.pressure)
	memoryFile = replaceByID(memoryFile, "id=\"visibility\">", w.visibility)
	memoryFile = replaceByID(memoryFile, "id=\"wind\">", w.wind)

	err = ioutil.WriteFile(HTMLFile, []byte(memoryFile), 0644)
	if err != nil {
		fmt.Printf("Error writing to file %s: %v\n", HTMLFile, err)
	}
}

func Weather(memoryFile string) {
    // Initial Weather load on startup
    getWeather(memoryFile)

    //==================================
    // Repeat Weather load every weatherdReloadInterval
    ticker := time.NewTicker(time.Minute * weatherReloadInterval)
    for range ticker.C {
        getWeather(memoryFile)
    }
}

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

func getWOTD(memoryFile string) {
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

    fieldName = "<summary>"
    w.def = extractWOTD(string(b), fieldName, 1)
    
    fmt.Printf("%v: Word - %s\n", time.Now(), w.word)
    fmt.Printf("%v: Def - %s\n\n", time.Now(), w.def)
    
    memoryFile = replaceByID(memoryFile, "id=\"word\">", w.word)
	memoryFile = replaceByID(memoryFile, "id=\"definition\">", w.def)
    
    err = ioutil.WriteFile(HTMLFile, []byte(memoryFile), 644)
    if err != nil {
		fmt.Printf("Error writing to file %s: %v\n", HTMLFile, err)
	}
}

func WOTD(memoryFile string) {
    // Initial WOTD load on startup
    getWOTD(memoryFile)

    //==================================
    // Repeat WOTD load every wotdReloadInterval
    ticker := time.NewTicker(time.Minute * wotdReloadInterval)
    for range ticker.C {
        getWOTD(memoryFile)
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

func getQOTD(memoryFile string) {
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

    fieldName := "<quote>"
	q.quote = "\"" + extractQOTD(string(b), fieldName, 1) + "\""
 
    fieldName = "<author>"
    q.source = extractQOTD(string(b), fieldName, 1)
	
    fmt.Printf("%v: Quote - %s\n", time.Now(), q.quote)
    fmt.Printf("%v: Src - %s\n\n", time.Now(), q.source)
    
    memoryFile = replaceByID(memoryFile, "id=\"quote\">", q.quote)
	memoryFile = replaceByID(memoryFile, "id=\"qsource\">", q.source)
    
    err = ioutil.WriteFile(HTMLFile,[]byte(memoryFile), 644)
    if err != nil {
		fmt.Printf("Error writing to file %s: %v\n", HTMLFile, err)
	}
}

func QOTD(memoryFile string) {
    // Initial WOTD load on startup
    fmt.Println("Start QOTD()")
    getQOTD(memoryFile)

    //==================================
    // Repeat WOTD load every qotdReloadInterval
    ticker := time.NewTicker(time.Minute * qotdReloadInterval)
    for range ticker.C {
        getQOTD(memoryFile)
    }
}

func replaceByID(src string, old string, new string) string {
	//fmt.Println("new = ", new)
	i := strings.Index(src, old)
	i = i + len(old)
	substr1 := src[:i]
	//fmt.Println("substr1[len(substr1) - 25:] = ", substr1[len(substr1)-25:])
	//fmt.Println()
	substr2 := src[i:]
	//fmt.Println("substr2[:25] = ", substr2[:25])
	i = strings.Index(substr2, "</")
	//fmt.Println("i = ", i)
	substr2 = substr2[i:]
	src = substr1 + new + substr2
	return src
}

func main() {
	file := HTMLFile
	src, err := ioutil.ReadFile(file)
	if err != nil {
		fmt.Printf("Error reading file %s: %v\n", file, err)
	}
    go Weather(string(src))
    //go QOTD(string(src))
    //go WOTD(string(src))
    
    select{}
/*
	quote := getQOTD()

	weather := getWeather()
	location := fmt.Sprintf("%s (%s N / %s W)", weather.location, weather.latitude, weather.longitude)
	//fmt.Printf("Observered:          %s\n", weather.updated)
	//fmt.Printf("Temperature:         %s\n", weather.temperature)
	//fmt.Printf("Humidity:            %s%%\n", weather.humidity)
	//fmt.Printf("Wind:                %s\n", weather.wind)
	//fmt.Printf("Barometric Presure:  %s\n", weather.pressure)
	//fmt.Printf("Visibilty:           %s mi\n", weather.visibility)

	file := "c:/users/lekrigbaum/Desktop/Family Planner/index.html"
	src, err := ioutil.ReadFile(file)
	if err != nil {
		fmt.Printf("Error reading file %s: %v\n", file, err)
	}

	newfile := replaceByID(string(src), "id=\"quote\">", quote.quote)
	newfile = replaceByID(newfile, "id=\"qsource\">", quote.source)
    
	newfile = replaceByID(newfile, "id=\"word\">", entry.word)
	newfile = replaceByID(newfile, "id=\"definition\">", entry.definition)

	newfile = replaceByID(newfile, "id=\"location\">", location)
	newfile = replaceByID(newfile, "id=\"updated\">", weather.updated)
	newfile = replaceByID(newfile, "id=\"cond\">", weather.conditions)
	newfile = replaceByID(newfile, "id=\"temp\">", weather.temperature)
	newfile = replaceByID(newfile, "id=\"humid\">", weather.humidity)
	newfile = replaceByID(newfile, "id=\"pressure\">", weather.pressure)
	newfile = replaceByID(newfile, "id=\"visibility\">", weather.visibility)
	newfile = replaceByID(newfile, "id=\"wind\">", weather.wind)
	err = ioutil.WriteFile(file, []byte(newfile), 0644)
	if err != nil {
		fmt.Printf("Error writing to file %s: %v\n", file, err)
	}
*/
}
