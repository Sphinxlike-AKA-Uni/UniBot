// Primary UniBot script
package Uni

import (
	"github.com/bwmarrin/discordgo"
	"os"
	"database/sql"
	crand "crypto/rand"
	"fmt"
	"encoding/binary"
	"github.com/BurntSushi/toml"
	"github.com/davecgh/go-spew/spew"
	"io"
	"time"
	"runtime"
)
// Decided to use roman numerals for versions
var versionstring string = "I.I.III"

type UniBot struct {
	S	*discordgo.Session
	Self	*discordgo.User
	Debug	uint8
	DB	*sql.DB
	Config 	UniConfig
	RNGChan	chan uint64
	messageGate bool
	MessageHandlers uint
	SC	chan os.Signal
	Restart	bool
	TempDir string
}

type UniConfig struct {
	Token 	  string
	CreatorID string
	DBDriver  string
	DBContent string
	ErrLogChannel string
}

func New() UniBot { // Creates a new Uni Bot
	var Uni UniBot
	Uni.SC = make(chan os.Signal, 1)
	Uni.RNGChan = make(chan uint64, 16)
	Uni.TempDir = os.TempDir()+"/UniBot"
	return Uni
}

func (Uni *UniBot) Startup(configfile string) error { // Start up the Uni Bot
	var err error
	
	// Decode Uni Config
	_, err = toml.DecodeFile(configfile, &Uni.Config)
	if err != nil {
		return err
	}
	
	// Open the SQL database
	err = Uni.OpenDB()
	if err != nil {
		return err
	}
	
	// Login Bot
	Uni.S, err = discordgo.New("Bot " + Uni.Config.Token)
	if err != nil {
		return err
	}

	// Adding the Handlers
	Uni.S.AddHandler(Uni.onMessageCreate)
	Uni.S.AddHandler(Uni.onReady)

	// Startup threads for Uni Bot
	go func() { // RNG
		// might consider making uni's RNG channel into a bytes.Buffer
		var ri uint64 // Random Int
		for {
			binary.Read(crand.Reader, binary.BigEndian, &ri)
			Uni.RNGChan <- ri
		}
	}()

	// every hour run the garbage collector
	go func() {
		for range time.Tick(time.Hour) {
			runtime.GC()
		}
	}()
	
	// Every hour change status on uni
	go func() {
		for range time.Tick(time.Hour) {
			chosenstatus := statuses[<-Uni.RNGChan%uint64(len(statuses))]
			if chosenstatus[0] == "Playing" { // Playing Game
				Uni.S.UpdateStatus(0, chosenstatus[1])
			} else if chosenstatus[0] == "Listening" { // Listening song
				Uni.S.UpdateListeningStatus(chosenstatus[1])
			}
		}
	}()
	
	// Create necessary directories
	for _, dirstr := range []string{
	Uni.TempDir,
	"errlog",
	"backups",} {
		if _, err := os.Stat(dirstr); os.IsNotExist(err) {os.MkdirAll(dirstr, 0755)}
	}
	
	
	// everything up to here should be fine unless "S.Open()" returns an error
	return Uni.S.Open()
}


// When bot is ready
func (Uni *UniBot) onReady(s *discordgo.Session, event *discordgo.Ready) {
	Uni.Self = event.User
	fmt.Println("UniBot is now ready")
	Uni.messageGate = true
}


// *sees message*
func (Uni *UniBot) onMessageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
	// If gate is open or message is from unibot herself, then don't process the message
	if !Uni.messageGate || m.Author.ID == Uni.Self.ID { return }
	//fmt.Println(m.Message) // seeing message?
	start1 := time.Now()
	
	//var stopwatch *StopWatch // TODO
	c, err := s.State.Channel(m.ChannelID)
	if err != nil {
		fmt.Println("Error occurred while getting channel: ", err)
		return
	}
	
	var (
	isDM bool
	g *discordgo.Guild
	name string
	modules uint64
	)
	
	if m.GuildID != "" { // "check if DM" statement
		// get guild variable
		g, err = s.State.Guild(c.GuildID)
		if err != nil {
			fmt.Println("Error occurred while getting guild(aka server): ", err)
			return
		}
		
		// Get author name/nickname
		if m.Member != nil { // handle webhook messages
			st, _ := s.GuildMember(g.ID, m.Author.ID)
			if st != nil {
				name = st.Nick
			}
		}
		
		// Check if server exists in database
		var s interface{}
		
		Uni.DBGetFirstVar(&s, "CheckGuild", g.ID)
		if s == nil { // if nothing was returned then create an index for it
			Uni.DBExec("CreateGuild", g.ID)
		}
		
		s = nil // Reset
		
		// Check if channel exists in database
		Uni.DBGetFirstVar(&s, "CheckChannelModules", c.ID)
		if s == nil { // if nothing was returned then create an index for it
			Uni.DBExec("CreateChannelModules", g.ID, c.ID)
		}
		// Grab Modules
		Uni.DBGetFirstVar(&modules, "GetChannelModules", c.ID)
	} else {
		isDM = true
	}
	
	if name == "" { // either is a direct message or something happened on line 145
		name = m.Author.Username
	}
	
	
	// Now that everything is setup, perform the magic
	var prefix string = "hey uni"
	e1 := time.Since(start1)
	start2 := time.Now()
	Uni.MessageHandlers += 1
	defer func() { Uni.MessageHandlers -= 1 }()
	prv, cmd := Uni.ProcessMessage(m, isDM, g, c, modules, name, prefix) // Magic happens here
	if Uni.Debug > 0 && prv != 0 {
		fmt.Printf("MessageCreate{%s} Handler took Init: %10s, Processing: %10s\n", cmd, e1, time.Since(start2))
	}

}

// send a message without having to wait until it sends
func (Uni *UniBot) Respond(cID, content string) {
	go Uni.S.ChannelMessageSend(cID, content)
}

// Something went horribly wrong and notify the main creator and log to file
// Will also exit the entire goroutine *if* err is not nil
func (Uni *UniBot) ErrRespond(err error, cID, action string, vars ...interface{}) {
	if err == nil {return}
	Uni.Respond(cID, fmt.Sprintf("A fatal error occured while %s, discontinuing operation.\nError: ```\n%v\n```", action, err))
	efn := fmt.Sprintf("errlog/%s.log", time.Now().Format(time.RFC3339Nano)) // Error File Name
	Uni.Respond(Uni.Config.ErrLogChannel, fmt.Sprintf("A fatal error occured while %s.\nErrFile: %s\nError: ```\n%v\n```", action, efn, err))
	var w io.Writer // writer
	// I hope to the gods that there was no error while opening this file
	ef, err := os.Create(efn)
	if err != nil { // somehow you got here but there is a failsafe and that's printing it to stdout
		w = os.Stdout
		fmt.Printf("Error occured while %s, printing to stdout\n", action)
	} else { // alright everything is fine
		w = ef
		fmt.Printf("Error occured while %s, printing to file %q\n", action, efn)
		defer ef.Close() // close file when finished with dumping
	}
	spew.Fdump(w, vars)
	runtime.Goexit()
}

// Only used for testing/debugging
func (Uni *UniBot) DebugPrintVars(vars ...interface{}) {
	spew.Fdump(os.Stdout, vars)
}