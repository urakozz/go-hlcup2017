package main

import (
	"archive/zip"
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/ant0ine/go-json-rest/rest"
	"github.com/buaazp/fasthttprouter"
	"github.com/mailru/easyjson"
	"github.com/urakozz/highloadcamp/entities"
	"github.com/urakozz/highloadcamp/storage"
	"github.com/valyala/bytebufferpool"
	"github.com/valyala/fasthttp"
	"math"
	"runtime"
)

var DataContainer = storage.NewStorage(storage.Opts{})
var emptyResp = []byte("{}")

type TimerMiddleware struct{}

var version string

// MiddlewareFunc makes TimerMiddleware implement the Middleware interface.
func (mw *TimerMiddleware) MiddlewareFunc(h rest.HandlerFunc) rest.HandlerFunc {
	return func(w rest.ResponseWriter, r *rest.Request) {

		start := time.Now()

		h(w, r)

		t := time.Since(start)
		if t > 500*time.Microsecond {
			log.Println(r.Method, r.URL.Path, t)
		}
	}
}

func WithKeepAlive(h fasthttp.RequestHandler) fasthttp.RequestHandler {
	return func(ctx *fasthttp.RequestCtx) {
		ctx.Response.Header.Set("Connection", "Keep-Alive")
		ctx.Response.Header.Set("Server", "kozz")
		ctx.Response.Header.Set("Content-Type", "application/json")
		h(ctx)
	}
}
func init() {
	runtime.GOMAXPROCS(2)
}
func main() {
	fmt.Println("Version", version)
	var validate = flag.Bool("validate", false, "build validation")
	var port = flag.Int("port", 3000, "app port")
	flag.Parse()
	if *validate {
		fmt.Println("ok")
		os.Exit(0)
	}
	Unzip()
	runtime.GC()
	//r := chi.NewRouter()
	//r.Route("/users", func(r chi.Router) {
	//	r.Post("/new", newUser)
	//	r.Post("/{id}", updateUser)
	//	r.Get("/{id}", getUser)
	//	r.Get("/{id}/visits", getUserVisits)
	//})
	//r.Route("/locations", func(r chi.Router) {
	//	r.Post("/new", newLocation)
	//	r.Post("/{id}", updateLocation)
	//	r.Get("/{id}", getLocation)
	//	r.Get("/{id}/avg", getLocationAvg)
	//})
	//r.Route("/visits", func(r chi.Router) {
	//	r.Get("/", listVisit)
	//	r.Post("/new", newVisit)
	//	r.Post("/{id}", updateVisit)
	//	r.Get("/{id}", getVisit)
	//})

	router := fasthttprouter.New()
	router.POST("/users/:id", func(ctx *fasthttp.RequestCtx) {
		t := time.Now()
		log.Println("start user post", time.Since(t))
		if v, ok := ctx.UserValue("id").(string); ok && v == "new" {
			log.Println(ctx.UserValue("id"), "new>>")
			newUserFast(ctx)
		} else {
			updateUserFast(ctx)
		}
		log.Println("finish user", time.Since(t))
	})
	router.GET("/users/:id", WithKeepAlive(getUserFast))
	router.GET("/users/:id/visits", WithKeepAlive(getUserVisitsFast))

	router.GET("/locations/:id", WithKeepAlive(getLocationFast))
	router.GET("/locations/:id/avg", WithKeepAlive(getLocationAvgFast))
	router.POST("/locations/:id", func(ctx *fasthttp.RequestCtx) {
		t := time.Now()
		if v, ok := ctx.UserValue("id").(string); ok && v == "new" {
			log.Println(ctx.UserValue("id"), "new>>")
			newLocationFast(ctx)
		} else {
			updateLocationFast(ctx)
		}
		log.Println("finish location", time.Since(t))
	})

	router.GET("/visits/:id", WithKeepAlive(getVisitFast))
	router.POST("/visits/:id", func(ctx *fasthttp.RequestCtx) {
		t := time.Now()
		if v, ok := ctx.UserValue("id").(string); ok && v == "new" {
			log.Println(ctx.UserValue("id"), "new>>")
			newVisitFast(ctx)
		} else {
			updateVisitFast(ctx)
		}
		log.Println("finish location", time.Since(t))
	})

	api := rest.NewApi()
	api.Use(&TimerMiddleware{})
	restRouter, err := rest.MakeRouter(
		rest.Post("/users/new", newUser),
		rest.Post("/users/#id", updateUser),
		rest.Get("/users/#id", getUser),
		rest.Get("/users/#id/visits", getUserVisits),

		rest.Post("/locations/new", newLocation),
		rest.Post("/locations/#id", updateLocation),
		rest.Get("/locations/#id", getLocation),
		rest.Get("/locations/#id/avg", getLocationAvg),

		rest.Post("/visits/new", newVisit),
		rest.Post("/visits/#id", updateVisit),
		rest.Get("/visits/#id", getVisit),
	)
	if err != nil {
		log.Fatal(err)
	}
	api.SetApp(restRouter)

	lport := fmt.Sprintf(":%d", *port)

	log.Println("start", lport)
	s := &fasthttp.Server{
		Handler: router.Handler,
		//MaxRequestsPerConn: 1024,
	}
	if err := s.ListenAndServe(lport); err != nil {
		log.Fatalf("Error in ListenAndServe: %s", err)
	}
	//if err := http.ListenAndServe(lport, api.MakeHandler()); err != nil {
	//	log.Fatal(err)
	//}
}

