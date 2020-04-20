package characters

import (
	"errors"
	"fmt"
	"wowstatistician/common"
	"wowstatistician/realms"

	"github.com/imroc/req"
)

// GenderDescription struct format
type GenderDescription struct {
	Male   string `json:"male"`
	Female string `json:"female"`
}

// Specialization struct format
type Specialization struct {
	Links             common.Links      `json:"_links"`
	Key               common.URL        `json:"key"`
	ID                int               `json:"id"`
	PlayableClass     Class             `json:"playable_class"`
	Name              string            `json:"name"`
	GenderDescription GenderDescription `json:"gender_description"`
	Media             common.Media      `json:"media"`
	Role              common.Value      `json:"role"`
	TalentTiers       []TalentTier      `json:"talent_tiers"`
	PvpTalents        []PvpTalent       `json:"pvp_talents"`
}

// TalentTier struct format
type TalentTier struct {
	Level     int      `json:"level"`
	Talents   []Talent `json:"talents"`
	TierIndex int      `json:"tier_index"`
}

// Talent struct format
type Talent struct {
	Talent       TalentDesc   `json:"talent"`
	SpellTooltip SpellTooltip `json:"spell_tooltip"`
	ColumnIndex  int          `json:"column_index"`
}

// PvpTalent struct format
type PvpTalent struct {
	Talent       TalentDesc   `json:"talent"`
	SpellTooltip SpellTooltip `json:"spell_tooltip"`
}

// TalentDesc struct format
type TalentDesc struct {
	Key  common.URL `json:"key"`
	Name string     `json:"name"`
	ID   int        `json:"id"`
}

// SpellTooltip struct format
type SpellTooltip struct {
	Description string `json:"description"`
	CastTime    string `json:"cast_time"`
	PowerCost   string `json:"power_cost"`
	Range       string `json:"range"`
	Cooldown    string `json:"cooldown"`
}

// Character struct format
type Character struct {
	Name          string       `json:"name"`
	ID            int          `json:"id"`
	Realm         realms.Realm `json:"realm"`
	Level         int          `json:"level"`
	PlayableClass Class        `json:"playable_class"`
	PlayableRace  Race         `json:"playable_race"`
}

// CharacterProfile struct format
type CharacterProfile struct {
	Links                  common.Links   `json:"_links"`
	ID                     int            `json:"id"`
	Name                   string         `json:"name"`
	Gender                 common.Value   `json:"gender"`
	Faction                common.Value   `json:"faction"`
	Race                   Race           `json:"race"`
	CharacterClass         Class          `json:"character_class"`
	ActiveSpec             Specialization `json:"active_spec"`
	Realm                  realms.Realm   `json:"realm"`
	Level                  int            `json:"level"`
	Experience             int            `json:"experience"`
	AchievementPoints      int            `json:"achievement_points"`
	Achievements           common.URL     `json:"achievements"`
	Titles                 common.URL     `json:"titles"`
	PvpSummary             common.URL     `json:"pvp_summary"`
	Encounters             common.URL     `json:"encounters"`
	Media                  common.URL     `json:"media"`
	LastLoginTimestamp     int            `json:"last_login_timestamp"`
	AverageItemLevel       int            `json:"average_item_level"`
	EquippedItemLevel      int            `json:"equipped_item_level"`
	Specializations        common.URL     `json:"specializations"`
	Statistics             common.URL     `json:"statistics"`
	MythicKeystoneProfile  common.URL     `json:"mythic_keystone_profile"`
	Equipment              common.URL     `json:"equipment"`
	Appearance             common.URL     `json:"appearance"`
	Collections            common.URL     `json:"collections"`
	ActiveTitle            Title          `json:"active_title"`
	Reputations            common.URL     `json:"reputations"`
	Quests                 common.URL     `json:"quests"`
	AchievementsStatistics common.URL     `json:"achievements_statistics"`
	Professions            common.URL     `json:"professions"`
}

// Race struct format
type Race struct {
	Key  common.URL `json:"key"`
	Name string     `json:"name"`
	ID   int        `json:"id"`
}

// Class struct format
type Class struct {
	Key  common.URL `json:"key"`
	Name string     `json:"name"`
	ID   int        `json:"id"`
}

// Title struct format
type Title struct {
	Key           common.URL `json:"key"`
	Name          string     `json:"name"`
	ID            int        `json:"id"`
	DisplayString string     `json:"display_string"`
}

// GetCharacterProfile return character profile for specified region, realm slug and character name slug
func GetCharacterProfile(token string, region string, realmSlug string, charName string) (*CharacterProfile, error) {
	authstr := fmt.Sprintf("Bearer %s", token)
	header := req.Header{
		"Authorization": authstr,
	}
	namespacestr := fmt.Sprintf("profile-%s", region)
	param := req.Param{
		"namespace": namespacestr,
		"locale":    "en_US",
	}
	urlstr := fmt.Sprintf("https://%s.api.blizzard.com/profile/wow/character/%s/%s", region, realmSlug, charName)
	request, err := req.Get(urlstr, header, param)
	if err != nil {
		return nil, errors.New("characters: could not retrieve character profile - " + err.Error())
	}
	var response CharacterProfile
	request.ToJSON(&response)
	return &response, nil
}
