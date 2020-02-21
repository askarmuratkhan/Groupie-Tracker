package main

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"text/template"
)

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

	if r.Method == "POST" {
		if r.URL.String() != "/" {
			w.WriteHeader(http.StatusBadRequest)
			tplAll.ExecuteTemplate(w, "error.html", 400)
			return
		}

		SearchHandler(w, r)
		return
	}
	if r.Method != "GET" {
		w.WriteHeader(http.StatusMethodNotAllowed)
		tplAll.ExecuteTemplate(w, "error.html", 405)
		return
	}

	GetArtistBase(w)
	// shablon := r.FormValue("toFind")
	// fmt.Println(shablon)
	//  TYPE := r.FormValue("searchType")
	//  fmt.Println(TYPE)

	if r.URL.String() == "/" { // условие для главной страницы
		err := tplAll.ExecuteTemplate(w, "index.html", DB)
		if err != nil {
			fmt.Println(err)
		}
	} else {

		name := r.URL.String()[1:]
		groupID, err := strconv.Atoi(name)
		if err != nil || groupID > 52 || groupID < 1 { // если нечего отображать
			w.WriteHeader(http.StatusNotFound)
			err2 := tplAll.ExecuteTemplate(w, "error.html", 404)
			if err2 != nil {
				fmt.Println(err2)
			}
			return
		}

		groupie := GetFullInfoForArtist(w, groupID)
		tplAll.ExecuteTemplate(w, "group.html", groupie)

	}

}

// obrabotchik poiska
func SearchHandler(w http.ResponseWriter, r *http.Request) {
	searchType := r.FormValue("searchType")
	searchString := r.FormValue("toFind")
	if strings.Contains(searchString, " // ") {
		searchType = searchString[strings.Index(searchString, " // ")+4:]
		searchString = searchString[:strings.Index(searchString, " // ")]
	}
	fmt.Println(searchString, searchType)
	var artistsFound []ArtistFull
	// fmt.Println(artistsFound)
	// fmt.Println(DB)
	for _, copy := range DB {
		switch searchType {
		case "artist", "Artist":
			if strings.Contains(strings.ToLower(copy.Name), strings.ToLower(searchString)) {
				artistsFound = append(artistsFound, copy)
				continue
			}
		case "member", "Members":
			for _, member := range copy.Members {
				if strings.Contains(strings.ToLower(member), strings.ToLower(searchString)) {
					artistsFound = append(artistsFound, copy)
					break
				}
			}
			continue
		case "creationDate", "Creation Date":
			if strconv.Itoa(int(copy.CreationDate)) == searchString {
				artistsFound = append(artistsFound, copy)
			}
			continue
		case "firstAlbum", "First Album":
			if copy.FirstAlbum == searchString {
				artistsFound = append(artistsFound, copy)
			}
			continue
		case "location", "Location":
			// fmt.Println(searchString)
			for _, place := range copy.Locations {
				// fmt.Println(place)
				if strings.Contains(strings.ToLower(place), strings.ToLower(searchString)) {
					artistsFound = append(artistsFound, copy)
					break
				}
			}
			continue
		}
	}
	// fmt.Println(artistsFound)
	err := tplAll.ExecuteTemplate(w, "found.html", artistsFound)
	if err != nil {
		fmt.Println(err)
	}
	return
}