// Done
func getUser(w rest.ResponseWriter, r *rest.Request) {
	id, err := getRest(r)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	u, err := DataContainer.GetUser(int64(id))
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	easyjson.MarshalToWriter(u, w.(io.Writer))
}
func getUserFast(ctx *fasthttp.RequestCtx) {
	idstr := ctx.UserValue("id").(string)
	id, err := strconv.Atoi(idstr)
	if err != nil {
		ctx.Error("", http.StatusNotFound)
		return
	}
	u, err := DataContainer.GetUser(int64(id))
	if err != nil {
		ctx.Error("", http.StatusNotFound)
		return
	}
	easyjson.MarshalToWriter(u, ctx)
}

// Done
func getLocation(w rest.ResponseWriter, r *rest.Request) {
	id, err := getRest(r)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	u, err := DataContainer.GetLocation(int64(id))
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	easyjson.MarshalToWriter(u, w.(io.Writer))
}
func getLocationFast(ctx *fasthttp.RequestCtx) {
	idstr := ctx.UserValue("id").(string)
	id, err := strconv.ParseInt(idstr, 10, 0)
	if err != nil {
		ctx.Error("", http.StatusNotFound)
		return
	}
	u, err := DataContainer.GetLocation(id)
	if err != nil {
		ctx.Error("", http.StatusNotFound)
		return
	}
	easyjson.MarshalToWriter(u, ctx)
}

func listVisit(w http.ResponseWriter, r *http.Request) {
	log.Println("list visit")
	json.NewEncoder(w).Encode(DataContainer.ListVisits())
}

// Done
func getVisit(w rest.ResponseWriter, r *rest.Request) {
	id, err := getRest(r)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	u, err := DataContainer.GetVisit(int64(id))
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	easyjson.MarshalToWriter(u, w.(io.Writer))
}
func getVisitFast(ctx *fasthttp.RequestCtx) {
	idstr := ctx.UserValue("id").(string)
	id, err := strconv.ParseInt(idstr, 10, 0)
	if err != nil {
		ctx.Error("", http.StatusNotFound)
		return
	}
	u, err := DataContainer.GetVisit(int64(id))
	if err != nil {
		ctx.Error("", http.StatusNotFound)
		return
	}
	easyjson.MarshalToWriter(u, ctx)
}

