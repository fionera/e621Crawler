package api

import (
	"encoding/json"
	"io/ioutil"
	"strconv"
	
	"github.com/sirupsen/logrus"
	"github.com/go-resty/resty"
)

const (
	listUrl  = "https://e621.net/posts.json"
	maxLimit = 320
)

type Posts struct {
	Posts []Post `json:"posts"`
}
type File struct {
	Width  int    `json:"width"`
	Height int    `json:"height"`
	Ext    string `json:"ext"`
	Size   int    `json:"size"`
	Md5    string `json:"md5"`
	URL    string `json:"url"`
}
type Preview struct {
	Width  int    `json:"width"`
	Height int    `json:"height"`
	URL    string `json:"url"`
}
type Sample struct {
	Has    bool   `json:"has"`
	Height int    `json:"height"`
	Width  int    `json:"width"`
	URL    string `json:"url"`
}
type Score struct {
	Up    int `json:"up"`
	Down  int `json:"down"`
	Total int `json:"total"`
}
type Tags struct {
	General   []string      `json:"general"`
	Species   []string      `json:"species"`
	Character []string      `json:"character"`
	Copyright []string      `json:"copyright"`
	Artist    []string      `json:"artist"`
	Invalid   []interface{} `json:"invalid"`
	Lore      []interface{} `json:"lore"`
	Meta      []string      `json:"meta"`
}
type Flags struct {
	Pending      bool `json:"pending"`
	Flagged      bool `json:"flagged"`
	NoteLocked   bool `json:"note_locked"`
	StatusLocked bool `json:"status_locked"`
	RatingLocked bool `json:"rating_locked"`
	Deleted      bool `json:"deleted"`
}
type Relationships struct {
	ParentID          interface{}   `json:"parent_id"`
	HasChildren       bool          `json:"has_children"`
	HasActiveChildren bool          `json:"has_active_children"`
	Children          []interface{} `json:"children"`
}
type Post struct {
	ID            int           `json:"id"`
	CreatedAt     string        `json:"created_at"`
	UpdatedAt     string        `json:"updated_at"`
	File          File          `json:"file"`
	Preview       Preview       `json:"preview"`
	Sample        Sample        `json:"sample"`
	Score         Score         `json:"score"`
	Tags          Tags          `json:"tags"`
	LockedTags    []interface{} `json:"locked_tags"`
	ChangeSeq     int           `json:"change_seq"`
	Flags         Flags         `json:"flags"`
	Rating        string        `json:"rating"`
	FavCount      int           `json:"fav_count"`
	Sources       []string      `json:"sources"`
	Pools         []interface{} `json:"pools"`
	Relationships Relationships `json:"relationships"`
	ApproverID    int           `json:"approver_id"`
	UploaderID    int           `json:"uploader_id"`
	Description   string        `json:"description"`
	CommentCount  int           `json:"comment_count"`
	IsFavorited   bool          `json:"is_favorited"`
	HasNotes      bool          `json:"has_notes"`
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
	client.SetHeaders(map[string]string{
		"user-agent": "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/73.0.3683.86 Safari/537.36",
		"Accept":     "application/json",
	})
	resp, err := client.R().SetQueryParams(queryParams).Get(listUrl)

	if err != nil {
		logrus.Fatal(err)
	}

	_ = ioutil.WriteFile("api/before_id_"+strconv.Itoa(beforeId)+".json", resp.Body(), 0755)

	var posts Posts
	err = json.Unmarshal(resp.Body(), &posts)
	if err != nil {
		logrus.Fatal(err)
	}

	return posts, nil
}
