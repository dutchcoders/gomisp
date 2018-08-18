package misp

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"time"
)

type Client struct {
	Key string

	*http.Client
	baseURL *url.URL
}

type SearchRequest struct {
	from  time.Time `json:"from"`
	to    time.Time `json:"to"`
	value string    `json:"value"`
	type_ string    `json:"type"`
}

func NewSearchRequest() *SearchRequest {
	return &SearchRequest{}
}

func (sr *SearchRequest) MarshalJSON() ([]byte, error) {
	m := map[string]interface{}{}
	m["from"] = sr.from.Format("2006-01-02")
	m["to"] = sr.to.Format("2006-01-02")

	if sr.value != "" {
		m["value"] = sr.value
	}

	if sr.type_ != "" {
		m["type"] = sr.type_
	}

	return json.Marshal(m)
}

func (sr *SearchRequest) From(t time.Time) *SearchRequest {
	sr.from = t
	return sr
}

func (sr *SearchRequest) To(t time.Time) *SearchRequest {
	sr.to = t
	return sr
}

func (sr *SearchRequest) Type(t string) *SearchRequest {
	sr.type_ = t
	return sr
}

func (sr *SearchRequest) Value(val string) *SearchRequest {
	sr.value = val
	return sr
}

type MISPResponse struct {
	Response json.RawMessage `json:"response"`
}

type SearchResult struct {
	Event Event `json:"Event"`
}

type Event struct {
	Analysis  string `json:"analysis"`
	Attribute []struct {
		Category           string `json:"category"`
		Comment            string `json:"comment"`
		Deleted            bool   `json:"deleted"`
		DisableCorrelation bool   `json:"disable_correlation"`
		Distribution       string `json:"distribution"`
		EventId            string `json:"event_id"`
		Galaxy             []interface{}
		Id                 string      `json:"id"`
		ObjectId           string      `json:"object_id"`
		ObjectRelation     interface{} `json:"object_relation"`
		ShadowAttribute    []interface{}
		SharingGroupId     string `json:"sharing_group_id"`
		Tag                []struct {
			Colour     string `json:"colour"`
			Exportable bool   `json:"exportable"`
			HideTag    bool   `json:"hide_tag"`
			Id         string `json:"id"`
			Name       string `json:"name"`
			UserId     string `json:"user_id"`
		}
		Timestamp string `json:"timestamp"`
		ToIds     bool   `json:"to_ids"`
		Type      string `json:"type"`
		Uuid      string `json:"uuid"`
		Value     string `json:"value"`
	}
	AttributeCount     string `json:"attribute_count"`
	Date               string `json:"date"`
	DisableCorrelation bool   `json:"disable_correlation"`
	Distribution       string `json:"distribution"`
	ExtendsUuid        string `json:"extends_uuid"`
	Galaxy             []interface{}
	Id                 string `json:"id"`
	Info               string `json:"info"`
	Locked             bool   `json:"locked"`
	Object             []interface{}
	Org                struct {
		Id   string `json:"id"`
		Name string `json:"name"`
		Uuid string `json:"uuid"`
	}
	OrgId string `json:"org_id"`
	Orgc  struct {
		Id   string `json:"id"`
		Name string `json:"name"`
		Uuid string `json:"uuid"`
	}
	OrgcId            string `json:"orgc_id"`
	ProposalEmailLock bool   `json:"proposal_email_lock"`
	PublishTimestamp  string `json:"publish_timestamp"`
	Published         bool   `json:"published"`
	RelatedEvent      []interface{}
	ShadowAttribute   []interface{}
	SharingGroupId    string `json:"sharing_group_id"`
	Tag               []struct {
		Colour     string `json:"colour"`
		Exportable bool   `json:"exportable"`
		HideTag    bool   `json:"hide_tag"`
		Id         string `json:"id"`
		Name       string `json:"name"`
		UserId     string `json:"user_id"`
	}
	ThreatLevelId string `json:"threat_level_id"`
	Timestamp     string `json:"timestamp"`
	Uuid          string `json:"uuid"`
}

func (c *Client) Search(sr *SearchRequest) ([]SearchResult, error) {
	request := struct {
		Request json.Marshaler `json:"request"`
	}{
		Request: sr,
	}

	req, err := c.NewRequest("POST", "/events/restSearch/download", request)
	if err != nil {
		return nil, err
	}

	resp := MISPResponse{}

	if err := c.Do(req, &resp); err != nil {
		return nil, err
	}

	result := []SearchResult{}
	if err := json.Unmarshal(resp.Response, &result); err != nil {
		return nil, err
	}

	return result, nil
}

type optionFn func(*Client)

// WithURL contains the MIPS target url
func WithURL(u url.URL) optionFn {
	return func(c *Client) {
		c.baseURL = &u
	}
}

// WithKey contains the MIPS API key
func WithKey(key string) optionFn {
	return func(c *Client) {
		c.Key = key
	}
}

// New returns a MIPS API client
func New(options ...optionFn) (*Client, error) {
	c := &Client{
		Client: http.DefaultClient,
	}

	for _, optionFn := range options {
		optionFn(c)
	}

	if c.baseURL == nil {
		return nil, fmt.Errorf("URL not set")
	}

	return c, nil
}
