package Uni
import "fmt"
var moduleslist []string = []string{"reddit", "derpi", "inspire", "minigames", "unibucks", "miscellaneous", "OwO"}

// Set module flag (setmodule = false means disable and setmodule = true means enable)
func (Uni *UniBot) ControlModule(setmodule bool, cID string, modulesID int, modules uint64) {
	var m uint64 = 2
	for a := 0; a < modulesID; a++ { // 2^modulesID
		m *= 2
	}
	
	if modules & m == m { // module found
		if setmodule { // enable
			Uni.Respond(cID, fmt.Sprintf("You appear to already have %q enabled", moduleslist[modulesID]))
		} else { // disable
			modules ^= m
			_, err := Uni.DBExec("UpdateChannelModules", modules, cID)
			if err == nil { // phew, nothing bad happened
				Uni.Respond(cID, fmt.Sprintf("%q is now disabled", moduleslist[modulesID]))
			} else { // uh oh something bad happened
				Uni.ErrRespond(err, cID, "updating modules variable", map[string]interface{}{"setmodule": setmodule, "cID": cID, "modulesID": modulesID, "modules": modules})
			}
		}
	} else { // module not found
		if setmodule { // enable
			modules |= m
			_, err := Uni.DBExec("UpdateChannelModules", modules, cID)
			if err == nil { // phew, nothing bad happened
				Uni.Respond(cID, fmt.Sprintf("%q is now enabled", moduleslist[modulesID]))
			} else { // uh oh something bad happened
				Uni.ErrRespond(err, cID, "updating modules variable", map[string]interface{}{"setmodule": setmodule, "cID": cID, "modulesID": modulesID, "modules": modules})
			}
		} else { // disable
			Uni.Respond(cID, fmt.Sprintf("You appear to already have %q disabled", moduleslist[modulesID]))
		}
	}
}