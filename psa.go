package main

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"strings"
	"time"

	"github.com/simononebyte/restup"
)

// PSAClient ...
type PSAClient struct {
	restup        *restup.RestUp
	excludes      []string
	reactiveSites []psaSite
}

type psaSite struct {
	Name     string
	SiteCode string
}

// NewPSAClient ...
func NewPSAClient(c config) *PSAClient {
	token := fmt.Sprintf("%s+%s:%s", c.ConnectWise.Company, c.ConnectWise.Username, c.ConnectWise.Password)
	psa := PSAClient{}
	psa.restup = restup.NewRestUp("https://api-eu.myconnectwise.net/v2019_4/apis/3.0/", token)
	psa.restup.AddHeader("clientId", c.ConnectWise.ClientID)
	fmt.Println("Using clientID", c.ConnectWise.ClientID)
	for _, v := range c.ReactiveSites {
		psa.reactiveSites = append(psa.reactiveSites, psaSite{v.Name, v.SiteCode})
	}
	psa.excludes = c.Excludes.Summary
	return &psa
}

// PSASiteTickets string = siteCode and int = ticket count
type PSASiteTickets map[string]int

// PSAStats ...
type PSAStats struct {
	MRRTickets int
	ORRTickets int
}

// PSATicket ...
type PSATicket struct {
	ID          int        `json:"id"`
	DateEntered time.Time  `json:"dateEntered"`
	Company     PSACompany `json:"company"`
	Board       PSABoard   `json:"board"`
	Summary     string     `josn:"summary"`
}

// PSACompany ...
type PSACompany struct {
	ID       int    `json:"id"`
	SiteCode string `json:"identifier"`
	Name     string `json:"name"`
}

// PSABoard ...
type PSABoard struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

// PSAQuery ...
type PSAQuery struct {
	OrderBy    string `json:"orderBy"`
	Conditions string `json:"conditions"`
}

// GetPSAStats ...
func (psa *PSAClient) GetPSAStats(start, end time.Time) PSAStats {

	stats := PSAStats{}

	fmt.Printf("PSA: Getting Tickets: \n")
	tickets, ticketsErr := psa.GetPSATicketsForWeek(start, end)
	if ticketsErr != nil {
		fmt.Printf("error: \n%s\n", ticketsErr)
		os.Exit(1)
	}
	reactive := psa.FilterReactiveTickets(tickets)
	fmt.Printf("Found : %d reactive tickets found out of %d raised\n", len(reactive), len(tickets))

	for _, r := range reactive {
		var mrr = false
		for _, s := range psa.reactiveSites {
			if r.Company.SiteCode == s.SiteCode {
				mrr = true
			}
		}

		if mrr == true {
			stats.MRRTickets++
		} else {
			stats.ORRTickets++
		}
	}
	return stats
}

// FilterReactiveTickets ..
func (psa *PSAClient) FilterReactiveTickets(tickets []PSATicket) []PSATicket {
	reactive := []PSATicket{}
	for _, v := range tickets {
		isReactive := false
		if strings.HasPrefix(v.Board.Name, "SD - Reactive") {
			if isSubjectExcluded(v.Summary, psa.excludes) == false {
				isReactive = true
			}
		}

		if isReactive == true {
			reactive = append(reactive, v)
		}
	}
	return reactive
}

// isSubjectExcluded excludes certain ticket subjects from the recorded stats
func isSubjectExcluded(summary string, excludes []string) bool {

	for _, v := range excludes {
		m, _ := regexp.MatchString(v, summary)
		if m == true {
			return true
		}
	}

	return false
}

func listTickets(tickets []PSATicket) {
	reader := bufio.NewReader(os.Stdin)

	maxLines := 25
	line := 0
	for _, v := range tickets {
		fmt.Printf("%6d  -  %s\n", v.ID, v.DateEntered)
		line++
		if line > maxLines {
			reader.ReadString('\n')
			line = 0
		}
	}
}

// GetPSATicketsForWeek ..
func (psa *PSAClient) GetPSATicketsForWeek(start, end time.Time) ([]PSATicket, error) {

	pageSize := 1000
	currentPage := 1
	filter := makePSADateFilter(start, end)
	tickets := []PSATicket{}

	fmt.Printf("Filter: %s\n", filter)

	for {
		cmd := fmt.Sprintf("service/tickets/search/?pageSize=%d&page=%d", pageSize, currentPage)
		page := []PSATicket{}

		if err := psa.restup.Post(cmd, filter, &page); err != nil {
			return tickets, err
		}
		if len(page) == 0 {
			break
		}
		tickets = append(tickets, page...)
		currentPage++
	}

	return tickets, nil
}

// makePSADateFilter returns the query filter needed by the PSA
func makePSADateFilter(start, end time.Time) PSAQuery {

	s := fmt.Sprintf("%d-%d-%d", start.Year(), start.Month(), start.Day())
	e := fmt.Sprintf("%d-%d-%d", end.Year(), end.Month(), end.Day())

	filterText := "recordType = 'ServiceTicket' AND dateEntered >= [%s] AND dateEntered < [%s]"
	return PSAQuery{
		OrderBy:    "dateEntered",
		Conditions: fmt.Sprintf(filterText, s, e),
	}
}
