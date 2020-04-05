// Main Script
package main
import (
	"fmt"
	"time"
	"os"
	"syscall"
	"flag"
	"./UniBot"
)

var (
	ucf string // Uni Config File
	err error
)


func init() {
	flag.StringVar(&ucf, "config", "UniConfig.inf", "Load Config File")
	flag.Parse()
}

func main() {
	uni := Uni.New()
	uni.Debug = 4 // {0 = No Debug, 1 = Light Debug, 2 = Medium Debug, >2 = All Debug}
	// Some things that can be set from the outside
	ust := time.Now() // Uni Startup Time
	err = uni.Startup(ucf)
	if err != nil {
		fmt.Println("Error starting up UniBot: ", err)
		os.Exit(0x1F)
	}
	fmt.Println("Time took to startup Uni: ", time.Since(ust))
	defer func() { // Shutdown uni properly
		fmt.Println("Uni Close Called")
		uni.DB.Close()
		uni.S.Close()
		if uni.Restart {
			os.Exit(1) // ExitCode 1 to signal a uni restart
		}
	}()

	for {
		cs := <-uni.SC // wait until a signal is captured
		fmt.Printf("Uni has caught signal %d(%s)\n", cs, cs)
		if cs != syscall.SIGPIPE && // Ignore broken pipe signals as uni can reconnect on her own
		cs != syscall.Signal(28) { // ignore "window changed" signals
			fmt.Println("Closing Uni")
			break
		}
	}

}
