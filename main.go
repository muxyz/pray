package main

import (
	"fmt"
	"net/http"
	"time"

	"mu.dev"
	"github.com/hablullah/go-prayer"
)

var (
	template = `
  <style>
  #name {
    text-align: right;
    width: 60px;
    margin-right: 20px;
    display: inline-block;
  }
  #times {
    padding-top: 100px;
  }
  </style>
  <div id="times">
  %s
  </div>
`
)

func dateFormat(v time.Time) string {
	return v.Format(time.Kitchen)
}

func printSchedule(sched prayer.Schedule) string {
	format := func(k string, v time.Time) string {
		return fmt.Sprintf(`<div><span id="name">%s</span><span id="time">%s</span></div>`, k, dateFormat(v))
	}

	var str string
	str += format("Fajr", sched.Fajr)
	str += format("🌅", sched.Sunrise)
	str += format("Zuhr", sched.Zuhr)
	str += format("Asr", sched.Asr)
	str += format("Maghrib", sched.Maghrib)
	str += format("Isha", sched.Isha)

	return str
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	date := time.Now().Format(time.DateOnly)

	// Calculate prayer schedule in London for 2023.
	// Since London in higher latitude, make sure to enable the adapter.
	europeLondon, _ := time.LoadLocation("Europe/London")
	londonSchedules, _ := prayer.Calculate(prayer.Config{
		Latitude:            51.41334,
		Longitude:           -0.36701,
		Timezone:            europeLondon,
		TwilightConvention:  prayer.MWL(),
		AsrConvention:       prayer.Shafii,
		HighLatitudeAdapter: prayer.NearestLatitude(),
		PreciseToSeconds:    true,
	}, 2024)

	head  := "I bear witness there is no deity but God, and I bear witness that Muhammad is the Messenger of God."
	content := ""

	for _, sched := range londonSchedules {
		if sched.Date != date {
			continue
		}
		content += printSchedule(sched)
		// schedule for today
		out := mu.Template("Pray", "Islamic Prayer Times", head, fmt.Sprintf(template, content))
		w.Write([]byte(out))
		return
	}
}

func main() {
	http.HandleFunc("/", indexHandler)

	http.ListenAndServe(":8082", nil)
}
