package main

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strings"
)

// структура идентичная API
type Artist struct {
	ID           float32
	Image        string
	Name         string
	Members      []string
	CreationDate float32
	FirstAlbum   string
	Locations    string
	ConcertDates string
	Relations    string
}

// структура, составленная из relation
// (составляется отдельным обработчиком)
type ConcDates struct {
	LocID    int
	Location string
	Dates    []string
}

type LocationsList struct {
	Index []struct {
		Id        int      `json:"id"`
		Locations []string `json:"locations"`
		Dates     string   `json:"dates"`
	} `json:"index"`
}

type DatesList struct {
	index struct {
		ID    int
		Dates []string
	}
}

// детальная структура по исполнителю для отображения на личной странице
type ArtistFull struct {
	ID           float32
	Image        string
	Name         string
	Members      []string
	CreationDate float32
	FirstAlbum   string
	ConcertDates []ConcDates
	Locations    []string
	CDates       []string
}

// основная база
var artists []Artist
var db []ArtistFull

// формируем основную базу исполнителей
func GetArtistBase(w http.ResponseWriter) {
	resp, err := http.Get("https://groupietrackers.herokuapp.com/api/artists")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		tplAll.ExecuteTemplate(w, "error.html", 500)
		return
	}
	data, err2 := ioutil.ReadAll(resp.Body)
	if err2 != nil {
		w.WriteHeader(http.StatusInternalServerError)
		tplAll.ExecuteTemplate(w, "error.html", 500)
		return
	}
	resp.Body.Close()
	json.Unmarshal(data, &artists)

	// delaem polnuyu bazu bez relats
	for _, group := range artists {
		var a ArtistFull
		a.ID = group.ID
		a.Image = group.Image
		a.Name = group.Name
		a.Members = group.Members
		a.CreationDate = group.CreationDate
		a.FirstAlbum = group.FirstAlbum
		db = append(db, a)
	}
	getLocationsForArtist(w)
	// getDatesForArtist(w)

	return
}

// формирует структуру ArtistFull для указанного испольнителя
func GetFullInfoForArtist(w http.ResponseWriter, groupID int) ArtistFull {
	relatslink := artists[groupID-1].Relations

	// запрашиваем Relations
	resp, err := http.Get(relatslink)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		tplAll.ExecuteTemplate(w, "error.html", 500)
		return db[groupID-1]
	}

	data, err2 := ioutil.ReadAll(resp.Body)
	if err2 != nil {
		w.WriteHeader(http.StatusInternalServerError)
		tplAll.ExecuteTemplate(w, "error.html", 500)
		return db[groupID-1]
	}
	resp.Body.Close()

	// обрабатываем полученные данные и формируем массив струтур ConcDates
	for index, value := range data {
		if index != 0 && value == '{' {
			data = data[index+1:]
			break
		}
	}
	data = data[:len(data)-4]
	stringData := string(data)
	uniqueLoc := strings.Split(stringData, "],")
	for index, location := range uniqueLoc {
		var loc ConcDates
		loc.LocID = index
		relation := strings.Split(location, ":[")
		loc.Location = relation[0][1 : len(relation[0])-1]
		loc.Dates = strings.Split(relation[1], ",")
		for i, val := range loc.Dates {
			loc.Dates[i] = val[1 : len(val)-1]
		}
		db[groupID-1].ConcertDates = append(db[groupID-1].ConcertDates, loc)
	}

	return db[groupID-1]
}
