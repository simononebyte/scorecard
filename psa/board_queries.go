package psa

import (
	"fmt"
	"time"
)

// Boards encapsulates queries run on a specific board
type Boards struct {
	client *Client
}

// NewBoards returns a new Boards ready to run queries on
func NewBoards(client *Client) *Boards {
	return &Boards{
		client: client,
	}
}

// GetBoards get the service boards currently active
func (b *Boards) GetBoards() ([]Board, error) {

	boardCmd := "/service/boards"
	boards, err := b.getCommand(boardCmd)
	if err != nil {
		return []Board{}, err
	}

	return boards, nil
}

// GetBoardID get ID for a service board
func (b *Boards) GetBoardID(boardName string) (int, error) {

	boardCmd := "/service/boards"
	boards, err := b.getCommand(boardCmd)
	if err != nil {
		return -1, err
	}

	for _, v := range boards {
		if v.Name == boardName {
			return v.ID, nil
		}
	}

	return -1, fmt.Errorf("error: Unable to find service board %v", boardName)
}

// GetOpenTickets get the tickets currently in an open state in the specified board
func (b *Boards) GetOpenTickets(boardName string) ([]Ticket, error) {

	boardCmd := "/service/boards"
	boards, berr := b.getCommand(boardCmd)
	if berr != nil {
		return []Ticket{}, berr
	}

	boardID := int(-1)
	for _, v := range boards {
		if v.Name == boardName {
			boardID = v.ID
			break
		}
	}
	if boardID == -1 {
		return []Ticket{}, fmt.Errorf("service board %s not found", boardName)
	}

	allTickets, terr := b.client.Tickets.GetOpenTickets()
	if terr != nil {
		return []Ticket{}, terr
	}
	tickets := make([]Ticket, 0)
	for _, ticket := range allTickets {
		if ticket.Board.ID == boardID {
			tickets = append(tickets, ticket)
		}
	}
	return tickets, nil
}

// GetOpenTicketsOlderThan all open tickets older than the specified days
func (b *Boards) GetOpenTicketsOlderThan(boardName string, days int) ([]Ticket, error) {

	boardID, err := b.GetBoardID(boardName)
	if err != nil {
		return []Ticket{}, err
	}
	if days > 0 {
		days = days * -1
	}
	date := time.Now().AddDate(0, 0, days)
	dateStr := fmt.Sprintf("%d-%d-%d", date.Year(), date.Month(), date.Day())

	ticketsCmd := "/service/tickets/search"

	conditions := make(map[string]string)
	conditions["conditions"] = fmt.Sprintf("ClosedFlag = False AND dateEntered < [%v] AND Board/ID = %v", dateStr, boardID)

	tickets, err := b.client.Tickets.postCommand(ticketsCmd, conditions)
	if err != nil {
		return []Ticket{}, err
	}

	return tickets, nil
}

// getCommand runs a getCommand
func (b *Boards) getCommand(cmd string) ([]Board, error) {

	pageSize := 1000
	currentPage := 1
	boards := []Board{}

	for {
		cmd = fmt.Sprintf("%s?pageSize=%d&page=%d", cmd, pageSize, currentPage)
		page := []Board{}

		if err := b.client.restup.Get(cmd, &page); err != nil {
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

// postCommand runs a POST API query
func (b *Boards) postCommand(cmd string, query map[string]string) ([]Board, error) {

	pageSize := 1000
	currentPage := 1
	boards := []Board{}

	for {
		cmd = fmt.Sprintf("%s?pageSize=%d&page=%d", cmd, pageSize, currentPage)
		page := []Board{}

		if err := b.client.restup.Post(cmd, query, &page); err != nil {
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
