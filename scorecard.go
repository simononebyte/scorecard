package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/simononebyte/scorecard/psa"
)

type config struct {
	Continuum     string       `json:"rmm_key"`
	ConnectWise   psa.Config   `json:"psa_key"`
	Excludes      psa.Excludes `json:"psa_excludes"`
	ReactiveSites []configSite `json:"reactive_endpoints"`
}

type configSite struct {
	Name     string `json:"name"`
	SiteCode string `json:"site_code"`
}

func main() {

	args := os.Args[1:]
	var week int

	switch len(args) {
	case 0:
		week = 1
	case 1:
		w, err := strconv.Atoi(args[0])
		if err != nil {
			usage()
		}
		week = w
	default:
		usage()
	}

	c, configErr := readConfig()

	if configErr != nil {
		fmt.Printf("error reading config: \n%s\n", configErr)
		os.Exit(1)
	}
	start, end := getDateRange(week)

	fmt.Printf("Gettings stats from Monday %v to Sunday %v", start, end)

	rmm := NewRMMClient(c)

	siteCodes := make([]string, 0)
	for _, v := range c.ReactiveSites {
		siteCodes = append(siteCodes, string(v.SiteCode))
	}

	excludes := make([]string, 0)
	for _, exclude := range c.Excludes.Summary {
		excludes = append(excludes, exclude)
	}

	psa := psa.NewClient(c.ConnectWise, siteCodes, excludes)

	rmmStats := rmm.GetRMMStats()
	psaStats := psa.GetStats(start, end)

	mrrRte := float64(psaStats.MRRTickets) / float64(rmmStats.TSCDevices)
	orrRte := float64(psaStats.ORRTickets) / float64(rmmStats.OtherDevices)

	fmt.Printf("\nType  - Devices - Tickets - RTE\n")
	fmt.Printf("Technology success   %3d       %3d       %f\n", rmmStats.TSCDevices, psaStats.MRRTickets, mrrRte)
	fmt.Printf("Other customers      %3d       %3d       %f\n", rmmStats.OtherDevices, psaStats.ORRTickets, orrRte)

	fmt.Printf("\n\nPress Enter to close window")
	bufio.NewReader(os.Stdin).ReadBytes('\n')
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

// getDateRange start and end dates to retrieve stats for
func getDateRange(week int) (startDate, endDate time.Time) {

	now := time.Now()

	offset := int(now.Weekday()) * -1

	if week == 0 {
		startDate = now.AddDate(0, 0, offset+1)
		endDate = now
		return
	}

	end := offset - ((week - 1) * 7)
	startDate = now.AddDate(0, 0, end-6)
	endDate = now.AddDate(0, 0, end)
	return
}
