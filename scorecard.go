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
	new         int
	open        int
	older7      int
	noUpdate7   int
	noUpdate31  int
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

	// new
	newTickets, err := psa.GetNewTicketsByBoardID(id, 7)
	if err != nil {
		panic("")
	}
	stats.new = len(newTickets)

	// open
	openTickets, err := psa.GetOpenTicketsByBoardID(id)
	if err != nil {
		panic("")
	}
	stats.open = len(openTickets)

	// older7
	older7Tickets, err := psa.GetOpenTicketsByBoardIDOlderThan(id, 7)
	if err != nil {
		panic("")
	}
	stats.older7 = len(older7Tickets)

	// noUpdate7
	noUpdate7Tickets, err := psa.GetOpenTicketsByBoardIDNotUpdatedIn(id, 31)
	if err != nil {
		panic("")
	}
	stats.noUpdate7 = len(noUpdate7Tickets)

	// noUpdate31
	noUpdate31Tickets, err := psa.GetOpenTicketsByBoardIDNotUpdatedIn(id, 31)
	if err != nil {
		panic("")
	}
	stats.noUpdate31 = len(noUpdate31Tickets)

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

// func printOpenTicketByPeson(psa *psa.Client, padding int, members []string) {
// 	fmtStr := " %-" + strconv.Itoa(padding) + "s : %3s\n"

// 	fmt.Printf("Open tickets by person\n")
// 	fmt.Printf("########################################################\n\n")

// 	for _, member := range members {
// 		tickets, err := psa.GetOpenOpenTicketsAssignedTo(member)
// 		if err != nil {
// 			fmt.Printf(fmtStr, err)
// 			continue
// 		}
// 		fmt.Printf(fmtStr, member, strconv.Itoa(len(tickets)))
// 	}
// 	fmt.Printf("\n\n")
// }

// func printOpenTicketCounts(psa *psa.Client, padding int, boards []string) {
// 	max := getWidest(boards)
// 	fmtStr := " %-" + strconv.Itoa(max) + "s : %3s\n"

// 	fmt.Printf("Open tickets by service board\n")
// 	fmt.Printf("########################################################\n\n")

// 	for _, board := range boards {
// 		tickets, err := psa.GetOpenTicketsOnBoard(board)
// 		if err != nil {
// 			fmt.Printf(fmtStr, board, "Board not found")
// 			continue
// 		}
// 		fmt.Printf(fmtStr, board, strconv.Itoa(len(tickets)))
// 	}
// 	fmt.Printf("\n\n")
// }

// func printOpenTicketsOlderThanDays(psa *psa.Client, padding int, boards []string, days int) {

// 	fmt.Printf("Open tickets by service board open more than %d days\n", days)
// 	fmt.Printf("########################################################\n\n")

// 	max := getWidest(boards)
// 	fmtStr := " %-" + strconv.Itoa(max) + "s : %3s\n"

// 	for _, board := range boards {
// 		tickets, err := psa.GetOpenTicketsOnBoardOlderThan(board, days)
// 		if err != nil {
// 			fmt.Printf(fmtStr, board, "Board not found")
// 			continue
// 		}
// 		fmt.Printf(fmtStr, board, strconv.Itoa(len(tickets)))
// 	}
// 	fmt.Printf("\n\n")
// }

// func printOpenAssignedTicketCounts(psa *psa.Client, padding int, boards []string) {
// 	max := getWidest(boards)
// 	fmtStr := " %-" + strconv.Itoa(max) + "s : %3s\n"

// 	fmt.Printf("Open tickets assigned to a resource by service board\n")
// 	fmt.Printf("########################################################\n\n")

// 	for _, board := range boards {
// 		tickets, err := psa.GetOpenAssignedTicketsOnBoard(board)
// 		if err != nil {
// 			fmt.Printf(fmtStr, board, "Board not found")
// 			continue
// 		}
// 		fmt.Printf(fmtStr, board, strconv.Itoa(len(tickets)))
// 	}
// 	fmt.Printf("\n\n")
// }

// func printOpenNotAssignedTicketCounts(psa *psa.Client, padding int, boards []string) {
// 	max := getWidest(boards)
// 	fmtStr := " %-" + strconv.Itoa(max) + "s : %3s\n"

// 	fmt.Printf("Open tickets not assigned to a resource by service board\n")
// 	fmt.Printf("########################################################\n\n")

// 	for _, board := range boards {
// 		tickets, err := psa.GetOpenNotAssignedTicketsOnBoard(board)
// 		if err != nil {
// 			fmt.Printf(fmtStr, board, "Board not found")
// 			continue
// 		}
// 		fmt.Printf(fmtStr, board, strconv.Itoa(len(tickets)))
// 	}
// 	fmt.Printf("\n\n")
// }

// func printOpenTicketsNotUpdatedInOnBoard(psa *psa.Client, padding int, boards []string, days int) {
// 	max := getWidest(boards)
// 	fmtStr := " %-" + strconv.Itoa(max) + "s : %3s\n"