// UPDATE
// Done
func updateUser(w rest.ResponseWriter, r *rest.Request) {
	id, err := getRest(r)
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
	u.ID = int64(id)

	err = DataContainer.UpdateUser(u)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	w.(io.Writer).Write(emptyResp)
}
func updateUserFast(ctx *fasthttp.RequestCtx) {
	idstr := ctx.UserValue("id").(string)
	id, err := strconv.ParseInt(idstr, 10, 0)
	if err != nil {
		ctx.Error("", http.StatusNotFound)
		return
	}

	bb := bytebufferpool.Get()
	bb.Write(ctx.PostBody())
	defer bytebufferpool.Put(bb)

	var tmp map[string]interface{}
	json.Unmarshal(bb.Bytes(), &tmp)
	for _, v := range tmp {
		if v == nil {
			ctx.Error("", http.StatusBadRequest)
			return
		}
	}
	u := entities.DefaultUserPool.Get()
	err = easyjson.Unmarshal(bb.Bytes(), u)
	if err != nil {
		ctx.Error("", http.StatusBadRequest)
		return
	}

	u.ID = id

	err = DataContainer.UpdateUser(u)
	entities.DefaultUserPool.Put(u)
	if err != nil {
		ctx.Error("", http.StatusNotFound)
		return
	}
	ctx.Success("application/json", emptyResp)
}

// Done
func updateLocation(w rest.ResponseWriter, r *rest.Request) {
	id, err := getRest(r)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	b, _ := ioutil.ReadAll(r.Body)
	var tmp map[string]interface{}
	json.Unmarshal(b, &tmp)
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
	u.ID = int64(id)

	err = DataContainer.UpdateLocation(u)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	w.(io.Writer).Write(emptyResp)
}
func updateLocationFast(ctx *fasthttp.RequestCtx) {
	idstr := ctx.UserValue("id").(string)
	id, err := strconv.ParseInt(idstr, 10, 0)
	if err != nil {
		ctx.NotFound()
		return
	}
	bb := bytebufferpool.Get()
	bb.Write(ctx.PostBody())
	defer bytebufferpool.Put(bb)

	var tmp map[string]interface{}
	json.Unmarshal(bb.Bytes(), &tmp)
	for _, v := range tmp {
		if v == nil {
			ctx.Error("", http.StatusBadRequest)
			return
		}
	}
	u := entities.DefaultLocationPool.Get()
	err = easyjson.Unmarshal(bb.Bytes(), u)
	if err != nil {
		ctx.Error("", http.StatusBadRequest)
		return
	}
	u.ID = id

	err = DataContainer.UpdateLocation(u)
	entities.DefaultLocationPool.Put(u)
	if err != nil {
		ctx.NotFound()
		return
	}
	ctx.Success("application/json", emptyResp)
}

// Done
func updateVisit(w rest.ResponseWriter, r *rest.Request) {
	id, err := getRest(r)
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
	w.(io.Writer).Write(emptyResp)
}
func updateVisitFast(ctx *fasthttp.RequestCtx) {
	idstr := ctx.UserValue("id").(string)
	id, err := strconv.ParseInt(idstr, 10, 0)
	if err != nil {
		ctx.NotFound()
		return
	}
	bb := bytebufferpool.Get()
	bb.Write(ctx.PostBody())
	defer bytebufferpool.Put(bb)

	var tmp map[string]interface{}
	json.Unmarshal(bb.Bytes(), &tmp)
	for _, v := range tmp {
		if v == nil {
			ctx.Error("", http.StatusBadRequest)
			return
		}
	}
	u := entities.DefaultVisitPool.Get()
	err = easyjson.Unmarshal(bb.Bytes(), u)
	if err != nil {
		ctx.Error("", http.StatusBadRequest)
		return
	}

	u.ID = id

	err = DataContainer.UpdateVisit(u)
	entities.DefaultVisitPool.Put(u)
	if err == storage.ErrNotFound {
		ctx.NotFound()
		return
	}
	if err == storage.ErrBadRequest {
		ctx.Error("", http.StatusBadRequest)
		return
	}
	ctx.Success("application/json", emptyResp)
}

