package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
	"text/template"
)

type Artist struct {
	Id           float32
	Image        string
	Name         string
	Members      []string
	CreationDate float32
	FirstAlbum   string
	Locations    string
	ConcertDates string
	Relations    string
}

type ConcDates struct {
	LocId    int
	Location string
	Dates    []string
}

type New struct {
	ID           int         `json:"id"`
	Image        string      `json:"image"`
	Name         string      `json:"name"`
	Members      []string    `json:"members"`
	CreationDate int         `json:"creationDate"`
	FirstAlbum   string      `json:"firstAlbum"`
	ConcertDates []ConcDates `json:"concertDates"`
}

var artists []Artist

var tplAll = template.Must(template.ParseFiles("templates/index.html"))
var tplGroup = template.Must(template.ParseFiles("templates/group.html"))
var eRRor = template.Must(template.ParseFiles("templates/404.html"))

func main() {

	mux := http.NewServeMux()
	fs := http.FileServer(http.Dir("CSS"))
	mux.HandleFunc("/", myHandlerMain)
	mux.Handle("/CSS/", http.StripPrefix("/CSS/", fs))
	http.ListenAndServe(":8080", mux)
}

func myHandlerMain(w http.ResponseWriter, r *http.Request) {
	resp, err := http.Get("https://groupietrackers.herokuapp.com/api/artists")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		eRRor.Execute(w, 500)
	}
	data, err2 := ioutil.ReadAll(resp.Body)
	if err2 != nil {
		fmt.Println(err2.Error)
	}
	resp.Body.Close()

	json.Unmarshal(data, &artists)

	if r.URL.String() == "/" {
		tplAll.Execute(w, artists)
	} else {
		var groupie New
		name := r.URL.String()[1:]
		groupID, err := strconv.Atoi(name)
		if err != nil || groupID > 52 || groupID < 1 {
			w.WriteHeader(http.StatusNotFound)
			eRRor.Execute(w, 404)
			return
		}
		relatslink := ""

		for _, group := range artists {

			if int(group.Id) == groupID {

				groupie.ID = int(group.Id)
				groupie.Image = group.Image
				groupie.Name = group.Name
				groupie.Members = group.Members
				groupie.CreationDate = int(group.CreationDate)
				groupie.FirstAlbum = group.FirstAlbum
				relatslink = group.Relations

				break
			}
		}

		resp, err = http.Get(relatslink)
		if err != nil {
			fmt.Println(err.Error)
		}
		data, err2 = ioutil.ReadAll(resp.Body)
		if err2 != nil {
			fmt.Println(err2.Error)
		}
		resp.Body.Close()

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
			loc.LocId = index
			relation := strings.Split(location, ":[")
			loc.Location = relation[0][1 : len(relation[0])-1]
			loc.Dates = strings.Split(relation[1], ",")
			for i, val := range loc.Dates {
				loc.Dates[i] = val[1 : len(val)-1]
			}
			groupie.ConcertDates = append(groupie.ConcertDates, loc)
		}

		tplGroup.Execute(w, groupie)

	}

}
