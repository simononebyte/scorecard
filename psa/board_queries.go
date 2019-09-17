package psa

import (
	"fmt"
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
