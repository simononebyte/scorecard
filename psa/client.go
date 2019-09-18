package psa

import (
	"fmt"

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

// OrderBy ..
type OrderBy string

const (
	// OrderByAsc sorts queries in ascending order
	OrderByAsc OrderBy = "asc"
	// OrderByDesc sorts queries in descending order
	OrderByDesc OrderBy = "desc"
)

// NewClient ...
func NewClient(c Config, reactiveSiteCodes []string, excludedSummaries []string) *Client {
	token := fmt.Sprintf("%s+%s:%s", c.Company, c.Username, c.Password)

	client := &Client{}
	client.restup = restup.NewRestUp("https://api-eu.myconnectwise.net/v2019_4/apis/3.0/", token)
	client.restup.AddHeader("clientId", c.ClientID)

	client.reactiveSites = reactiveSiteCodes
	client.excludes = excludedSummaries

	return client
}

// getBoardCommand runs a Service Board GET API query
func (c *Client) getBoardCommand(cmd string) ([]Board, error) {

	pageSize := 1000
	currentPage := 1
	boards := []Board{}

	for {
		cmd = fmt.Sprintf("%s?pageSize=%d&page=%d", cmd, pageSize, currentPage)
		page := []Board{}

		if err := c.restup.Get(cmd, &page); err != nil {
			return []Board{}, err
		}
		if len(page) == 0 {
			break
		}
		boards = append(boards, page...)
		if len(page) <= pageSize {
			break
		}
		currentPage++
	}

	return boards, nil
}

// postBoardCommand runs a Service Board POST API query
func (c *Client) postBoardCommand(cmd string, query map[string]string) ([]Board, error) {

	pageSize := 1000
	currentPage := 1
	boards := []Board{}

	for {
		cmd = fmt.Sprintf("%s?pageSize=%d&page=%d", cmd, pageSize, currentPage)
		page := []Board{}

		if err := c.restup.Post(cmd, query, &page); err != nil {
			return []Board{}, err
		}
		if len(page) == 0 {
			break
		}
		boards = append(boards, page...)
		if len(page) <= pageSize {
			break
		}
		currentPage++
	}

	return boards, nil
}

// getCommand runs a getCommand
func (c *Client) getTicketsCommand(cmd string) ([]Ticket, error) {

	pageSize := 1000
	currentPage := 1
	tickets := []Ticket{}

	for {
		cmd = fmt.Sprintf("%s?pageSize=%d&page=%d", cmd, pageSize, currentPage)
		page := []Ticket{}

		if err := c.restup.Get(cmd, &page); err != nil {
			return []Ticket{}, err
		}
		if len(page) == 0 {
			break
		}
		tickets = append(tickets, page...)
		if len(page) <= pageSize {
			break
		}
		currentPage++
	}

	return tickets, nil
}

// postTicketsCommand runs a POST API query
func (c *Client) postTicketsCommand(cmd string, query map[string]string) ([]Ticket, error) {

	pageSize := 1000
	currentPage := 1
	tickets := []Ticket{}

	for {
		cmd = fmt.Sprintf("%s?pageSize=%d&page=%d", cmd, pageSize, currentPage)
		page := []Ticket{}

		if err := c.restup.Post(cmd, query, &page); err != nil {
			return []Ticket{}, err
		}
		if len(page) == 0 {
			break
		}
		tickets = append(tickets, page...)
		if len(page) <= pageSize {
			break
		}
		currentPage++
	}

	return tickets, nil
}

// GetStats ...
// func (client *Client) GetStats(start, end time.Time) Stats {

// 	stats := Stats{}

// 	fmt.Printf("PSA: Getting Tickets: \n")
// 	tickets, ticketsErr := client.GetTicketsForWeek(start, end)
// 	if ticketsErr != nil {
// 		fmt.Printf("error: \n%s\n", ticketsErr)
// 		os.Exit(1)
// 	}
// 	reactive := client.FilterReactiveTickets(tickets)
// 	fmt.Printf("Found : %d reactive tickets found out of %d raised\n", len(reactive), len(tickets))

// 	for _, r := range reactive {
// 		var mrr = false
// 		for _, siteCode := range client.reactiveSites {
// 			if r.Company.SiteCode == siteCode {
// 				mrr = true
// 			}
// 		}

// 		if mrr == true {
// 			stats.MRRTickets++
// 		} else {
// 			stats.ORRTickets++
// 		}
// 	}
// 	return stats
// }

// FilterReactiveTickets ..
// func (client *Client) FilterReactiveTickets(tickets []Ticket) []Ticket {
// 	reactive := []Ticket{}
// 	for _, v := range tickets {
// 		isReactive := false
// 		if strings.HasPrefix(v.Board.Name, "SD - Reactive") {
// 			if isSubjectExcluded(v.Summary, client.excludes) == false {
// 				isReactive = true
// 			}
// 		}

// 		if isReactive == true {
// 			reactive = append(reactive, v)
// 		}
// 	}
// 	return reactive
// }

// isSubjectExcluded excludes certain ticket subjects from the recorded stats
// func isSubjectExcluded(summary string, excludes []string) bool {

// 	for _, v := range excludes {
// 		m, _ := regexp.MatchString(v, summary)
// 		if m == true {
// 			return true
// 		}
// 	}

// 	return false
// }

// func listTickets(tickets []Ticket) {
// 	reader := bufio.NewReader(os.Stdin)

// 	maxLines := 25
// 	line := 0
// 	for _, v := range tickets {
// 		fmt.Printf("%6d  -  %s\n", v.ID, v.DateEntered)
// 		line++
// 		if line > maxLines {
// 			reader.ReadString('\n')
// 			line = 0
// 		}
// 	}
// }

// // GetTicketsForWeek ..
// func (client *Client) GetTicketsForWeek(start, end time.Time) ([]Ticket, error) {

// 	pageSize := 1000
// 	currentPage := 1
// 	filter := makePSADateFilter(start, end)
// 	tickets := []Ticket{}

// 	fmt.Printf("Filter: %s\n", filter)

// 	for {
// 		cmd := fmt.Sprintf("service/tickets/search/?orderBy=%s&pageSize=%d&page=%d", OrderByAsc, pageSize, currentPage)
// 		page := []Ticket{}

// 		if err := client.restup.Post(cmd, filter, &page); err != nil {
// 			return tickets, err
// 		}
// 		if len(page) == 0 {
// 			break
// 		}
// 		tickets = append(tickets, page...)
// 		currentPage++
// 	}

// 	return tickets, nil
// }

// // makePSADateFilter returns the query filter needed by the PSA
// func makePSADateFilter(start, end time.Time) Query {

// 	s := fmt.Sprintf("%d-%d-%d", start.Year(), start.Month(), start.Day())
// 	e := fmt.Sprintf("%d-%d-%d", end.Year(), end.Month(), end.Day())

// 	filterText := "recordType = 'ServiceTicket' AND dateEntered >= [%s] AND dateEntered < [%s]"
// 	return Query{
// 		OrderBy:    "dateEntered",
// 		Conditions: fmt.Sprintf(filterText, s, e),
// 	}
// }

// // runTicketCommand ..
// func (client *Client) runTicketCommand(cmd string, query Query) ([]Ticket, error) {

// 	pageSize := 1000
// 	currentPage := 1
// 	tickets := []Ticket{}

// 	for {
// 		cmd = fmt.Sprintf("service/tickets/search/?pageSize=%d&page=%d", pageSize, currentPage)
// 		page := []Ticket{}

// 		if err := client.restup.Post(cmd, query, &page); err != nil {
// 			return []Ticket{}, err
// 		}
// 		if len(page) == 0 {
// 			break
// 		}
// 		tickets = append(tickets, page...)
// 		currentPage++
// 	}

// 	return tickets, nil
// }

// // getCommand runs a getCommand
// func (client *Client) getCommand(cmd string, orderBy OrderBy) ([]Ticket, error) {

// 	pageSize := 1000
// 	currentPage := 1
// 	tickets := []Ticket{}

// 	for {
// 		cmd = fmt.Sprintf("service/tickets/search/?pageSize=%d&page=%d", pageSize, currentPage)
// 		page := []Ticket{}

// 		if err := client.restup.Get(cmd, &page); err != nil {
// 			return []Ticket{}, err
// 		}
// 		if len(page) == 0 {
// 			break
// 		}
// 		tickets = append(tickets, page...)
// 		currentPage++
// 	}

// 	return tickets, nil
// }
