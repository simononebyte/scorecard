package main

import (
	"fmt"
	"os"
)

// RMMClient encapsulates the RMM API Client
type RMMClient struct {
	apiClient     *APIClient
	reactiveSites []rmmSite
}

type rmmSite struct {
	Name     string
	SiteCode string
}

// NewRMMClient ...
func NewRMMClient(c config) *RMMClient {
	rmm := RMMClient{}
	rmm.apiClient = NewAPICLient(c.Continuum)
	for _, v := range c.ReactiveSites {
		rmm.reactiveSites = append(rmm.reactiveSites, rmmSite{v.Name, v.SiteCode})
	}
	return &rmm
}

// RMMStats ...
type RMMStats struct {
	TSCDevices   int
	OtherDevices int
}

// RMMSite ...
type RMMSite struct {
	Name     string `json:"name"`
	SiteCode string `json:"siteCode"`
}

// RMMDevice ...
type RMMDevice struct {
	MachineID           string `json:"machineID"`
	MachineName         string `json:"machineName"`
	FriendlyName        string `json:"friendlyName"`
	HostName            string `json:"hostName"`
	LastUpdated         string `json:"lastUpdated "`
	LastSeenOnline      string `json:"lastSeenOnline"`
	LastReboot          string `json:"lastReboot"`
	AssetType           string `json:"assetType"`
	CompanyID           string `json:"companyId"`
	CompanyName         string `json:"companyName"`
	SiteID              string `json:"siteId"`
	SiteCode            string `json:"siteCode"`
	SiteName            string `json:"siteName"`
	LastLoginUser       string `json:"lastLoginUser"`
	ManufacturerName    string `json:"manufacturerName"`
	ModelNumber         string `json:"modelNumber"`
	ModelSKU            string `json:"modelSKU"`
	SupportSerialNumber string `json:"supportSerialNumber"`
	OperatingSystem     string `json:"operatingSystem"`
	CPU                 string `json:"cpu"`
	MmemoryTotal        string `json:"memoryTotal"`
}

// GetRMMStats ...
func (rmm *RMMClient) GetRMMStats() RMMStats {

	stats := RMMStats{}

	fmt.Printf("RMM: Getting Sites: ")
	sites, sitesErr := rmm.GetRMMSites()
	if sitesErr != nil {
		fmt.Printf("error: \n%s\n", sitesErr)
		os.Exit(1)
	}
	fmt.Printf("%d returned\n", len(sites))

	fmt.Printf("RMM: Getting device counts: ")
	for _, v := range sites {
		devs, devsErr := rmm.GetRMMEndpoints(v.SiteCode)
		if devsErr != nil {
			fmt.Printf("error: %s\n%s\n", v.Name, devsErr)
			os.Exit(1)
		}
		fmt.Printf(".")
		if rmm.IsTSCSite(v.Name) == true {
			stats.TSCDevices += len(devs)
			continue
		}
		stats.OtherDevices += len(devs)
	}
	fmt.Printf("\n")
	return stats
}

// GetRMMSites ..
func (rmm *RMMClient) GetRMMSites() ([]RMMSite, error) {

	url := "https://api.itsupport247.net/reporting/v1/sites"
	sites := []RMMSite{}

	if err := rmm.apiClient.Get(url, &sites); err != nil {
		return sites, err
	}

	return sites, nil
}

// GetRMMEndpoints ...
func (rmm *RMMClient) GetRMMEndpoints(siteCode string) ([]RMMDevice, error) {

	url := fmt.Sprintf("https://itsapi.itsupport247.net/reporting/v1/sites/%s/devices/", siteCode)
	devices := []RMMDevice{}

	if err := rmm.apiClient.Get(url, &devices); err != nil {
		return devices, err
	}

	return devices, nil
}

// IsTSCSite determines if site is a Technology Success Customer site
func (rmm *RMMClient) IsTSCSite(site string) bool {
	for _, v := range rmm.reactiveSites {
		if v.Name == site {
			return true
		}
	}
	return false
}

// func doRMMGet(url string, key string, out interface{}) error {

// 	b64key := b64.StdEncoding.EncodeToString([]byte(key))

// 	req, reqErr := http.NewRequest(http.MethodGet, url, nil)
// 	if reqErr != nil {
// 		return reqErr
// 	}

// 	req.Header.Set("Authorization", fmt.Sprintf("Basic %s", b64key))

// 	doErr := doJSONRequest(req, out)
// 	if doErr != nil {
// 		return doErr
// 	}

// 	return nil
// }
