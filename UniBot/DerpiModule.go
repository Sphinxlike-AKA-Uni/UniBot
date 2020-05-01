package Uni
import (
	"fmt"
	"time"
	"strings"
	"net/url"
	"encoding/json"
	"github.com/bwmarrin/discordgo"
)

// https://derpibooru.org/pages/api

// Derpi Structs
type Image struct { // Derpibooru's way to store image info
	Updated_At string `json:"updated_at"`
	Spoilered bool `json:"spoilered"`
	Created_At string `json:"created_at"`
	Deletion_Reason string `json:"deletion_reason"`
	Sha512Hash string `json:"sha512_hash"`
	Duplicate_Of string `json:"duplicate_of"`
	Description string `json:"description"`
	AspectRatio float64 `json:"aspect_ratio"`
	Tag_Ids []int64 `json:"tag_ids"`
	Downvotes int64 `json:"downvotes"`
	Mime_Type string `json:"mime_type"`
	Tag_Count int64 `json:"tag_count"`
	Upvotes int64 `json:"upvotes"`
	UploaderID int64 `json:"uploader_id"`
	Name string `json:"name"`
	Representations struct {
		ThumbTiny string `json:"thumb_tiny"`
		ThumbSmall string `json:"thumb_small"`
		Thumb string `json:"thumb"`
		Small string `json:"small"`
		Medium string `json:"medium"`
		Large string `json:"large"`
		Tall string `json:"tall"`
		Full string `json:"full"`
	} `json:"representations"`
	Uploader string `json:"uploader"`
	Faves int64 `json:"faves"`
	ID int64 `json:"id"`
	Source_URL string `json:"source_url"`
	Height int `json:"height"`
	Score int64 `json:"score"`
	Hidden_From_Users bool `json:"hidden_from_users"`
	Tags []string `json:"tags"`
	Width int `json:"width"`
	First_Seen_At string `json:"first_seen_at"`
	Comment_Count int64 `json:"comment_count"`
	Orig_Sha512_Hash string `json:"orig_sha512_hash"`
	Wilson_Score float64 `json:"wilson_score"`
	Format string `json:"format"`
	Thumbnails_Generated bool `json:"thumbnails_generated"`
	View_Url string `json:"view_url"`
	Intensities struct {
		NE float64 `json:"ne"`
		NW float64 `json:"nw"`
		SE float64 `json:"se"`
		SW float64 `json:"sw"`
	} `json:"intensities"`
	//Interactions []interface{} `json:"interactions"`// still dunno what this is
}

// Derpi Search Results
type Search struct {
	Images []Image `json:"images"`
	Total int64 `json:"total"`
	Interactions []interface{} `json:"interactions"`// still dunno what this is
}


// For derpi filters
type Filter struct {
	ID int64
	Name string
	Description string
	Hidden_Tag_IDs []int64
	Spoilered_Tag_IDs []int64
	Spoilered_Tags []string
	Hidden_Complex string
	Spoilered_Complex string
	Public bool
	System bool
	User_Count int64
	User_ID int64
}

type GetFilter struct {
	Filter *Filter `json:"filter"`
}

// Oembed
type Oembed struct {
	Author_Name string `json:"author_name"`
	Author_URL string `json:"author_url"`
	Cache_Age int64 `json:"cache_age"`
	Derpibooru_Comments int64 `json:"derpibooru_comments"`
	Derpibooru_ID int64 `json:"derpibooru_id"`
	Derpibooru_Score int64 `json:"derpibooru_score"`
	Derpibooru_Tags []string `json:"derpibooru_tags"`
	Provider_Name string `json:"provider_name"`
	Provider_URL string `json:"provider_url"`
	Title string `json:"title"`
	Type string `json:"type"`
	Version string `json:"version"`
}

// Derpi Comment
type Comment struct {
	Author string `json:"author"`
	Avatar string `json:"avatar"`
	Body string `json:"body"`
	Created_At time.Time `json:"created_at"`
	Edit_Reason string `json:"edit_reason"`
	Edit_At string `json:"edit_at"`
	ID int64 `json:"id"`
	Image_ID int64 `json:"image_id"`
	Updated_At time.Time `json:"updated_at"`
	User_ID int64 `json:"user_id"`
}

type GetComment struct {
	Comment *Comment `json:"comment"`
}


// Derpi API Stuff here

// Get Filter data from derpi
func (Uni *UniBot) GetDerpiFilter(filterid string) (*Filter, error) {
	resp, err := Uni.HTTPRequest("GET", fmt.Sprintf("https://derpibooru.org/api/v1/json/filters/%d", filterid), map[string]interface{}{"User-Agent": GrabUserAgent(),}, nil)
	if err != nil { return nil, err }
	var f *Filter
	json.NewDecoder(resp.Body).Decode(&f)
	return f, nil
}