// 	fmt.Printf("Open tickets not updated in %v days by service board\n", days)
// 	fmt.Printf("########################################################\n\n")

// 	for _, board := range boards {
// 		tickets, err := psa.GetOpenTicketsNotUpdatedInOnBoard(board, days)
// 		if err != nil {
// 			fmt.Printf(fmtStr, board, "Board not found")
// 			continue
// 		}
// 		fmt.Printf(fmtStr, board, strconv.Itoa(len(tickets)))
// 	}
// 	fmt.Printf("\n\n")
// }

// func printNewTicketsInLastOnBoard(psa *psa.Client, padding int, boards []string, days int) {
// 	max := getWidest(boards)
// 	fmtStr := " %-" + strconv.Itoa(max) + "s : %3v\n"

// 	fmt.Printf("New tickets in last %v days by service board\n", days)
// 	fmt.Printf("########################################################\n\n")

// 	for _, board := range boards {
// 		tickets, err := psa.GetNewTicketsInLastOnBoard(board, days)
// 		if err != nil {
// 			fmt.Printf(fmtStr, board, "Board not found")
// 			continue
// 		}
// 		fmt.Printf(fmtStr, board, len(tickets))
// 	}
// 	fmt.Printf("\n\n")
// }

// func printEscaltedAndReferredTicketsInLast(psa *psa.Client, padding int, days int) {

// 	fmtStr := " %-" + strconv.Itoa(padding) + "s : %3v\n"

// 	fmt.Printf("Tickets escalated or referred in the last %v days\n", days)
// 	fmt.Printf("########################################################\n\n")

// 	if escalated, err := psa.GetEscalatedTicketsInLast(days); err != nil {
// 		fmt.Printf(fmtStr, "Escalated", err)
// 	} else {
// 		fmt.Printf(fmtStr, "Escalated", len(escalated))
// 	}

// 	// if escalated, err := psa.GetEscalatedTicketsInLast(days); err != nil {
// 	// 	fmt.Printf(fmtStr, "Escalated", err )
// 	// } else {
// 	// 	fmt.Printf(fmtStr, "Escalated", escalated )
// 	// }

// 	fmt.Printf("\n\n")
// }

// func getWidest(list []string) int {
// 	max := 0
// 	for _, v := range list {
// 		if len(v) > max {
// 			max = len(v)
// 		}
// 	}
// 	return max
// }

// func scoreCard1() {

// 	args := os.Args[1:]
// 	var week int

// 	switch len(args) {
// 	case 0:
// 		week = 1
// 	case 1:
// 		w, err := strconv.Atoi(args[0])
// 		if err != nil {
// 			usage()
// 		}
// 		week = w
// 	default:
// 		usage()
// 	}

// 	c, configErr := readConfig()

// 	if configErr != nil {
// 		fmt.Printf("error reading config: \n%s\n", configErr)
// 		os.Exit(1)
// 	}
// 	start, end := getDateRange(week)

// 	fmt.Printf("Gettings stats from Monday %v to Sunday %v", start, end)

// 	rmm := NewRMMClient(c)

// 	siteCodes := make([]string, 0)
// 	for _, v := range c.ReactiveSites {
// 		siteCodes = append(siteCodes, string(v.SiteCode))
// 	}

// 	excludes := make([]string, 0)
// 	for _, exclude := range c.Excludes.Summary {
// 		excludes = append(excludes, exclude)
// 	}

// 	psa := psa.NewClient(c.ConnectWise, siteCodes, excludes)

// 	rmmStats := rmm.GetRMMStats()
// 	psaStats := psa.GetStats(start, end)

// 	mrrRte := float64(psaStats.MRRTickets) / float64(rmmStats.TSCDevices)
// 	orrRte := float64(psaStats.ORRTickets) / float64(rmmStats.OtherDevices)

// 	fmt.Printf("\nType  - Devices - Tickets - RTE\n")
// 	fmt.Printf("Technology success   %3d       %3d       %f\n", rmmStats.TSCDevices, psaStats.MRRTickets, mrrRte)
// 	fmt.Printf("Other customers      %3d       %3d       %f\n", rmmStats.OtherDevices, psaStats.ORRTickets, orrRte)

// 	fmt.Printf("\n\nPress Enter to close window")
// 	bufio.NewReader(os.Stdin).ReadBytes('\n')
// }

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

// // getDateRange start and end dates to retrieve stats for
// func getDateRange(week int) (startDate, endDate time.Time) {

// 	now := time.Now()

// 	offset := int(now.Weekday()) * -1

// 	if week == 0 {
// 		startDate = now.AddDate(0, 0, offset+1)
// 		endDate = now
// 		return
// 	}

// 	end := offset - ((week - 1) * 7)
// 	startDate = now.AddDate(0, 0, end-6)
// 	endDate = now.AddDate(0, 0, end)
// 	return
// }
