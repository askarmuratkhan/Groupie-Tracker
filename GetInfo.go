package main

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
)

func getLocationsForArtist(w http.ResponseWriter) {
	var loc LocationsList
	resp, err := http.Get("https://groupietrackers.herokuapp.com/api/locations")
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
	json.Unmarshal(data, &loc)

	for index, structure := range loc.Index {
		db[index].Locations = structure.Locations
	}

	return
}

// func getDatesForArtist(w http.ResponseWriter, link string) []string {
// 	var dates DatesList
// 	resp, err := http.Get(link)
// 	if err != nil {
// 		w.WriteHeader(http.StatusInternalServerError)
// 		tplAll.ExecuteTemplate(w, "error.html", 500)
// 		return nil
// 	}
// 	data, err2 := ioutil.ReadAll(resp.Body)
// 	if err2 != nil {
// 		w.WriteHeader(http.StatusInternalServerError)
// 		tplAll.ExecuteTemplate(w, "error.html", 500)
// 		return nil
// 	}
// 	resp.Body.Close()
// 	json.Unmarshal(data, dates)

// 	return dates.Dates
// }
