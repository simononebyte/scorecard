package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"

	"github.com/simononebyte/scorecard/psa"
)

type config struct {
	Continuum     string         `json:"rmm_key"`
	ConnectWise   psa.Config     `json:"psa_key"`
	Boards        []configBoards `json:"psa_boards"`
	Excludes      psa.Excludes   `json:"psa_excludes"`
	ReactiveSites []configSite   `json:"reactive_endpoints"`
}

type configSite struct {
	Name     string `json:"name"`
	SiteCode string `json:"site_code"`
}

type configBoards struct {
	ID        int    `json:"id"`
	Name      string `json:"name"`
	Worksheet string `json:"worksheet"`
}

var excludeBoards = []string{
	"Planned Time Off",
}

type boardStats struct {
	open        int
	new         int
	noUpdate7   int
	older7      int
	older31     int
	assigned    int
	notAssigned int
}

func main() {

	c, configErr := readConfig()
	if configErr != nil {
		fmt.Printf("error reading config: \n%s\n", configErr)
		os.Exit(1)
	}

	psa, err := psa.NewClient(c.ConnectWise, excludeBoards)
	if err != nil {
		panic(err)
	}

	for _, board := range c.Boards {
		stats := getStatsforBoard(psa, board.ID)
		fmt.Println(board.Name, " : ", stats)
	}

	fmt.Printf("\n\nPress Enter to close window")

	bufio.NewReader(os.Stdin).ReadBytes('\n')
}

func getStatsforBoard(psa *psa.Client, id int) boardStats {
	stats := boardStats{}

	// open
	openTickets, err := psa.GetOpenTicketsByBoardID(id)
	if err != nil {
		panic("")
	}
	stats.open = len(openTickets)

	// new
	newTickets, err := psa.GetNewTicketsByBoardID(id, 7)
	if err != nil {
		panic("")
	}
	stats.new = len(newTickets)

	// noUpdate7
	noUpdate7Tickets, err := psa.GetOpenTicketsByBoardIDNotUpdatedIn(id, 7)
	if err != nil {
		panic("")
	}
	stats.noUpdate7 = len(noUpdate7Tickets)

	// older7
	older7Tickets, err := psa.GetOpenTicketsByBoardIDOlderThan(id, 7)
	if err != nil {
		panic("")
	}
	stats.older7 = len(older7Tickets)

	// older31
	older31Tickets, err := psa.GetOpenTicketsByBoardIDOlderThan(id, 31)
	if err != nil {
		panic("")
	}
	stats.older31 = len(older31Tickets)

	// assigned
	assignedTickets, err := psa.GetOpenAssignedTicketsByBoardID(id)
	if err != nil {
		panic("")
	}
	stats.assigned = len(assignedTickets)

	// notAssigned
	notAssignedTickets, err := psa.GetOpenNotAssignedTicketsByBoardID(id)
	if err != nil {
		panic("")
	}
	stats.notAssigned = len(notAssignedTickets)

	return stats
}

func usage() {
	fmt.Println("Displays rective tickets per endpoint per week.")
	fmt.Println("")
	fmt.Println("    scrorecard [week]")
	fmt.Println("")
	fmt.Println("    week - The number of weeks ago to resturn stats for")
	fmt.Println("           Default is 1 which returns the previous week")
	fmt.Println("           0 returns the current week")
	os.Exit(0)
}

func readConfig() (config, error) {
	c := config{}

	f, _ := os.Open("scorecard.json")
	defer f.Close()

	d := json.NewDecoder(f)
	err := d.Decode(&c)

	return c, err
}
