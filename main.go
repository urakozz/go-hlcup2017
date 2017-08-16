package main

import (
	"archive/zip"
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"./entities"
	"./storage"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/mailru/easyjson"
)

var DataContainer = storage.NewStorage(storage.Opts{})
var emptyResp = []byte("{}")

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
	b, err := u.MarshalJSON()
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

	b, err := u.MarshalJSON()
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
	b, err := u.MarshalJSON()
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
	err = easyjson.Unmarshal(b, u)
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
	w.Write(emptyResp)
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
	err = easyjson.Unmarshal(b, u)
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
	w.Write(emptyResp)
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
	err = easyjson.Unmarshal(b, u)
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
	w.Write(emptyResp)
}

// NEW
func newUser(w http.ResponseWriter, r *http.Request) {
	log.Println("newUser")
	u := &entities.User{}

	err := easyjson.UnmarshalFromReader(r.Body, u)
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
	err := easyjson.UnmarshalFromReader(r.Body, u)
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
	w.Write(emptyResp)
}
func newVisit(w http.ResponseWriter, r *http.Request) {
	log.Println("newVisit")
	u := &entities.Visit{}
	err := easyjson.UnmarshalFromReader(r.Body, u)
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
	w.Write(emptyResp)
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

	//if len(visits.Visits) == 0 {
	//	w.Write([]byte(`{"visits":[]}`))
	//	return
	//}
	b, err := easyjson.Marshal(visits)
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

	avg := DataContainer.GetLocationVisitsFilteredAvg(int64(id), opts)
	if avg == 0 {
		w.Write([]byte(`{"avg":0.0}`))
		return
	}

	w.Write([]byte(fmt.Sprintf(`{"avg":%.5f}`, avg)))
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

		if strings.HasPrefix(f.Name, "locations") {
			locations := &entities.LocationContainer{}
			err := easyjson.UnmarshalFromReader(rc, locations)
			if err != nil {
				log.Fatal(err)
			}
			DataContainer.LoadLocations(locations.Locations)
			log.Println("loaded locations", len(locations.Locations))
		} else if strings.HasPrefix(f.Name, "users") {
			users := &entities.UserContainer{}
			err := easyjson.UnmarshalFromReader(rc, users)
			if err != nil {
				log.Fatal(err)
			}
			DataContainer.LoadUsers(users.Users)
			log.Println("loaded users", len(users.Users))
		} else if strings.HasPrefix(f.Name, "visits") {
			visits := &entities.VisitContainer{}
			err := easyjson.UnmarshalFromReader(rc, visits)
			if err != nil {
				log.Fatal(err)
			}
			DataContainer.LoadVisits(visits.Visits)
			log.Println("loaded visits", len(visits.Visits))
		}
		rc.Close()
	}
	DataContainer.ProcessLoad()
	DataContainer.WarmUp()
}
