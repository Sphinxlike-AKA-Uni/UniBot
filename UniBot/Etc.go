// Script that contains extra functions that is useful in some areas
package Uni
import (
	"net/http"
	"io"
	"io/ioutil"
	"fmt"
	"runtime"
	"encoding/json"
)

// Checks if a string is inside of a string array
func IsStringInArray(str string, strarray []string) bool {
	for _, v := range strarray {
		if str == v { return true }
	}
	return false
}

// Easier way to do HTTP Requests
func (Uni *UniBot) HTTPRequest(method, link string, headers map[string]interface{}, body io.Reader) (*http.Response, error) {
	req, err := http.NewRequest(method, link, body)
	if err != nil { return nil, err }
	for k, v := range headers {
		req.Header.Set(k, fmt.Sprintf("%v", v))
	}
	return Uni.S.Client.Do(req)
}

// Just get everything from request all at once
func (Uni *UniBot) HTTPRequestBytes(method, link string, headers map[string]interface{}, body io.Reader) ([]byte, error){
	r, err := Uni.HTTPRequest(method, link, headers, body)
	if err != nil { return nil, err }
	return ioutil.ReadAll(r.Body)
}

// Transform the recieved json into struct provided // Mainly so i didn't have to check errs twice
func (Uni *UniBot) HTTPRequestJSON(method, link string, headers map[string]interface{}, body io.Reader, GivenStruct interface{}) (error) {
	resp, err := Uni.HTTPRequest(method, link, headers, body)
	if err != nil { return err }
	return json.NewDecoder(resp.Body).Decode(&GivenStruct)
}

// a function to generate a user agent
func GrabUserAgent() string {
	return fmt.Sprintf("UniBot/%s (%s; %s; %s)", versionstring, runtime.GOARCH, runtime.GOOS, runtime.Version())
}