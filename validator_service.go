package main

import (
	"fmt"
	"time"

	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

func startValidatorService(db *mgo.Collection) {

	for true {
		ReportBatch := getPendingReportsBatch(db)
		for len(ReportBatch) > 0 {

			for _, rep := range ReportBatch {

				language := getLanguage(rep.Message)
				sentment := getSentiment(rep.Message, language)

				fmt.Println("Report " + rep.FacebookId + " -> " + language + " " + sentment)

				if sentment != "negative" {
					rep.Status = ReportStateAIDiscarted

				} else {
					rep.Status = ReportStateAISelected
				}

				rep.save(db)
			}

			ReportBatch = getPendingReportsBatch(db)
		}
		time.Sleep(10 * time.Second)
	}

}

func getPendingReportsBatch(db *mgo.Collection) []Report {

	var ReportBatch []Report

	if err := db.Find(bson.M{"status": ReportStateUnchecked}).Sort("-timestamp").Batch(24).All(&ReportBatch); err != nil {
		panic(err)
	}

	return ReportBatch

}
