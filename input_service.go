package main

import (
	"fmt"
	"time"

	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

type User struct {
	Name string `json:"username"`
	Id   string `json:"id"`
}
type Report struct {
	Message    string    `json:"message"`
	Timestamp  time.Time `json:"timestamp"`
	FacebookId string    `json:"facebookid"`
}

func startInputService(db *mgo.Collection, pageId string, pageAccessToken string) {

	FacebookPosts, _ := getFacebookPosts(pageId, pageAccessToken, "")

	for _, post := range FacebookPosts {

		rep := post.getReport()

		rep.save(db)
		/*
			FacebookComments, FacebookCommentsPagination := post.getComments(pageAccessToken)

			for _, comment := range FacebookComments {
				rep = comment.getReport()
				rep.save(db)
			}

			for ; FacebookCommentsPagination.hasNext(); FacebookComments, FacebookCommentsPagination = post.getComments(pageAccessToken, FacebookCommentsPagination.next()) {

				for _, comment := range FacebookComments {
					rep = comment.getReport(db)
					rep.save(db)
				}

			}*/
	}
}

func (p *Report) save(db *mgo.Collection) {

	if _, err := db.Upsert(bson.M{"facebookid": p.FacebookId}, p); err != nil {
		panic(err)
	}
	fmt.Println("Added new report to database")
}
