package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
)

type wotd struct {
	word       string
	definition string
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
	//==============================================
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

	end := strings.Index(text[start:], "</")
	end = start + end
	return text[start:end]
}

func getWOTD() wotd {
	var w wotd
	url := "https://wordsmith.org/awad/rss1.xml"
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

	fieldName := "<title>"
	wotd := extractWOTD(string(b), fieldName, 2)
	w.word = wotd

	fieldName = "<description>"
	definition := extractWOTD(string(b), fieldName, 2)
	w.definition = definition

	return w
}

func extractQOTD(text string, str string, rep int) string {
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

func getQOTD() qotd {
	var q qotd
	url := "https://feeds.feedburner.com/brainyquote/QUOTEBR"
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

	fieldName := "<title>"
	source := extractQOTD(string(b), fieldName, 3)
	q.source = source

	fieldName = "<description>"
	quote := extractQOTD(string(b), fieldName, 3)
	q.quote = quote

	return q
}

func write(file string, word string, definition string) int {
	f, err := os.Create(file)
	if err != nil {
		//fmt.Printf("Error opening file \"%s\": %d\n", file)
		return 1
	}
	defer f.Close()

	_, err = f.WriteString(word + "\n")
	if err != nil {
		//fmt.Printf("Error writing to file \"%s\": %d\n", file)
		return 2
	}

	//fmt.Printf("wrote %d bytes\n", n)

	_, err = f.WriteString(definition + "\n")
	if err != nil {
		//fmt.Printf("Error writing to file \"%s\": %d\n", file)
		return 2
	}
	//fmt.Printf("wrote %d bytes\n", n)
	return 0
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
	entry := getWOTD()
	//fmt.Println(entry.word)
	//fmt.Println(entry.definition)
	write("wotd.txt", entry.word, entry.definition)

	quote := getQOTD()
	//fmt.Println(quote.source)
	//fmt.Println("" + quote.quote + "")
	write("qotd.txt", quote.quote, quote.source)

	weather := getWeather()
	location := fmt.Sprintf("%s (%s N / %s W)", weather.location, weather.latitude, weather.longitude)
	//fmt.Printf("Observered:          %s\n", weather.updated)
	//fmt.Printf("Temperature:         %s\n", weather.temperature)
	//fmt.Printf("Humidity:            %s%%\n", weather.humidity)
	//fmt.Printf("Wind:                %s\n", weather.wind)
	//fmt.Printf("Barometric Presure:  %s\n", weather.pressure)
	//fmt.Printf("Visibilty:           %s mi\n", weather.visibility)
	write("weather.txt", weather.location, weather.updated)

	file := "c:/users/lekrigbaum/Desktop/Family Planner/index.html"
	src, err := ioutil.ReadFile(file)
	if err != nil {
		fmt.Printf("Error reading file %s: %v\n", file, err)
	}

	newfile := replaceByID(string(src), "id=\"quote\">", quote.quote)
	//fmt.Printf("\n\n=== First: replaced quote ========================================================\n")
	//fmt.Println(newfile)
	newfile = replaceByID(newfile, "id=\"qsource\">", quote.source)
	//fmt.Printf("\n\n=== Second: replaced source ======================================================\n")
	//fmt.Println(newfile)
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
}
