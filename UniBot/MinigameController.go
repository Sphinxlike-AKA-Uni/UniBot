package Uni

import (
	//"fmt"
	"strings"
	"./Minigames/sweeper"
	"github.com/bwmarrin/discordgo"
)

func (Uni *UniBot) PlaySweep(s *discordgo.Session, m *discordgo.MessageCreate) {
	//fmt.Println(Sweeper.SweepWithRNG(10, 10, 5, Uni.Rng))
	Respond(s, m, strings.Replace(Sweeper.SweepWithRNG(10, 10, 5, Uni.Rng), "B", "<:bomb:>", -1))
}