package Uni
import (
	"fmt"
	"os"
	"strconv"
	"github.com/jmank88/ubjson"
	"time"
)


// Create a Uni Bucks profile if user doesn't have one
func (Uni *UniBot) CheckIfProfileExists(cID, uID string) bool {
	u, err := Uni.DBGetFirst("GrabUniBucks", uID)
	if u == nil { // user does not have a unibucks profile
		if err != nil {
			Uni.ErrRespond(err, cID, "checking unibucks profile", map[string]interface{}{"err": err, "cID": cID})
			return false
		} else { // create uni bucks profile for user
			Uni.DBExec("InsertUniBucksProfile", uID)
		}
	}
	return true // everything should be a okie
}

////// Slot Machine

// No BS slot machine but even if it's all chance you still might lose a bit tbh
func (Uni *UniBot) SlotRoll(cID, uID, name string) {
	var slots [3][3]int
	items := []string{":seven:", ":cherries:", ":banana:", ":grapes:", ":watermelon:", ":apple:", ":eggplant:", ":chocolate_bar:", ":ice_cream:", ":doughnut:"}
	for i := 0; i < 9; i++ {
		slots[i/3][i%3] = int(<-Uni.RNGChan%uint64(len(items)))
	}
	r := ""
	for i := 0; i < 3; i++ {
		for j := 0; j < 3; j++ {
			r = fmt.Sprintf("%s%s", r, items[slots[i][j]])
		}
		r = fmt.Sprintf("%s\n", r)
	}
	
	inc := CheckRows(slots)
	var moreinfo string
	u, err := Uni.DBGetFirst("GrabUniBucks", uID)
	if err != nil {
		Uni.ErrRespond(err, cID, "grabbing unibucks profile", map[string]interface{}{"err": err, "cID": cID})
		return
	}
	
	var ti float64
	
	if len(inc) == 0 {
		moreinfo = "`LOSS (-25)`"
		ti -= 25
	} else {
		for _, inci := range inc {
			ti += float64(inci)
		}
		moreinfo = fmt.Sprintf("`WIN(+%g)`", ti)
	}
	
	
	Uni.Respond(cID, fmt.Sprintf("%s\n%s\n**%s's Uni Bucks: %g**", r, moreinfo, name, u.(float64)+ti))
	Uni.DBExec("AddUniBucks", ti, uID)
}

// Check to see if user has won
func CheckRows(slots [3][3]int) []int {
	/*
	:seven: = 777
	:cherries: = 100
	:banana: = 85
	:grapes: = 222
	:watermelon: = 127
	:apple: = 55
	:eggplant: = 10
	:chocolate_bar: = 333
	:ice_cream: = 444
	:doughnut: = 654
	*/
	ia := []int{777, 100, 85, 222, 127, 55, 10, 333, 444, 654}
	incarray := []int{}
	// 3 in a row, horizional
	for i := 0; i < 3; i++ {
		if (slots[i][0] == slots[i][1]) && (slots[i][0] == slots[i][2]) {
			incarray = append(incarray, ia[slots[i][0]])
		}
	}
	
	// 3 in a row, vertical
	for i := 0; i < 3; i++ {
		if (slots[0][i] == slots[1][i]) && (slots[0][i] == slots[2][i]) {
			incarray = append(incarray, ia[slots[0][i]])
		}
	}
	
	// 3 in a row, diagonal
	if ((slots[0][0] == slots[1][1]) && (slots[0][0] == slots[2][2])) {
		incarray = append(incarray, ia[slots[0][0]])
	}
	
	if ((slots[0][2] == slots[1][1]) && (slots[0][2] == slots[2][0])) {
		incarray = append(incarray, ia[slots[0][2]])
	}
	
	return incarray

}


////// Blackjack

type BlackjackGame struct {
	Bet float64
	DealerCards []Card
	UserCards []Card
	Deck []Card
}

var SuitStrings []string = []string{"♣", "♠", "♦", "♥"}
var CardStrings []string = []string{"A", "2", "3", "4", "5", "6", "7", "8", "9", "10", "J", "Q", "K"}
type Card struct {
	Suit byte
	Num byte
}

