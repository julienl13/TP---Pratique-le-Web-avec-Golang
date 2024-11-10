package main

import (
	"html/template"
	"net/http"
	"strconv"
	"sync"
	"regexp"
)

var (
	promoData = struct {
		Name      string
		Filiere   string
		Level     string
		Students  []struct {
			FirstName string
			LastName  string
			Age       int
			Gender    string
		}
	}{
		Name:    "B1 Informatique",
		Filiere: "Informatique",
		Level:   "Bachelor 1",
		Students: []struct {
			FirstName string
			LastName  string
			Age       int
			Gender    string
		}{
			{"Jean", "Dupont", 19, "masculin"},
			{"Marie", "Durand", 20, "féminin"},
			{"Paul", "Martin", 18, "masculin"},
		},
	}
	viewCount int
	viewMutex sync.Mutex
)

func main() {
	http.HandleFunc("/promo", promoHandler)
	http.HandleFunc("/change", changeHandler)
	http.HandleFunc("/user/form", userFormHandler)
	http.HandleFunc("/user/treatment", userTreatmentHandler)
	http.HandleFunc("/user/display", userDisplayHandler)

	http.ListenAndServe(":8080", nil)
}

func promoHandler(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("templates/promo.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	tmpl.Execute(w, promoData)
}

func changeHandler(w http.ResponseWriter, r *http.Request) {
	viewMutex.Lock()
	viewCount++
	currentCount := viewCount
	viewMutex.Unlock()

	msg := ""
	if currentCount%2 == 0 {
		msg = "Le nombre de vues est pair : " + strconv.Itoa(currentCount)
	} else {
		msg = "Le nombre de vues est impair : " + strconv.Itoa(currentCount)
	}
	tmpl, err := template.New("change").Parse("<html><body><h1>{{.}}</h1></body></html>")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	tmpl.Execute(w, msg)
}

func userFormHandler(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("templates/user_form.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	tmpl.Execute(w, nil)
}

func userTreatmentHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Redirect(w, r, "/user/form", http.StatusSeeOther)
		return
	}

	firstName := r.FormValue("firstName")
	lastName := r.FormValue("lastName")
	birthDate := r.FormValue("birthDate")
	gender := r.FormValue("gender")

	nameRegex := regexp.MustCompile(`^[A-Za-z]{1,32}$`)
	if !nameRegex.MatchString(firstName) || !nameRegex.MatchString(lastName) || (gender != "masculin" && gender != "féminin" && gender != "autre") {
		http.Redirect(w, r, "/user/form?error=invalid", http.StatusSeeOther)
		return
	}

	http.Redirect(w, r, "/user/display?firstName="+firstName+"&lastName="+lastName+"&birthDate="+birthDate+"&gender="+gender, http.StatusSeeOther)
}

func userDisplayHandler(w http.ResponseWriter, r *http.Request) {
	firstName := r.URL.Query().Get("firstName")
	lastName := r.URL.Query().Get("lastName")
	birthDate := r.URL.Query().Get("birthDate")
	gender := r.URL.Query().Get("gender")

	if firstName == "" || lastName == "" || birthDate == "" || gender == "" {
		http.Error(w, "Veuillez renseigner toutes les informations personnelles.", http.StatusBadRequest)
		return
	}

	tmpl, err := template.New("display").Parse("<html><body><h1>Informations de l'utilisateur</h1><p>Nom: {{.LastName}}</p><p>Prénom: {{.FirstName}}</p><p>Date de naissance: {{.BirthDate}}</p><p>Sexe: {{.Gender}}</p></body></html>")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	tmpl.Execute(w, struct {
		FirstName string
		LastName  string
		BirthDate string
		Gender    string
	}{
		FirstName: firstName,
		LastName:  lastName,
		BirthDate: birthDate,
		Gender:    gender,
	})
}
