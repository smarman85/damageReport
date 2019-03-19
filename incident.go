package main

import (
        "github.com/smarman85/dateRange"
        "encoding/json"
        "fmt"
        "io/ioutil"
        "net/http"
        "os"
				"log"
        "strconv"
)


type Incident struct {
        Incident_info []Incidents `json:"incidents"`
        MoreIncidents bool        `json:"more"`
}

type Incidents struct {
        Number int    `json:"incident_number"`
        Info   string `json:"description"`
        Creation string `json:"created_at"`
        LastUpdate string `json:"last_status_change_at"`
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

func reportIncident(responseBody []byte) (bool, int) {

        // initialize incidnet var
        var incident Incident
        err := json.Unmarshal(responseBody, &incident)
        check(err)

        var incidentCount int = 0
        for i := 0; i < len(incident.Incident_info); i++ {
                if incident.Incident_info[i].Team.TeamName == "DevOps Team" {
                        fmt.Println("Incident Number:       " + strconv.Itoa(incident.Incident_info[i].Number))
                        fmt.Println("Incident Created at:   " + incident.Incident_info[i].Creation)
                        fmt.Println("Last Status Change:    " + incident.Incident_info[i].LastUpdate)
                        fmt.Println("Incident INFO:         " + incident.Incident_info[i].Info)
                        fmt.Println("Incident Team:         " + incident.Incident_info[i].Team.TeamName)
                        fmt.Println("Learn More:            " + incident.Incident_info[i].Detail)
                        fmt.Println("**************************************")
                        incidentCount += 1
                }
        }

        // more results?
        var keepGoing bool
        if incident.MoreIncidents == true {
                keepGoing = true
        } else {
                keepGoing = false
        }

        return keepGoing, incidentCount

} // close func reportIncident

func apiRequest(queryString, authToken, offset string) []byte {

        // set up api req
        // pagination offset has to increase by limit number with every request
        request, _ := http.NewRequest("GET", "https://api.pagerduty.com/incidents?"+queryString+"&limit=100&offset="+offset, nil)
        request.Header.Set("Accept", "application/vnd.pagerduty+json;version=2")
        request.Header.Set("Authorization", "Token token="+authToken)

        resp, err := http.DefaultClient.Do(request)
        check(err)

        body, _ := ioutil.ReadAll(resp.Body)
        check(err)

        // close api request
        defer resp.Body.Close()

        return body

}

func main() {

        if os.Getenv("PAGER_DUTY_TOKEN") == "" {
                log.Fatal("Error loading PagerDuty environment variable")
        }

        // set pagerDuty token
        authToken := os.Getenv("PAGER_DUTY_TOKEN")
        queryString := dateRange.Run()

        body := apiRequest(queryString, authToken, "0")

        moreResults, incidentCount := reportIncident(body)

        var alertCount int = incidentCount
        offset := 0
        for moreResults {
                offset += 100
                body := apiRequest(queryString, authToken, strconv.Itoa(offset))
                levelDeeper, incidentCount := reportIncident(body)
                alertCount += incidentCount
                if levelDeeper == false {
                        break
                }
        }
        fmt.Println("Number of Incidents: " + strconv.Itoa(alertCount))

} // close func main() 
