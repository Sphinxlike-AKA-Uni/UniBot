package Uni

import (
	"fmt"
	"strings"
	"time"
	"strconv"
	"syscall"
	"math/big"
	"math/rand"
	"bytes"
	"io/ioutil"
	"crypto/sha256"
	"github.com/BurntSushi/toml"
	"github.com/bwmarrin/discordgo"
)





func (Uni *UniBot) ProcessMessage(s *discordgo.Session, m *discordgo.MessageCreate, isDM bool, g *discordgo.Guild , c *discordgo.Channel, modules int, prefix string) (int, string) {
	
	if m.Author.ID == Uni.CreatorID {
		if strings.HasPrefix(m.Content, "$~restart") {
			Uni.Msggate = true
			for Uni.Msghandlers > 1 {
				time.Sleep(44 * time.Millisecond)
			}
			Uni.Restart = true
			Uni.SC <- syscall.SIGTERM //syscall.SIGINT
		} else if strings.HasPrefix(m.Content, "$~shutdown") {
			Uni.Msggate = true
			for Uni.Msghandlers > 1 {
				time.Sleep(44 * time.Millisecond)
			}
			//os.Exit(2147483647)
			Uni.SC <- syscall.SIGTERM //syscall.SIGINT
		} else if strings.HasPrefix(m.Content, "$~reply ") {
			content := m.Content[8:]
			a := strings.Split(content, "|")
			s.ChannelMessageSend(a[0], a[1])
		}
		
		
	} 	
	

	if strings.HasPrefix(strings.ToLower(m.Content), prefix) {
		//if (!isDM && modules & 2 == 2) || isDM { // Reddit search
		if (!isDM && modules & 2 == 2) { // Reddit search
			
			nsfw := false
			if isDM {
				nsfw = Uni.CheckNSFW(s, m)
			} else {
				nsfw = c.NSFW
			}
			
			
			for _, f := range []string{
			" find a post on r/", 		" find a post on /r/",		" find a post on ", 
			" find a post in r/", 		" find a post in /r/", 		" find a post in ", 
			" find a post from r/", 	" find a post from /r/", 	" find a post from ", 
			" find a post within r/", 	" find a post within /r/", 	" find a post within ", 
			" find me a post on r/", 	" find me a post on  /r/", 	" find me a post on ", 
			" find a post r/",			" find a post /r/",	 		" find a post ", 
			} {
				if strings.HasPrefix(strings.ToLower(m.Content), prefix+f) {
					Uni.GrabRedditPost(s, m, m.Content[len(prefix)+len(f):], nsfw, "", "", 0)
					return 1, "Reddit"			
				}
			}
			
			
			for _, f := range []string{
			" find a top post on r/", " find a top post on /r/", " find a top post on ", 
			" find a top post in r/", " find a top post in /r/", " find a top post in ", 
			" find a top post from r/", " find a top post from /r/", " find a top post from ", 
			" find a top post within r/", " find a top post within /r/", " find a top post within ", 
			" find a top post r/", " find a top post /r/", " find a top post ", 
			} {
				if strings.HasPrefix(strings.ToLower(m.Content), prefix+f) {
					Uni.GrabRedditPost(s, m, m.Content[len(prefix)+len(f):], nsfw, "/top", "", 0)
					return 1, "Reddit"			
				}
			}
						
			for _, f := range []string{
			" find a new post on r/", " find a new post on /r/", " find a new post on ", 
			" find a new post in r/", " find a new post in /r/", " find a new post in ", 
			" find a new post from r/", " find a new post from /r/", " find a new post from ", 
			" find a new post within r/", " find a new post within /r/", " find a new post within ", 
			" find a new post r/", " find a new post /r/", " find a new post ", 
			} {
				if strings.HasPrefix(strings.ToLower(m.Content), prefix+f) {
					Uni.GrabRedditPost(s, m, m.Content[len(prefix)+len(f):], nsfw, "/new", "", 0)
					return 1, "Reddit"			
				}
			} 
			
		}

		if (!isDM && modules & 4 == 4) || isDM { // Derpibooru search
			// Mispell detection
			if ( !strings.Contains(strings.ToLower(m.Content), "search") &&
			(strings.Contains(strings.ToLower(m.Content), "on derpi") || 
			strings.Contains(strings.ToLower(m.Content), "on derpibooru"))) {
				Respond(s, m, "Do what to derpibooru?")
				return 0, ""
			}
			
			if strings.Contains(strings.ToLower(m.Content), "seach") || strings.Contains(strings.ToLower(m.Content), "serch") || strings.Contains(strings.ToLower(m.Content), "sarch") || (strings.Contains(strings.ToLower(m.Content), "earch") && !strings.Contains(strings.ToLower(m.Content), "search")) { Respond(s, m, "Hehe you mistyped \"search\" you silly >w<"); return 0, ""; }
			if strings.Contains(strings.ToLower(m.Content), "depi") || strings.Contains(strings.ToLower(m.Content), "deri") { Respond(s, m, "Hehe you mistyped \"derpi\" you silly >w<"); return 0, ""; }
			
			// Ok you seem fine
			
			if strings.HasPrefix(strings.ToLower(m.Content), prefix+" search on derpibooru for ") {
				Uni.SearchOnDerpi(s, m , m.Content[len(prefix)+26:])
				return 1, "Derpi Search"
			} else if strings.HasPrefix(strings.ToLower(m.Content), prefix+" search on derpibooru ") {
				Uni.SearchOnDerpi(s, m , m.Content[len(prefix)+22:])
				return 1, "Derpi Search"
			} else if strings.HasPrefix(strings.ToLower(m.Content), prefix+" search on derpi for ") {
				Uni.SearchOnDerpi(s, m , m.Content[len(prefix)+21:])
				return 1, "Derpi Search"
			} else if strings.HasPrefix(strings.ToLower(m.Content), prefix+" search on derpi ") {
				Uni.SearchOnDerpi(s, m , m.Content[len(prefix)+17:])
				return 1, "Derpi Search"
			} else if strings.HasPrefix(strings.ToLower(m.Content), prefix+" search derpibooru for ") {
				Uni.SearchOnDerpi(s, m , m.Content[len(prefix)+23:])
				return 1, "Derpi Search"
			} else if strings.HasPrefix(strings.ToLower(m.Content), prefix+" search derpibooru ") {
				Uni.SearchOnDerpi(s, m , m.Content[len(prefix)+19:])
				return 1, "Derpi Search"
			} else if strings.HasPrefix(strings.ToLower(m.Content), prefix+" search derpi for ") {
				Uni.SearchOnDerpi(s, m , m.Content[len(prefix)+18:])
				return 1, "Derpi Search"
			} else if strings.HasPrefix(strings.ToLower(m.Content), prefix+" search derpi ") {
				Uni.SearchOnDerpi(s, m , m.Content[len(prefix)+14:])
				return 1, "Derpi Search"
			}  else if strings.HasPrefix(strings.ToLower(m.Content), prefix+" derpi image ") {
				Uni.ImageInfo(s, m, m.Content[len(prefix)+13:])
				return 1, "Derpi Image"
			} else if strings.HasPrefix(strings.ToLower(m.Content), prefix+" derpibooru image ") {
				Uni.ImageInfo(s, m, m.Content[len(prefix)+18:])
				return 1, "Derpi Image"
			}
		}
		if (!isDM && modules & 8 == 8) || isDM { // "Inspire me"
			if strings.HasPrefix(strings.ToLower(m.Content), prefix+" inspire me") {
				Inspire(s, m)
				return 1, "Inspiro"
			}
		}

	}

	
	
	if !isDM { // From server and server only
		if strings.HasPrefix(strings.ToLower(m.Content), prefix+" set my nick from javascript ") {
			dc, err := (strconv.Unquote("\"" + m.Content[len(prefix)+29:] + "\""))
			if err != nil {
				Respond(s, m, fmt.Sprintf("```%s```", err.Error()))
				return 2, "nickname error"
			}
			RenameUser(s, m, g.ID, dc)
			return 1, "Nickname"
		} else if strings.HasPrefix(strings.ToLower(m.Content), prefix+" set my nick from js ") {
			dc, err := (strconv.Unquote("\"" + m.Content[len(prefix)+21:] + "\""))
			if err != nil {
				Respond(s, m, fmt.Sprintf("```%s```", err.Error()))
				return 2, "nickname error"
			}
			RenameUser(s, m, g.ID, dc)
			return 1, "Nickname"
		} else if strings.HasPrefix(strings.ToLower(m.Content), prefix+" set my nick as ") {
			RenameUser(s, m, g.ID, m.Content[len(prefix)+15:])
			return 1, "Nickname"
		} else if strings.HasPrefix(strings.ToLower(m.Content), prefix+" set my nick ") {
			RenameUser(s, m, g.ID, m.Content[len(prefix)+13:])
			return 1, "Nickname"
		} else if strings.HasPrefix(strings.ToLower(m.Content), prefix+" am i admin") {
			if a := Uni.Admin_Detect(s, m); a {
				s.ChannelMessageSend(m.ChannelID, "Ye")
			} else {
				s.ChannelMessageSend(m.ChannelID, "Nope")
			}
			return 1, "Am I admin"
		}
		
		if (
		strings.HasPrefix(strings.ToLower(m.Content), "*hugs") ||
		strings.HasPrefix(strings.ToLower(m.Content), "*boops") || 
		strings.HasPrefix(strings.ToLower(m.Content), "*snugs") ||
		strings.HasPrefix(strings.ToLower(m.Content), "*snuggles")) {
			AssistAction(s, m)
			return 1, "Assist"
		}
		
		
		if a := Uni.Admin_Detect(s, m); a { // Server Admin Commands
			if strings.HasPrefix(strings.ToLower(m.Content), prefix+" set admin role ") {
				Uni.UpdateAdminRole(s, m, g, m.Content[len(prefix)+16:])
				return 1, "Update Admin Role"
			}

			
			// Enable Modules
			if strings.HasPrefix(strings.ToLower(m.Content), prefix+" enable module derpi") {
				Uni.EnableModule(s, m, c.ID, 4, modules)
				return 1, "Enable Module"
			} else if strings.HasPrefix(strings.ToLower(m.Content), prefix+" enable module reddit") {
				Uni.EnableModule(s, m, c.ID, 2, modules)
				return 1, "Enable Module"
			} else if strings.HasPrefix(strings.ToLower(m.Content), prefix+" enable module inspire") {
				Uni.EnableModule(s, m, c.ID, 8, modules)
				return 1, "Enable Module"
			}
			
			
			// Disable Modules
			if strings.HasPrefix(strings.ToLower(m.Content), prefix+" disable module derpi") {
				Uni.DisableModule(s, m, c.ID, 4, modules)
				return 1, "Disable Module"
			} else if strings.HasPrefix(strings.ToLower(m.Content), prefix+" disable module reddit") {
				Uni.DisableModule(s, m, c.ID, 2, modules)
				return 1, "Disable Module"
			} else if strings.HasPrefix(strings.ToLower(m.Content), prefix+" disable module inspire") {
				Uni.DisableModule(s, m, c.ID, 8, modules)
				return 1, "Disable Module"
			}
			
			// Perish/Ban
			if strings.HasPrefix(strings.ToLower(m.Content), prefix+" perish ") {
				a := strings.Split(m.Content, "\n")
				Uni.Perish(s, m, g, a[0][len(prefix)+8:])
				return 1, "Perish"
			} else if strings.HasPrefix(strings.ToLower(m.Content), prefix+" ban ") {
				a := strings.Split(m.Content, "\n")
				Uni.Perish(s, m, g, a[0][len(prefix)+5:])
				return 1, "Perish"
			} else if strings.HasPrefix(strings.ToLower(m.Content), prefix+" purge ") || strings.HasPrefix(strings.ToLower(m.Content), prefix+" clear ") {
				a := m.Content[len(prefix)+7:]
				i, err := strconv.Atoi(a)
				if err != nil {
					Respond(s, m, err.Error())
					return 0, ""
				}
				Uni.Purge(s, m, i)
				return 1, "Purge chat"
			}
					
		
			// Lua things
			if strings.HasPrefix(strings.ToLower(m.Content), prefix+" enable lua") {
				Uni.EnableLua(s, m, g.ID)
				return 1, "Enable Lua"
			}
			
			if strings.HasPrefix(strings.ToLower(m.Content), prefix+" rewrite lua") {
				Uni.RewriteLua(s, m, g.ID)
				return 1, "Rewrite Lua"
			}
			
			if strings.HasPrefix(strings.ToLower(m.Content), prefix+" view lua") {
				Uni.ViewLua(s, m, g.ID)
				return 1, "View Lua"
			}
			
			
			// Admin part of modules
			if modules & 4 == 4 { // Derpi
				if strings.HasPrefix(strings.ToLower(m.Content), prefix+" set derpi filter ") {
					Uni.DerpiSetFilter(s, m, g.ID, c.ID, m.Content[len(prefix)+18:])
					return 1, "Derpi Set Filter"
				}
			}
		}
	}
	
	
	if strings.HasPrefix(strings.ToLower(m.Content), prefix+" give me nsfw") {
		Uni.GiveNSFW(s, m)
		return 1, "Give NSFW"
	} else if strings.HasPrefix(strings.ToLower(m.Content), prefix+" revoke my nsfw") {
		Uni.RevokeNSFW(s, m)
		return 1, "Revoke NSFW"
	}
	
	if strings.ToLower(m.Content) == prefix {
		Respond(s, m, "Hmm?")
		return 1, "Hmm?"
	}
	
	
	
	// For fun
	if strings.HasPrefix(strings.ToLower(m.Content), prefix+" roll ") {
		a := m.Content[len(prefix)+6:]
		d := strings.Split(strings.ToLower(a), "d")
		
		amountofdice := 1
		if len(d) < 2 {
			Respond(s, m, "How much dice?")
			return 0, ""
		}
		if len(d[0]) != 0 {
			amountofdice, _ = strconv.Atoi(d[0])
		}		
		
		dicelimit := 6
		if len(d[1]) != 0 {
			dicelimit, _ = strconv.Atoi(d[1])
		}
		
		if amountofdice == 0 {
			Respond(s, m, "Cannot roll 0 dice")
			return 1, "Roll dice"
		}
		
		if dicelimit == 0 {
			Respond(s, m, "Dice limit cannot be 0")
			return 1, "Roll dice"
		}
		
		// Limit checker
		
		if amountofdice < 0 {
			amountofdice = 2
		}
		
		if amountofdice > 50 {
			listofresponses := []string{
				"Too many things to roll >~<",
				"Das alot of dice OwO",
				"Nu",
				"No",
				"B-but 50 is enough ;~;",
				"Errrr......."}
			Respond(s, m, listofresponses[Uni.Rng.Intn(len(listofresponses))])
			return 0, ""
		}
		
		dice := []int64{}
		
		for i := 0; i < amountofdice; i++ {
			dice = append(dice, Uni.Rng.Int63n(int64(dicelimit)))
		}
		
		r := ""
		total := big.NewInt(0)
		for index, dice := range dice {
			if index != 0 {
				r = fmt.Sprintf("%s, ", r)
			}
			r = fmt.Sprintf("%s%d", r, dice+1)
			//total += int64(dice+1)
			total.Add(total, big.NewInt(dice+1))
		}
		if amountofdice < 15 {
			Respond(s, m, fmt.Sprintf("Dice rolled: (%s)\nTotal: %d", r, total))
		} else {
			Respond(s, m, fmt.Sprintf("Total: %d", total))
		}
		return 1, "Roll dice"
	}
	
	
	return 0, ""
	
}




