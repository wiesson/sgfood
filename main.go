package main

import (
	"net/http"
	"fmt"
	"encoding/json"
	"time"
	"os"
	"strconv"
)

const apiUrl string = "http://altepost.sipgate.net/api.php"

type Meals struct {
	Date   string `json:"date"`
	Day    string `json:"day"`
	Future bool `json:"future"`
	Meals  []Meal `json:"meals"`
}

type Meal struct {
	Id       string `json:"id"`
	Type     string `json:"type"`
	Name     string `json:"name"`
	Calories int `json:"calories"`
	// Rating   []Rating `json:"rating"`
}

type Rating struct {
	Good int `json:"good"`
	Bad  int `json:"bad"`
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
	t := time.Now()
	meals := &Meals{}
	isoYear, isoWeek := t.ISOWeek()
	isoWeekDay := int(t.Weekday())

	if len(os.Args) > 1 {
		if os.Args[1] == "tomorrow" {
			isoWeekDay += 1
		}

		if os.Args[1] == "yesterday" && isoWeekDay > 0 {
			isoWeekDay -= 1
		}
	}

	s := fmt.Sprintf("%s?kw=%02d/%d&day=%d", apiUrl, isoWeek, isoYear, isoWeekDay)
	response, err := http.Get(s)
	defer response.Body.Close()

	if err != nil {
		panic(err)
	}

	err = json.NewDecoder(response.Body).Decode(meals)

	if err != nil {
		panic(err)
	}

	for _, value := range meals.Meals {
		s := fmt.Sprintf("%s  %s\n", EmojiByType(value.Type), value.Name)
		fmt.Print(s)
	}
}
