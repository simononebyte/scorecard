package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
)

type config struct {
	Continuum     string       `json:"rmm_key"`
	ConnectWise   PSAAuth      `json:"psa_key"`
	Excludes      PSAExcludes  `json:"psa_excludes"`
	ReactiveSites []configSite `json:"reactive_endpoints"`
}

type configSite struct {
	Name     string `json:"name"`
	SiteCode string `json:"site_code"`
}

// PSAAuth hold the APIP credentials for the PSA system
type PSAAuth struct {
	Company  string `json:"company"`
	Username string `json:"public"`
	Password string `json:"private"`
}

// PSAExcludes detail tickets that should be exluded from scorecard
type PSAExcludes struct {
	Summary []string `json:"summary"`
}

func main() {

	c, configErr := readConfig()

	if configErr != nil {
		fmt.Printf("error reading config: \n%s\n", configErr)
		os.Exit(1)
	}
	rmm := NewRMMClient(c)
	psa := NewPSAClient(c)

	rmmStats := rmm.GetRMMStats()
	psaStats := psa.GetPSAStats()

	mrrRte := float64(psaStats.MRRTickets) / float64(rmmStats.TSCDevices)
	orrRte := float64(psaStats.ORRTickets) / float64(rmmStats.OtherDevices)

	fmt.Printf("\nType  - Devices - Tickets - RTE\n")
	fmt.Printf("Technology success   %3d       %3d       %f\n", rmmStats.TSCDevices, psaStats.MRRTickets, mrrRte)
	fmt.Printf("Other customers      %3d       %3d       %f\n", rmmStats.OtherDevices, psaStats.ORRTickets, orrRte)

	fmt.Printf("\n\nPress Enter to close window")
	bufio.NewReader(os.Stdin).ReadBytes('\n')
}

func readConfig() (config, error) {
	c := config{}

	f, _ := os.Open("scorecard.json")
	defer f.Close()

	d := json.NewDecoder(f)
	err := d.Decode(&c)

	return c, err
}
