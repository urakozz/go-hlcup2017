package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"./entities"
	"./storage"
	"archive/zip"
	"bytes"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"io/ioutil"
	"os"
	"strings"
	"time"
)

var DataContainer = storage.NewStorage(storage.Opts{})

func main() {
	var validate = flag.Bool("validate", false, "build validation")
	var port = flag.Int("port", 3000, "app port")
	flag.Parse()
	if *validate {
		fmt.Println("ok")
		os.Exit(0)
	}

	Unzip()
	r := chi.NewRouter()
	r.Use(middleware.Timeout(2 * time.Minute))
	r.Route("/users", func(r chi.Router) {
		r.Post("/new", newUser)
		r.Post("/{id}", updateUser)
		r.Get("/{id}", getUser)
		r.Get("/{id}/visits", getUserVisits)
	})
	r.Route("/locations", func(r chi.Router) {
		r.Post("/new", newLocation)
		r.Post("/{id}", updateLocation)
		r.Get("/{id}", getLocation)
		r.Get("/{id}/avg", getLocationAvg)
	})
	r.Route("/visits", func(r chi.Router) {
		r.Get("/", listVisit)
		r.Post("/new", newVisit)
		r.Post("/{id}", updateVisit)
		r.Get("/{id}", getVisit)
	})
	log.Println("start", *port)
	if err := http.ListenAndServe(fmt.Sprintf(":%d", *port), r); err != nil {
		log.Fatal(err)
	}
}

// Done
func getUser(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	log.Println("getUser", "id", id)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	u, err := DataContainer.GetUser(int64(id))
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	b, err := json.Marshal(u)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Write(b)
}

// Done
func getLocation(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	log.Println("getUser", "id", id)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	u, err := DataContainer.GetLocation(int64(id))
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	b, err := json.Marshal(u)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Write(b)
}

func listVisit(w http.ResponseWriter, r *http.Request) {
	log.Println("list visit")
	json.NewEncoder(w).Encode(DataContainer.ListVisits())
}

// Done
func getVisit(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	log.Println("getVisit", "id", id)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	u, err := DataContainer.GetVisit(int64(id))
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	b, err := json.Marshal(u)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Write(b)
}

// UPDATE
// Done
func updateUser(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	log.Println("updateUser", "id", id)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	b, _ := ioutil.ReadAll(r.Body)
	var tmp map[string]interface{}
	json.NewDecoder(bytes.NewReader(b)).Decode(&tmp)
	for _, v := range tmp {
		if v == nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
	}
	u := &entities.User{}
	err = json.NewDecoder(bytes.NewReader(b)).Decode(u)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	//if u.ID >= 0 && u.ID != int64(id) {
	//	w.WriteHeader(http.StatusBadRequest)
	//	return
	//}
	u.ID = int64(id)

	err = DataContainer.UpdateUser(u)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	w.Write([]byte("{}"))
}

// Done
func updateLocation(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	log.Println("updateLocation", "id", id)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	b, _ := ioutil.ReadAll(r.Body)
	var tmp map[string]interface{}
	json.NewDecoder(bytes.NewReader(b)).Decode(&tmp)
	for _, v := range tmp {
		if v == nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
	}
	u := &entities.Location{}
	err = json.NewDecoder(bytes.NewReader(b)).Decode(u)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	//if u.ID >= 0 && u.ID != int64(id) {
	//	w.WriteHeader(http.StatusBadRequest)
	//	return
	//}
	u.ID = int64(id)

	err = DataContainer.UpdateLocation(u)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	w.Write([]byte("{}"))
}

