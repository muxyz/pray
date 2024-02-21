package main

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/hablullah/go-prayer"
	"mu.dev"
)

var (
	template = `
  <style>
  #name {
    text-align: right;
    width: 100px;
    margin-right: 20px;
    display: inline-block;
  }
  #shahada {
    padding: 100px 20px;
    text-align: center;
    max-width: 400px;
  }
  #times {
    max-width: 400px;
    text-align: center;
    font-size: 2em;
  }
  </style>
  <div id="shahada">
  %s
  </div>
  <div id="times">
  %s
  </div>
`
)

type City struct {
	Name     string
	Lat, Lon float64
	Location string
}

var cities = []City{
	City{Name: "London", Lat: 51.41334, Lon: -0.36701, Location: "Europe/London"},
	City{Name: "Doha", Lat: 25.286106, Lon: 51.534817, Location: "Asia/Qatar"},
	City{Name: "Mecca", Lat: 21.422487, Lon: 39.826206, Location: "Asia/Riyadh"},
	City{Name: "New York", Lat: 40.730610, Lon: -73.935242, Location: "America/New_York"},
	City{Name: "Sydney", Lat: -33.865143, Lon: 151.209900, Location: "Australia/Sydney"},
	City{Name: "Tokyo", Lat: 35.672855, Lon: 139.817413, Location: "Asia/Tokyo"},
}

func dateFormat(v time.Time) string {
	return v.Format(time.Kitchen)
}

func printSchedule(city string, sched prayer.Schedule) string {
	format := func(k string, v time.Time) string {
		return fmt.Sprintf(`<div><span id="name">%s</span><span id="time">%s</span></div>`, k, dateFormat(v))
	}

	str := fmt.Sprintf(`<h2 id="%s">%s</h2>`, strings.ReplaceAll(city, " ", ""), city)
	str += format("Fajr", sched.Fajr)
	str += format("🌅", sched.Sunrise)
	str += format("Zuhr", sched.Zuhr)
	str += format("Asr", sched.Asr)
	str += format("Maghrib", sched.Maghrib)
	str += format("Isha", sched.Isha)
	str += "<br><br>"
	return str
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	date := time.Now().Format(time.DateOnly)

	head := "<p>I bear witness there is no deity but God, and I bear witness that Muhammad is the Messenger of God.</p>"
	nav := ""
	content := ""

	for _, city := range cities {
		// Calculate prayer schedule in London for 2023.
		// Since London in higher latitude, make sure to enable the adapter.
		tz, _ := time.LoadLocation(city.Location)
		schedules, _ := prayer.Calculate(prayer.Config{
			Latitude:            city.Lat,
			Longitude:           city.Lon,
			Timezone:            tz,
			TwilightConvention:  prayer.MWL(),
			AsrConvention:       prayer.Shafii,
			HighLatitudeAdapter: prayer.NearestLatitude(),
			PreciseToSeconds:    true,
		}, 2024)

		nav += fmt.Sprintf(`<a href="#%s" class="head">%s</a>`, strings.ReplaceAll(city.Name, " ", ""), city.Name)

		for _, sched := range schedules {
			if sched.Date != date {
				continue
			}

			content += printSchedule(city.Name, sched)
		}
	}

	out := mu.Template("Pray", "Islamic Prayer Times", nav, fmt.Sprintf(template, head, content))
	w.Write([]byte(out))
	return
}

func main() {
	http.HandleFunc("/", indexHandler)

	http.ListenAndServe(":8082", nil)
}
