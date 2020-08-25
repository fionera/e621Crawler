package api

import (
	"encoding/json"
	"github.com/go-resty/resty"
	"io/ioutil"
	"strconv"
)

const (
	listUrl  = "https://e621.net/post/index.json"
	maxLimit = 320
)

type Posts []Post

type Post struct {
	Artist    []string `json:"artist"`
	Author    string   `json:"author"`
	Change    int      `json:"change"`
	Children  string   `json:"children"`
	CreatedAt struct {
		JSONClass string `json:"json_class"`
		N         int    `json:"n"`
		S         int    `json:"s"`
	} `json:"created_at"`
	CreatorID     int         `json:"creator_id"`
	Description   string      `json:"description"`
	FavCount      int         `json:"fav_count"`
	FileExt       string      `json:"file_ext"`
	FileSize      int         `json:"file_size"`
	FileURL       string      `json:"file_url"`
	HasChildren   bool        `json:"has_children"`
	HasComments   bool        `json:"has_comments"`
	HasNotes      bool        `json:"has_notes"`
	Height        int         `json:"height"`
	ID            int         `json:"id"`
	LockedTags    interface{} `json:"locked_tags"`
	Md5           string      `json:"md5"`
	ParentID      int         `json:"parent_id"`
	PreviewHeight int         `json:"preview_height"`
	PreviewURL    string      `json:"preview_url"`
	PreviewWidth  int         `json:"preview_width"`
	Rating        string      `json:"rating"`
	SampleHeight  int         `json:"sample_height"`
	SampleURL     string      `json:"sample_url"`
	SampleWidth   int         `json:"sample_width"`
	Score         int         `json:"score"`
	Source        string      `json:"source"`
	Sources       []string    `json:"sources"`
	Status        string      `json:"status"`
	Tags          string      `json:"tags"`
	Width         int         `json:"width"`
}

func List(limit int, beforeId int, page int, tags string, typedTags bool) (Posts, error) {
	queryParams := map[string]string{}
	if limit != 0 {
		if limit > maxLimit {
			limit = maxLimit
		}
		queryParams["limit"] = strconv.Itoa(limit)
	}
	if beforeId != 0 {
		queryParams["before_id"] = strconv.Itoa(beforeId)
	}
	if page != 0 {
		queryParams["page"] = strconv.Itoa(page)
	}
	if tags != "" {
		queryParams["tags"] = tags
	}
	if typedTags != false {
		queryParams["typed_tags"] = strconv.FormatBool(typedTags)
	}
	client := resty.New()
	resp, err := client.R().
		SetQueryParams(queryParams).
		SetHeader("user-agent", "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/73.0.3683.86 Safari/537.36").
		SetHeader("cookie", "__cfduid=dfdd28cc2534b56b8e878dd540944bf221555100794; blacklist_avatars=true; blacklist_users=false; css=hexagon%3Boverrides%2Fchristmas%3B1462320000; mode=view; e621=BAh7CDoPc2Vzc2lvbl9pZCIlMzZhMzcwMjc4YmMxMmEzZDM0YmRlNTA4ZmU2ZGI2MzA6EF9jc3JmX3Rva2VuSSIxR0pkMXBqMWxJMVZlcWFHQmk1bkNHNlcwY2JzcGtkQ0tUc0xVaTltZGhHaz0GOgZFRkkiCmZsYXNoBjsHRklDOidBY3Rpb25Db250cm9sbGVyOjpGbGFzaDo6Rmxhc2hIYXNoewAGOgpAdXNlZHsA--62ec4afec78a2df7a9934f6dc6ea8c20c8731da4; cf_clearance=42717a4c94b201952cb6f6b4bb38a874a6808c26-1555114721-2592000-250").
		SetHeader("Accept", "application/json").
		Get(listUrl)

	if err != nil {
		return nil, err
	}

	_ = ioutil.WriteFile("api/before_id_"+strconv.Itoa(beforeId)+".json", resp.Body(), 0755)

	var posts Posts
	err = json.Unmarshal(resp.Body(), &posts)
	if err != nil {
		return nil, err
	}

	return posts, nil
}
