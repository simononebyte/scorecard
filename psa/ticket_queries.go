package psa

import (
	"fmt"
	"time"
)

const (
	ticketsEndpoint      string = "/service/tickets"
	ticketSearchEndpoint string = "/service/tickets/search"
)

// GetTickets get the service boards currently active
func (c *Client) GetTickets() ([]Ticket, error) {

	tickets, err := c.getTicketsCommand(ticketsEndpoint)
	if err != nil {
		return []Ticket{}, err
	}

	return tickets, nil
}

// GetOpenTickets gets all open tickets
func (c *Client) GetOpenTickets() ([]Ticket, error) {

	conditions := make(map[string]string)
	conditions["conditions"] = "ClosedFlag = False"

	tickets, err := c.postTicketsCommand(ticketSearchEndpoint, conditions)
	if err != nil {
		return []Ticket{}, err
	}

	return tickets, nil
}

// GetOpenTicketsOnBoard gets all open tickets on a service board
func (c *Client) GetOpenTicketsOnBoard(name string) ([]Ticket, error) {

	boardID, err := c.GetBoardID(name)
	if err != nil {
		return []Ticket{}, err
	}

	conditions := make(map[string]string)
	conditions["conditions"] = fmt.Sprintf("ClosedFlag = False AND Board/ID = %v", boardID)

	tickets, err := c.postTicketsCommand(ticketSearchEndpoint, conditions)
	if err != nil {
		return []Ticket{}, err
	}

	return tickets, nil
}

// GetOpenAssignedTickets gets all open & assigned tickets
func (c *Client) GetOpenAssignedTickets() ([]Ticket, error) {

	conditions := make(map[string]string)
	conditions["conditions"] = "ClosedFlag = False AND resources LIKE '*'"

	tickets, err := c.postTicketsCommand(ticketSearchEndpoint, conditions)
	if err != nil {
		return []Ticket{}, err
	}

	return tickets, nil
}

// GetOpenAssignedTicketsOnBoard gets all open & assigned tickets
func (c *Client) GetOpenAssignedTicketsOnBoard(name string) ([]Ticket, error) {

	boardID, err := c.GetBoardID(name)
	if err != nil {
		return []Ticket{}, err
	}

	conditions := make(map[string]string)
	conditions["conditions"] = fmt.Sprintf("ClosedFlag = False AND resources LIKE '*' AND Board/ID = %v", boardID)

	tickets, err := c.postTicketsCommand(ticketSearchEndpoint, conditions)
	if err != nil {
		return []Ticket{}, err
	}

	return tickets, nil
}

// GetOpenTicketsOlderThan all open tickets older than the specified days
func (c *Client) GetOpenTicketsOlderThan(days int) ([]Ticket, error) {

	if days > 0 {
		days = days * -1
	}
	date := time.Now().AddDate(0, 0, days)
	dateStr := fmt.Sprintf("%d-%d-%d", date.Year(), date.Month(), date.Day())

	conditions := make(map[string]string)
	conditions["conditions"] = fmt.Sprintf("ClosedFlag = False AND dateEntered < [%v]", dateStr)

	tickets, err := c.postTicketsCommand(ticketSearchEndpoint, conditions)
	if err != nil {
		return []Ticket{}, err
	}

	return tickets, nil
}

// GetOpenTicketsObBoardOlderThan all open tickets older than the specified dayson a baord
func (c *Client) GetOpenTicketsObBoardOlderThan(name string, days int) ([]Ticket, error) {

	boardID, err := c.GetBoardID(name)
	if err != nil {
		return []Ticket{}, err
	}

	if days > 0 {
		days = days * -1
	}
	date := time.Now().AddDate(0, 0, days)
	dateStr := fmt.Sprintf("%d-%d-%d", date.Year(), date.Month(), date.Day())

	conditions := make(map[string]string)
	conditions["conditions"] = fmt.Sprintf("ClosedFlag = False AND dateEntered < [%v]  AND Board/ID = %v", dateStr, boardID)

	tickets, err := c.postTicketsCommand(ticketSearchEndpoint, conditions)
	if err != nil {
		return []Ticket{}, err
	}

	return tickets, nil
}

// GetOpenNotAssignedTickets gets all open & assigned tickets
func (c *Client) GetOpenNotAssignedTickets() ([]Ticket, error) {

	conditions := make(map[string]string)
	conditions["conditions"] = fmt.Sprintf("ClosedFlag = False AND resources = NULL")

	tickets, err := c.postTicketsCommand(ticketSearchEndpoint, conditions)
	if err != nil {
		return []Ticket{}, err
	}

	return tickets, nil
}

// GetOpenNotAssignedTicketsOnBoard gets all open & assigned tickets
func (c *Client) GetOpenNotAssignedTicketsOnBoard(name string) ([]Ticket, error) {

	boardID, err := c.GetBoardID(name)
	if err != nil {
		return []Ticket{}, err
	}

	conditions := make(map[string]string)
	conditions["conditions"] = fmt.Sprintf("ClosedFlag = False AND resources = NULL AND Board/ID = %v", boardID)

	tickets, err := c.postTicketsCommand(ticketSearchEndpoint, conditions)
	if err != nil {
		return []Ticket{}, err
	}

	return tickets, nil
}
