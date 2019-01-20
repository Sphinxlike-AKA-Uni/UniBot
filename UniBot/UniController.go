package Uni


import (
	"fmt"
	"time"
	"os"
	"math"
	"math/rand"
	"database/sql"
	"io"
	"runtime"
	"encoding/binary"
	"io/ioutil"
	"strconv"
	"strings"
	//"net/http"
	_ "github.com/mattn/go-sqlite3"
	"./UniLua"//"github.com/yuin/gopher-lua"
	"github.com/bwmarrin/discordgo"
	"bufio"
)


type UniBot struct {
	DG *discordgo.Session
	User *discordgo.User
	Rng *rand.Rand
	Debug bool
	DBLocation string
	Database *sql.DB
	Token string
	Email string
	Password string
	IsBot bool
	Msggate bool
	Restart bool
	Msghandlers int
	CreatorID string
	SC chan os.Signal
	RC chan LC
	LuaDir string
	APIPressure map[string]float64
	RecentRedditPosts RecentRedditPostsStruct
}

type RecentRedditPostsStruct struct {
	IDs [25]string
	Index int
}

type LC struct {
	RC chan lua.LValue
	LV lua.LValue
}

func New() (UniBot) {
	var Uni UniBot
	Uni.SC = make(chan os.Signal, 1)
	Uni.APIPressure = make(map[string]float64)
	Uni.IsBot = true
	Uni.DBLocation = "Uni.db"
	Uni.Msggate = false
	return Uni
}

func (Uni *UniBot) Startup() (error) {
	Uni.Rng = rand.New(rand.NewSource(99))
	Uni.Rng.Seed(time.Now().UnixNano() * 333)
	SwapUniSeed(Uni.Rng)
	Uni.Database, _ = sql.Open("sqlite3", Uni.DBLocation)
	Uni.Database.Exec("CREATE TABLE IF NOT EXISTS ServerData(id varchar(25) not null, assistancecooldown bigint not null default 0, adminrole varchar(25) default null)")
	Uni.Database.Exec("CREATE TABLE IF NOT EXISTS NSFWList(userid varchar(25) not null)")
	Uni.Database.Exec("CREATE TABLE IF NOT EXISTS Modules(gID varchar(25) not null, cID varchar(25) not null, modules int not null default 0)")
	Uni.Database.Exec("CREATE TABLE IF NOT EXISTS DerpiFilters(gID varchar(25) not null, cID varchar(25) not null, filterID bigint not null)")
	Uni.Database.Exec("PRAGMA journal_mode = WAL")
	Uni.Database.Exec("PRAGMA synchronous = EXTRA")
	Uni.Database.Exec("PRAGMA temp_store = 2")
	Uni.Database.SetMaxOpenConns(1)
	if Uni.IsBot {
		dg, err := discordgo.New("Bot " + Uni.Token)
		Uni.DG = dg
		if err != nil {
			return err
		}
		
	} else {
		dg, err := discordgo.New(Uni.Email, Uni.Password)
		Uni.DG = dg
		if err != nil {
			return err
		}
	}
	
	
	Uni.DG.AddHandler(Uni.onMessage)
	Uni.DG.AddHandler(Uni.onReady)
	Uni.DG.AddHandler(Uni.onChannelCreate)
	Uni.DG.AddHandler(Uni.onChannelRemove)
	
	for _, dirstr := range []string{
	"tmp",
	"errlog",
	"backups",
	Uni.LuaDir} {
		if _, err := os.Stat(dirstr); os.IsNotExist(err) {_ = os.MkdirAll(dirstr, 0755)}
	}
	
	
	go Uni.BackupDB()
	
	
	Uni.RC = make(chan LC)
	go Uni.LuaReact(Uni.DG, Uni.RC)
	go Uni.DG.Open()

	return nil
}









/************************************************************************/

func (Uni *UniBot) onChannelCreate(s *discordgo.Session, c *discordgo.ChannelCreate) {
	if c.Name != "" {
		Uni.Database.Exec(fmt.Sprintf("INSERT INTO Modules (gID, cID) VALUES ('%s', '%s');", c.GuildID, c.ID))
	}
}


func (Uni *UniBot) onChannelRemove(s *discordgo.Session, c *discordgo.ChannelDelete) {
	if c.Name != "" {
		Uni.Database.Exec(fmt.Sprintf("DELETE FROM Modules WHERE cID = '%s';", c.ID))
	}
}

func IsStringInsideArray(sub string, sa []string) bool {
	for _, i := range sa {
		if i == sub {
			return true
		}
	}
	return false
}