// Get server ID of a message
func getGuild(m *discordgo.MessageCreate, s *discordgo.Session) (*discordgo.Guild, error) {
	c, err := s.State.Channel(m.ChannelID)
	if err != nil {
		// Could not find channel.
		return nil, err
	}
	g, err := s.State.Guild(c.GuildID)
	if err != nil {
		// Unable to get server
		return nil, err
	}
	return g, nil
}






func AssistAction(s *discordgo.Session, m *discordgo.MessageCreate) {
	Respond(s, m, fmt.Sprintf("*also %s", m.Content[1:]))
}


func ErrRespond(session *discordgo.Session, message *discordgo.MessageCreate, response string) {
	t := make([]byte, 10)
	for i := range t {
		t[i] = byte(rand.Intn(256))
	}
	elh := fmt.Sprintf("%x", sha256.Sum256(t))
	var buf bytes.Buffer
	enc := toml.NewEncoder(&buf)
	_ = enc.Encode(message)
	_ = ioutil.WriteFile(fmt.Sprintf("errlog/%s.inf", elh), buf.Bytes(), 0644)
	session.ChannelMessageSend("371838361658195969", fmt.Sprintf("%s\n\nErr log message file: %s", response, elh))
	Respond(session, message, response)
}



func Respond(s *discordgo.Session, m *discordgo.MessageCreate, response string) {
	go s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("%s", response))
}


