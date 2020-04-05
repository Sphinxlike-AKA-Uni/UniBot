package Uni

func (Uni *UniBot) Inspire(cID string) {
	inspiration, err := Uni.HTTPRequestBytes("GET", "http://inspirobot.me/api?generate=true", map[string]interface{}{"User-Agent": GrabUserAgent()}, nil)
	if err != nil {
		Uni.ErrRespond(err, cID, "getting from inspirobot", map[string]interface{}{"err": err, "cID": cID})
		return
	}
	Uni.Respond(cID, string(inspiration))
}