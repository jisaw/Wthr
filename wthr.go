package main

import (
	"encoding/json"
	"fmt"
	"github.com/codegangsta/cli"
	"io/ioutil"
	"net/http"
	"os"
)

type WeatherJSON struct {
	Coord   map[string]float64
	Weather []struct {
		Id          int
		Main        string
		Description string
		icon        string
	}
	Base string
	Main struct {
		Temp     float32
		Pressure int
		Humidity int
		Temp_min float64
		Temp_max float64
	}
	Wind struct {
		Speed float64
		Deg   int
	}
	Clouds map[string]int
	Rain   map[string]int
	Dt     int
	Sys    struct {
		Type    int
		Id      int
		Message float64
		Country string
		Sunrise int
		Sunset  int
	}
	Id   int
	Name string
	Cod  int
}

const queryURL = "http://api.openweathermap.org/data/2.5/weather?units=imperial&q="

func retrieveWeather(area string) {
	data := WeatherJSON{}
	r, _ := http.Get(queryURL + area)
	defer r.Body.Close()

	body, _ := ioutil.ReadAll(r.Body)
	json.Unmarshal(body, &data)
	fmt.Printf("City: %s | ", data.Name)
	fmt.Printf("Temp: %-3.2ff | ", data.Main.Temp)
	fmt.Printf("Description: %s | ", data.Weather[0].Description)
}

func main() {
	app := cli.NewApp()
	app.Name = "wthr"
	app.Version = "0.0.1"
	app.Usage = "Enter a city name or Zip code and get some weather information"
	app.Action = func(c *cli.Context) {
		//println("This is working")
		if len(c.Args()) >= 1 {
			retrieveWeather(c.Args()[0])
		} else {
			println("Please pass in a city name or zip code!")
		}
	}

	app.Run(os.Args)
}
