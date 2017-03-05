package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/context"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

type Adapter func(http.Handler) http.Handler

func Adapt(h http.Handler, adapters ...Adapter) http.Handler {
	for _, adapter := range adapters {
		h = adapter(h)
	}
	return h
}

func withDB(db *mgo.Collection) Adapter {
	// return the Adapter
	return func(h http.Handler) http.Handler {
		// the adapter (when called) should return a new handler
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// copy the database session
			dbsession := db
			// save it in the mux context
			context.Set(r, "database", dbsession)
			// pass execution to the original handler
			h.ServeHTTP(w, r)
		})
	}
}

func main() {

	// connect to the database
	mongoUrl := os.Getenv("SCALINGO_MONGO_URL")
	fmt.Println(mongoUrl)
	if mongoUrl == "" {
		mongoUrl = "localhost/harambot"
	}
	session, err := mgo.Dial(mongoUrl)
	db := session.DB("").C("report_info")
	fmt.Println("DB initialized")

	if err != nil {
		log.Fatal("cannot dial mongo", err)
	}

	defer session.Close() // clean up when we’re done
	// Adapt our handle function using withDB

	pageId := os.Getenv("PAGE_ID")
	if pageId == "" {
		pageId = "informer.upc"
	}

	pageAccessToken := os.Getenv("ACCESS_TOKEN")
	if pageAccessToken == "" {
		pageAccessToken = "EAACEdEose0cBANWw65HIuRi7Hvrj4TauhqVGo9HxmvnSysDXZCYZAZAiEis5y8eOovHX5grWxNRj5A16XJDmU3HOw6eTN8D6mvwRN9NRnHQ0t8ghuN5sj5Q4Vg2BBoVZBtHyxeech9sIa7Ms03GyfNcyiaN92mNXvLiZCAZAxcd2yc6CyjvC4HLcFzEOVT7YYZD"
	}

	go startInputService(db, pageId, pageAccessToken)
	go startValidatorService(db)

	fmt.Println("Services initialized")

	// add the handler
	http.Handle("/potentiallist", Adapt(http.HandlerFunc(potentialListHandler), withDB(db), context.ClearHandler))
	http.Handle("/report", Adapt(http.HandlerFunc(reportHandler), withDB(db), context.ClearHandler))
	http.Handle("/", http.FileServer(http.Dir("./static/")))
	// start the server
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatal(err)
	}
}

type listRequest struct {
	Offset int `json="offset"`
}

func potentialListHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "POST":
		handleServePotentialList(w, r)
	default:
		http.Error(w, "Not supported", http.StatusMethodNotAllowed)
	}
}

func handleServePotentialList(w http.ResponseWriter, r *http.Request) {

	db := context.Get(r, "database").(*mgo.Collection)

	var req listRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {

		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	var ReportBatch []Report

	fmt.Print("Offset: ")
	fmt.Println(req.Offset)

	if err := db.Find(bson.M{"status": ReportStateAISelected}).Sort("+timestamp").Skip(req.Offset).Limit(25).All(&ReportBatch); err != nil {
		panic(err)
	}

	if err := json.NewEncoder(w).Encode(ReportBatch); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

}

func reportHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "POST":
		handleReportChange(w, r)
	default:
		http.Error(w, "Not supported", http.StatusMethodNotAllowed)
	}
}

func handleReportChange(w http.ResponseWriter, r *http.Request) {

	db := context.Get(r, "database").(*mgo.Collection)

	var change ReportModification

	if err := json.NewDecoder(r.Body).Decode(&change); err != nil {

		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	var rep Report
	if err := db.Find(bson.M{"facebookid": change.FacebookId}).One(&rep); err != nil {

		http.Error(w, err.Error(), http.StatusBadRequest)
		return

	}

	var newStatus int

	switch change.StatusName {
	case "confirmed":
		newStatus = ReportStateUserConfirmed
	case "discarted":
		newStatus = ReportStateUserDiscarted
	default:
		newStatus = ReportStateUnchecked

	}

	rep.Status = newStatus

	rep.update(db)
}

/*
func handleInsert(w http.ResponseWriter, r *http.Request) {
	db := context.Get(r, "database").(*mgo.Collection)
	// decode the request body

	var c comment

	log.Println("On insert")

	if err := json.NewDecoder(r.Body).Decode(&c); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	log.Println("Parsed")
	// give the comment a unique ID and set the time
	c.ID = bson.NewObjectId()
	c.When = time.Now()
	// insert it into the database
	if err := db.Insert(&c); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	log.Println("Added")

	// redirect to it
	http.Redirect(w, r, "/comments/"+c.ID.Hex(), http.StatusTemporaryRedirect)
}

func handleRead(w http.ResponseWriter, r *http.Request) {
	db := context.Get(r, "database").(*mgo.Collection)
	// load the comments
	var comments []*comment
	if err := db.
		Find(nil).Sort("-when").Limit(100).All(&comments); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	// write it out
	if err := json.NewEncoder(w).Encode(comments); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}*/
