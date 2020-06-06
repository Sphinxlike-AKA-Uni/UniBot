package Uni

import (
	"fmt"
	"strings"
	"github.com/bwmarrin/discordgo"
	"syscall"
	"time"
	"os"
)



func (Uni *UniBot) ProcessMessage(m *discordgo.MessageCreate, isDM bool, g *discordgo.Guild , c *discordgo.Channel, modules uint64, name, prefix string) (int, string) {
	// Creator Only Commands
	if m.Author.ID == Uni.Config.CreatorID {
		if strings.HasPrefix(m.Content, "$~restart") {
			Uni.messageGate = true
			for Uni.MessageHandlers > 1 {
				time.Sleep(444 * time.Millisecond)
			}
			Uni.Restart = true
			Uni.SC <- syscall.SIGTERM // force terminate self
		} else if strings.HasPrefix(m.Content, "$~shutdown") {
			Uni.messageGate = true
			for Uni.MessageHandlers > 1 {
				time.Sleep(444 * time.Millisecond)
			}
			Uni.SC <- syscall.SIGTERM // force terminate self
		} else if strings.HasPrefix(m.Content, "$~reply ") { // *squee*
			content := m.Content[8:]
			a := strings.Split(content, "|")
			Uni.Respond(a[0], a[1])
		}
	}


	if strings.HasPrefix(strings.ToLower(m.Content), prefix) {
		// Want help? (i'm sorry if the README is not good ^^")
		if strings.HasPrefix(strings.ToLower(m.Content), prefix+" help") ||
		strings.HasPrefix(strings.ToLower(m.Content), prefix+" halp") ||
		strings.HasPrefix(strings.ToLower(m.Content), prefix+" hlep") {
			Uni.Respond(m.ChannelID, "https://github.com/Sphinxlike-AKA-Uni/UniBot/blob/master/README.md")
			return 1, "Help"
		}
		
		// Current version of uni?
		if strings.HasPrefix(strings.ToLower(m.Content), prefix+" ver") {
			// Semantic Versioning https://semver.org/
			Uni.Respond(m.ChannelID, fmt.Sprintf("`%s`", versionstring))
		}
		
		// Module stuff begin here
		if (!isDM && modules & 2 == 2) || isDM { // Reddit search
			var nsfw bool
			if isDM {
				// did user request for nsfw?
				nsfwstr := ""
				Uni.DBGetFirstVar(&nsfwstr, "CheckNSFW", m.Author.ID)
				if nsfwstr != "" { nsfw = true } // user index was returned
			} else {
				nsfw = c.NSFW // channel specified NSFW tag
			}
			
			if strings.HasPrefix(strings.ToLower(m.Content), prefix+" find a ") { // skip if not close to reddit command
				// for loops to prevent a bunch of "else ifs"
				for _, RedditType := range []string{"", "top ", "new "} {
				for _, Adverb := range []string{"on ", "in ", "from ", "within ", ""} {
					for _, RedditFormat := range []string{"", "r/", "/r/"} {
						rs := fmt.Sprintf("%s find a %spost %s%s", prefix, RedditType, Adverb, RedditFormat)
						if strings.HasPrefix(strings.ToLower(m.Content), rs) {
							Uni.GrabRedditPost(m.ChannelID, m.Content[len(rs):], nsfw, RedditType)
							return 1, "Reddit"
						}
					}
				}
				}
			}
			
		}
		
		if (!isDM && modules & 4 == 4) || isDM { // Derpibooru search
			for _, f := range []string{
			" search on derpibooru for ", " search on derpibooru ", " search on derpi for ", " search on derpi ",
			" search derpibooru for ", " search derpibooru ", " search derpi for ", " search derpi ",
			} {
				if strings.HasPrefix(strings.ToLower(m.Content), prefix+f) {
					Uni.SearchOnDerpi(c.ID, m.Content[len(prefix)+len(f):])
					return 1, "Derpi Search"
				}
			}
			
			if strings.HasPrefix(strings.ToLower(m.Content), prefix+" derpi total ") {
				f, err := Uni.GetChannelDerpiFilter(c.ID)
				Uni.ErrRespond(err, c.ID, "grabbing channel derpi filter", map[string]interface{}{"err": err, "cID": c.ID})
				s, err := Uni.DerpiSearch(m.Content[len(prefix)+13:], f, map[string]interface{}{"per_page": 1})
				Uni.ErrRespond(err, c.ID, "acquiring derpi search", map[string]interface{}{"err": err, "cID": c.ID})
				Uni.Respond(c.ID, fmt.Sprintf("`%s` returned with %d image%s", m.Content[len(prefix)+13:], s.Total, map[bool]string{true: "s", false: ""}[s.Total != 1]))
				return 1, "Derpi Total"
			}
		}
		
		if (!isDM && modules & 8 == 8) || isDM { // "Inspire me"
			if strings.HasPrefix(strings.ToLower(m.Content), prefix+" inspire me") {
				Uni.Inspire(m.ChannelID)
				return 1, "Inspiro"
			}
		}
		
		if (!isDM && modules & 16 == 16) || isDM { // Minigames
			// UNO
			if strings.HasPrefix(strings.ToLower(m.Content), prefix+" start uno") {
				/*uf, err := os.OpenFile(fmt.Sprintf("%s/uno_%s", Uni.TempDir, m.Author.ID), os.O_RDWR, 0644) // UNO File
				if !(os.IsNotExist(err)) { // reaches here if file exists
					Uni.UNOCreate(cID, m.Author.ID)
				} else {
					Uni.Respond(cID, "There is currently an UNO game going on, if you wish for more than one UNO game to start I would recommend having another text channel.")
				}*/
				/*
bf, err := os.OpenFile(fmt.Sprintf("%s/b_%s", Uni.TempDir, uID), os.O_RDWR, 0644) // Blackjack File
	if !(os.IsNotExist(err)) { // reaches here if file exists
		Uni.Respond(cID, "You appear to already have a game in session, here lemme show you the cards")
		return true
	}
				*/
			}
			
		}
		
		if (!isDM && modules & 32 == 32) || isDM { // Uni Bucks Minigames
			Uni.CheckIfProfileExists(c.ID, m.Author.ID) // Check if user has Uni Bucks( will create one if they don't)
			if strings.HasPrefix(strings.ToLower(m.Content), prefix+" bank") ||
			strings.HasPrefix(strings.ToLower(m.Content), prefix+" balance") || 
			strings.HasPrefix(strings.ToLower(m.Content), prefix+" wallet") {
				u, err := Uni.DBGetFirst("GrabUniBucks", m.Author.ID)
				Uni.ErrRespond(err, c.ID, "grabbing unibucks profile", map[string]interface{}{"err": err, "cID": c.ID})
				Uni.Respond(c.ID, fmt.Sprintf("**%s's Uni Bucks: %.2f**", name, u))
				return 1, "Uni Bucks Bank"
			}
			
			if strings.HasPrefix(strings.ToLower(m.Content), prefix+" slot roll") {
				Uni.SlotRoll(c.ID, m.Author.ID, name)
				return 1, "Slot Roll"
			} else if strings.HasPrefix(strings.ToLower(m.Content), prefix+" play blackjack ") {
				if Uni.StartBlackjack(m.Author.ID, c.ID, m.Content[len(prefix)+16:]) {
					Uni.RepresentCards(c.ID, m.Author.ID, name)
				}
				return 1, "Blackjack Start"
			} else if strings.HasPrefix(strings.ToLower(m.Content), prefix+" daily") {
				Uni.Daily(c.ID, m.Author.ID, name)
				return 1, "UniBucks Daily"
			}
			
			if _, err := os.Stat(fmt.Sprintf("%s/b_%s", Uni.TempDir, m.Author.ID)); !os.IsNotExist(err) { // Can't hit or stay without actually playing blackjack silly
				if strings.HasPrefix(strings.ToLower(m.Content), prefix+" hit") {
					Uni.BlackjackHit(c.ID, m.Author.ID)
					Uni.RepresentCards(c.ID, m.Author.ID, name)
					return 1, "Blackjack Hit"
				} else if strings.HasPrefix(strings.ToLower(m.Content), prefix+" stay") {
					Uni.BlackjackStay(c.ID, m.Author.ID, name)
					return 1, "Blackjack Stay"
				}
			}
			
			
		}
		// Module stuff end here
		
		// Both DM and server commands
		if strings.HasPrefix(strings.ToLower(m.Content), prefix+" give nsfw") ||
		strings.HasPrefix(strings.ToLower(m.Content), prefix+" give me nsfw") {
			_, err := Uni.DBExec("GiveNSFW", m.Author.ID)
			if err != nil {
				Uni.ErrRespond(err, c.ID, "putting user in NSFWList", map[string]interface{}{"err": err, "cID": c.ID,})
			} else {
				Uni.Respond(c.ID, "You now have the ability to summon NSFW in DMs")
			}
			return 1, "Give NSFW"
		} else if strings.HasPrefix(strings.ToLower(m.Content), prefix+" revoke my nsfw") {
			_, err := Uni.DBExec("RevokeNSFW", m.Author.ID)
			if err != nil {
				Uni.ErrRespond(err, c.ID, "removing user in NSFWList", map[string]interface{}{"err": err, "cID": c.ID,})
			} else {
				Uni.Respond(c.ID, "You now no longer have the ability to summon NSFW in DMs")
			}
			return 1, "Revoke NSFW"
		}
	
		
		// Admin commands
		if Uni.Admin_Detect(m.ChannelID, m.Author.ID, g) {
			// Enabling/Disabling of a module
			if strings.HasPrefix(strings.ToLower(m.Content), prefix+" enable module ") ||
			strings.HasPrefix(strings.ToLower(m.Content), prefix+" disable module ") {
				var validity bool
				for index, module := range moduleslist {
					if strings.HasPrefix(strings.ToLower(m.Content), prefix+" enable module "+module) {
						Uni.ControlModule(true, m.ChannelID, index, modules)
						validity = true
					} else if strings.HasPrefix(strings.ToLower(m.Content), prefix+" disable module "+module) {
						Uni.ControlModule(false, m.ChannelID, index, modules)
						validity = true
					}
				}
				if !validity {
					Uni.Respond(m.ChannelID, "Module does not appear to be valid")
				}
			}
			
			// Set Derpi Filter
			if modules & 4 == 4 { // Derpi
				if strings.HasPrefix(strings.ToLower(m.Content), prefix+" set derpi filter ") {
					Uni.SetChannelDerpiFilter(g.ID, c.ID, m.Content[len(prefix)+18:])
					return 1, "Derpi Set Filter"
				}
			}
		}
		
		// Hmm?
		if strings.ToLower(m.Content) == prefix {
			Uni.Respond(c.ID, "Hmm?")
			return 1, "Hmm?"
		}
	}
	
	
	
	// Playful assistance
	if (
	strings.HasPrefix(strings.ToLower(m.Content), "*hugs") ||
	strings.HasPrefix(strings.ToLower(m.Content), "*boops") || 
	strings.HasPrefix(strings.ToLower(m.Content), "*snugs") ||
	strings.HasPrefix(strings.ToLower(m.Content), "*snuggles")) {
		Uni.Respond(m.ChannelID, fmt.Sprintf("*also %s", m.Content[1:]))
		return 1, "Assist"
	}
	
	
	// nothing happened
	return 0, ""
}