// Done
func updateVisit(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	b, _ := ioutil.ReadAll(r.Body)
	var tmp map[string]interface{}
	json.NewDecoder(bytes.NewReader(b)).Decode(&tmp)
	for _, v := range tmp {
		if v == nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
	}
	u := &entities.Visit{}
	err = json.NewDecoder(bytes.NewReader(b)).Decode(u)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	//if u.ID >= 0 && u.ID != int64(id) {
	//	w.WriteHeader(http.StatusBadRequest)
	//	return
	//}
	u.ID = int64(id)

	err = DataContainer.UpdateVisit(u)
	if err == storage.ErrNotFound {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	if err == storage.ErrBadRequest {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	w.Write([]byte("{}"))
}

// NEW
func newUser(w http.ResponseWriter, r *http.Request) {
	log.Println("newUser")
	u := &entities.User{}
	err := json.NewDecoder(r.Body).Decode(u)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	if u.ID == 0 {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	if err = u.Validate(); err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	DataContainer.NewUser(u)
	w.Write([]byte("{}"))
}
func newLocation(w http.ResponseWriter, r *http.Request) {
	log.Println("newLocation")
	u := &entities.Location{}
	err := json.NewDecoder(r.Body).Decode(u)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	if u.ID == 0 {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	if err = u.Validate(); err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	DataContainer.NewLocation(u)
	w.Write([]byte("{}"))
}
func newVisit(w http.ResponseWriter, r *http.Request) {
	log.Println("newVisit")
	u := &entities.Visit{}
	err := json.NewDecoder(r.Body).Decode(u)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	if u.ID == 0 {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	if err = u.Validate(); err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	if err := DataContainer.NewVisit(u); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	w.Write([]byte("{}"))
}

// CUSTOM
func getUserVisits(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	_, err = DataContainer.GetUser(int64(id))
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	opts := storage.GetUserVisitsOpts{}
	if fromDate := r.URL.Query().Get("fromDate"); fromDate != "" {
		if v, err := strconv.Atoi(fromDate); err == nil {
			v32 := int64(v)
			opts.FromDate = &v32
		} else {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
	}
	if toDate := r.URL.Query().Get("toDate"); toDate != "" {
		if v, err := strconv.Atoi(toDate); err == nil {
			v32 := int64(v)
			opts.ToDate = &v32
		} else {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
	}

	if toDistance := r.URL.Query().Get("toDistance"); toDistance != "" {
		if v, err := strconv.Atoi(toDistance); err == nil {
			v32 := int64(v)
			opts.ToDistance = &v32
		} else {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
	}

	if v := r.URL.Query().Get("country"); v != "" {
		opts.Country = &v
	}

	visits := DataContainer.GetUserVisitsFiltered(int64(id), opts)
	if len(visits) == 0 {
		w.Write([]byte(`{"visits":[]}`))
		return
	}
	res := &ShortVisitContainer{Visits: visits}
	b, err := json.Marshal(res)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Println(err)
		return
	}
	w.Write(b)
}
func getLocationAvg(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	_, err = DataContainer.GetLocation(int64(id))
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	opts := storage.GetLocationVisitsOpts{}
	if fromDate := r.URL.Query().Get("fromDate"); fromDate != "" {
		if v, err := strconv.Atoi(fromDate); err == nil {
			v32 := int64(v)
			opts.FromDate = &v32
		} else {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
	}
	if toDate := r.URL.Query().Get("toDate"); toDate != "" {
		if v, err := strconv.Atoi(toDate); err == nil {
			v32 := int64(v)
			opts.ToDate = &v32
		} else {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
	}
	if fromAge := r.URL.Query().Get("fromAge"); fromAge != "" {
		if v, err := strconv.Atoi(fromAge); err == nil {
			v32 := int64(v)
			opts.FromAge = &v32
		} else {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
	}
	if toAge := r.URL.Query().Get("toAge"); toAge != "" {
		if v, err := strconv.Atoi(toAge); err == nil {
			v32 := int64(v)
			opts.ToAge = &v32
		} else {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
	}
	if gen := r.URL.Query().Get("gender"); gen != "" {
		if gen == "m" || gen == "f" {
			opts.Gender = &gen
		} else {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
	}

	visits := DataContainer.GetLocationVisitsFiltered(int64(id), opts)
	if len(visits) == 0 {
		w.Write([]byte(`{"avg":0.0}`))
		return
	}
	var sum int64
	for _, v := range visits {
		sum += int64(*v.Mark)
	}
	avg := float64(sum) / float64(len(visits))

	w.Write([]byte(fmt.Sprintf(`{"avg":%.5f}`, avg)))
}

type ShortVisitContainer struct {
	Visits []*entities.ShortVisit `json:"visits"`
}
type VisitContainer struct {
	Visits []*entities.Visit `json:"visits"`
}
type UserContainer struct {
	Users []*entities.User `json:"users"`
}
type LocationContainer struct {
	Locations []*entities.Location `json:"locations"`
}

func Unzip() {
	r, err := zip.OpenReader("tmp/data/data.zip")
	if err != nil {
		log.Fatal(err)
	}

	defer r.Close()

	for _, f := range r.File {
		if f.FileInfo().IsDir() {
			continue
		}
		log.Println(f.Name)

		rc, err := f.Open()
		if err != nil {
			log.Fatal(err)
			return
		}

		dec := json.NewDecoder(rc)
		if strings.HasPrefix(f.Name, "locations") {
			locations := &LocationContainer{}
			err := dec.Decode(locations)
			if err != nil {
				log.Fatal(err)
			}
			DataContainer.LoadLocations(locations.Locations)
			log.Println("loaded locations", len(locations.Locations))
		} else if strings.HasPrefix(f.Name, "users") {
			users := &UserContainer{}
			err := dec.Decode(users)
			if err != nil {
				log.Fatal(err)
			}
			DataContainer.LoadUsers(users.Users)
			log.Println("loaded users", len(users.Users))
		} else if strings.HasPrefix(f.Name, "visits") {
			visits := &VisitContainer{}
			err := dec.Decode(visits)
			if err != nil {
				log.Fatal(err)
			}
			DataContainer.LoadVisits(visits.Visits)
			log.Println("loaded visits", len(visits.Visits))
		}
		rc.Close()
	}
	DataContainer.ProcessLoad()
}
