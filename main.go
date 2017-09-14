package main

import (
	"net/http"
	"fmt"
	"encoding/json"
	"time"
	"os"
	"flag"
)

type Meals struct {
	Date    string `json:"date"`
	Day     string `json:"day"`
	Future  bool   `json:"future"`
	Meals   []Meal `json:"meals"`
	Weekday int    `json:"dayNumber"`
}

type Meal struct {
	Id       string `json:"id"`
	Type     string `json:"type"`
	Name     string `json:"name"`
	Calories int    `json:"calories"`
}

func EmojiByType(typeOfFood string) string {
	switch typeOfFood {
	case "carne":
		return "\U0001f356"
	case "fisch":
		return "\U0001f41f"
	case "salat":
		return "\U0001f331"
	case "vegi":
		return "\U0001f338"
	default:
		return "\U0001f374"
	}
}

func main() {
	apiBaseUrl := "http://altepost.sipgate.net/api.php"

	t := time.Now()
	isoYear, isoWeek := t.ISOWeek()
	isoWeekDay := int(t.Weekday())

	when := flag.String("when", "today", "Which date would you like to check? Yesterday, today or tomorrow?")
	flag.Parse()

	switch *when {
	case "yesterday":
		isoWeekDay += -1

		if isoWeekDay < 1 {
			fmt.Fprint(os.Stderr, "yesterday was in the last week. Can't do that :-(")
			return
		}
	case "tomorrow":
		isoWeekDay += 1

		if isoWeekDay > 7 {
			fmt.Fprint(os.Stderr, "tomorrow seems to be in the next week. Can't do that :-(")
			return
		}
	}

	response, err := http.Get(fmt.Sprintf("%s?kw=%02d/%d&day=%d", apiBaseUrl, isoWeek, isoYear, isoWeekDay))
	if err != nil {
		fmt.Fprintf(os.Stderr, "an occured during meal request: %v\n", err)
		return
	}
	defer response.Body.Close()

	meal := &Meals{}
	err = json.NewDecoder(response.Body).Decode(meal)
	if err != nil {
		fmt.Fprintf(os.Stderr, "the meal response could not be decoded: %v\n", err)
		return
	}

	if len(meal.Meals) == 0 {
		fmt.Fprint(os.Stderr, "nothing to eat today :-(\n")
		return
	}

	for _, value := range meal.Meals {
		fmt.Printf("%s  %s\n", EmojiByType(value.Type), value.Name)
	}
}
