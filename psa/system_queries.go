package psa

import (
	"fmt"
	"strconv"
)

const (
	activeMembersEndpoint string = "/system/members?conditions=disableOnlineFlag=false AND type/ID!=NULL"
	auditTrailEndpoint    string = "system/audittrail?type=Ticket&id=%v"
)

// GetMembers get active members
func (c *Client) GetMembers() ([]Member, error) {

	members, err := c.getMembersCommand(activeMembersEndpoint)
	if err != nil {
		return []Member{}, err
	}

	return members, nil
}

// GetMemberName get the name from a member identifier
func (c *Client) GetMemberName(identifier string) (string, error) {

	members, err := c.getMembersCommand(activeMembersEndpoint)
	if err != nil {
		return "", err
	}

	for _, v := range members {
		if v.Identifier == identifier {
			return v.Name, nil
		}
	}

	return "", fmt.Errorf("member identifier not found")
}

// GetTicketAuditTrail gets the audit trail entries for a specific ticket
func (c *Client) GetTicketAuditTrail(ticketID int) ([]Audit, error) {

	cmd := fmt.Sprintf(auditTrailEndpoint, strconv.Itoa(ticketID))

	audit, err := c.getAuditTrailCommand(cmd)
	if err != nil {
		return []Audit{}, err
	}

	return audit, nil
}
