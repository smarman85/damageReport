package main

import (
        "os"
        "github.com/joho/godotenv"
        "fmt"
        "encoding/json"
        "io/ioutil"
        "log"
        "net/http"
        "strconv"
        "dateRange"
)

type Incident struct {
        Incident_info []Incidents `json:"incidents"`
}

type Incidents struct {
        Number int    `json:"incident_number"`
        Info   string `json:"description"`
        Team   Team   `json:"escalation_policy"`
        Detail string `json:"html_url"`
}

type Team struct {
        TeamName string `json:"summary"`
}

func check(e error) {
        if e != nil {
          panic(e)
        }
}

// load .env file
func main() {
        homeDir := os.Getenv("HOME")
        err := godotenv.Load(homeDir + "/go/src/github.com/damageReport/.env")
        if err != nil {
          log.Fatal("Error loading .env file")
        }

        // set pagerDuty token
        auth_token := os.Getenv("PAGER_DUTY_TOKEN")
        queryString := dateRange.Run()

        // set up api req
        // request, _ := http.NewRequest("GET", "https://api.pagerduty.com/incidents?since=2018-12-03T10:00&until=2018-12-10T10:00&limit=100", nil)
        request, _ := http.NewRequest("GET", "https://api.pagerduty.com/incidents?" + queryString + "&limit=100", nil)
        request.Header.Set("Accept", "application/vnd.pagerduty+json;version=2")
        request.Header.Set("Authorization", "Token token=" + auth_token)

        resp, err := http.DefaultClient.Do(request)
        check(err)

        // close api request
        defer resp.Body.Close()

        body, _ := ioutil.ReadAll(resp.Body)
        check(err)

        // initialize incidnet var
        var incident Incident
        err = json.Unmarshal(body, &incident)
        check(err)

        for i := 0; i < len(incident.Incident_info); i ++ {
                if incident.Incident_info[i].Team.TeamName == "DevOps Team" {
                        fmt.Println("Incident Number: " + strconv.Itoa(incident.Incident_info[i].Number))
                        fmt.Println("Incident INFO:   " + incident.Incident_info[i].Info)
                        fmt.Println("Incident Team:   " + incident.Incident_info[i].Team.TeamName)
                        fmt.Println("Learn More:      " + incident.Incident_info[i].Detail)
                        fmt.Println("**************************************")
                }
        }

        fmt.Println("Number of Incidents last week: " + strconv.Itoa(len(incident.Incident_info)))

} // close func main() 
