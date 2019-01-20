package Uni
import (
	"fmt"
	"io/ioutil"
	"net/http"
	"encoding/json"
	/*"bytes"
	"github.com/BurntSushi/toml"*/
	"github.com/bwmarrin/discordgo"
)

type RedditSearchResult struct {
	Kind string
	Data RedditData
}

type RedditData struct {
	Modhash string
	Dist int
	Children []RedditPost
}

type RedditPost struct {
	Kind string
	Data struct {
		Author string
		Clicked bool
		ID string
		Link_flair_text string
		Mod_reason_title string
		Name string
		Num_crossposts int
		Num_comments int
		Over_18 bool
		SubReddit string
		Selftext string
		Saved bool
		Title string
		Ups int
		URL string
	}
}


func Request(link string) ([]byte, error) {
	tmpClient := http.Client{}
	req, err := http.NewRequest(http.MethodGet, link, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("User-Agent", "UniBot_Discord_Reddit-Search")
	fmt.Println(link)
	res, err := tmpClient.Do(req)
	if err != nil {
		return nil, err
	}
	//fmt.Println(res)
	
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	return body, nil
}

func RequestRedditDataFrom(link string, limit int, after string) (*RedditSearchResult, error) {
	rlink := fmt.Sprintf("%s", link)
	if limit != 0 {
		rlink = fmt.Sprintf("%s?limit=%d", rlink, limit)
	}
	
	if after != "" {
		rlink = fmt.Sprintf("%s?after=%s", rlink, after)
	}
	
	
	rd, err := Request(rlink)
	if err != nil {
		return nil, err	
	}
	
	var rSearch *RedditSearchResult
	
	err = json.Unmarshal(rd, &rSearch)
	if err != nil {
		return nil, err	
	}
	
	return rSearch, nil
	
}

func (Uni *UniBot) GrabRedditPost(s *discordgo.Session, m *discordgo.MessageCreate, subreddit string, nsfw bool, sortby, after string, retries int) {
	rlink := "https://reddit.com/"
	rlink = fmt.Sprintf("%sr/%s", rlink, subreddit)
	rrd, err := RequestRedditDataFrom(fmt.Sprintf("%s%s.json", rlink, sortby), 12, after)
	if err != nil {
		fmt.Println(err)
		return
	}
	
	
	
	if rrd.Data.Dist == 0 {
		if after == "" {
			Respond(s, m, "Returned no results for that")
			return
		} else {
			Uni.GrabRedditPost(s, m, subreddit, nsfw, sortby, "", retries+1)
		}
	}
	
	
	//if (Uni.Rng.Int63n(2) == 0) {
	if (Uni.Rng.Int63n(7) < 4+(Uni.Rng.Int63()%2)) { // 4 or 5
		Uni.GrabRedditPost(s, m, subreddit, nsfw, sortby, rrd.Data.Children[(int(Uni.Rng.Int63()+1) % int(rrd.Data.Dist))].Data.Name, retries+1)
	} else {
		rp := rrd.Data.Children[int(Uni.Rng.Int63() % int64(rrd.Data.Dist))]
		if rp.Data.Over_18 && !nsfw {
			if retries > 3 {
				Respond(s, m, "Returned post is nsfw, refusing to send reddit post")
				return
			} else {
				Uni.GrabRedditPost(s, m, subreddit, nsfw, sortby, rrd.Data.Children[(int(Uni.Rng.Int63()+1) % int(rrd.Data.Dist))].Data.Name, retries+1)
			}
		}
		
		lft := ""
		if len(rp.Data.Link_flair_text) != 0 {
			lft = fmt.Sprintf("\nFlair text: %s,\n\n", rp.Data.Link_flair_text)
		}
		
		for _, index := range Uni.RecentRedditPosts.IDs {
			if index == "" {
				break
			}
			
			if rp.Data.ID == index {
				Uni.GrabRedditPost(s, m, subreddit, nsfw, sortby, rrd.Data.Children[(int(Uni.Rng.Int63()+1) % int(rrd.Data.Dist))].Data.Name, retries+1)
				return
			}
		}
		
		Uni.RecentRedditPosts.IDs[Uni.RecentRedditPosts.Index] = rp.Data.ID
		Uni.RecentRedditPosts.Index = (Uni.RecentRedditPosts.Index+1) % 25
		
		Respond(s, m, fmt.Sprintf("%s\n\n\"%s\" on r/%s\n\n%d ↑, %d X-posts, %d Comments\n%s<https://redd.it/%s>", rp.Data.URL, rp.Data.Title, rp.Data.SubReddit, rp.Data.Ups, rp.Data.Num_crossposts, rp.Data.Num_comments, lft, rp.Data.ID))
	}
	
}

