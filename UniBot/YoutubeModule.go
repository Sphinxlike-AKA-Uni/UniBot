package Uni
import (
	"fmt"
	"os"
	"os/exec"
	"syscall"
	"time"
	//"io"
	"github.com/bwmarrin/discordgo"
)

// Unused youtube-dl code that uni no longer uses

func (Uni *UniBot) DownloadMP3(s *discordgo.Session, m *discordgo.MessageCreate, link string) {
	// youtube-dl -x --audio-format flac --audio-quality 0 
	filename := fmt.Sprintf("tmp/%s.mp3", m.ID)
	ytcmd := exec.Command("youtube-dl", "-x", "--audio-format", "mp3", "--audio-quality", "0", "-o", filename, link)
	exitCode := 0
	var err error
	sc := make(chan int)
	go func(err *error, sc chan int) {
		*err = ytcmd.Run()
		sc <- 1
	}(&err, sc)
	
	go func(sc chan int) {
		time.Sleep(25 * time.Second)
		sc <- 0
	}(sc)
	
	a := <-sc
	if a == 0 {
		Respond(s, m, "Command took more than 25 seconds, aborting process")
		ytcmd.Process.Kill()
		return
	}
	
	if err != nil {
        if exitError, ok := err.(*exec.ExitError); ok {
            ws := exitError.Sys().(syscall.WaitStatus)
            exitCode = ws.ExitStatus()
        }
    } else {
        ws := ytcmd.ProcessState.Sys().(syscall.WaitStatus)
        exitCode = ws.ExitStatus()
    }
	if exitCode == 0 {
		r, _ := os.Open(filename)
		rs, _ := r.Stat()
		if rs.Size() > 8000000 {
			Respond(s, m, "File size over 8MB")
			return
		}
		Respond(s, m, "Uploading MP3...")
		s.ChannelFileSend(m.ChannelID, filename, r)
	} else {
		Respond(s, m, fmt.Sprintf("Exit code returned: %d", exitCode))
	}
}