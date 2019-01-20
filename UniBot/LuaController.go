package Uni
import (
	"fmt"
	"os"
	"context"
	"time"
	"strings"
	"io/ioutil"
	"./UniLua"//"github.com/yuin/gopher-lua"
	"net/http"
	"github.com/bwmarrin/discordgo"
)

func (Uni *UniBot) EnableLua(s *discordgo.Session, m *discordgo.MessageCreate, gID string) {
	cld := fmt.Sprintf("%s/%s", Uni.LuaDir, gID)
	_ = os.MkdirAll(cld, 0755)
	_ = ioutil.WriteFile(cld+"/main.lua", []byte{}, 0755)
	Respond(s, m, "I have prepared lua script for this server")
}


func (Uni *UniBot) RewriteLua(s *discordgo.Session, m *discordgo.MessageCreate, gID string) {
	_, err := s.Guild(gID)
	if err != nil {
		Respond(s, m, err.Error())
		return
	}
	
	
	a := strings.Split(m.Content, "\n")
	lc := ""
	for i, line := range a {
		if !(i == 0 ||
		strings.HasPrefix(line, "```")) {
			lc = fmt.Sprintf("%s%s\n", lc, line)
		}
	}
	
	err = ioutil.WriteFile(fmt.Sprintf("%s/%s/main.lua", Uni.LuaDir, gID), []byte(lc), 0444)
	if err != nil {
		Respond(s, m, err.Error())
		return
	}
	
	Respond(s, m, "Lua script modified")
	
}

func (Uni *UniBot) ViewLua(s *discordgo.Session, m *discordgo.MessageCreate, gID string) {
	lc, err := ioutil.ReadFile(fmt.Sprintf("%s/%s/main.lua", Uni.LuaDir, gID))
	if err != nil {
		Respond(s, m, err.Error())
		return
	}
	Respond(s, m, fmt.Sprintf("```lua\n%s```", string(lc)))
}

func (Uni *UniBot) LPrint(L *lua.LState) int {
	fmt.Println(L.UniVars)
	t := L.NewTable()
	t.RawSet(lua.LString("type"), lua.LString("MESSAGE_SEND"))
	t.RawSet(lua.LString("channelID"), lua.LString(L.UniVars.ChannelID))
	t.RawSet(lua.LString("message"), L.Get(1))
	c := make(chan lua.LValue)
	l := &LC{LV: t, RC: c}
	Uni.RC <- *l
	r := <-c
	if r != nil {
		L.Push(r)
	} else {
		L.Push(lua.LNil)
	}
	return 1
}

func (Uni *UniBot) LChannelPrint(L *lua.LState) int {
	t := L.NewTable()
	t.RawSet(lua.LString("type"), lua.LString("MESSAGE_SEND"))
	t.RawSet(lua.LString("channelID"), lua.LString(L.ToString(1)))
	t.RawSet(lua.LString("message"), lua.LString(L.ToString(2)))
	c := make(chan lua.LValue)
	l := &LC{LV: t, RC: c}
	Uni.RC <- *l
	r := <-c
	if r != nil {
		L.Push(r)
	} else {
		L.Push(lua.LNil)
	}
	return 1
}



func (Uni *UniBot) LMessageDelete(L *lua.LState) int {
	t := L.NewTable()
	t.RawSet(lua.LString("type"), lua.LString("MESSAGE_DELETE"))
	t.RawSet(lua.LString("channelID"), lua.LString(L.UniVars.ChannelID))
	t.RawSet(lua.LString("ID"), lua.LString(L.UniVars.ID))
	c := make(chan lua.LValue)
	l := &LC{LV: t, RC: c}
	Uni.RC <- *l
	r := <-c
	if r != nil {
		L.Push(r)
	} else {
		L.Push(lua.LNil)
	}
	return 1
}

func LStringFind(L *lua.LState) int {
	a := L.ToString(1)
	b := L.ToString(2)
	L.Push(lua.LNumber(strings.Index(a, b)))
	return 1
}

func LStringLower(L *lua.LState) int {
	a := L.ToString(1)
	L.Push(lua.LString(strings.ToLower(a)))
	return 1
}

func LNilFunction(L *lua.LState) int {
	L.Push(lua.LNumber(1))
	return 1
}



func (Uni *UniBot) LuaHTTPGet(L *lua.LState) int {
	t := L.NewTable()
	h := L.NewTable()
	t.RawSet(lua.LString("type"), lua.LString("HTTP_GET"))
	t.RawSet(lua.LString("channelID"), lua.LString(L.UniVars.ChannelID))
	t.RawSet(lua.LString("link"), lua.LString(L.ToString(1)))
	t.RawSet(lua.LString("httptable"), h)
	c := make(chan lua.LValue)
	l := &LC{LV: t, RC: c}
	Uni.RC <- *l
	r := <-c
	if r != nil {
		L.Push(r)
	} else {
		L.Push(lua.LNil)
	}
	return 1
}