// Starting a blackjack deck
func (Uni *UniBot) StartBlackjack(uID, cID, betstr string) bool {
	var bjg BlackjackGame
	bf, err := os.OpenFile(fmt.Sprintf("%s/b_%s", Uni.TempDir, uID), os.O_RDWR, 0644) // Blackjack File
	if !(os.IsNotExist(err)) { // reaches here if file exists
		Uni.Respond(cID, "You appear to already have a game in session, here lemme show you the cards")
		return true
	} else {
		goto FileCreate
	}
	
FileCreate:
	bet, err := strconv.ParseFloat(betstr, 64)
	if err != nil {
		Uni.Respond(cID, "Error parsing bet text")
		return false
	}
	
	if bet < 0 { // can't have it where you lose and gain stuff
		Uni.Respond(cID, "Plz input a value greater than 0")
		return false
	}
	
	u, err := Uni.DBGetFirst("GrabUniBucks", uID)
	if err != nil {
		Uni.ErrRespond(err, cID, "grabbing unibucks profile", map[string]interface{}{"err": err, "cID": cID})
		return false
	}
	
	
	if u.(float64) < bet {
		Uni.Respond(cID, "Your bet is too high. Try something lower")
		return false
	}
	
	
	bjg.Bet = bet
	// Generate deck
	for i := 0; i < 20; i += 0 {
		tgc := Card{byte(<-Uni.RNGChan)%4, byte(<-Uni.RNGChan)%13}
		ce := false // Card Exists
		for _, ic := range bjg.Deck {
			if ic == tgc { ce = true }
		}
		
		if !ce {
			i++
			bjg.Deck = append(bjg.Deck, tgc)
		}
	}
	
	// Dealer draw
	bjg.DeckDraw(false)
	bjg.DeckDraw(false)
	
	// User draw
	bjg.DeckDraw(true)
	bjg.DeckDraw(true)
	
	// Write to file
	bf, err = os.Create(fmt.Sprintf("%s/b_%s", Uni.TempDir, uID))
	if err != nil {
		Uni.ErrRespond(err, cID, "creating blackjack file", map[string]interface{}{"err": err, "cID": cID})
		return false
	}
	defer bf.Close() // close so another function will read it
	ubjson.NewEncoder(bf).Encode(bjg)
	return true
}

// Draw top card from deck ( true = user, false = dealer )
func (bjg *BlackjackGame) DeckDraw(whom bool) {
	topcard := bjg.Deck[0]
	bjg.Deck = bjg.Deck[1:]
	if whom { // user
		bjg.UserCards = append(bjg.UserCards, topcard)
	} else { // dealer
		bjg.DealerCards = append(bjg.DealerCards, topcard)
	}
}

// Get the amount which person has (bool param same as "DeckDraw")
func (bjg *BlackjackGame) GetCardValues(cards []Card) ( v int ) {
	for i := 0; i < 2; i++ {
		for _, ic := range cards {
			switch ic.Num {
				case 0: v += []int{11, 1}[i] // Ace
				case 1: v += 2 // 2
				case 2: v += 3 // 3
				case 3: v += 4 // 4
				case 4: v += 5 // 5
				case 5: v += 6 // 6
				case 6: v += 7 // 7
				case 7: v += 8 // 8
				case 8: v += 9 // 9
				case 9: v += 10 // 10
				case 10, 11, 12: v += 10 // Jack, Queen, King
			}
		}
		if v > 21 && i == 0 { // retry calculation (where ace is benefitted as 1)
			v = 0
			continue
		} else {
			return
		}
	}
	return
}

