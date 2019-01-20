package Uni

import (
	"io/ioutil"
	"net/http"
	"github.com/bwmarrin/discordgo"
)





func Inspire(s *discordgo.Session, m *discordgo.MessageCreate) {
	tmpClient := http.Client{}
	
	req, _ := http.NewRequest(http.MethodGet, "http://inspirobot.me/api?generate=true", nil)
	req.Header.Set("User-Agent", "Uni_Inspiration")
	res, _ := tmpClient.Do(req)
	
	inspiration, _ := ioutil.ReadAll(res.Body)
	Respond(s, m, string(inspiration))
}