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

type FacebookComment struct {
	Message      string      `json:"message"`
	Created_time ISO8601Time `json:"created_time"`
	Id           string      `json:"id"`
}

type FacebookPost struct {
	Message      string      `json:"message"`
	Created_time ISO8601Time `json:"created_time"`
	Id           string      `json:"id"`
}

type FacebookPagination struct {
	Prev string `json:"previous"`
	Next string `json:"next"`
}

type FacebookCommentListResponse struct {
	Comments   []FacebookComment  `json:"data"`
	Pagination FacebookPagination `json:"paging"`
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

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)

	FbPostList := FacebookPostListResponse{}

	json.Unmarshal([]byte(string(body)), &FbPostList)

	return FbPostList.Posts, FbPostList.Pagination

}

func (p *FacebookPost) getReport() Report {

	t, _ := time.Parse("2006-01-02T15:04:05.999999999-0700", string(p.Created_time))
	return Report{Message: p.Message, Timestamp: t, FacebookId: p.Id, Status: ReportStateUnchecked}

}

func (p *FacebookPost) getComments(pageAccessToken string) ([]FacebookComment, FacebookPagination) {

	fmt.Println("Getting comments for post: " + p.Id)

	url := "https://graph.facebook.com/v2.8/" + p.Id + "/comments?access_token=" + pageAccessToken

	resp, err := http.Get(url)

	if err != nil {
		log.Fatal("Could not fetch post comments", err)
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)

	FbCommentList := FacebookCommentListResponse{}

	json.Unmarshal([]byte(string(body)), &FbCommentList)

	return FbCommentList.Comments, FbCommentList.Pagination

}

func (c *FacebookComment) getReport() Report {

	t, _ := time.Parse("2006-01-02T15:04:05.999999999-0700", string(c.Created_time))
	return Report{Message: c.Message, Timestamp: t, FacebookId: c.Id, Status: ReportStateUnchecked}

}
