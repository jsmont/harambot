package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

type ISO8601Time string

type FacebookPost struct {
	Message      string      `json:"message"`
	Created_time ISO8601Time `json:"created_time"`
	Id           string      `json:"id"`
}

type FacebookPagination struct {
	Prev string `json:"previous"`
	Next string `json:"next"`
}

type FacebookPostListResponse struct {
	Posts      []FacebookPost     `json:"data"`
	Pagination FacebookPagination `json:"paging"`
}

func getFacebookPosts(pageId string, pageAccessToken string, postsPage string) ([]FacebookPost, FacebookPagination) {

	url := "https://graph.facebook.com/v2.8/" + pageId + "/posts?access_token=" + pageAccessToken

	if postsPage != "" {
		url += "&page=" + postsPage
	}

	resp, err := http.Get(url)

	if err != nil {
		log.Fatal("Could not fetch group posts", err)
	}

	fmt.Println(resp)

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)

	fmt.Println(string(body))

	FbPostList := FacebookPostListResponse{}

	json.Unmarshal([]byte(string(body)), &FbPostList)

	return FbPostList.Posts, FbPostList.Pagination

}

func (p *FacebookPost) getReport() Report {
	t, _ := time.Parse("2006-01-02T15:04:05.999999999-0700", string(p.Created_time))
	return Report{Message: p.Message, Timestamp: t, FacebookId: p.Id}

}
