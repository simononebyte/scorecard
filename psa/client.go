package psa

import (
	"fmt"
	"strconv"
	"time"

	"github.com/simononebyte/restup"
)

// Config holds the API credentials for the PSA system
type Config struct {
	Company  string `json:"company"`
	Username string `json:"public"`
	Password string `json:"private"`
	ClientID string `json:"client_id"`
	APIBase  string `json:"api_base"`
}

// Client ...
type Client struct {
	restup        *restup.RestUp
	excludeBoards []Board
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

// NewClient creates a new PSA Client.
//  Config contians the ConnectWise API and Client Keys
//  globalBoardExcludes is a list of service boards that will be
//  excluded from all qureies
func NewClient(c Config, globalBoardExcludes []string) (*Client, error) {
	token := fmt.Sprintf("%s+%s:%s", c.Company, c.Username, c.Password)

	client := &Client{}
	client.restup = restup.NewRestUp(c.APIBase, token)
	client.restup.AddHeader("clientId", c.ClientID)

	if len(globalBoardExcludes) > 0 {
		if err := client.populateExcludes(globalBoardExcludes); err != nil {
			return &Client{}, err
		}
	}
	return client, nil
}

func (c *Client) populateExcludes(excludes []string) error {
	boards, err := c.GetBoards()
	if err != nil {
		return err
	}
	c.excludeBoards = make([]Board, 0)
	for _, b := range boards {
		for _, e := range excludes {
			if b.Name == e {
				c.excludeBoards = append(c.excludeBoards, b)
				break
			}
		}
	}
	if len(excludes) != len(c.excludeBoards) {
		return fmt.Errorf("not all boards to be excluded were found")
	}
	return nil
}

func newCondition(condition string, a ...interface{}) map[string]string {
	conditions := make(map[string]string)
	conditions["conditions"] = fmt.Sprintf(condition, a...)
	return conditions
}

func dateStringFromDays(days int) string {
	if days > 0 {
		days = days * -1
	}
	date := time.Now().AddDate(0, 0, days)
	return fmt.Sprintf("%d-%d-%d", date.Year(), date.Month(), date.Day())
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

// getMembersCommand runs a getCommand
func (c *Client) getMembersCommand(cmd string) ([]Member, error) {

	pageSize := 1000
	currentPage := 1
	members := []Member{}

	for {
		cmd = fmt.Sprintf("%s?pageSize=%d&page=%d", cmd, pageSize, currentPage)
		page := []Member{}

		if err := c.restup.Get(cmd, &page); err != nil {
			return []Member{}, err
		}
		if len(page) == 0 {
			break
		}
		members = append(members, page...)
		if len(page) <= pageSize {
			break
		}
		currentPage++
	}

	return members, nil
}

func (c *Client) wrapExcludedBoards(condition string) string {
	newCondition := "((" + condition + ")"
	for _, v := range c.excludeBoards {
		newCondition += " AND Board/ID != " + strconv.Itoa(v.ID) + ")"
	}
	return newCondition
}

// getTicketSourceCommand runs a getCommand
func (c *Client) getTicketSourceCommand(cmd string) ([]TicketSource, error) {

	pageSize := 1000
	currentPage := 1
	sources := []TicketSource{}

	for {
		cmd = fmt.Sprintf("%s?pageSize=%d&page=%d", cmd, pageSize, currentPage)
		page := []TicketSource{}

		if err := c.restup.Get(cmd, &page); err != nil {
			return []TicketSource{}, err
		}
		if len(page) == 0 {
			break
		}
		sources = append(sources, page...)
		if len(page) <= pageSize {
			break
		}
		currentPage++
	}

	return sources, nil
}

// getMembersCommand runs a getCommand
func (c *Client) getAuditTrailCommand(cmd string) ([]Audit, error) {

	pageSize := 1000
	currentPage := 1
	audit := []Audit{}

	for {
		page := []Audit{}

		if err := c.restup.Get(cmd, &page); err != nil {
			return []Audit{}, err
		}
		if len(page) == 0 {
			break
		}
		audit = append(audit, page...)
		if len(page) <= pageSize {
			break
		}
		currentPage++
	}

	return audit, nil
}