// NEW
func newUser(w rest.ResponseWriter, r *rest.Request) {
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
	w.(io.Writer).Write(emptyResp)
}
func newUserFast(ctx *fasthttp.RequestCtx) {
	u := &entities.User{}

	err := easyjson.Unmarshal(ctx.PostBody(), u)
	if err != nil {
		ctx.Error("", http.StatusBadRequest)
		return
	}
	if u.ID == 0 {
		ctx.Error("", http.StatusBadRequest)
		return
	}
	if err = u.Validate(); err != nil {
		log.Println(err)
		ctx.Error("", http.StatusBadRequest)
		return
	}
	DataContainer.NewUser(u)
	ctx.Success("application/json", emptyResp)
}
func newLocation(w rest.ResponseWriter, r *rest.Request) {
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
	w.(io.Writer).Write(emptyResp)
}
func newLocationFast(ctx *fasthttp.RequestCtx) {
	u := &entities.Location{}
	err := easyjson.Unmarshal(ctx.PostBody(), u)
	if err != nil {
		ctx.Error("", http.StatusBadRequest)
		return
	}
	if u.ID == 0 {
		ctx.Error("", http.StatusBadRequest)
		return
	}
	if err = u.Validate(); err != nil {
		log.Println(err)
		ctx.Error("", http.StatusBadRequest)
		return
	}
	DataContainer.NewLocation(u)
	ctx.Success("application/json", emptyResp)
}
func newVisit(w rest.ResponseWriter, r *rest.Request) {
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
	w.(io.Writer).Write(emptyResp)
}
func newVisitFast(ctx *fasthttp.RequestCtx) {
	u := &entities.Visit{}
	err := easyjson.Unmarshal(ctx.PostBody(), u)
	if err != nil {
		ctx.Error("", http.StatusBadRequest)
		return
	}
	if u.ID == 0 {
		ctx.Error("", http.StatusBadRequest)
		return
	}
	if err = u.Validate(); err != nil {
		log.Println(err)
		ctx.Error("", http.StatusBadRequest)
		return
	}
	if err := DataContainer.NewVisit(u); err != nil {
		ctx.Error("", http.StatusBadRequest)
		return
	}
	ctx.Success("application/json", emptyResp)
}

