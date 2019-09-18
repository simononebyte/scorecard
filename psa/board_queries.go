package psa

import (
	"fmt"
)

const (
	boardsEndpoint string = "/service/boards"
)

// GetBoards get the service boards currently active
func (c *Client) GetBoards() ([]Board, error) {

	boards, err := c.getBoardCommand(boardsEndpoint)
	if err != nil {
		return []Board{}, err
	}

	return boards, nil
}

// GetBoardID get ID for a service board
func (c *Client) GetBoardID(boardName string) (int, error) {

	boards, err := c.getBoardCommand(boardsEndpoint)
	if err != nil {
		return -1, err
	}

	for _, v := range boards {
		if v.Name == boardName {
			return v.ID, nil
		}
	}

	return -1, fmt.Errorf("service board not found")
}
