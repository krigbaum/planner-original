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

func getWeather() weather {
	var w weather
	url := "http://w1.weather.gov/xml/current_obs/KLAF.xml"
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
	//fmt.Printf("\n\n\n%s\n", b)

	fieldName := "<location>"
	location := extractWeather(string(b), fieldName, 1)
	w.location = location

	fieldName = "<latitude>"
	latitude := extractWeather(string(b), fieldName, 1)
	w.latitude = latitude

	fieldName = "<longitude>"
	longitude := extractWeather(string(b), fieldName, 1)
	w.longitude = longitude[1:]

	fieldName = "<observation_time>"
	updated := extractWeather(string(b), fieldName, 1)
	w.updated = updated

	fieldName = "<weather>"
	conditions := extractWeather(string(b), fieldName, 1)
	w.conditions = conditions

	fieldName = "<temperature_string>"
	temperature := extractWeather(string(b), fieldName, 1)
	w.temperature = temperature

	fieldName = "<relative_humidity>"
	humidity := extractWeather(string(b), fieldName, 1)
	w.humidity = humidity

	fieldName = "<wind_string>"
	wind := extractWeather(string(b), fieldName, 1)
	w.wind = wind

	fieldName = "<pressure_in>"
	pressure := extractWeather(string(b), fieldName, 1)
	w.pressure = pressure

	fieldName = "<visibility_mi>"
	visibility := extractWeather(string(b), fieldName, 1)
	w.visibility = visibility

	return w
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

func getWOTD() wotd {
    var w wotd
    url := wotdURL
   
   // Initial Run
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
    
    fmt.Printf("%v: %s\n", time.Now(), w.word)
    //==================================
   
    ticker := time.NewTicker(time.Minute * 5)
    for _ = range ticker.C {
        
        resp, err = http.Get(url)
        if err != nil {
            fmt.Fprintf(os.Stderr, "fetch: %v\n", err)
            os.Exit(1)
        }
        b, err = ioutil.ReadAll(resp.Body)
        resp.Body.Close()
        if err != nil {
            fmt.Fprintf(os.Stderr, "fetch: reading %s: %v\n", url, err)
            os.Exit(1)
        }
        fieldName = "<title>"
        w.word = extractWOTD(string(b), fieldName, 2)
    
        fieldName = "<summary>"
        w.def = extractWOTD(string(b), fieldName, 1)
        
        fmt.Printf("%v: %s\n", time.Now(), w.word)
    
    }
    return w
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

func getQOTD() qotd {
	var q qotd
	url := "http://quotes.rest/qod.xml"
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
	
    return q
}

func replaceByID(src string, old string, new string) string {
	//fmt.Println("new = ", new)
	i := strings.Index(src, old)
	i = i + len(old)
	substr1 := src[:i]
	fmt.Println("substr1[len(substr1) - 25:] = ", substr1[len(substr1)-25:])
	fmt.Println()
	substr2 := src[i:]
	fmt.Println("substr2[:25] = ", substr2[:25])
	i = strings.Index(substr2, "</")
	fmt.Println("i = ", i)
	substr2 = substr2[i:]
	src = substr1 + new + substr2
	return src
}

func main() {
	//entry := getWOTD()
    getWOTD()
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
