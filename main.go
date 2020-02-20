package main

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
	"text/template"
	// "fmt"
)

// структура идентичная API
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

// структура, составленная из relation
// (составляется отдельным обработчиком)
type ConcDates struct {
	LocId    int
	Location string
	Dates    []string
}

// детальная структура по исполнителю для отображения на личной странице
type ArtistFull struct {
	ID           int         `json:"id"`
	Image        string      `json:"image"`
	Name         string      `json:"name"`
	Members      []string    `json:"members"`
	CreationDate int         `json:"creationDate"`
	FirstAlbum   string      `json:"firstAlbum"`
	ConcertDates []ConcDates `json:"concertDates"`
}

// основная база
var artists []Artist

// страницы, используемые в проекте
var tplAll = template.Must(template.ParseGlob("templates/*.html"))

// программа использует API "https://groupietrackers.herokuapp.com/api"
// и представляет данные в двух видах:
// 1. Общий список групп с названиями и фото
// 2. Индивидуальная страница с полными данными по группе, имеющимися на указанном выше ресурсе
func main() {

	mux := http.NewServeMux()
	fs := http.FileServer(http.Dir("CSS"))
	mux.HandleFunc("/", myHandlerMain)
	mux.Handle("/CSS/", http.StripPrefix("/CSS/", fs))
	http.ListenAndServe(":8080", mux)
}

// Основной обработчик
func myHandlerMain(w http.ResponseWriter, r *http.Request) {
	
	
	// shablon := r.FormValue("toFind")
	// fmt.Println(shablon)
	//  TYPE := r.FormValue("searchType")
	//  fmt.Println(TYPE)

	if r.Method == "POST" {
		if r.URL.String() != "/" {
			w.WriteHeader(http.StatusBadRequest)
			tplAll.ExecuteTemplate(w, "error.html", 400)
			return
		}
		http.HandleFunc("/", SearchHandler)

	}
	if r.Method != "GET" {
		w.WriteHeader(http.StatusBadRequest)
		tplAll.ExecuteTemplate(w, "error.html", 400)
		return
	}

	GetArtistBase(w) 
	
	if r.URL.String() == "/" { // условие для главной страницы
	
		tplAll.ExecuteTemplate(w, "index.html", artists)
	} else {

		name := r.URL.String()[1:]
		groupID, err := strconv.Atoi(name)
		if err != nil && groupID > 52 && groupID < 1 { // если нечего отображать
			w.WriteHeader(http.StatusNotFound)
			tplAll.ExecuteTemplate(w, "error.html", 404)
			return
		}
		
			groupie := GetFullInfoForArtist(w, groupID)
			tplAll.ExecuteTemplate(w, "group.html", groupie)
		
	} 
	


}

func SearchHandler(w http.ResponseWriter, r *http.Request) {

}

// формирует структуру ArtistFull для указанного испольнителя
func GetFullInfoForArtist(w http.ResponseWriter, groupID int) ArtistFull {
	var groupie ArtistFull
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

	// запрашиваем Relations
	resp, err := http.Get(relatslink)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		tplAll.ExecuteTemplate(w, "error.html", 500)
		return groupie
	}
	
	data, err2 := ioutil.ReadAll(resp.Body)
	if err2 != nil {
		w.WriteHeader(http.StatusInternalServerError)
		tplAll.ExecuteTemplate(w, "error.html", 500)
		return groupie
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
		loc.LocId = index
		relation := strings.Split(location, ":[")
		loc.Location = relation[0][1 : len(relation[0])-1]
		loc.Dates = strings.Split(relation[1], ",")
		for i, val := range loc.Dates {
			loc.Dates[i] = val[1 : len(val)-1]
		}
		groupie.ConcertDates = append(groupie.ConcertDates, loc)
	}

	return groupie
}

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
	
	return 
}
