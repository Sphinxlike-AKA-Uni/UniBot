package Uni
import (
	"fmt"
	"strings"
)
// Reddit Structs
type RedditSearchResult struct {
	Kind string
	Data RedditData
}

type RedditData struct {
	Modhash string
	Dist int
	Children []RedditPost
	After string
	Before string
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
		Stickied bool
		Title string
		Ups int
		URL string
	}
}

// Send a random reddit post from requested subreddit
func (Uni *UniBot) GrabRedditPost(cID, subreddit string, nsfw bool, sortby string) {
	if strings.HasPrefix(subreddit, "/") { subreddit = subreddit[1:]
	} else if !strings.HasPrefix(subreddit, "r/") { subreddit = "r/"+subreddit }
	
	if strings.HasSuffix(sortby, " ") { sortby = sortby[:len(sortby)-1] } // get rid of the space at the end
	//fmt.Printf("%q %q\n", subreddit, sortby)
	var (
	pagelimit uint64 = 10
	limit int = 20
	)
	Beginning:
	var (
	rpi uint64 = <-Uni.RNGChan%pagelimit // Reddit Post Index
	srp *RedditPost // Selected Reddit Post
	after string = ""
	indexes []int
	)
	for i := uint64(0); srp == nil; i++ { // Loop through indexes
		// Grab Reddit Page
		var srs RedditSearchResult
		var link string = fmt.Sprintf("https://www.reddit.com/%s.json?limit=%d", subreddit, limit)
		if after != "" {
			link = fmt.Sprintf("%s&after=%s", link, after)
		}
		if sortby != "" {
			link = fmt.Sprintf("%s&sort=%s", link, sortby)
		}
		
		err := Uni.HTTPRequestJSON("GET", link, map[string]interface{}{"User-Agent": GrabUserAgent()}, nil, &srs)
		if err != nil { // Somehow got an error
			Uni.ErrRespond(err, cID, "getting reddit results", map[string]interface{}{"err": err, "cID": cID, "subreddit": subreddit, "RedditSearchResultStruct": srs, "link": link})
			return
		}
		
		// There was no error, proceed
		
		// Size Detection
		if len(srs.Data.Children) == 0 {
			if i == 0 { // because it was the first attempt and got nothing, probably meant the subreddit doesn't exist or has nothing
				Uni.Respond(cID, "Subreddit does not appear to exist")
				return
			}
			pagelimit = i
			goto Beginning // Subreddit apparently too small, retry but with that in mind
		}
		
		if i != rpi { // do next page
			after = srs.Data.Children[srs.Data.Dist-1].Data.Name
			continue
		}
		
		
		// Filter out NSFW posts for SFW channels
		for i, post := range srs.Data.Children {
			if !post.Data.Stickied {
				if !nsfw {
					if !post.Data.Over_18 {
						indexes = append(indexes, i)
					}
				} else {
					indexes = append(indexes, i)
				}
			}
		}
		
		if len(indexes) == 0 {
			Uni.Respond(cID, "All posts appear to be tagged NSFW")
			return
		}
		// Finally, a random post can be selected
		srp = &srs.Data.Children[indexes[<-Uni.RNGChan%uint64(len(indexes))]]
	}
	
	lft := ""
	if len(srp.Data.Link_flair_text) != 0 {
		lft = fmt.Sprintf("\nFlair text: %s,\n\n", srp.Data.Link_flair_text)
	}
	
	Uni.Respond(cID, fmt.Sprintf("%s\n\n\"%s\" on r/%s\n\n%d ↑, %d X-posts, %d Comments\n%s<https://redd.it/%s>", srp.Data.URL, srp.Data.Title, srp.Data.SubReddit, srp.Data.Ups, srp.Data.Num_crossposts, srp.Data.Num_comments, lft, srp.Data.ID))
	
}