// When bot is ready
func (Uni *UniBot) onReady(s *discordgo.Session, event *discordgo.Ready) {
	Uni.User = event.User
	fmt.Println("Uni is now ready.")
	// DB Check
	
	serverIDs := []string{}
	
	rows, err := Uni.Database.Query("SELECT gid FROM Modules")
	if err != nil {
		fmt.Println("Error on startup DB check: ", err)
		return
	}
	for rows.Next() {
		var id string
		err = rows.Scan(&id)
		if err != nil {
			fmt.Println("Error on startup DB check: ", err)
			return
		}
		
		if !IsStringInsideArray(id, serverIDs) {
			serverIDs = append(serverIDs, id)
		}
		
	}
	rows.Close()
	
	gIDList := []string{}
	for _, g := range s.State.Guilds {
		gIDList = append(gIDList, g.ID)
	}
	
	for _, id := range serverIDs {
		if !IsStringInsideArray(id, gIDList) {
			Uni.Database.Exec(fmt.Sprintf("DELETE FROM Modules WHERE gID = '%s';", id))
			Uni.Database.Exec(fmt.Sprintf("DELETE FROM ServerData WHERE id = '%s';", id))
			fmt.Println(id)
		}
	}
	
	Uni.Database.Exec("VACUUM")
	go func() {
		for {
			SwapUniSeed(Uni.Rng)
			randomtimer := time.Duration(time.Duration(int64(math.Abs(float64(Uni.Rng.Int31())))*333))
			if Uni.Debug {
				fmt.Println("Swapping random seed, waiting ", randomtimer)
			}
			time.Sleep(randomtimer)
		}
	}()
	
	go Uni.LuaPressureDecay()
}





// *sees message*
func (Uni *UniBot) onMessage(s *discordgo.Session, m *discordgo.MessageCreate) {
	if Uni.Msggate {return}
	
	if m.Author.ID == Uni.User.ID { return }
	
	start1 := time.Now()
	
	c, err := s.State.Channel(m.ChannelID)
	if err != nil {
		// Could not find channel.
		return
	}
	
	g, _ := s.State.Guild(c.GuildID)
	cID := "NULL"
	
	isDM := false
	
	if (g != nil) {
		Uni.ServerDatabaseExistCheck(s, m, c.GuildID)
		Uni.ChannelExistCheck(s, m, c.GuildID, m.ChannelID) // To catch channels that were created when uni was asleep
		cID = c.ID
	} else {
		isDM = true
	}
	
	rows, err := Uni.Database.Query(fmt.Sprintf("SELECT modules FROM Modules WHERE cid = '%s'", cID))
	if err != nil {
		ErrRespond(s, m, fmt.Sprintf(" (Query failure while checking channel modules)	 %s", err))
		return
	}
	modules := 0
	for rows.Next() {
		_ = rows.Scan(&modules)
	}
	go rows.Close()
	
	var prefix string = "hey uni"
	
	/*if !isDM {
		for _, member := range g.Members {
			if member.User.ID == Uni.User.ID {
				if member.Nick != "" {
					prefix = fmt.Sprintf("hey %s", strings.ToLower(member.Nick))
				}
			}
		}
	}*/
	
	
	e1 := time.Since(start1)
	start2 := time.Now()
	Uni.Msghandlers += 1
	prv, cmd := Uni.ProcessMessage(s, m, isDM, g, c, modules, prefix) // Magic happens (mostly) here
	if Uni.Debug && prv != 0 {
		fmt.Printf("MessageCreate{%s} Handler took Init: %10s, Processing: %10s\n", cmd, e1, time.Since(start2))
	}

	if !isDM { // Prevents uni from crashing
		if _, err := os.Stat(Uni.LuaDir+"/"+g.ID+"/main.lua"); !os.IsNotExist(err) {
			Uni.ParseServerLua(s, m, g.ID)
		}
	}
	
	
	Uni.Msghandlers -= 1
}

/****************************************************/


func SwapUniSeed(rngobject *rand.Rand) {
	if runtime.GOOS == "linux" {
		// /dev/urandom as int64
		urd, err := os.Open("/dev/urandom")
		if err != nil {
			fmt.Println("Error opening /dev/urandom, leaving rng seed as is")
			return
		}
		urb := make([]byte, 10)
		_, err = urd.Read(urb)
		if err != nil {
			fmt.Println("Error reading /dev/urandom, leaving rng seed as is")
			return
		}
		rngobject.Seed(int64(binary.BigEndian.Uint64(urb)))
	} else { // then my own rng thingy which is kinda bad
		rngobject.Seed(time.Now().UnixNano()*rngobject.Int63())
	}
}