// Grab a search of certain tags and other parameters
func (Uni *UniBot) DerpiSearch(tags, filterid string, AdditionalParameters map[string]interface{}) (*Search, error) {
	// https://derpibooru.org/api/v1/json/search/images?q=safe&per_page=1
	// filter_id 	Assuming the user can access the filter ID given by the parameter, overrides the current filter for this request. This is primarily useful for unauthenticated API access.
	// key 	An optional authentication token. If omitted, no user will be authenticated.

	// You can find your authentication token in your account settings.
	// page 	Controls the current page of the response, if the response is paginated. Empty values default to the first page.
	// per_page 	Controls the number of results per page, up to a limit of 50, if the response is paginated. The default is 25.
	// q 	The current search query, if the request is a search request.
	// sd 	The current sort direction, if the request is a search request.
	// sf 	The current sort field, if the request is a search request.
	// {Those are notes from the API and i just copied them here}
	
	var link string = fmt.Sprintf("https://derpibooru.org/api/v1/json/search/images?q=%s&filter_id=%s", tags, filterid)
	if AdditionalParameters == nil { goto ProcessLink } // Scanning "AdditionalParameters" since it would crash if it did
	for k, v := range AdditionalParameters {
		link = fmt.Sprintf("%s&%s=%v", link, k, v)
	}
	
ProcessLink:
	var s *Search
	err := Uni.HTTPRequestJSON("GET", link, map[string]interface{}{"User-Agent": GrabUserAgent()}, nil, &s)
	return s, err
}

// Discord stuff here

// Get a random image with the following tags(and filter if set)
func (Uni *UniBot) SearchOnDerpi(cID, tags string) {
	tags = url.QueryEscape(tags)
	tags = strings.Replace(strings.Replace(tags, " ", "+", -1), "&", "%26", -1)
	f, err := Uni.GetChannelDerpiFilter(cID)
	if err != nil {
		Uni.ErrRespond(err, cID, "getting derpibooru filter", map[string]interface{}{"err": err, "cID": cID, "tags": tags, "filter": f})
		return
	}
	s, err := Uni.DerpiSearch(tags, f, map[string]interface{}{"sf": "random", "per_page": 1})
	Uni.ErrRespond(err, cID, "searching derpibooru", map[string]interface{}{"err": err, "cID": cID, "tags": tags, "filter": f})
	
	if len(s.Images) == 0 {
		Uni.Respond(cID, "Search has returned empty")
		return
	}
	
	var imagetags string = strings.Join(s.Images[0].Tags, ", ")
	if len(imagetags) > 2047 { // description capped at 2048 characters
		imagetags = imagetags[:2043]+"...."
	}
	
	embed := &discordgo.MessageEmbed{
		Type: "image",
		/*Author: &discordgo.MessageEmbedAuthor{
			URL:	
			Name:	"",
		},*/
		Color:  int(<-Uni.RNGChan%(1<<24)), //random color picker
		Fields: []*discordgo.MessageEmbedField{
			&discordgo.MessageEmbedField{Name: "Upload Date&Time", Value: s.Images[0].Created_At, Inline: true},
			&discordgo.MessageEmbedField{Name: "Score", Value: fmt.Sprintf("%d - %d = %d", s.Images[0].Upvotes, s.Images[0].Downvotes, s.Images[0].Upvotes-s.Images[0].Downvotes, ), Inline: true},
			&discordgo.MessageEmbedField{Name: "Image Info", Value: fmt.Sprintf("%s Format, (%dx%d)", s.Images[0].Format, s.Images[0].Width, s.Images[0].Height,), Inline: true},
		},
		Image: &discordgo.MessageEmbedImage{
			URL: s.Images[0].View_Url,
			Width: s.Images[0].Width,
			Height: s.Images[0].Height,
		},
		Title: "Click here for the derpibooru image page",
		URL: fmt.Sprintf("https://derpibooru.org/images/%d", s.Images[0].ID),
		Description: imagetags,
	}
	
	if s.Images[0].Format == "gif" {
		embed.Type = "gifv"
	} else if s.Images[0].Format == "webm" {
		embed.Type = "video"
		embed.Video = &discordgo.MessageEmbedVideo{URL: s.Images[0].View_Url, Width: s.Images[0].Width, Height: s.Images[0].Height,}
		embed.Image = nil
	}
	// Am going to clean most of this up soon
	Uni.S.ChannelMessageSendEmbed(cID, embed)
}

// Grab a channel's derpi filter (set to unibot's public default if none applied)
func (Uni *UniBot) GetChannelDerpiFilter(cID string) (string, error) {
	var fstr string = "157679" // https://www.derpibooru.org/filters/157679 being the default filter
	err := Uni.DBGetFirstVar(&fstr, "GetDerpiFilter", cID)
	return fstr, err
}

// Set the channel's derpi filter
func (Uni *UniBot) SetChannelDerpiFilter(gID, cID, filterid string) {
	f, err := Uni.GetDerpiFilter(filterid)
	Uni.ErrRespond(err, cID, "requesting filter data", map[string]interface{}{"gID": gID, "cID": cID, "err": err, "filterid": filterid})
	if f == nil { // redirected
		Uni.Respond(cID, "Filter seems to have returned nil, is the filter public?")
		return
	}
	
	// proceed
	fstr := ""
	Uni.DBGetFirstVar(&fstr, "GetDerpiFilter", cID)
	if fstr == "" { // no such index exists, create index
		_, err = Uni.DBExec("InsertDerpiFilter", gID, cID, filterid)
	} else { // index for channel already exists, update index
		_, err = Uni.DBExec("UpdateDerpiFilter", filterid, cID)
	}
	
	Uni.ErrRespond(err, cID, "setting channel derpibooru filter", map[string]interface{}{"gID": gID, "cID": cID, "err": err, "filterid": filterid, "fstr": fstr})
	Uni.Respond(cID, fmt.Sprintf("Filter set to ID: %d, %q", f.ID, f.Name))
}