// Present cards
func (Uni *UniBot) RepresentCards(cID, uID, name string) {
	bf, err := os.OpenFile(fmt.Sprintf("%s/b_%s", Uni.TempDir, uID), os.O_RDWR, 0644) // Blackjack File
	if err != nil {
		Uni.ErrRespond(err, cID, "opening blackjack file", map[string]interface{}{"err": err, "cID": cID,})
		return
	}
	var bjg BlackjackGame
	ubjson.NewDecoder(bf).Decode(&bjg)
	
	dcv := bjg.GetCardValues(bjg.DealerCards) // Dealer Card Value
	ucv := bjg.GetCardValues(bjg.UserCards) // User Card Value
	
	rs := "Dealer Cards: " // Response Speech
	if dcv == 21 || ucv == 21 || ucv > 21 {	
		for _, ic := range bjg.DealerCards {
			rs = fmt.Sprintf("%s, %s%s", rs, SuitStrings[ic.Suit], CardStrings[ic.Num])
		}
	} else {
		for dci, ic := range bjg.DealerCards {
			if dci == 0 {
				rs = fmt.Sprintf("%s ?", rs)
			} else {
				rs = fmt.Sprintf("%s, %s%s", rs, SuitStrings[ic.Suit], CardStrings[ic.Num])
			}
		}
	}
	rs = fmt.Sprintf("%s\n%s's Cards: ", rs, name)
	
	for _, ic := range bjg.UserCards {
		rs = fmt.Sprintf("%s, %s%s", rs, SuitStrings[ic.Suit], CardStrings[ic.Num])
	}
	
	rs = fmt.Sprintf("%s ( %d )", rs, ucv)
	
	// Win Detection
	var fb bool // Finish Blackjack
	if dcv == 21 {
		rs = fmt.Sprintf("%s\n%s", rs, fmt.Sprintf("Dealer blackjack, dealer wins `-%g`", bjg.Bet))
		Uni.DBExec("AddUniBucks", -bjg.Bet, uID)
		fb = true
	}
	
	if ucv == 21 {
		rs = fmt.Sprintf("%s\n%s", rs, fmt.Sprintf("%s blackjack, %s wins `+%g`", name, name, bjg.Bet))
		Uni.DBExec("AddUniBucks", bjg.Bet, uID)
		fb = true
	}
	
	if ucv > 21 {
		rs = fmt.Sprintf("%s\n%s", rs, fmt.Sprintf("%s bust, dealer wins `-%g`", name, bjg.Bet))
		Uni.DBExec("AddUniBucks", -bjg.Bet, uID)
		fb = true
	}
	
	if fb { // print user's unibucks on finish
		usersunibucks, err := Uni.DBGetFirst("GrabUniBucks", uID)
		if err != nil { // fook
			return
		}
		rs = fmt.Sprintf("%s\n**%s's Uni Bucks: %g**", rs, name, usersunibucks)
		os.Remove(fmt.Sprintf("%s/b_%s", Uni.TempDir, uID))
	}
	Uni.Respond(cID, rs)
	
}

// Hit me
func (Uni *UniBot) BlackjackHit(cID, uID string) {
	bf, err := os.OpenFile(fmt.Sprintf("%s/b_%s", Uni.TempDir, uID), os.O_RDWR, 0644) // Blackjack File
	if err != nil {
		Uni.ErrRespond(err, cID, "opening blackjack file", map[string]interface{}{"err": err, "cID": cID,})
		return
	}
	
	var bjg BlackjackGame
	ubjson.NewDecoder(bf).Decode(&bjg)
	
	if (len(bjg.DealerCards) == 0 || len(bjg.UserCards) == 0 || len(bjg.Deck) == 0) {
		os.Remove(fmt.Sprintf("%s/b_%s", Uni.TempDir, uID))
		Uni.Respond(cID, "Error reading blackjack file, closing session")
		return
	}
	bjg.DeckDraw(true)
	// reopen file
	bf.Close()
	bf, _ = os.OpenFile(fmt.Sprintf("%s/b_%s", Uni.TempDir, uID), os.O_RDWR, 0644) // Blackjack File
	ubjson.NewEncoder(bf).Encode(&bjg)
}

