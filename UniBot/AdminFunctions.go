package Uni
import (
	"github.com/bwmarrin/discordgo"
	//"fmt"
)

// is user an admin in server?
func (Uni *UniBot) Admin_Detect(cID, uID string, g *discordgo.Guild) bool {

	if uID == g.OwnerID { // user is guild owner
		return true
	}
	
	var s string
	Uni.DBGetFirstVar(&s, "GetGuildAdminRole", g.ID)
	if s != "" { // something did get returned
		u, err := Uni.S.GuildMember(g.ID, uID)
		if err != nil { // an error occurred and it automatically assumes false
			Uni.ErrRespond(err, cID, "grabbing guild member", map[string]string{"cID": cID, "uID": uID}, g)
			goto End
		}
		// everything is fine
		for _, role := range u.Roles {
			if role == s { return true } // a detected match
		}
		
		// if reached here, no match
	}
	
	End:
	return false
}