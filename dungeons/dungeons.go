package dungeons

import (
	"errors"
	"fmt"
	"wowstatistician/common"

	"github.com/imroc/req"
)

// MythicDungeonsIndex struct format
type MythicDungeonsIndex struct {
	Links    common.Links `json:"_links"`
	Dungeons []Dungeon    `json:"dungeons"`
}

// Dungeon struct format
type Dungeon struct {
	Key  common.URL `json:"key"`
	Name string     `json:"name"`
	ID   int        `json:"id"`
}

// MythicDungeon struct format
type MythicDungeon struct {
	Links            common.Links      `json:"_links"`
	ID               int               `json:"id"`
	Name             string            `json:"name"`
	Map              Map               `json:"map"`
	Zone             Zone              `json:"zone"`
	Dungeon          Dungeon           `json:"dungeon"`
	KeystoneUpgrades []KeystoneUpgrade `json:"keystone_upgrades"`
}

// KeystoneUpgrade struct format
type KeystoneUpgrade struct {
	UpgradeLevel       int `json:"upgrade_level"`
	QualifyingDuration int `json:"qualifying_duration"`
}

// Zone struct format
type Zone struct {
	Key  common.URL `json:"key"`
	Name string     `json:"name"`
}

// Map struct format
type Map struct {
	Name string `json:"name"`
	ID   int    `json:"id"`
}

// GetMythicDungeonsIndex return mythic dungeons index for specified region
func GetMythicDungeonsIndex(token string, region string) (*MythicDungeonsIndex, error) {
	authstr := fmt.Sprintf("Bearer %s", token)
	header := req.Header{
		"Authorization": authstr,
	}
	namespacestr := fmt.Sprintf("dynamic-%s", region)
	param := req.Param{
		"namespace": namespacestr,
		"locale":    "en_US",
	}
	urlstr := fmt.Sprintf("https://%s.api.blizzard.com/data/wow/mythic-keystone/dungeon/index", region)
	request, err := req.Get(urlstr, header, param)
	if err != nil {
		return nil, errors.New("dungeons: could not retrieve mythic dungeons index - " + err.Error())
	}
	var response MythicDungeonsIndex
	request.ToJSON(&response)
	return &response, nil
}

// GetMythicDungeon return mythic dungeon for provided dungeons index url
func GetMythicDungeon(token string, url common.URL) (*MythicDungeon, error) {
	authstr := fmt.Sprintf("Bearer %s", token)
	header := req.Header{
		"Authorization": authstr,
	}
	param := req.Param{
		"locale": "en_US",
	}
	request, err := req.Get(url.Href, header, param)
	if err != nil {
		return nil, errors.New("dungeons: could not retrieve mythic dungeon - " + err.Error())
	}
	var response MythicDungeon
	request.ToJSON(&response)
	return &response, nil
}
