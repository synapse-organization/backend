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
var Cities = []City{}
var ProvinceCities = map[int][]City{}

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
	err = decoder.Decode(&Cities)
	if err != nil {
		panic(err)
	}

	for _, city := range Cities {
		ProvinceCities[city.ProvinceID] = append(ProvinceCities[city.ProvinceID], city)
	}
	log.GetLog().Info("Provinces loaded successfully")
	log.GetLog().Info("Cities loaded successfully")
}
