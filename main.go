package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net"
	"net/http"
	"strings"
)

type UserResponse struct {
	Ip       string `json:"client_ip"`
	Location string `json:"location"`
	Greeting string `json:"greeting"`
}

type Locationdetails struct {
	Location string `json:"city"`
	Weather  string `json:"temp"`
}

func GetUserResponse(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, "NON-GET REQUEST", 400)
	}

	urlquery := r.URL.Query().Get("visitor_name")
	if urlquery == "" {
		http.Error(w, "could not find url query param", 400)
	}

	ip := getIP(r)

	location := getGeoLocation(ip)
	weather := getWeather(location)
	cleanValue := strings.Replace(queryValue, "\"", "", -1)

	
	result := &UserResponse{
		Ip:       ip,
		Location: location,
		Greeting: fmt.Sprintf("Hello %s!, the temperature is %s degrees Celsius in %s", cleanValue, weather, location),
	}

	jsonResponse, err := json.Marshal(result)
	if err != nil {
		log.Println(err)
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonResponse)

}

func getIP(r *http.Request) string {
	// Check the X-Forwarded-For header first
	xff := r.Header.Get("X-Forwarded-For")
	if xff != "" {
		ips := strings.Split(xff, ",")
		if len(ips) > 0 {
			return strings.TrimSpace(ips[0])
		}
	}

	// If no X-Forwarded-For, check the X-Real-IP header
	xri := r.Header.Get("X-Real-IP")
	if xri != "" {
		return strings.TrimSpace(xri)
	}

	// If no headers, fall back to RemoteAddr
	ip, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		return ""
	}

	return ip
}

func getGeoLocation(ip string) string {
	// Replace with the actual URL and your API key for the geolocation service
	resp, err := http.Get(fmt.Sprintf("https://api.ipgeolocation.io/ipgeo?apiKey=ae608f3e55fe4c7d8621bffc87e4fcc0&ip=%s", ip))
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	var location Locationdetails
	if err := json.NewDecoder(resp.Body).Decode(&location); err != nil {
		panic(err)
	}

	return location.Location
}

func getWeather(city string) string {
	// Replace with the actual URL and your API key for the weather service
	resp, err := http.Get("http://api.openweathermap.org/data/2.5/weather?q=" + city + "&appid=fe6e2300fbeb4a7d98bbe5d141fe9685&units=metric")
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	var weather Locationdetails
	if err := json.NewDecoder(resp.Body).Decode(&weather); err != nil {
		panic(err)
	}

	return weather.Weather
}

func main() {
	http.HandleFunc("/api/hello", GetUserResponse)
	fmt.Printf("starting server at 5008/n")
	if err := http.ListenAndServe(":5008", nil); err != nil {
		log.Fatal(err)
	}
}

/*func getLocation(ip string) (string, error) {
    url := "https://ipinfo.io/" + ip + "/json"
    resp, err := http.Get(url)
    if err != nil {
        return "", err
    }
    defer resp.Body.Close()

    body, err := ioutil.ReadAll(resp.Body)
    if err != nil {
        return "", err
    }

    var data map[string]interface{}
    err = json.Unmarshal(body, &data)
    if err != nil {
        return "", err
    }

    location := data["city"].(string) + ", " + data["country"].(string)
    return location, nil
}
*/
