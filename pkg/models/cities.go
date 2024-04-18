package models

import (
	"barista/pkg/log"
	"encoding/json"
	"os"
)

type Province struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

type City struct {
	ID         int    `json:"id"`
	Name       string `json:"name"`
	ProvinceID int    `json:"ostan"`
}

var Provinces = []Province{}
var cities = []City{}
var Cities = map[int][]City{}

func init() {
	file, err := os.Open("assets/ostan.json")
	if err != nil {
		panic(err)
	}

	defer file.Close()

	decoder := json.NewDecoder(file)
	err = decoder.Decode(&Provinces)
	if err != nil {
		panic(err)
	}

	file, err = os.Open("assets/shahr.json")
	if err != nil {
		panic(err)
	}

	defer file.Close()

	decoder = json.NewDecoder(file)
	err = decoder.Decode(&cities)
	if err != nil {
		panic(err)
	}

	for _, city := range cities {
		Cities[city.ProvinceID] = append(Cities[city.ProvinceID], city)
	}
	log.GetLog().Info("Provinces loaded successfully")
	log.GetLog().Info("Cities loaded successfully")
}