func Cheer(s *discordgo.Session, m *discordgo.MessageCreate) {
	cq := []string{"Okie ^^", "Will do ^^", "Tada!", "Hehe sure thing ^w^", "*squee*"}
	Respond(s, m, cq[rand.Intn(len(cq))])
}


func (Uni *UniBot) CheckNSFW(session *discordgo.Session, message *discordgo.MessageCreate) bool {
	
	rows, err := Uni.Database.Query("select * from NSFWList")
	if err != nil {
		ErrRespond(session, message, fmt.Sprintf(" (Query failure while checking NSFW List) Error occurred; %s", err))
		return false
	}
	defer rows.Close()
	for rows.Next() {
		var id string
		err = rows.Scan(&id)
		if err != nil {
			ErrRespond(session, message, fmt.Sprintf(" (NSFW userlist row scan failure) Error occurred; %s", err))
			return false
		}
		
		if id == message.Author.ID {
			return true
		}
	}
	return false
}

func (Uni *UniBot) GiveNSFW(s *discordgo.Session, m *discordgo.MessageCreate) {
	if a := Uni.CheckNSFW(s, m); a {
		Respond(s, m, "You seem to already have NSFW")
	} else {
		Uni.Database.Exec(fmt.Sprintf("INSERT INTO NSFWList VALUES (%s)", m.Author.ID))
		Respond(s, m, "You now have the ability to summon NSFW with searches in my DMs")
	}
}


func (Uni *UniBot) RevokeNSFW(s *discordgo.Session, m *discordgo.MessageCreate) {
	if a := Uni.CheckNSFW(s, m); a {
		Uni.Database.Exec(fmt.Sprintf("DELETE FROM NSFWList WHERE userid IS '%s'", m.Author.ID))
		Respond(s, m, "Your NSFW is revoked, You no longer have the ability to access NSFW in DMs")
	} else {
		Respond(s, m, "You seem to not have NSFW")
	}
}


func RenameUser(s *discordgo.Session, m *discordgo.MessageCreate, gID, nick string) {
	
	err := s.GuildMemberNickname(gID, m.Author.ID, nick)
	if err != nil {
		Respond(s, m, fmt.Sprintf("```%s```", err.Error()))
		return
	}
	Cheer(s, m)
}

