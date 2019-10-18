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
	conditions := newCondition("ClosedFlag = False AND dateEntered >= [%v] AND Board/ID = %v",dateStr, boardID)
	return c.postTicketsCommand(ticketSearchEndpoint, conditions)
}

// GetOpenTicketsByBoardIDNotUpdatedIn gets all open tickets on a service board
// boardID: The PSA board ID
// days: Tickets that have not been upated in x days
func (c *Client) GetOpenTicketsByBoardIDNotUpdatedIn(boardID int, days int) ([]Ticket, error) {

	dateStr := dateStringFromDays(days)
	conditions := newCondition("ClosedFlag = False AND dateEntered >= [%v] AND Board/ID = %v", dateStr, boardID)
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

/*
##############################################################################
  Initial naive query implementations
##############################################################################
*/

// // GetTickets get the service boards currently active
// func (c *Client) GetTickets() ([]Ticket, error) {

// 	tickets, err := c.getTicketsCommand(ticketsEndpoint)
// 	if err != nil {
// 		return []Ticket{}, err
// 	}

// 	return tickets, nil
// }

// // GetOpenTickets gets all open tickets
// func (c *Client) GetOpenTickets() ([]Ticket, error) {

// 	conditions := make(map[string]string)
// 	conditions["conditions"] = "ClosedFlag = False"

// 	tickets, err := c.postTicketsCommand(ticketSearchEndpoint, conditions)
// 	if err != nil {
// 		return []Ticket{}, err
// 	}

// 	return tickets, nil
// }

// // GetOpenTicketsOnBoard gets all open tickets on a service board
// func (c *Client) GetOpenTicketsOnBoard(name string) ([]Ticket, error) {

// 	boardID, err := c.GetBoardID(name)
// 	if err != nil {
// 		return []Ticket{}, err
// 	}

// 	conditions := newCondition("ClosedFlag = False AND Board/ID = %v", boardID)

// 	tickets, err := c.postTicketsCommand(ticketSearchEndpoint, conditions)
// 	if err != nil {
// 		return []Ticket{}, err
// 	}

// 	return tickets, nil
// }

// // GetOpenAssignedTickets gets all open & assigned tickets
// func (c *Client) GetOpenAssignedTickets() ([]Ticket, error) {

// 	conditions := make(map[string]string)
// 	conditions["conditions"] = "ClosedFlag = False AND resources LIKE '*'"

// 	tickets, err := c.postTicketsCommand(ticketSearchEndpoint, conditions)
// 	if err != nil {
// 		return []Ticket{}, err
// 	}

// 	return tickets, nil
// }

// // GetOpenAssignedTicketsOnBoard gets all open & assigned tickets
// func (c *Client) GetOpenAssignedTicketsOnBoard(name string) ([]Ticket, error) {

// 	boardID, err := c.GetBoardID(name)
// 	if err != nil {
// 		return []Ticket{}, err
// 	}

// 	conditions := newCondition("ClosedFlag = False AND resources LIKE '*' AND Board/ID = %v", boardID)

// 	tickets, err := c.postTicketsCommand(ticketSearchEndpoint, conditions)
// 	if err != nil {
// 		return []Ticket{}, err
// 	}

// 	return tickets, nil
// }

// // GetOpenTicketsOlderThan all open tickets older than the specified days
// func (c *Client) GetOpenTicketsOlderThan(days int) ([]Ticket, error) {

// 	if days > 0 {
// 		days = days * -1
// 	}
// 	date := time.Now().AddDate(0, 0, days)
// 	dateStr := fmt.Sprintf("%d-%d-%d", date.Year(), date.Month(), date.Day())

// 	conditions := newCondition("ClosedFlag = False AND dateEntered < [%v]", dateStr)

// 	tickets, err := c.postTicketsCommand(ticketSearchEndpoint, conditions)
// 	if err != nil {
// 		return []Ticket{}, err
// 	}

// 	return tickets, nil
// }

// // GetOpenTicketsOnBoardOlderThan all open tickets older than the specified dayson a baord
// func (c *Client) GetOpenTicketsOnBoardOlderThan(name string, days int) ([]Ticket, error) {

// 	boardID, err := c.GetBoardID(name)
// 	if err != nil {
// 		return []Ticket{}, err
// 	}

// 	if days > 0 {
// 		days = days * -1
// 	}
// 	date := time.Now().AddDate(0, 0, days)
// 	dateStr := fmt.Sprintf("%d-%d-%d", date.Year(), date.Month(), date.Day())

// 	conditions := newCondition("ClosedFlag = False AND dateEntered < [%v]  AND Board/ID = %v", dateStr, boardID)

// 	tickets, err := c.postTicketsCommand(ticketSearchEndpoint, conditions)
// 	if err != nil {
// 		return []Ticket{}, err
// 	}

// 	return tickets, nil
// }

// // GetOpenNotAssignedTickets gets all open & assigned tickets
// func (c *Client) GetOpenNotAssignedTickets() ([]Ticket, error) {

// 	conditions := newCondition("ClosedFlag = False AND resources = NULL")

// 	tickets, err := c.postTicketsCommand(ticketSearchEndpoint, conditions)
// 	if err != nil {
// 		return []Ticket{}, err
// 	}

// 	return tickets, nil
// }

// // GetOpenNotAssignedTicketsOnBoard gets all open & assigned tickets
// func (c *Client) GetOpenNotAssignedTicketsOnBoard(name string) ([]Ticket, error) {

// 	boardID, err := c.GetBoardID(name)
// 	if err != nil {
// 		return []Ticket{}, err
// 	}

// 	conditions := newCondition("ClosedFlag = False AND resources = NULL AND Board/ID = %v", boardID)

// 	tickets, err := c.postTicketsCommand(ticketSearchEndpoint, conditions)
// 	if err != nil {
// 		return []Ticket{}, err
// 	}

// 	return tickets, nil
// }

// // GetOpenTicketsNotUpdatedIn all open tickets not updtaed in specified number of days
// func (c *Client) GetOpenTicketsNotUpdatedIn(days int) ([]Ticket, error) {

// 	if days > 0 {
// 		days = days * -1
// 	}
// 	date := time.Now().AddDate(0, 0, days)
// 	dateStr := fmt.Sprintf("%d-%d-%d", date.Year(), date.Month(), date.Day())

// 	conditions := newCondition("ClosedFlag = False AND _info/LastUpdated < [%v]", dateStr)

// 	tickets, err := c.postTicketsCommand(ticketSearchEndpoint, conditions)
// 	if err != nil {
// 		return []Ticket{}, err
// 	}

// 	return tickets, nil
// }

// // GetOpenTicketsNotUpdatedInOnBoard all open tickets not updtaed in specified number of days
// func (c *Client) GetOpenTicketsNotUpdatedInOnBoard(name string, days int) ([]Ticket, error) {

// 	boardID, err := c.GetBoardID(name)
// 	if err != nil {
// 		return []Ticket{}, err
// 	}

// 	if days > 0 {
// 		days = days * -1
// 	}
// 	date := time.Now().AddDate(0, 0, days)
// 	dateStr := fmt.Sprintf("%d-%d-%d", date.Year(), date.Month(), date.Day())

// 	conditions := newCondition("ClosedFlag = False AND _info/LastUpdated < [%v] AND Board/ID = %v", dateStr, boardID)

// 	tickets, err := c.postTicketsCommand(ticketSearchEndpoint, conditions)
// 	if err != nil {
// 		return []Ticket{}, err
// 	}

// 	return tickets, nil
// }

// // GetNewTicketsInLast all open tickets not updtaed in specified number of days
// func (c *Client) GetNewTicketsInLast(days int) ([]Ticket, error) {

// 	if days > 0 {
// 		days = days * -1
// 	}
// 	date := time.Now().AddDate(0, 0, days)
// 	dateStr := fmt.Sprintf("%d-%d-%d", date.Year(), date.Month(), date.Day())

// 	conditions := newCondition("dateEntered >= [%v]", dateStr)

// 	tickets, err := c.postTicketsCommand(ticketSearchEndpoint, conditions)
// 	if err != nil {
// 		return []Ticket{}, err
// 	}

// 	return tickets, nil
// }

// // GetNewTicketsInLastOnBoard all open tickets not updtaed in specified number of days
// func (c *Client) GetNewTicketsInLastOnBoard(name string, days int) ([]Ticket, error) {

// 	boardID, err := c.GetBoardID(name)
// 	if err != nil {
// 		return []Ticket{}, err
// 	}

// 	if days > 0 {
// 		days = days * -1
// 	}
// 	date := time.Now().AddDate(0, 0, days)
// 	dateStr := fmt.Sprintf("%d-%d-%d", date.Year(), date.Month(), date.Day())

// 	conditions := newCondition("dateEntered >= [%v] AND Board/ID = %v", dateStr, boardID)

// 	tickets, err := c.postTicketsCommand(ticketSearchEndpoint, conditions)
// 	if err != nil {
// 		return []Ticket{}, err
// 	}

// 	return tickets, nil
// }

// // GetOpenOpenTicketsAssignedTo gets all open assigned assigned to specific person
// func (c *Client) GetOpenOpenTicketsAssignedTo(identifier string) ([]Ticket, error) {

// 	condition := fmt.Sprintf("ClosedFlag = False AND resources LIKE '*%v*'", identifier)
// 	conditions := make(map[string]string)
// 	conditions["conditions"] = c.wrapExcludedBoards(condition)

// 	tickets, err := c.postTicketsCommand(ticketSearchEndpoint, conditions)
// 	if err != nil {
// 		return []Ticket{}, err
// 	}

// 	return tickets, nil
// }

// // GetTicketSources gets the various methods by which a ticket can be raised
// func (c *Client) GetTicketSources() ([]TicketSource, error) {

// 	sources, err := c.getTicketSourceCommand(ticketSourceEndpoint)
// 	if err != nil {
// 		return []TicketSource{}, err
// 	}

// 	return sources, nil
// }

// // GetEscalatedTicketsInLast all tickets escalated in last x days
// func (c *Client) GetEscalatedTicketsInLast(days int) ([]Ticket, error) {

// 	if days > 0 {
// 		days = days * -1
// 	}
// 	date := time.Now().AddDate(0, 0, days)

// 	dateStr := fmt.Sprintf("%d-%d-%d", date.Year(), date.Month(), date.Day())

// 	dateCond := fmt.Sprintf("(dateEntered >= [%v] OR _info/lastUpdated >= [%v])", dateStr, dateStr)
// 	sourceCond := "(source/id = 8 OR source/id = 9)"

// 	fullCond := dateCond + " AND " + sourceCond
// 	fmt.Println(fullCond)

// 	conditions := make(map[string]string)
// 	conditions["conditions"] = fullCond

// 	tickets, err := c.postTicketsCommand(ticketSearchEndpoint, conditions)
// 	if err != nil {
// 		return []Ticket{}, err
// 	}

// 	escalated := make([]Ticket, 0)
// ticketEscalated:
// 	for _, t := range tickets {
// 		audit, err := c.GetTicketAuditTrail(t.ID)
// 		if err != nil {
// 			fmt.Println(err)
// 			continue
// 		}
// 		for _, a := range audit {
// 			if a.Text == escalatedText {
// 				if date.Before(a.EnteredDate) {
// 					fmt.Println(date)
// 					escalated = append(escalated, t)
// 					break ticketEscalated
// 				}
// 			}
// 		}

// 	}
// 	return escalated, nil
// }