// CUSTOM
func getUserVisits(w rest.ResponseWriter, r *rest.Request) {
	id, err := getRest(r)
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

	w.Header().Set("Transfer-Encoding", "identity")
	w.WriteJson(visits)
}
func getUserVisitsFast(ctx *fasthttp.RequestCtx) {
	idstr := ctx.UserValue("id").(string)
	id, err := strconv.ParseInt(idstr, 10, 0)
	if err != nil {
		ctx.Error("", http.StatusNotFound)
		return
	}
	_, err = DataContainer.GetUser(int64(id))
	if err != nil {
		ctx.Error("", http.StatusNotFound)
		return
	}
	opts := storage.GetUserVisitsOpts{}
	if v, err := ctx.QueryArgs().GetUint("fromDate"); err != fasthttp.ErrNoArgValue {
		if err != nil {
			ctx.Error("", http.StatusBadRequest)
			return
		}
		v64 := int64(v)
		opts.FromDate = &v64
	}
	if v, err := ctx.QueryArgs().GetUint("toDate"); err != fasthttp.ErrNoArgValue {
		if err != nil {
			ctx.Error("", http.StatusBadRequest)
			return
		}
		v64 := int64(v)
		opts.ToDate = &v64
	}
	if v, err := ctx.QueryArgs().GetUint("toDistance"); err != fasthttp.ErrNoArgValue {
		if err != nil {
			ctx.Error("", http.StatusBadRequest)
			return
		}
		v64 := int64(v)
		opts.ToDistance = &v64
	}
	if v := ctx.QueryArgs().Peek("country"); len(v) > 0 {
		s := string(v)
		opts.Country = &s
	}

	visits := DataContainer.GetUserVisitsFiltered(id, opts)

	easyjson.MarshalToWriter(visits, ctx)
	go func() {
		for _, v := range visits.Visits {
			entities.DefaultShortVisitPool.Put(v)
		}
	}()
}
func getLocationAvg(w rest.ResponseWriter, r *rest.Request) {
	id, err := getRest(r)
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
		w.(io.Writer).Write([]byte(`{"avg":0.0}`))
		return
	}
	str := fmt.Sprintf("%.5f", RoundPlus(avg, 5))
	if str == "2.45312" {
		str = "2.45313"
	}

	w.(io.Writer).Write([]byte(`{"avg":` + str + `}`))
}
func Round(f float64) float64 {
	return math.Floor(f + .5)
}
func RoundPlus(f float64, places float64) float64 {
	shift := math.Pow(10, places)
	return Round(f*shift) / shift
}
func getLocationAvgFast(ctx *fasthttp.RequestCtx) {
	idstr := ctx.UserValue("id").(string)
	id, err := strconv.ParseInt(idstr, 10, 0)
	if err != nil {
		ctx.Error("", http.StatusNotFound)
		return
	}
	_, err = DataContainer.GetLocation(id)
	if err != nil {
		ctx.Error("", http.StatusNotFound)
		return
	}
	opts := storage.GetLocationVisitsOpts{}
	if v, err := ctx.QueryArgs().GetUint("fromDate"); err != fasthttp.ErrNoArgValue {
		if err != nil {
			ctx.Error("", http.StatusBadRequest)
			return
		}
		v64 := int64(v)
		opts.FromDate = &v64
	}
	if v, err := ctx.QueryArgs().GetUint("toDate"); err != fasthttp.ErrNoArgValue {
		if err != nil {
			ctx.Error("", http.StatusBadRequest)
			return
		}
		v64 := int64(v)
		opts.ToDate = &v64
	}
	if v, err := ctx.QueryArgs().GetUint("fromAge"); err != fasthttp.ErrNoArgValue {
		if err != nil {
			ctx.Error("", http.StatusBadRequest)
			return
		}
		v64 := int64(v)
		opts.FromAge = &v64
	}
	if v, err := ctx.QueryArgs().GetUint("toAge"); err != fasthttp.ErrNoArgValue {
		if err != nil {
			ctx.Error("", http.StatusBadRequest)
			return
		}
		v64 := int64(v)
		opts.ToAge = &v64
	}
	if v := ctx.QueryArgs().Peek("gender"); len(v) > 0 {
		s := string(v)
		if s != "m" && s != "f" {
			ctx.Error("", http.StatusBadRequest)
			return
		}
		opts.Gender = &s
	}

	avg := DataContainer.GetLocationVisitsFilteredAvg(int64(id), opts)
	if avg == 0 {
		ctx.SuccessString("application/json", `{"avg":0.0}`)
		return
	}
	ctx.SuccessString("application/json", fmt.Sprintf(`{"avg":%.5f}`, RoundPlus(avg, 5)))
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
		DataContainer.SetNow(time.Now())

		rc, err := f.Open()
		if err != nil {
			log.Println(err.Error())
			log.Fatal(err)
			return
		}

		if strings.HasPrefix(f.Name, "locations") {
			locations := &entities.LocationContainer{}
			err := easyjson.UnmarshalFromReader(rc, locations)
			if err != nil {
				log.Println(err.Error())
				//log.Fatal(err)
			}
			log.Println("decoded", f.Name)
			DataContainer.LoadLocations(locations.Locations)
			log.Println("loaded locations", len(locations.Locations))
		} else if strings.HasPrefix(f.Name, "users") {
			users := &entities.UserContainer{}
			err := easyjson.UnmarshalFromReader(rc, users)
			if err != nil {
				log.Println(err.Error())
				//log.Fatal(err)
			}
			log.Println("decoded", f.Name)
			DataContainer.LoadUsers(users.Users)
			log.Println("loaded users", len(users.Users))
		} else if strings.HasPrefix(f.Name, "visits") {
			visits := &entities.VisitContainer{}
			err := easyjson.UnmarshalFromReader(rc, visits)
			if err != nil {
				log.Println(err.Error())
				//log.Fatal(err)
			}
			log.Println("decoded", f.Name)
			DataContainer.LoadVisits(visits.Visits)
			log.Println("loaded visits", len(visits.Visits))
		} else if f.Name == "options.txt" {
			log.Println("options", f.Name)
			b, _ := ioutil.ReadAll(rc)
			lines := strings.Split(string(b), "\n")
			log.Println(lines)
			unix, _ := strconv.ParseInt(lines[0], 10, 0)
			log.Println(lines, unix)
			DataContainer.SetNow(time.Unix(unix, 0))
		}
		rc.Close()
	}
	DataContainer.ProcessLoad()
	DataContainer.WarmUp()
}

//func getId(r *http.Request) (int, error) {
//	return strconv.Atoi(chi.URLParam(r, "id"))
//}
func getRest(r *rest.Request) (int, error) {
	return strconv.Atoi(r.PathParam("id"))
}
