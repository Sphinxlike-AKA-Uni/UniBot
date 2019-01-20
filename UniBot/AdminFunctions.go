package Uni
import (
	"fmt"
	"strings"
	"github.com/bwmarrin/discordgo"
)




// Update An Admin Role From A Server
func (Uni *UniBot) UpdateAdminRole(s *discordgo.Session, m *discordgo.MessageCreate, g *discordgo.Guild, c string) {
	if m.Author.ID != g.OwnerID {return}
	var rr []*discordgo.Role
	for _, role := range g.Roles {
		if strings.Contains(strings.ToLower(role.Name), strings.ToLower(c)) || strings.Contains(role.ID, c) {
			rr = append(rr, role)
		}
		
	}
	
	if len(rr) > 1 {
		rt := "Returned roles:```\n"
		for _, role := range rr {
			rt = rt+fmt.Sprintf("%s,  %s\n", role.ID, role.Name)
		}
		Respond(s, m, rt+"```")
	} else if len(rr) == 0 {
		Respond(s, m, fmt.Sprintf("No roles found for \"%s\"", c))
	} else {
		Uni.Database.Exec(fmt.Sprintf("UPDATE ServerData SET adminrole = \"%s\" WHERE id = \"%s\"", rr[0].ID, g.ID))
		Respond(s, m, fmt.Sprintf("Admin role has now been set to \"%s\" \\<@%s\\>", rr[0].Name, rr[0].ID,))
	}
	
}



func (Uni *UniBot) Admin_Detect(session *discordgo.Session, message *discordgo.MessageCreate) bool {
	g, err := getGuild(message, session)
	if err != nil {
		session.ChannelMessageSend(message.ChannelID, "Error receiving guild assuming you are not admin")
		session.ChannelMessageSend("371838361658195969", "Message ID: "+message.ID+"\n"+err.Error())
		return false
	}
	if message.Author.ID == g.OwnerID || message.Author.ID == Uni.CreatorID {
		return true
	}
	rows, err := Uni.Database.Query(fmt.Sprintf("SELECT adminrole FROM ServerData WHERE id = %s", g.ID))
	if err != nil {
		ErrRespond(session, message, fmt.Sprintf(" (Query failure while checking adminroles) %s", err))
	}
	defer rows.Close()
	for rows.Next() {
		adminrole := ""
		err = rows.Scan(&adminrole)
		if err != nil {
			return false
		}
		
	    for _, member := range g.Members {
			if member.User.ID == message.Author.ID {
				for _, rolestr := range member.Roles {
					if rolestr == adminrole { return true }
				}
			}
		}
	}
	return false
	
}


func (Uni *UniBot) Perish(s *discordgo.Session, m *discordgo.MessageCreate, g *discordgo.Guild, ID string) {	
	var ul []*discordgo.User
	if len(m.Mentions) > 0 {
		ul = m.Mentions
	} else {
		for _, user := range g.Members {
			if strings.Contains(strings.ToLower(user.User.Username), strings.ToLower(ID)) || strings.Contains(strings.ToLower(user.Nick), strings.ToLower(ID)) || strings.Contains(user.User.ID, ID) {
				ul = append(ul, user.User)
			}
		}
	}
	
	if len(ul) == 1 {
		a := strings.Split(m.Content, "\n")
		reason := ""
		if len(a) == 2 {
			reason = a[1]
		}
		err := s.GuildBanCreateWithReason(g.ID, ul[0].ID, reason, 7)
		if err != nil {
			Respond(s, m, fmt.Sprintf("```%s```", err))
		} else {
			Cheer(s, m)
		}
	} else if len(ul) > 1 {
		rt := "Returned users:```\n"
		for _,  user := range ul {
			rt = rt+fmt.Sprintf("%s,  %s#%s\n", user.ID, user.Username, user.Discriminator)
		}
		Respond(s, m, rt+"```")
	} else {
		Respond(s, m, fmt.Sprintf("No users found for \"%s\"", ID))
	}
}

func (Uni *UniBot) Purge(s *discordgo.Session, m *discordgo.MessageCreate, amount int) {	
	for {
		cm, err := s.ChannelMessages(m.ChannelID, 100, "", "", "")
		if err != nil {
			Respond(s, m, err.Error())
			return
		}
		
		if len(cm) == 0 || amount == 0 {
			return
		}
		
		/*for _, message := range cm {
			amount -= 1
			if amount == 0 {
				return
			}
			err := s.ChannelMessageDelete(m.ChannelID, message.ID)
			if err != nil {
				fmt.Println(err)
				return
			}
		}*/
		mIDs := []string{}
		for _, message := range cm {
			amount -= 1
			if amount == 0 {
				break
			}
			mIDs = append(mIDs, message.ID)
		}
		
		s.ChannelMessagesBulkDelete(m.ChannelID, mIDs)
		
		fmt.Println(amount)
	}
}
