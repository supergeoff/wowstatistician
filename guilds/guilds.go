package guilds

import (
	"errors"
	"fmt"
	"strings"
	"wowstatistician/characters"
	"wowstatistician/common"
	"wowstatistician/realms"

	"github.com/imroc/req"
)

// GuildRoster struct format
type GuildRoster struct {
	Links   common.Links `json:"_links"`
	Guild   Guild        `json:"guild"`
	Members []Member     `json:"members"`
}

// Guild struct format
type Guild struct {
	Key     common.URL   `json:"key"`
	Name    string       `json:"name"`
	ID      int          `json:"id"`
	Realm   realms.Realm `json:"realm"`
	Faction common.Value `json:"faction"`
	Slug    string       `json:"slug"`
}

// Member struct format
type Member struct {
	Character characters.Character `json:"character"`
	Rank      int                  `json:"rank"`
}

// MakeGuildSlug return guild slug for specified guild name
func MakeGuildSlug(guildName string) string {
	guildSlug := strings.ToLower(guildName)
	guildSlug = strings.ReplaceAll(guildSlug, " ", "-")
	return guildSlug
}

// GetGuildRoster return guild roster for specified region, realm slug and guild name slug
func GetGuildRoster(token string, region string, realmSlug string, guildSlug string) (*GuildRoster, error) {
	authstr := fmt.Sprintf("Bearer %s", token)
	header := req.Header{
		"Authorization": authstr,
	}
	namespacestr := fmt.Sprintf("profile-%s", region)
	param := req.Param{
		"namespace": namespacestr,
		"locale":    "en_US",
	}
	urlstr := fmt.Sprintf("https://%s.api.blizzard.com/data/wow/guild/%s/%s/roster", region, realmSlug, guildSlug)
	request, err := req.Get(urlstr, header, param)
	if err != nil {
		return nil, errors.New("guilds: could not retrieve guild roster - " + err.Error())
	}
	var response GuildRoster
	request.ToJSON(&response)
	return &response, nil
}
