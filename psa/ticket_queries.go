package psa

const (
	ticketsEndpoint      string = "/service/tickets"
	ticketSearchEndpoint string = "/service/tickets/search"
	ticketSourceEndpoint string = "/service/sources"

	escalatedText string = "Status has been updated from \"Needs-Info\" to \"Escalated from Helpdesk\"."
)

// GetNewTicketsByBoardID gets all new tickets on a service board
// boardID: The PSA board ID
// days: New tickets with the last x days
func (c *Client) GetNewTicketsByBoardID(boardID int, days int) ([]Ticket, error) {

	dateStr := dateStringFromDays(days)
	conditions := newCondition("dateEntered >= [%v] AND Board/ID = %v", dateStr, boardID)
	return c.postTicketsCommand(ticketSearchEndpoint, conditions)
}

// GetOpenTicketsByBoardID gets all open tickets on a service board
// boardID: The PSA board ID
func (c *Client) GetOpenTicketsByBoardID(boardID int) ([]Ticket, error) {

	conditions := newCondition("ClosedFlag = False AND Board/ID = %v", boardID)
	return c.postTicketsCommand(ticketSearchEndpoint, conditions)
}

// GetOpenTicketsByBoardIDOlderThan gets all open tickets on a service board
// boardID: The PSA board ID
// days:
func (c *Client) GetOpenTicketsByBoardIDOlderThan(boardID int, days int) ([]Ticket, error) {

	dateStr := dateStringFromDays(days)
	conditions := newCondition("ClosedFlag = False AND dateEntered <= [%v] AND Board/ID = %v", dateStr, boardID)
	return c.postTicketsCommand(ticketSearchEndpoint, conditions)
}

// GetOpenTicketsByBoardIDNotUpdatedIn gets all open tickets on a service board
// boardID: The PSA board ID
// days: Tickets that have not been upated in x days
func (c *Client) GetOpenTicketsByBoardIDNotUpdatedIn(boardID int, days int) ([]Ticket, error) {

	dateStr := dateStringFromDays(days)
	conditions := newCondition("ClosedFlag = False AND _info/LastUpdated <= [%v] AND Board/ID = %v", dateStr, boardID)
	return c.postTicketsCommand(ticketSearchEndpoint, conditions)
}

// GetOpenAssignedTicketsByBoardID gets all open tickets on a service board
// boardID: The PSA board ID
func (c *Client) GetOpenAssignedTicketsByBoardID(boardID int) ([]Ticket, error) {

	conditions := newCondition("ClosedFlag = False AND Board/ID = %v AND resources LIKE '*'", boardID)
	return c.postTicketsCommand(ticketSearchEndpoint, conditions)
}

// GetOpenNotAssignedTicketsByBoardID gets all open tickets on a service board
// boardID: The PSA board ID
func (c *Client) GetOpenNotAssignedTicketsByBoardID(boardID int) ([]Ticket, error) {

	conditions := newCondition("ClosedFlag = False AND Board/ID = %v AND resources = NULL", boardID)
	return c.postTicketsCommand(ticketSearchEndpoint, conditions)
}
