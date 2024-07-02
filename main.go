package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"text/template"
)

const API_KEY = "1c26f3d561d521c4174ff1f8c341d0da"
const URL = "https://api.openweathermap.org/data/3.0/onecall"

type Weather struct {
	ID          int    `json:"id"`
	Main        string `json:"main"`
	Description string `json:"description"`
	Icon        string `json:"icon"`
}

type Rain struct {
	OneHour float64 `json:"1h"`
}

type PageData struct {
	Lat     string
	Lng     string
	Weather *Weather
	Error   string
}

type Current struct {
	Dt         int       `json:"dt"`
	Sunrise    int       `json:"sunrise"`
	Sunset     int       `json:"sunset"`
	Temp       float64   `json:"temp"`
	FeelsLike  float64   `json:"feels_like"`
	Pressure   int       `json:"pressure"`
	Humidity   int       `json:"humidity"`
	DewPoint   float64   `json:"dew_point"`
	UVI        float64   `json:"uvi"`
	Clouds     int       `json:"clouds"`
	Visibility int       `json:"visibility"`
	WindSpeed  float64   `json:"wind_speed"`
	WindDeg    int       `json:"wind_deg"`
	WindGust   float64   `json:"wind_gust"`
	Weather    []Weather `json:"weather"`
	Rain       Rain      `json:"rain"`
}

func weatherHandler(w http.ResponseWriter, r *http.Request) {
	lat := r.FormValue("lat")
	lng := r.FormValue("lng")
	tmpl, err := template.ParseFiles("weather.html")
	pageData := PageData{
		Lat:     lat,
		Lng:     lng,
		Weather: nil,
		Error:   "",
	}

	if lat == "" || lng == "" {
		tmpl.Execute(w, pageData)

		return
	}

	if err != nil {
		pageData.Error = err.Error()
		tmpl.Execute(w, pageData)
		return
	}

	url := fmt.Sprintf(URL+"?lat=%s&lon=%s&appid=%s", lat, lng, API_KEY)

	response, err := http.Get(url)

	fmt.Println(url)

	if err != nil {
		pageData.Error = err.Error()
		tmpl.Execute(w, pageData)
		return
	}

	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		pageData.Error = "Weather API has issue"
		tmpl.Execute(w, pageData)
		return

	}

	body, err := ioutil.ReadAll(response.Body)

	if err != nil {
		pageData.Error = err.Error()
		tmpl.Execute(w, pageData)
		return
	}

	var weatherData struct {
		Current Current `json:"current"`
	}

	err = json.Unmarshal(body, &weatherData)

	if err != nil {
		pageData.Error = err.Error()
		tmpl.Execute(w, pageData)
		return
	}

	weather := weatherData.Current.Weather[0]
	pageData.Weather = &weather

	err = tmpl.Execute(w, pageData)

	if err != nil {
		pageData.Error = err.Error()
		tmpl.Execute(w, pageData)
		return
	}
}

func main() {
	http.HandleFunc("/", weatherHandler)
	fmt.Println("Server is listening on port 8080...")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		fmt.Println("Error starting server:", err)
	}
}
