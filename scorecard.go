package main

import (
	"bufio"
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/simononebyte/scorecard/psa"
	"github.com/tealeg/xlsx"
)

type config struct {
	Continuum     string         `json:"rmm_key"`
	ConnectWise   psa.Config     `json:"psa_key"`
	Boards        []configBoards `json:"psa_boards"`
	Excludes      psa.Excludes   `json:"psa_excludes"`
	ReactiveSites []configSite   `json:"reactive_endpoints"`
	StatsFile     string         `json:"stats_file"`
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

type boardStatsMap map[string]boardStats

func (m boardStatsMap) getKeys() []string {
	i := 0
	keys := make([]string, len(m))
	for k := range m {
		keys[i] = k
	}
	return keys
}

const statCount = 7

func main() {

	batchFlag := flag.Bool("batch", false, "Run in batch mode")
	flag.Parse()
	fmt.Println("Batch? ", *batchFlag)
	if isBatch(batchFlag) {
		fmt.Println("Running in Batch mode")
	} else {
		fmt.Println("Running in interractive mode")
	}
	c, configErr := readConfig()
	if configErr != nil {
		fmt.Printf("error reading config: \n%s\n", configErr)
		os.Exit(1)
	}

	psa, err := psa.NewClient(c.ConnectWise, excludeBoards)
	if err != nil {
		panic(err)
	}

	stats := boardStatsMap{}

	for _, board := range c.Boards {
		stats[board.Name] = getStatsforBoard(psa, board.ID)
	}

	// TODO
	// Wrap in if clause - only print if in interactive mode
	if isBatch(batchFlag) {
		saveStats(c, stats)
		os.Exit(0)
	}

	// Interactive mode
	printStats(c, stats)

}

func printStats(c config, stats boardStatsMap) {

	// boardWidth := maxStringLen(stats.)
	for name, stat := range stats {
		fmt.Println(name)
		fmt.Printf("  Open                : %3d\n", stat.open)
		fmt.Printf("  New                 : %3d\n", stat.new)
		fmt.Printf("  No Update in 7 days : %3d\n", stat.noUpdate7)
		fmt.Printf("  Older 7 days        : %3d\n", stat.older7)
		fmt.Printf("  Older 31 days       : %3d\n", stat.older31)
		fmt.Printf("  Assigned            : %3d\n", stat.assigned)
		fmt.Printf("  Not Assigned        : %3d\n", stat.notAssigned)
		fmt.Println("---------------------------")
	}
	fmt.Printf("\n\nPress Enter to close window")
	bufio.NewReader(os.Stdin).ReadBytes('\n')
}

func rightPadString(s string, w int) string {
	fStr := "%" + strconv.Itoa(w) + "s"
	return fmt.Sprintf(fStr, s)
}

func leftPadString(s string, w int) string {
	fStr := "%-" + strconv.Itoa(w) + "s"
	return fmt.Sprintf(fStr, s)
}

func isBatch(flag *bool) bool {
	return flag != nil && *flag == true
}

func saveStats(c config, stats boardStatsMap) error {

	// open excel file
	f, err := xlsx.OpenFile(c.StatsFile)
	if err != nil {
		panic(err)
	}

	for _, board := range c.Boards {
		sheet := getSheet(f, board.Worksheet)
		if sheet == nil {
			return fmt.Errorf("error: unable to find worksheet %v", board.Worksheet)
		}

		stat := stats[board.Name]
		var row *xlsx.Row

		if isLastRowToday(sheet) {
			row = sheet.Rows[sheet.MaxRow-1]
		} else {
			row = sheet.AddRow()
		}

		if len(row.Cells) < statCount+1 {
			for i := len(row.Cells); i < statCount+1; i++ {
				row.AddCell()
			}
		}

		row.Cells[0].SetValue(time.Now().UTC().Truncate(24 * time.Hour))
		row.Cells[1].SetValue(stat.open)
		row.Cells[2].SetValue(stat.new)
		row.Cells[3].SetValue(stat.noUpdate7)
		row.Cells[4].SetValue(stat.older7)
		row.Cells[5].SetValue(stat.older31)
		row.Cells[6].SetValue(stat.assigned)
		row.Cells[7].SetValue(stat.notAssigned)

	}

	return f.Save(c.StatsFile)
}

func isLastRowToday(sheet *xlsx.Sheet) bool {
	today := time.Now().UTC().Truncate(24 * time.Hour)
	lastRow := sheet.Rows[sheet.MaxRow-1]
	lastTime, _ := lastRow.Cells[0].GetTime(false)
	return today == lastTime
}

func getSheet(file *xlsx.File, name string) *xlsx.Sheet {
	for _, sheet := range file.Sheets {
		if sheet.Name == name {
			return sheet
		}
	}
	return nil
}

func getDDMMYYYYString(date time.Time) string {
	return fmt.Sprintf("%d//%d//%d", date.Day(), date.Month(), date.Year())
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
	fmt.Println("    scrorecard -batch")
	fmt.Println("")
	fmt.Println("    batch - Saves stats to the Excel spreadsheet specified")
	fmt.Println("            in the config file")
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

func getMaxStringLen(list []string) int {

	max := 0
	for _, v := range list {
		if len(v) > max {
			max = len(v)
		}
	}
	return max
}