func (Uni *UniBot) UpdateGameStatuses(waittime int, getfile string) {
	time.Sleep(time.Second * 4)
	for {
		gfile, err := os.Open(getfile)
		if err != nil {
			fmt.Println("An error occured getting game list, ", err)
			fmt.Println(err)
			fmt.Println("UniBot will not update status")
			return
		}
		var lines []string
		scanner := bufio.NewScanner(gfile)
		for scanner.Scan() {
			lines = append(lines, scanner.Text())
		}
		gfile.Close()
		chosenline := lines[(int(Uni.Rng.Int63()) % len(lines))]
		if strings.HasPrefix(chosenline, "Playing,") {
			Uni.DG.UpdateStatus(0, chosenline[8:])
		} else if strings.HasPrefix(chosenline, "Listening,") {
			Uni.DG.UpdateListeningStatus(chosenline[10:])
		}
		time.Sleep(time.Second * time.Duration(waittime))
	}
}


func (Uni *UniBot) LuaPressureDecay() {
	allisnormal := true
	for {
		allisnormal = true
		for k := range Uni.APIPressure {
			time.Sleep(144447 * time.Microsecond)
			if Uni.APIPressure[k] <= .5 {
				Uni.APIPressure[k] = .5
				continue
			}
			
			allisnormal = false
			
			Uni.APIPressure[k] -= math.Atan2(math.Atan(77.777), 444)
			fmt.Println("cooldown time for ", k, ": ", Uni.APIPressure[k])
			if Uni.APIPressure[k] <= .5 {
				Uni.APIPressure[k] = .5
			}
			
		}
		if len(Uni.APIPressure) == 0 || allisnormal {
			time.Sleep(7 * time.Second)
		}
		
	}
}



func (Uni *UniBot) ServerDatabaseExistCheck(s *discordgo.Session, m *discordgo.MessageCreate, id string) {
	rows, err := Uni.Database.Query(fmt.Sprintf("SELECT * FROM ServerData WHERE id = '%s'", id))
	if err != nil {
		ErrRespond(s, m, fmt.Sprintf(" (Query failure while checking server) %s", err))
	}
	
	if rows == nil {
		return
	}
	
	defer rows.Close()
	
	for rows.Next() {return}
	Uni.Database.Exec(fmt.Sprintf("INSERT INTO ServerData (id) VALUES ('%s');", id))
	chs, _ := s.GuildChannels(id)
	cl := []string{}
	for _, ch := range chs {
		//if ch.ParentID == "" {
			cl = append(cl, ch.ID)
		//}
	}
	
	for _, cli := range cl {
		Uni.Database.Exec(fmt.Sprintf("INSERT INTO Modules (gID, cID) VALUES ('%s', '%s');", id, cli))
	}
	
}



func (Uni *UniBot) ChannelExistCheck(s *discordgo.Session, m *discordgo.MessageCreate, gID, cID string) {
	rows, err := Uni.Database.Query(fmt.Sprintf("SELECT * FROM Modules WHERE cID = '%s'", cID))
	if err != nil {
		ErrRespond(s, m, fmt.Sprintf(" (Query failure while checking channel) %s", err))
	}
	defer rows.Close()
	for rows.Next() {
		return
	}
	Uni.Database.Exec(fmt.Sprintf("INSERT INTO Modules (gID, cID) VALUES ('%s', '%s');", gID, cID))
	fmt.Sprintf("INSERT INTO Modules (gID, cID) VALUES ('%s', '%s');", gID, cID)
}


func (Uni *UniBot) BackupDB() {
	timestr, err := ioutil.ReadFile("BackupDBTimer.inf")
	if err != nil {
		fmt.Println("Error occured inside BackupDB: ", err)
		return
	}
	timerint, err := strconv.Atoi(string(timestr))
	if err != nil {
		fmt.Println("Error occured inside BackupDB: ", err)
		return
	}
	
	for {
		time.Sleep(time.Duration(time.Duration(timerint)) * time.Second)
		unidb, err := os.Open(Uni.DBLocation)
		if err != nil {
			fmt.Println("Error occured inside BackupDB: ", err)
			return
		}
  
		unidbcopy, err := os.OpenFile(strings.Replace(fmt.Sprintf("backups/%s.db", time.Now().Format(time.RFC3339)), ":", "__", -1), os.O_RDWR|os.O_CREATE, 0666)
		if err != nil {
			fmt.Println("Error occured inside BackupDB: ", err)
			return
		}

		_, err = io.Copy(unidbcopy, unidb)
		if err != nil {
			fmt.Println("Error occured inside BackupDB: ", err)
			return
		}
		unidb.Close()
		unidbcopy.Close()
	}
	
}
