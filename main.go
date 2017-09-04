package main

import (
	"archive/zip"
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

	"runtime"

	"bytes"

	"github.com/ant0ine/go-json-rest/rest"
	"github.com/buaazp/fasthttprouter"
	"github.com/mailru/easyjson"
	"github.com/urakozz/highloadcamp/entities"
	"github.com/urakozz/highloadcamp/storage"
	"github.com/valyala/bytebufferpool"
	"github.com/valyala/fasthttp"
)

var DataContainer = storage.NewStorage(storage.Opts{})
var emptyResp = []byte("{}")
var emptyAvg = []byte(`{"avg":0.0}`)

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

	lport := fmt.Sprintf(":%d", *port)

	log.Println("start", lport)
	s := &fasthttp.Server{
		Handler: router.Handler,
		//MaxRequestsPerConn: 1024,
	}
	if err := s.ListenAndServe(lport); err != nil {
		log.Fatalf("Error in ListenAndServe: %s", err)
	}
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
		ctx.Success("application/json", emptyAvg)
		return
	}
	buf := bytes.NewBufferString(`{"avg":`)
	buf.WriteString(strconv.FormatFloat(avg, 'f', 5, 64))
	buf.WriteRune('}')
	ctx.Success("application/json", buf.Bytes())
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
}

