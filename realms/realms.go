package realms

import (
	"errors"
	"fmt"
	"wowstatistician/common"

	"github.com/imroc/req"
)

// Region struct format
type Region struct {
	Key  common.URL `json:"key"`
	Name string     `json:"name"`
	ID   int        `json:"id"`
}

// Realm struct format
type Realm struct {
	Links          common.Links `json:"_links"`
	Key            common.URL   `json:"key"`
	ID             int          `json:"id"`
	Region         Region       `json:"region"`
	ConnectedRealm common.URL   `json:"connected_realm"`
	Name           string       `json:"name"`
	Category       string       `json:"category"`
	Locale         string       `json:"locale"`
	Timezone       string       `json:"timezone"`
	Type           common.Value `json:"type"`
	IsTournament   bool         `json:"is_tournament"`
	Slug           string       `json:"slug"`
}

// ConnectedRealms struct format
type ConnectedRealms struct {
	Links              common.Links `json:"_links"`
	ID                 int          `json:"id"`
	HasQueue           bool         `json:"has_queue"`
	Status             common.Value `json:"status"`
	Population         common.Value `json:"population"`
	Realms             []Realm      `json:"realms"`
	MythicLeaderboards common.URL   `json:"mythic_leaderboards"`
	Auctions           common.URL   `json:"auctions"`
}

// RealmsIndex struct format
type RealmsIndex struct {
	Links  common.Links `json:"_links"`
	Realms []Realm      `json:"realms"`
}

// ConnectedRealmsIndex struct format
type ConnectedRealmsIndex struct {
	Links           common.Links `json:"_links"`
	ConnectedRealms []common.URL `json:"connected_realms"`
}

// GetConnectedRealmsIndex return connected realms index for specified region
func GetConnectedRealmsIndex(token string, region string) (*ConnectedRealmsIndex, error) {
	namespacestr := fmt.Sprintf("dynamic-%s", region)
	authstr := fmt.Sprintf("Bearer %s", token)
	header := req.Header{
		"Authorization": authstr,
	}
	param := req.Param{
		"namespace": namespacestr,
		"locale":    "en_US",
	}
	urlstr := fmt.Sprintf("https://%s.api.blizzard.com/data/wow/connected-realm/index", region)
	request, err := req.Get(urlstr, header, param)
	if err != nil {
		return nil, errors.New("realms: could not retrieve connected realms index - " + err.Error())
	}
	var response ConnectedRealmsIndex
	request.ToJSON(&response)
	return &response, nil
}

// GetConnectedRealms return connected realms provided index url
func GetConnectedRealms(token string, url common.URL) (*ConnectedRealms, error) {
	authstr := fmt.Sprintf("Bearer %s", token)
	header := req.Header{
		"Authorization": authstr,
	}
	param := req.Param{
		"locale": "en_US",
	}
	request, err := req.Get(url.Href, header, param)
	if err != nil {
		return nil, errors.New("realms: could not retrieve connected realms - " + err.Error())
	}
	var response ConnectedRealms
	request.ToJSON(&response)
	return &response, nil
}

// GetRealmsIndex return realm index for specified region
func GetRealmsIndex(token string, region string) (*RealmsIndex, error) {
	authstr := fmt.Sprintf("Bearer %s", token)
	header := req.Header{
		"Authorization": authstr,
	}
	namespacestr := fmt.Sprintf("dynamic-%s", region)
	param := req.Param{
		"namespace": namespacestr,
		"locale":    "en_US",
	}
	urlstr := fmt.Sprintf("https://%s.api.blizzard.com/data/wow/realm/index", region)
	request, err := req.Get(urlstr, header, param)
	if err != nil {
		return nil, errors.New("realms: could not retrieve realms index - " + err.Error())
	}
	var response RealmsIndex
	request.ToJSON(&response)
	return &response, nil
}

// GetRealm return realm provided index url
func GetRealm(token string, url common.URL) (*Realm, error) {
	authstr := fmt.Sprintf("Bearer %s", token)
	header := req.Header{
		"Authorization": authstr,
	}
	param := req.Param{
		"locale": "en_US",
	}
	request, err := req.Get(url.Href, header, param)
	if err != nil {
		return nil, errors.New("realms: could not retrieve realm - " + err.Error())
	}
	var response Realm
	request.ToJSON(&response)
	return &response, nil
}
