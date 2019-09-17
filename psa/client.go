package psa

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"strings"
	"time"

	"github.com/simononebyte/restup"
)

// Config holds the API credentials for the PSA system
type Config struct {
	Company  string `json:"company"`
	Username string `json:"public"`
	Password string `json:"private"`
	ClientID string `json:"client_id"`
}

// Client ...
type Client struct {
	restup        *restup.RestUp
	excludes      []string
	reactiveSites []string
}

// SiteTickets string = siteCode and int = ticket count
type SiteTickets map[string]int

// Stats ...
type Stats struct {
	MRRTickets int
	ORRTickets int
}

// NewClient ...
func NewClient(c Config, reactiveSiteCodes []string, excludedSummaries []string) *Client {
	token := fmt.Sprintf("%s+%s:%s", c.Company, c.Username, c.Password)
	psa := Client{}
	psa.restup = restup.NewRestUp("https://api-eu.myconnectwise.net/v2019_4/apis/3.0/", token)
	psa.restup.AddHeader("clientId", c.ClientID)

	psa.reactiveSites = reactiveSiteCodes
	psa.excludes = excludedSummaries

	return &psa
}

// GetStats ...
func (client *Client) GetStats(start, end time.Time) Stats {

	stats := Stats{}

	fmt.Printf("PSA: Getting Tickets: \n")
	tickets, ticketsErr := client.GetTicketsForWeek(start, end)
	if ticketsErr != nil {
		fmt.Printf("error: \n%s\n", ticketsErr)
		os.Exit(1)
	}
	reactive := client.FilterReactiveTickets(tickets)
	fmt.Printf("Found : %d reactive tickets found out of %d raised\n", len(reactive), len(tickets))

	for _, r := range reactive {
		var mrr = false
		for _, siteCode := range client.reactiveSites {
			if r.Company.SiteCode == siteCode {
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
func (client *Client) FilterReactiveTickets(tickets []Ticket) []Ticket {
	reactive := []Ticket{}
	for _, v := range tickets {
		isReactive := false
		if strings.HasPrefix(v.Board.Name, "SD - Reactive") {
			if isSubjectExcluded(v.Summary, client.excludes) == false {
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

func listTickets(tickets []Ticket) {
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

// GetTicketsForWeek ..
func (client *Client) GetTicketsForWeek(start, end time.Time) ([]Ticket, error) {

	pageSize := 1000
	currentPage := 1
	filter := makePSADateFilter(start, end)
	tickets := []Ticket{}

	fmt.Printf("Filter: %s\n", filter)

	for {
		cmd := fmt.Sprintf("service/tickets/search/?pageSize=%d&page=%d", pageSize, currentPage)
		page := []Ticket{}

		if err := client.restup.Post(cmd, filter, &page); err != nil {
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
func makePSADateFilter(start, end time.Time) Query {

	s := fmt.Sprintf("%d-%d-%d", start.Year(), start.Month(), start.Day())
	e := fmt.Sprintf("%d-%d-%d", end.Year(), end.Month(), end.Day())

	filterText := "recordType = 'ServiceTicket' AND dateEntered >= [%s] AND dateEntered < [%s]"
	return Query{
		OrderBy:    "dateEntered",
		Conditions: fmt.Sprintf(filterText, s, e),
	}
}
