package main

import (
	"encoding/json"
	"fmt"
	"github.com/codegangsta/cli"
	. "github.com/jisaw/goUtils"
	"io/ioutil"
	"net/http"
	"os"
)

// Current Weather API Response Struct
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

// 5day Forecast API Struct
type FiveDayJSON struct {
	City struct {
		Id   int
		Name string
	}
	Coord struct {
		Lon float64
		Lat float64
	}
	Country string
	Cod     string
	Message string
	Cnt     int
	List    []struct {
		Dt   int
		Main struct {
			Temp       float32
			Temp_min   float32
			Temp_max   float32
			Pressure   float32
			Sea_level  float32
			Grnd_level float32
			Humidity   int
			Temp_kf    float32
		}
		Weather []struct {
			Id          int
			Main        string
			Description string
			Icon        string
		}
		Clouds struct {
			All int
		}
		Wind struct {
			Speed float32
			Deg   float32
		}
		Sys struct {
			Pod string
		}
		Dt_txt string
	}
}

type Config struct {
	Country string
	City    string
	Unit    string
}

// Open Weather API Endpoint
const weatherQueryURL = "http://api.openweathermap.org/data/2.5/weather?mode=json&"

// Open Weather API 5 Day Forecast Endpoint
const fiveDayQueryURL = "http://api.openweathermap.org/data/2.5/forecast?mode=json&"

// Writes config to json file and returns Config struct
func writeConfig(c Config) Config {
	result, err := json.Marshal(c)
	CheckErr(err)
	err2 := ioutil.WriteFile("/tmp/config.json", result, 0644)
	CheckErr(err2)
	return c
}

// Reads config from json file and returns Config struct
func readConfig() Config {
	file, err := ioutil.ReadFile("/tmp/config.json")
	if err != nil {
		return Config{}
	}
	config := Config{}
	json.Unmarshal(file, &config)
	return config
}

// Returns a string to be amended to the queryWeatherURL
func (c *Config) urlAmendment() string {
	return string("&q=" + c.City + "," + c.Country + "&unit=" + c.Unit)
}

// Retrieves, parses, and prints current weather
func retrieveWeather(c Config) {
	// TODO rework with c Config
	data := WeatherJSON{}
	r, _ := http.Get(weatherQueryURL + c.urlAmendment())
	defer r.Body.Close()

	body, _ := ioutil.ReadAll(r.Body)
	json.Unmarshal(body, &data)
	printCurrentWeather(data)
}

func printCurrentWeather(data WeatherJSON) {
	fmt.Printf("City: %s | ", data.Name)
	fmt.Printf("Temp: %.ff | ", data.Main.Temp)
	fmt.Printf("High: %.ff | ", data.Main.Temp_max)
	fmt.Printf("Low: %.ff | ", data.Main.Temp_min)
	fmt.Printf("Wind Speed: %.fmph | ", data.Wind.Speed)
	fmt.Printf("Description: %s\n", Capitalize(data.Weather[0].Description))
}

func retrieveFiveDay(c Config) {
	data := FiveDayJSON{}
	r, _ := http.Get(fiveDayQueryURL + c.urlAmendment())
	defer r.Body.Close()

	body, _ := ioutil.ReadAll(r.Body)
	json.Unmarshal(body, &data)
	printFiveDayWeather(data)
}

func printFiveDayWeather(data FiveDayJSON) {
	//TODO Print 5 day forecast
	fmt.Printf("WORKING TEST WORKING")
}

func main() {
	app := cli.NewApp()
	app.Name = "wthr"
	app.Version = "0.0.1"
	app.Usage = "Enter a city name or Zip code and get some weather information"
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:  "town, t",
			Value: " ",
			Usage: "Sets the town to recieve weather information about",
		},
		cli.StringFlag{
			Name:  "country, c",
			Value: "us",
			Usage: "Sets the country to use",
		},
		cli.StringFlag{
			Name:  "unit, u",
			Value: "imperial",
			Usage: "Sets the unit of measurement",
		},
		cli.BoolFlag{
			Name:  "5day, 5",
			Usage: "returns a 5 day forecast",
		},
	}
	// Main Logic
	app.Action = func(c *cli.Context) {
		config := readConfig()
		if len(c.String("town")) > 1 {
			config.City = c.String("town")
		}
		if len(c.String("country")) > 1 {
			config.Country = c.String("country")
		}
		if len(c.String("unit")) > 1 {
			config.Unit = c.String("unit")
		}

		writeConfig(config)

		if c.Bool("5day") {
			retrieveFiveDay(config)
		} else {
			retrieveWeather(config)
		}
		//
		// TODO REWORK LOGIC
		//
	}
	app.Run(os.Args)
}