func (Uni *UniBot) ParseServerLua(s *discordgo.Session, m *discordgo.MessageCreate, gID string) {
	t := time.Now()
	L := lua.NewState(lua.Options{CallStackSize: 14, RegistrySize:  14*20,})
	
	// Hardwiring the variables
	L.UniVars.ChannelID = m.ChannelID
	L.UniVars.GuildID = gID
	L.UniVars.ID = m.ID
	
	L.SetGlobal("ID", lua.LString(m.ID))
	L.SetGlobal("content", lua.LString(m.Content))
	L.SetGlobal("channelID", lua.LString(m.ChannelID))
	L.SetGlobal("timestamp", lua.LString(fmt.Sprintf("%s", m.Timestamp)))
	L.SetGlobal("author_ID", lua.LString(m.Author.ID))
	L.SetGlobal("author_username", lua.LString(m.Author.Username))
	L.SetGlobal("author_discriminator", lua.LString(m.Author.Discriminator))
	L.SetGlobal("author_bot", lua.LBool(m.Author.Bot))
	
	L.SetGlobal("print", L.NewFunction(Uni.LPrint))
	L.SetGlobal("channelprint", L.NewFunction(Uni.LChannelPrint))
	L.SetGlobal("string.find", L.NewFunction(LStringFind))
	L.SetGlobal("string.lower", L.NewFunction(LStringLower))
	L.SetGlobal("delete", L.NewFunction(Uni.LMessageDelete))
	
	
	L.SetGlobal("http_get", L.NewFunction(Uni.LuaHTTPGet))

	ctx, cancel := context.WithTimeout(context.Background(), 7*time.Second)
	L.SetContext(ctx)
	fmt.Println("Time took to create lua state machine for ID (", m.ID, "): ", time.Since(t))
	t = time.Now()
	err := L.DoFile(fmt.Sprintf("%s/%s/main.lua", Uni.LuaDir, gID))
	fmt.Println("Time took to parse lua for ID (", m.ID, "): ", time.Since(t))
	if err != nil {
		Respond(s, m, fmt.Sprintf("Lua state machine error: \n```%s```", err))
	}
	cancel()
	L.Close()
}



func (Uni *UniBot) LuaReact(s *discordgo.Session, rc chan LC) {
	for {
		a := <- rc
		
		if Uni.APIPressure[fmt.Sprintf("%s", a.LV.(*lua.LTable).RawGet(lua.LString("channelID")))] == 0 {
			Uni.APIPressure[fmt.Sprintf("%s", a.LV.(*lua.LTable).RawGet(lua.LString("channelID")))] = .5
		} else {
			Uni.APIPressure[fmt.Sprintf("%s", a.LV.(*lua.LTable).RawGet(lua.LString("channelID")))] *= 1.5
		}
		fmt.Println(Uni.APIPressure[fmt.Sprintf("%s", a.LV.(*lua.LTable).RawGet(lua.LString("channelID")))])
		
		go func() {
			time.Sleep(time.Duration(Uni.APIPressure[fmt.Sprintf("%s", a.LV.(*lua.LTable).RawGet(lua.LString("channelID")))]) * time.Second)
			switch fmt.Sprintf("%s", a.LV.(*lua.LTable).RawGet(lua.LString("type"))) {
			case "MESSAGE_SEND":
				_, err := s.ChannelMessageSend(fmt.Sprintf("%s", a.LV.(*lua.LTable).RawGet(lua.LString("channelID"))), fmt.Sprintf("%s", a.LV.(*lua.LTable).RawGet(lua.LString("message"))))
				if err != nil {	
					a.RC <- lua.LString(err.Error())
					return
				}
			case "MESSAGE_DELETE":
				err := s.ChannelMessageDelete(fmt.Sprintf("%s", a.LV.(*lua.LTable).RawGet(lua.LString("channelID"))), fmt.Sprintf("%s", a.LV.(*lua.LTable).RawGet(lua.LString("ID"))))
				if err != nil {	
					a.RC <- lua.LString(err.Error())
					return
				}
			case "HTTP_GET":
				LuaHandleHTTPGet(a.RC, a.LV.(*lua.LTable).RawGet(lua.LString("httptable")).(*lua.LTable), fmt.Sprintf("%s", a.LV.(*lua.LTable).RawGet(lua.LString("link"))))
				return
			}
			a.RC <- nil
		}()
	}
}


func LuaHandleHTTPGet(a chan lua.LValue, h *lua.LTable, link string) {
	tmpClient := &http.Client{Timeout: time.Second * 2,}
	var (
		err error
		req *http.Request
		res *http.Response
		resp []byte = []byte{}
	)
	req, _ = http.NewRequest(http.MethodGet, link, nil)
	
	res, err = tmpClient.Do(req)

	if err != nil {
		goto finish
	}
		
	resp, err = ioutil.ReadAll(res.Body)
	if err != nil {
		goto finish
	}
	
	finish:
	
	if res == nil {
		h.RawSet(lua.LString("Status"), lua.LNil)
		h.RawSet(lua.LString("StatusCode"), lua.LNil)
	} else {
		h.RawSet(lua.LString("Status"), lua.LString(res.Status))
		h.RawSet(lua.LString("StatusCode"), lua.LNumber(res.StatusCode))
	}
	if err != nil {
		h.RawSet(lua.LString("error"), lua.LString(err.Error()))
		h.RawSet(lua.LString("resp"), lua.LNil)
	} else {
		h.RawSet(lua.LString("error"), lua.LNil)
		h.RawSet(lua.LString("resp"), lua.LString(resp))
	}
	a <- h
}
