package main

import (
	"fmt"
	"time"

	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

type User struct {
	Name string `json:"name"`
	Id   string `json:"id"`
}

const (
	ReportStateUnchecked     = iota
	ReportStateAIDiscarted   = iota
	ReportStateAISelected    = iota
	ReportStateUserDiscarted = iota
	ReportStateUserConfirmed = iota
)

type Report struct {
	Owner      User      `json:"owner"`
	Message    string    `json:"message"`
	Timestamp  time.Time `json:"timestamp"`
	FacebookId string    `json:"facebookid"`
	Status     int       `json:"status"`
}

type ReportModification struct {
	FacebookId string `json:"id"`
	StatusName string `json:"status_name"`
}

func startInputService(db *mgo.Collection, pageId string, pageAccessToken string) {

	FacebookPosts, _ := getFacebookPosts(pageId, pageAccessToken, "")

	for _, post := range FacebookPosts {

		rep := post.getReport()

		rep.save(db)

		FacebookComments, _ := post.getComments(pageAccessToken)

		for _, comment := range FacebookComments {
			rep = comment.getReport()
			rep.save(db)
		}
		/*
			for ; FacebookCommentsPagination.hasNext(); FacebookComments, FacebookCommentsPagination = post.getComments(pageAccessToken, FacebookCommentsPagination.next()) {

				for _, comment := range FacebookComments {
					rep = comment.getReport(db)
					rep.save(db)
				}

			}*/
	}
}

func (p *Report) save(db *mgo.Collection) {

	exists := Report{}

	if err := db.Find(bson.M{"facebookid": p.FacebookId}).One(&exists); err != nil {
		fmt.Println(err)
		if err := db.Insert(&p); err != nil {
			panic(err)
		}

		fmt.Println("Added new report to database")
	} else {
		fmt.Println("Report repeated")
	}

}

func (p *Report) update(db *mgo.Collection) {

	if _, err := db.Upsert(bson.M{"facebookid": p.FacebookId}, p); err != nil {
		panic(err)
	}
	fmt.Println("Added new report to database: " + p.FacebookId)
}
