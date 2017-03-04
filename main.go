package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"time"

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

func withDB(db *mgo.Session) Adapter {
	// return the Adapter
	return func(h http.Handler) http.Handler {
		// the adapter (when called) should return a new handler
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// copy the database session
			dbsession := db.Copy()
			defer dbsession.Close() // clean up
			// save it in the mux context
			context.Set(r, "database", dbsession)
			// pass execution to the original handler
			h.ServeHTTP(w, r)
		})
	}
}

func handle(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		handleRead(w, r)
	case "POST":
		handleInsert(w, r)
	default:
		http.Error(w, "Not supported", http.StatusMethodNotAllowed)
	}
}

func main() {
	// connect to the database
	mongoUrl := os.Getenv("MONGO_URL")
	if mongoUrl == "" {
		mongoUrl = "localhost"
	}
	db, err := mgo.Dial(mongoUrl)
	if err != nil {
		log.Fatal("cannot dial mongo", err)
	}

	defer db.Close() // clean up when we’re done
	// Adapt our handle function using withDB

	pageId := "informer.upc"
	pageAccessToken := "EAACEdEose0cBAK74yl9gHSB4Yfxz2rJXtvD49TBZCc1qTeKkGKC1jWOeZAygPDQlzGcDwtPcOambChqJkJRLiviFtM3eZCEFPV5NzYWZBa7ESQ4GXFZAHFpvjZBzYFghbUOwX5J7SoohBQdH7jtiEcihUclUZAb4C8G7qAMzvYpbmakQTUt9ZBTvUWPX4lIdM8w8TA59eMGIvQZDZD"

	go startInputService(pageId, pageAccessToken)

	h := Adapt(http.HandlerFunc(handle), withDB(db))
	// add the handler
	http.Handle("/comments", context.ClearHandler(h))
	// start the server
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatal(err)
	}
}

type comment struct {
	ID     bson.ObjectId `json:"id" bson:"_id"`
	Author string        `json:"author" bson:"author"`
	Text   string        `json:"text" bson:"text"`
	When   time.Time     `json:"when" bson:"when"`
}

func handleInsert(w http.ResponseWriter, r *http.Request) {
	db := context.Get(r, "database").(*mgo.Session)
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
	if err := db.DB("harambot").C("comments").Insert(&c); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	log.Println("Added")

	// redirect to it
	http.Redirect(w, r, "/comments/"+c.ID.Hex(), http.StatusTemporaryRedirect)
}

func handleRead(w http.ResponseWriter, r *http.Request) {
	db := context.Get(r, "database").(*mgo.Session)
	// load the comments
	var comments []*comment
	if err := db.DB("harambot").C("comments").
		Find(nil).Sort("-when").Limit(100).All(&comments); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	// write it out
	if err := json.NewEncoder(w).Encode(comments); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}