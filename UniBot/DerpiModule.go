package Uni

import (
	"fmt"
	"net/http"
	"time"
	"math"
	"strings"
	"net/url"
	"io/ioutil"
	"encoding/json"
	"github.com/bwmarrin/discordgo"
)








type Image struct { // Derpibooru's way to store image info
	ID               int             `json:"id"`
	CreatedAt        time.Time          `json:"created_at"`
	UpdatedAt        time.Time          `json:"updated_at"`
	DuplicateReports []*DuplicateReport `json:"duplicate_reports"`
	FirstSeenAt      time.Time          `json:"first_seen_at"`
	UploaderID       interface{}        `json:"uploader_id"`
	FileName         string             `json:"file_name"`
	Description      string             `json:"description"`
	Uploader         string             `json:"uploader"`
	Image            string             `json:"image"`
	Score            int                `json:"score"`
	Upvotes          int                `json:"upvotes"`
	Downvotes        int                `json:"downvotes"`
	Faves            int                `json:"faves"`
	CommentCount     int                `json:"comment_count"`
	Tags             string             `json:"tags"`
	TagIds           []int              `json:"tag_ids"`
	Width            int                `json:"width"`
	Height           int                `json:"height"`
	AspectRatio      float64            `json:"aspect_ratio"`
	OriginalFormat   string             `json:"original_format"`
	MimeType         string             `json:"mime_type"`
	Sha512Hash       string             `json:"sha512_hash"`
	OrigSha512Hash   string             `json:"orig_sha512_hash"`
	SourceURL        string             `json:"source_url"`
	Representations  struct {
		ThumbTiny  string `json:"thumb_tiny"`
		ThumbSmall string `json:"thumb_small"`
		Thumb      string `json:"thumb"`
		Small      string `json:"small"`
		Medium     string `json:"medium"`
		Large      string `json:"large"`
		Tall       string `json:"tall"`
		Full       string `json:"full"`
	} `json:"representations"`
	IsRendered   bool          `json:"is_rendered"`
	IsOptimized  bool          `json:"is_optimized"`
	Interactions []interface{} `json:"interactions"`
}


type DuplicateReport struct {
	ID                  int64               `json:"id"`
	State               string              `json:"state"`
	Reason              string              `json:"reason"`
	ImageIDNumber       int                 `json:"image_id_number"`
	TargetImageIDNumber int                 `json:"target_image_id_number"`
	User                interface{}         `json:"user"`
	CreatedAt           string              `json:"created_at"`
	Modifier            *DupeReportModifier `json:"modifier"`
}



type DupeReportModifier struct {
	ID           int64  `json:"id"`
	Name         string `json:"name"`
	Avatar       string `json:"avatar"`
	CommentCount int    `json:"comment_count"`
	UploadCount  int    `json:"upload_count"`
	PostCount    int    `json:"post_count"`
	TopicCount   int    `json:"topic_count"`
}

type Search struct {
	Search []Image `json:"search"`
	Total int `json:"total"`
	Interactions []interface{} `json:"interactions"`
}

// For derpi filters
type Filter struct {
	ID int
	Name string
	Description string
	Hidden_Tag_IDs []int
	Spoilered_Tag_IDs []int
	Spoilered_Tags []string
	Hidden_Complex string
	Spoilered_Complex string
	Public bool
	System bool
	User_Count int
	User_ID int
}


func DerpiSearch(tags string, filterID, page, perpage int) (*Search, error)  {
	tmpClient := http.Client{}
	
	req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("https://derpibooru.org/search.json?q=%s&filter_id=%d&page=%d&perpage=%d", tags, filterID, page, perpage), nil)
	if err != nil {
		return nil, err
	}
	
	req.Header.Set("User-Agent", "Uni_Derpi_Search")
	
	res, getErr := tmpClient.Do(req)
	if getErr != nil {
		return nil, getErr
	}
	
	body, readErr := ioutil.ReadAll(res.Body)
	if readErr != nil {
		return nil, readErr
	}
	
	var tSearch *Search
	
	jsonErr := json.Unmarshal(body, &tSearch)	
	return tSearch, jsonErr
}



func GetDerpiFilter(filterid string) (*Filter, error) {
	tmpClient := http.Client{}
	
	req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("https://derpibooru.org/filters/%s.json", filterid), nil)
	if err != nil {
		return nil, err
	}
	
	req.Header.Set("User-Agent", "Uni_Derpi_Search")
	
	res, getErr := tmpClient.Do(req)
	if getErr != nil {
		return nil, getErr
	}
	
	body, readErr := ioutil.ReadAll(res.Body)
	if readErr != nil {
		return nil, readErr
	}
	
	var tFil *Filter
	
	json.Unmarshal(body, &tFil)	
	return tFil, nil
}

func GetVSLinkOfImage(i Image) string {
	return fmt.Sprintf("https:%s", strings.Replace(strings.Replace(i.Representations.Tall, "/img/", "/img/view/", -1), "/tall", "", -1))
}


