package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"time"
)

type WeatherData struct {
	Success string `json:"success"`
	Result  struct {
		ResourceID string `json:"resource_id"`
		Fields     []struct {
			ID   string `json:"id"`
			Type string `json:"type"`
		} `json:"fields"`
	} `json:"result"`
	Records struct {
		Location []struct {
			LocationName   string `json:"locationName"`
			WeatherElement []struct {
				ElementName string `json:"elementName"`
				Time        []struct {
					StartTime string `json:"startTime"`
					EndTime   string `json:"endTime"`
					Parameter struct {
						ParameterName  string `json:"parameterName"`
						ParameterValue string `json:"parameterValue"`
						ParameterUnit  string `json:"parameterUnit"`
					} `json:"parameter"`
				} `json:"time"`
			} `json:"weatherElement"`
		} `json:"location"`
	} `json:"records"`
}

func handleWeatherAPI(w http.ResponseWriter, r *http.Request) {
	url := "https://opendata.cwa.gov.tw/api/v1/rest/datastore/F-C0032-001?Authorization=CWA-40B819F7-7B59-44CA-95EE-51AEA3FF6381&locationName=%E8%87%BA%E4%B8%AD%E5%B8%82"
	resp, err := http.Get(url)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var weatherData WeatherData
	err = json.Unmarshal(body, &weatherData)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	locationName := weatherData.Records.Location[0].LocationName
	location := weatherData.Records.Location[0].WeatherElement[0].Time[1].Parameter.ParameterName
	PoP := weatherData.Records.Location[0].WeatherElement[1].Time[1].Parameter.ParameterName
	MinT := weatherData.Records.Location[0].WeatherElement[2].Time[1].Parameter.ParameterName
	Ci := weatherData.Records.Location[0].WeatherElement[3].Time[1].Parameter.ParameterName
	MaxT := weatherData.Records.Location[0].WeatherElement[4].Time[1].Parameter.ParameterName

	startTime := weatherData.Records.Location[0].WeatherElement[0].Time[1].StartTime
	endTime := weatherData.Records.Location[0].WeatherElement[0].Time[1].EndTime

	startTimeFormatted, err := time.Parse("2006-01-02 15:04:05", startTime)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	endTimeFormatted, err := time.Parse("2006-01-02 15:04:05", endTime)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	message := fmt.Sprintf("\n今天%s的天氣：%s\n溫度：%s°C ~ %s°C\n降雨機率：%s%%\n舒適度：%s\n時間：%s ~ %s\n",
		locationName, location, MinT, MaxT, PoP, Ci, startTimeFormatted.Format("01-02 15:04"), endTimeFormatted.Format("01-02 15:04"))
	lineNotify(message)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": message})
}

func lineNotify(message string) {
	token := "iStUiDHiODsJQkjbmWqWL0i6AbGG6qrtw0l8IvUF0F3"
	lineUrl := "https://notify-api.line.me/api/notify"
	req, err := http.NewRequest("POST", lineUrl, strings.NewReader(url.Values{"message": {message}}.Encode()))
	if err != nil {
		fmt.Printf("Error creating Line Notify request: %v", err)
		return
	}
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	_, err = http.DefaultClient.Do(req)
	if err != nil {
		fmt.Printf("Error sending Line Notify: %v", err)
	}
}

func main() {
	http.HandleFunc("/weather", handleWeatherAPI)
	fmt.Println("Server started on port 8080")
	http.ListenAndServe(":8080", nil)
}
