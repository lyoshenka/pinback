package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/weppos/publicsuffix-go/publicsuffix"
	// _ "github.com/motemen/go-loghttp/global" // log http requests
)

type Client struct {
	BaseURL *url.URL
}

func NewClient(authToken string) Client {
	u, _ := url.Parse("https://api.pinboard.in/?format=json")
	q := u.Query()
	q.Set("auth_token", authToken)
	u.RawQuery = q.Encode()
	return Client{
		BaseURL: u,
	}
}

func (c Client) Query(path string, values url.Values, data interface{}) error {
	u := *c.BaseURL
	u.Path = path

	q := u.Query()
	for key, vals := range values {
		q[key] = vals
	}
	u.RawQuery = q.Encode()

	rsp, err := http.Get(u.String())
	if err != nil {
		return err
	}
	defer rsp.Body.Close()

	if rsp.StatusCode != http.StatusOK {
		return fmt.Errorf("bad response: %s", rsp.Status)
	}

	body, err := ioutil.ReadAll(rsp.Body)
	if err != nil {
		return err
	}

	trimmed := bytes.TrimLeft(body, "\xef\xbb\xbf") // remove UTF-8 BOM

	return json.Unmarshal(trimmed, &data)
}

func (c Client) Recent() ([]Post, error) {
	var data struct {
		Date  time.Time `json:"date"`
		User  string    `json:"user"`
		Posts []Post    `json:"posts"`
	}
	q := url.Values{}
	q.Add("count", "100")
	if err := c.Query("/v1/posts/recent", q, &data); err != nil {
		return nil, err
	}
	return data.Posts, nil
}

func (c Client) Since(t time.Time) ([]Post, error) {
	q := url.Values{}
	q.Add("fromdt", t.UTC().Format("2006-01-02T15:04:05Z"))
	var posts []Post
	if err := c.Query("/v1/posts/all", q, &posts); err != nil {
		return nil, err
	}
	return posts, nil
}

type Post struct {
	Title       string    `json:"description"`
	Description string    `json:"extended"`
	Hash        string    `json:"hash"`
	URL         string    `json:"href"`
	Time        time.Time `json:"time"`
	Tags        []string  `json:"tags"`
	Shared      bool      `json:"shared"`
	Toread      bool      `json:"toread"`
	Meta        string    `json:"meta"`
}

func (p Post) Domain() string {
	u, err := url.Parse(p.URL)
	if err != nil {
		panic(err)
	}
	d, err := publicsuffix.Domain(u.Hostname())
	if err != nil {
		panic(err)
	}
	return strings.ToLower(d)
}

func (p *Post) UnmarshalJSON(b []byte) error {
	type tmpPost Post
	tmp := struct {
		tmpPost
		Tags   string `json:"tags"`
		Shared string `json:"shared"`
		Toread string `json:"toread"`
	}{}

	err := json.Unmarshal(b, &tmp)
	if err != nil {
		return err
	}

	p.Title = tmp.Title
	p.Description = tmp.Description
	p.Hash = tmp.Hash
	p.URL = tmp.URL
	p.Time = tmp.Time
	p.Tags = strings.Fields(tmp.Tags)
	p.Shared = tmp.Shared == "yes"
	p.Toread = tmp.Toread == "yes"
	p.Meta = tmp.Meta
	return nil
}
