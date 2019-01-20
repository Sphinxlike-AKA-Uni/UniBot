package Uni

import (
	"fmt"
	"github.com/bwmarrin/discordgo"

)


func (Uni *UniBot) EnableModule(s *discordgo.Session, m *discordgo.MessageCreate, cID string, which, modules int) {
	if which == 2 {
		if modules & 2 == 2 {// Already has reddit search module
			Respond(s, m, "You seem to already have the reddit search module enabled")
			return
		} else { // Give reddit search module
			modules = modules | 2
			Respond(s, m, "Reddit search module enabled")
		}
	}
	if which == 4 {
		if modules & 4 == 4 {// Already has derpibooru search module
			Respond(s, m, "You seem to already have the derpibooru search module enabled")
			return
		} else { // Give derpibooru search module
			modules = modules | 4
			Respond(s, m, "Derpibooru search module enabled")
		}
	}
	if which == 8 {
		if modules & 8 == 8 {// Already has inspire module
			Respond(s, m, "You seem to already have the inspiration module enabled")
			return
		} else {
			modules = modules | 8
			Respond(s, m, "Inspiration module enabled")
		}
	}
	
	
	Uni.Database.Exec(fmt.Sprintf("UPDATE Modules SET modules = %d WHERE cID = '%s';", modules, cID))
}

func (Uni *UniBot) DisableModule(s *discordgo.Session, m *discordgo.MessageCreate, cID string, which, modules int) {
	if which == 2 {
		if modules & 2 == 2 {// Already has reddit search module
			modules = modules ^ 2
			Respond(s, m, "Reddit search module disabled")
		} else { // Revoke reddit search module
			Respond(s, m, "You seem to already have the reddit search module disabled")
			return
		}
	}
	if which == 4 {
		if modules & 4 == 4 {// Already has derpibooru search module
			modules = modules ^ 4
			Respond(s, m, "Derpibooru search module disabled")
		} else { // Revoke derpibooru search module
			Respond(s, m, "You seem to already have the derpibooru search module disabled")
			return
		}
	}
	if which == 8 {
		if modules & 8 == 8 {// Already has inspire module
			modules = modules ^ 8
			Respond(s, m, "Inspiration module disabled")
		} else {
			Respond(s, m, "You seem to already have the inspiration module disabled")
			return
		}
	}
	Uni.Database.Exec(fmt.Sprintf("UPDATE Modules SET modules = %d WHERE cID = '%s';", modules, cID))
}