func (Uni *UniBot) GetFilter(s *discordgo.Session, m *discordgo.MessageCreate) int {
	c, _ := s.State.Channel(m.ChannelID)
	f := 56027 // The spooky everything filter
	
	if c.GuildID == "" {
		if a := Uni.CheckNSFW(s, m); !a {
			f = 157679
		}
	} else {
		
		rows, err := Uni.Database.Query(fmt.Sprintf("SELECT filterID FROM DerpiFilters WHERE cID = '%s'", c.ID))
		if err != nil {
			ErrRespond(s, m, fmt.Sprintf(" (Query failure while checking database) %s", err))
		}
		
		defer rows.Close()
		
		for rows.Next() { // Index already exists
			rows.Scan(&f)
			return f
		}
		
		
		if !c.NSFW {
			f = 157679
		}
	}
	
	return f
	
}


func (Uni *UniBot) DerpiSetFilter(s *discordgo.Session, m *discordgo.MessageCreate, gID, cID, filterID string) {
	a, err := GetDerpiFilter(filterID)
	if err != nil {
		Respond(s, m, fmt.Sprintf("Error while grabbing filter: ", err))
		return
	}
	
	if a == nil {
		Respond(s, m, "Filter seems to have returned nil, is the filter public?")
		return
	}
	
	rows, err := Uni.Database.Query(fmt.Sprintf("SELECT * FROM DerpiFilters WHERE gID = '%s' AND cID = '%s'", gID, cID))
	if err != nil {
		ErrRespond(s, m, fmt.Sprintf(" (Query failure while checking database) %s", err))
	}
	
	overwrite := 0
	
	for rows.Next() { // Index already exists
		overwrite = 1
	}
	
	rows.Close()
	
	if overwrite == 1 {
		_, err := Uni.Database.Exec(fmt.Sprintf("UPDATE DerpiFilters SET filterID = \"%s\" WHERE cID = \"%s\";", filterID, cID))
		if err != nil {
			Respond(s, m, fmt.Sprintf("%s", err))
			return
		}
	} else {
		Uni.Database.Exec(fmt.Sprintf("INSERT INTO DerpiFilters VALUES ('%s','%s','%s');", gID, cID, filterID))
	}
	
	Respond(s, m, fmt.Sprintf("Filter %s to ID: %d, \"%s\"", []string{"set", "overwritten"}[overwrite], a.ID, a.Name))
	
	
}


func (Uni *UniBot) getsearch(s *discordgo.Session, m *discordgo.MessageCreate, tags string) (*Search, error) {
	f := Uni.GetFilter(s, m)
	t := strings.Replace(strings.Replace(tags, " ", "+", -1), "&", "%26", -1)
	tSearch, err := DerpiSearch(t, f, 1, 1)
	
	if err != nil {
		return nil, err
	}

	if tSearch.Total == 0 {
		return DerpiSearch(t, f, 1, 1)
	} else {
		return DerpiSearch(t, f, int(math.Abs(float64(Uni.Rng.Intn(tSearch.Total))))+1, 1)
	}
}




func (Uni *UniBot) SearchOnDerpi(s *discordgo.Session, m *discordgo.MessageCreate, tags string) {
	tags = url.QueryEscape(tags)
	sr, err := Uni.getsearch(s, m, tags)
	if err != nil {
		ErrRespond(s, m, fmt.Sprintf(" (Searching error) Error occurred; %s", err))
	} else {
		if len(sr.Search) != 0 {
			Respond(s, m, GetVSLinkOfImage(sr.Search[0]))
		} else {
			Respond(s, m, fmt.Sprintf("Returned inquiry is empty for \"%s\"", tags))
		}
	}
}



func (Uni *UniBot) ImageInfo(s *discordgo.Session, m *discordgo.MessageCreate, ID string) {
	i, err := DerpiSearch(fmt.Sprintf("id:%s", ID), Uni.GetFilter(s, m), 1, 1)
	if err != nil {
		ErrRespond(s, m, fmt.Sprintf(" (Getting Image error) Error occurred; %s", err))
		return
	}
	if len(i.Search) == 0 {
		Respond(s, m, "The image ID either doesn't exist, deleted or has been filtered")
		return
	}
	r := fmt.Sprintf("Image results of <https://derpibooru.org/%s>\n", ID)
	r = fmt.Sprintf("%sTags: `%s`\n", r, i.Search[0].Tags)
	r = fmt.Sprintf("%sScore: %d ↑  %d ↓\n", r, i.Search[0].Upvotes, i.Search[0].Downvotes)
	r = fmt.Sprintf("%soriginal_format: %s\n", r, i.Search[0].OriginalFormat)
	if i.Search[0].SourceURL != "" {
		r = fmt.Sprintf("%sSource URL: <%s>\n", r, i.Search[0].SourceURL)
	}
	r = fmt.Sprintf("%sVS Link: %s\n", r, GetVSLinkOfImage(i.Search[0]))
	Respond(s, m, r)
}