// Blackjack stay, (which dealer will make his moves here)
func (Uni *UniBot) BlackjackStay(cID, uID, name string) {
	bf, err := os.OpenFile(fmt.Sprintf("%s/b_%s", Uni.TempDir, uID), os.O_RDWR, 0644) // Blackjack File
	if err != nil {
		Uni.ErrRespond(err, cID, "opening blackjack file", map[string]interface{}{"err": err, "cID": cID,})
		return
	}
	
	var bjg BlackjackGame
	ubjson.NewDecoder(bf).Decode(&bjg)
	
	if (len(bjg.DealerCards) == 0 || len(bjg.UserCards) == 0 || len(bjg.Deck) == 0) {
		os.Remove(fmt.Sprintf("%s/b_%s", Uni.TempDir, uID))
		Uni.Respond(cID, "Error reading blackjack file, closing session")
		return
	}
	
	for { // Dealer "AI"
		an := int(<-Uni.RNGChan%10)
		dcv := bjg.GetCardValues(bjg.DealerCards)
		if an+10 > dcv {
			bjg.DeckDraw(false)
		} else {
			break // finish AI
		}
	}
	
	dcv := bjg.GetCardValues(bjg.DealerCards)
	ucv := bjg.GetCardValues(bjg.UserCards)
	
	rs := "Dealer Cards: "
	for _, ic := range bjg.DealerCards {
		rs = fmt.Sprintf("%s, %s%s", rs, SuitStrings[ic.Suit], CardStrings[ic.Num])
	}
	
	rs = fmt.Sprintf("%s ( %d )\n%s's Cards: ", rs, dcv, name)
	
	for _, ic := range bjg.UserCards {
		rs = fmt.Sprintf("%s, %s%s", rs, SuitStrings[ic.Suit], CardStrings[ic.Num])
	}
	rs = fmt.Sprintf("%s ( %d )", rs, ucv)
	if dcv > 21 { // Dealer bust
		rs = fmt.Sprintf("%s\nDealer bust, %s wins `+%g`", rs, name, bjg.Bet)
		Uni.DBExec("AddUniBucks", bjg.Bet, uID)
	} else if dcv > ucv { // Dealer was greater than user
		rs = fmt.Sprintf("%s\nDealer wins `-%g`", rs, bjg.Bet)
		Uni.DBExec("AddUniBucks", -bjg.Bet, uID)
	} else if dcv < ucv { // User was greater than dealer
		rs = fmt.Sprintf("%s\n%s wins `+%g`", rs, name, bjg.Bet)
		Uni.DBExec("AddUniBucks", bjg.Bet, uID)
	} else if dcv == ucv {
		rs = fmt.Sprintf("%s\nDealer and %s has tied `0`", rs, name)
	}
	
	usersunibucks, err := Uni.DBGetFirst("GrabUniBucks", uID)
	if err != nil { // fook
		return
	}
	
	rs = fmt.Sprintf("%s\n**%s's Uni Bucks: %g**", rs, name, usersunibucks)
	
	Uni.Respond(cID, rs)
	
	
	os.Remove(fmt.Sprintf("%s/b_%s", Uni.TempDir, uID))
}

// other uni bucks stuff
func (Uni *UniBot) Daily(cID, uID, name string) {
	var nanoseconds int64
	n, err := Uni.DBGetFirst("GrabDailyUniBucksTime", uID)
	if err != nil {
		Uni.ErrRespond(err, cID, "getting daily time", map[string]interface{}{"err": err, "cID": cID, })
		return
	}
	
	if n == nil { // no entry from user
		Uni.DBExec("InsertDailyUniBucksTime", uID, time.Now().UnixNano())
	}
	
	if nanoseconds+int64(time.Hour*24) < time.Now().UnixNano() { // give daily
		dinc := (<-Uni.RNGChan%(<-Uni.RNGChan%2480))+20 //Daily increase
		Uni.DBExec("AddUniBucks", dinc, uID)
		u, err := Uni.DBGetFirst("GrabUniBucks", uID)
		if err != nil { return }
		Uni.Respond(cID, fmt.Sprintf("You have recieved %d Uni Bucks for the day\n**%s's Uni Bucks: %g**", dinc, name, u))
	} else {
		Uni.Respond(cID, fmt.Sprintf("Plz wait %v", time.Duration((nanoseconds+int64(time.Hour*24))-time.Now().UnixNano())))
	}
}
