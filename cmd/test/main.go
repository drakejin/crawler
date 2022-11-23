package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
)

func main() {

	url := "https://www.reddit.com/r/aww/top/?t=all"
	method := "GET"

	payload := strings.NewReader(``)

	client := &http.Client{}
	req, err := http.NewRequest(method, url, payload)

	if err != nil {
		fmt.Println(err)
		return
	}
	req.Header.Add("Cookie", "csv=2; edgebucket=wZEwiYsD9cYrhrX3mT; loid=0000000000uiaxjdvf.2.1669222828000.Z0FBQUFBQmpmbEdzaG9kZGdyaTNfcUNCbGp2NWh0ZXQtZUVUWXBCWGVzTW1ZdHBGN0xEbm1oUzgxbFJibEQ0d3d5Q1Brd2tLa1p2ekRjY2QtWGd2b05VRTVpRy1IblRRSi1CXzlvbEJfdmxEOG5GMWdaTGstdTV2ZklnZ1lnRUtiYmhDNVhTVkpDeVU; session_tracker=ghckafnjdddmbodqrm.0.1669222828436.Z0FBQUFBQmpmbEdzUjNDZ3J1UTZidm8wSW9VeXhJS01jYnBfOFNhRmRQSGtOU3BOOEoxaWNPaDl0ZDJDT1g4RGdWaVNGSUtoTXZnVnluY2dkTlNpUUNsT3ZWd1p5TkdDbXItMzRabG5qS2dkZzM3U3ZHQXFfdzlSbHdWS1NmclRaUWtMU2ZrbVBzTjg; token_v2=eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE2NjkzMDkxMDgsInN1YiI6Ii0yb3Npb3BXYTNTUmNablVLaEREV3Z6eWhfckFqM2ciLCJsb2dnZWRJbiI6ZmFsc2UsInNjb3BlcyI6WyIqIiwiZW1haWwiLCJwaWkiXX0.cdlCeUfkL0h4Z21T6HhyJu4Jv1iUR0Ffhwuye3uoeMg")

	res, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(string(body))
}
