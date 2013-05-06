package adn

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"time"
)

type Source struct {
	Name     string `json:"name,omitempty"`
	Link     string `json:"link,omitempty"`
	ClientID string `json:"client_id,omitempty"`
}

type Mention struct {
	Name     string `json:"name,omitempty"`
	ID       string `json:"id,omitempty"`
	Position int    `json:"pos,omitempty"`
	Length   int    `json:"len,omitempty"`
}

type Hashtag struct {
	Name     string `json:"name,omitempty"`
	Position int    `json:"pos,omitempty"`
	Length   int    `json:"len,omitempty"`
}

type Link struct {
	Test          string `json:"text,omitempty"`
	URL           string `json:"url,omitempty"`
	Position      int    `json:"pos,omitempty"`
	Length        int    `json:"len,omitempty"`
	AmendedLength int    `json:"amended_len,omitempty"`
}

type Entities struct {
	Mentions []Mention `json:"mentions,omitempty"`
	Hashtags []Hashtag `json:"hashtags,omitempty"`
	Links    []Link    `json:"links,omitempty"`
}

type Annotation struct {
	Type  string                 `json:"type,omitempty"`
	Value map[string]interface{} `json:"value,omitempty"`
}

type Post struct {
	ID           string       `json:"id,omitempty"`
	User         *User        `json:"user,omitempty"`
	CreatedAt    time.Time    `json:"created_at,omitempty"`
	Text         string       `json:"text,omitempty"`
	HTML         string       `json:"html,omitempty"`
	Source       Source       `json:"source,omitempty"`
	ReplyTo      *string      `json:"reply_to,omitempty"`
	CanonicalURL string       `json:"canonical_url,omitempty"`
	ThreadID     string       `json:"thread_id,omitempty"`
	NumReplies   int          `json:"num_replies,omitempty"`
	NumStars     int          `json:"num_stars,omitempty"`
	NumReposts   int          `json:"num_reposts,omitempty"`
	Annotations  []Annotation `json:"annotations,omitempty"`
	Entities     *Entities    `json:"entities,omitempty"`
	IsDeleted    bool         `json:"is_deleted,omitempty"`
	MachineOnly  bool         `json:"machine_only,omitempty"`
	YouStarred   bool         `json:"you_starred,omitempty"`
	StarredBy    []User       `json:"starred_by,omitempty"`
	YouReposted  bool         `json:"you_reposted,omitempty"`
	Reposters    []User       `json:"reposters,omitempty"`
	RepostOf     *Post        `json:"repost_of,omitempty"`
}

func (a *ADN) CreatePost(post Post) (Post, error) {
	url := a.makeURL("/stream/0/posts")
	data, err := json.Marshal(post)
	if err != nil {
		return Post{}, err
	}
	req, err := a.makeRequest("POST", url, bytes.NewBuffer(data))
	if err != nil {
		return Post{}, err
	}
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return Post{}, err
	}
	defer resp.Body.Close()
	respData, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return Post{}, err
	}
	var postResp Post
	err = json.Unmarshal(respData, &postResp)
	if err != nil {
		return Post{}, err
	}
	return postResp, nil
}
