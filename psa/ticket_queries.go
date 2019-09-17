package psa

import "fmt"

// Tickets encapsulates queries run on a specific board
type Tickets struct {
	client *Client
}

// NewTickets returns a new Tickets ready to run queries on
func NewTickets(client *Client) *Tickets {
	return &Tickets{
		client: client,
	}
}

// GetTickets get the service boards currently active
func (t *Tickets) GetTickets() ([]Ticket, error) {

	ticketsCmd := "/service/tickets"
	tickets, err := t.getCommand(ticketsCmd)
	if err != nil {
		return []Ticket{}, err
	}

	return tickets, nil
}

// GetOpenTickets get the service boards currently active
func (t *Tickets) GetOpenTickets() ([]Ticket, error) {

	ticketsCmd := "/service/tickets/search"

	conditions := make(map[string]string)
	conditions["conditions"] = "ClosedFlag = False"

	tickets, err := t.postCommand(ticketsCmd, conditions)
	if err != nil {
		return []Ticket{}, err
	}

	return tickets, nil
}

// getCommand runs a getCommand
func (t *Tickets) getCommand(cmd string) ([]Ticket, error) {

	pageSize := 1000
	currentPage := 1
	tickets := []Ticket{}

	for {
		cmd = fmt.Sprintf("%s?pageSize=%d&page=%d", cmd, pageSize, currentPage)
		page := []Ticket{}

		if err := t.client.restup.Get(cmd, &page); err != nil {
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

// postCommand runs a POST API query
func (t *Tickets) postCommand(cmd string, query map[string]string) ([]Ticket, error) {

	pageSize := 1000
	currentPage := 1
	tickets := []Ticket{}

	for {
		cmd = fmt.Sprintf("%s?pageSize=%d&page=%d", cmd, pageSize, currentPage)
		page := []Ticket{}

		if err := t.client.restup.Post(cmd, query, &page); err != nil {
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
