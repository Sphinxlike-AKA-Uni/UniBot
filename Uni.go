package main

import (
	"fmt"
	"time"
	"os"
	"os/signal"
	"syscall"
	"io/ioutil"
	"./UniBot"
)








// Start everything

func main() {
	defer fmt.Println("Uni shutdown") //defer test
	pui := time.Now() // Pre-Uni Init
	uni := Uni.New()
	uni.Debug = true
	uniownerid, _ := ioutil.ReadFile("../OwnerID.inf")
	uni.CreatorID = string(uniownerid)
	unitoken, _ := ioutil.ReadFile("../token.inf")
	uni.Token = string(unitoken)
	uni.DBLocation = "../Uni.db"
	uni.LuaDir = "../Lua"
	err := uni.Startup()
	if err != nil {
		fmt.Println(err)
		os.Exit(2147483646)
	}
	go uni.UpdateGameStatuses(1800, "GameList.inf")
	defer UniClose(uni) // Just in case
	fmt.Println("Time took to startup Uni: ", time.Since(pui))
	signal.Notify(uni.SC)
	for {
		cs := <-uni.SC
		if cs != syscall.SIGPIPE { // IDK why ubuntu server keeps giving me this whenever I host uni for a while
			if cs != syscall.Signal(0x1C) { // So I can SSH properly with JuiceSSH
				fmt.Println("Uni caught signal: ", cs, "\nClosing Uni....")
				fmt.Printf("Signal ID: %d\n", cs)
				break
			}
		}
	}
	fmt.Println("Shutting down Uni")
	UniClose(uni)
	fmt.Println("Uni close")
	if uni.Restart {
		os.Exit(1)
	}
}


func UniClose(uni Uni.UniBot) {
	uni.Database.Close()
	uni.DG.Close()
	fmt.Println("Uni Close Called")
}