package leatherboards

import (
	"errors"
	"fmt"
	"wowstatistician/characters"
	"wowstatistician/common"
	"wowstatistician/dungeons"
	"wowstatistician/guilds"
	"wowstatistician/realms"

	"github.com/imroc/req"
)

// MythicLeatherboardsIndex struct format
type MythicLeatherboardsIndex struct {
	Links          common.Links  `json:"_links"`
	ID             int           `json:"id"`
	Region         realms.Region `json:"region"`
	ConnectedRealm common.URL    `json:"connected_realm"`
	Name           string        `json:"name"`
	Category       string        `json:"category"`
	Locale         string        `json:"locale"`
	Timezone       string        `json:"timezone"`
	Type           common.Value  `json:"type"`
	IsTournament   bool          `json:"is_tournament"`
	Slug           string        `json:"slug"`
}

// MythicLeatherboard struct format
type MythicLeatherboard struct {
	Links                common.Links    `json:"_links"`
	Map                  dungeons.Map    `json:"map"`
	Period               int             `json:"period"`
	PeriodStartTimestamp int             `json:"period_start_timestamp"`
	PeriodEndTimestamp   int             `json:"period_end_timestamp"`
	ConnectedRealm       common.URL      `json:"connected_realm"`
	LeadingGroups        []LeadingGroup  `json:"leading_groups"`
	KeystoneAffixes      []KeystoneAffix `json:"keystone_affixes"`
	MapChallengeModeID   int             `json:"map_challenge_mode_id"`
	Name                 string          `json:"name"`
}

// RealmsMythicLeatherboards struct format
type RealmsMythicLeatherboards struct {
	Links               common.Links  `json:"_links"`
	CurrentLeaderboards []Leaderboard `json:"current_leaderboards"`
}

// Leaderboard struct format
type Leaderboard struct {
	Key  common.URL `json:"key"`
	Name string     `json:"name"`
	ID   int        `json:"id"`
}

// KeystoneAffix struct format
type KeystoneAffix struct {
	Affix         AffixDetail `json:"keystone_affix"`
	StartingLevel int         `json:"starting_level"`
}

// AffixDetail struct format
type AffixDetail struct {
	Key  common.URL `json:"key"`
	Name string     `json:"name"`
	ID   int        `json:"id"`
}

// LeadingGroup struct format
type LeadingGroup struct {
	Ranking            int      `json:"ranking"`
	Duration           int      `json:"duration"`
	CompletedTimestamp int      `json:"completed_timestamp"`
	KeystoneLevel      int      `json:"keystone_level"`
	Members            []Member `json:"members"`
}

// Member struct format
type Member struct {
	Profile        characters.Character `json:"profile"`
	Faction        common.Value         `json:"faction"`
	Specialization Specialization       `json:"specialization"`
}

// Specialization struct format
type Specialization struct {
	Key common.URL `json:"key"`
	ID  int        `json:"id"`
}

// RaidLeatherboard struct format
type RaidLeatherboard struct {
	Links        common.Links  `json:"_links"`
	Slug         string        `json:"slug"`
	CriteriaType string        `json:"criteria_type"`
	Zone         dungeons.Zone `json:"zone"`
	Entries      []Entry       `json:"entries"`
}

// Entry struct format
type Entry struct {
	Guild                 guilds.Guild          `json:"guild"`
	Character             characters.Character  `json:"character"`
	Rating                int                   `json:"rating"`
	Faction               common.Value          `json:"faction"`
	Timestamp             int                   `json:"timestamp"`
	Region                string                `json:"region"`
	Rank                  int                   `json:"rank"`
	SeasonMatchStatistics SeasonMatchStatistics `json:"season_match_statistics"`
	Tier                  Season                `json:"tier"`
}

// PvpSeasonsIndex struct format
type PvpSeasonsIndex struct {
	Links         common.Links `json:"_links"`
	Seasons       []Season     `json:"seasons"`
	CurrentSeason Season       `json:"current_season"`
}

// Season struct format
type Season struct {
	Key common.URL `json:"key"`
	ID  int        `json:"id"`
}

// PvpSeason struct format
type PvpSeason struct {
	Links                common.Links `json:"_links"`
	ID                   int          `json:"id"`
	Leaderboards         common.URL   `json:"leaderboards"`
	Rewards              common.URL   `json:"rewards"`
	SeasonStartTimestamp int          `json:"season_start_timestamp"`
}

// PvpLeatherboards struct format
type PvpLeatherboards struct {
	Links        common.Links  `json:"_links"`
	Season       Season        `json:"season"`
	Leaderboards []Leaderboard `json:"leaderboards"`
}

// PvpLeatherboard struct format
type PvpLeatherboard struct {
	Links   common.Links `json:"_links"`
	Season  Season       `json:"season"`
	Name    string       `json:"name"`
	Bracket Bracket      `json:"bracket"`
	Entries []Entry      `json:"entries"`
}

// Bracket struct format
type Bracket struct {
	ID   int    `json:"id"`
	Type string `json:"type"`
}

// SeasonMatchStatistics struct format
type SeasonMatchStatistics struct {
	Played int `json:"played"`
	Won    int `json:"won"`
	Lost   int `json:"lost"`
}

// MythicSeasonsIndex struct format
type MythicSeasonsIndex struct {
	Links         common.Links `json:"_links"`
	Seasons       []Season     `json:"seasons"`
	CurrentSeason Season       `json:"current_season"`
}

// MythicSeason struct format
type MythicSeason struct {
	Links          common.Links `json:"_links"`
	ID             int          `json:"id"`
	StartTimestamp int          `json:"start_timestamp"`
	Periods        []Period     `json:"periods"`
}

// KeystonePeriodsIndex struct format
type KeystonePeriodsIndex struct {
	Links         common.Links `json:"_links"`
	Periods       []Period     `json:"periods"`
	CurrentPeriod Period       `json:"current_period"`
}

// Period struct format
type Period struct {
	Key common.URL `json:"key"`
	ID  int        `json:"id"`
}

// KeystonePeriod struct format
type KeystonePeriod struct {
	Links          common.Links `json:"_links"`
	ID             int          `json:"id"`
	StartTimestamp int          `json:"start_timestamp"`
	EndTimestamp   int          `json:"end_timestamp"`
}

// GetMythicKeystonePeriodsIndex return mythic keystone periods index for specified region
func GetMythicKeystonePeriodsIndex(token string, region string) (*KeystonePeriodsIndex, error) {
	authStr := fmt.Sprintf("Bearer %s", token)
	header := req.Header{
		"Authorization": authStr,
	}
	namespaceStr := fmt.Sprintf("dynamic-%s", region)
	param := req.Param{
		"namespace": namespaceStr,
		"locale":    "en_US",
	}
	urlStr := fmt.Sprintf("https://%s.api.blizzard.com/data/wow/mythic-keystone/period/index", region)
	request, err := req.Get(urlStr, header, param)
	if err != nil {
		return nil, errors.New("leatherboards: could not retrieve keystone period index - " + err.Error())
	}
	var response KeystonePeriodsIndex
	request.ToJSON(&response)
	return &response, nil
}

// GetMythicKeystonePeriod return keystone period for provided periods index url
func GetMythicKeystonePeriod(token string, url common.URL) (*KeystonePeriod, error) {
	authStr := fmt.Sprintf("Bearer %s", token)
	header := req.Header{
		"Authorization": authStr,
	}
	param := req.Param{
		"locale": "en_US",
	}
	request, err := req.Get(url.Href, header, param)
	if err != nil {
		return nil, errors.New("leatherboards: could not retrieve keystone period - " + err.Error())
	}
	var response KeystonePeriod
	request.ToJSON(&response)
	return &response, nil
}

// GetMythicSeasonsIndex return mythic seasons index for specified region
func GetMythicSeasonsIndex(token string, region string) (*MythicSeasonsIndex, error) {
	namespaceStr := fmt.Sprintf("dynamic-%s", region)
	authStr := fmt.Sprintf("Bearer %s", token)

	header := req.Header{
		"Authorization": authStr,
	}
	param := req.Param{
		"namespace": namespaceStr,
		"locale":    "en_US",
	}
	urlStr := fmt.Sprintf("https://%s.api.blizzard.com/data/wow/mythic-keystone/season/index", region)
	request, err := req.Get(urlStr, header, param)
	if err != nil {
		return nil, errors.New("leatherboards: could not retrieve mythic seasons index - " + err.Error())
	}
	var response MythicSeasonsIndex
	request.ToJSON(&response)
	return &response, nil
}

// GetMythicSeason return mythic season for provided seasons index url
func GetMythicSeason(token string, url common.URL) (*MythicSeason, error) {
	authStr := fmt.Sprintf("Bearer %s", token)
	header := req.Header{
		"Authorization": authStr,
	}
	param := req.Param{
		"locale": "en_US",
	}
	request, err := req.Get(url.Href, header, param)
	if err != nil {
		return nil, errors.New("leatherboards: could not retrieve mythic season - " + err.Error())
	}
	var response MythicSeason
	request.ToJSON(&response)
	return &response, nil
}

// GetRealmsMythicLeatherboards return mythic leatherboards list for provided connected realms url
func GetRealmsMythicLeatherboards(token string, url common.URL) (*RealmsMythicLeatherboards, error) {
	authStr := fmt.Sprintf("Bearer %s", token)
	header := req.Header{
		"Authorization": authStr,
	}
	param := req.Param{
		"locale": "en_US",
	}
	request, err := req.Get(url.Href, header, param)
	if err != nil {
		return nil, errors.New("leatherboards: could not retrieve realms mythic leatherboards list - " + err.Error())
	}
	var response RealmsMythicLeatherboards
	request.ToJSON(&response)
	return &response, nil
}

// GetMythicLeatherboard return myhtic leatherboard for provided connected realms leatherboards list url
func GetMythicLeatherboard(token string, url common.URL) (*MythicLeatherboard, error) {
	authStr := fmt.Sprintf("Bearer %s", token)
	header := req.Header{
		"Authorization": authStr,
	}
	param := req.Param{
		"locale": "en_US",
	}
	request, err := req.Get(url.Href, header, param)
	if err != nil {
		return nil, errors.New("leatherboards: could not retrieve mythic leatherboard - " + err.Error())
	}
	var response MythicLeatherboard
	request.ToJSON(&response)
	return &response, nil
}

// GetSpecifiMythicLeatherboard return mythic leatherboard for provided region, period and connected realms
func GetSpecifiMythicLeatherboard(token string, region string, period KeystonePeriod, realms realms.ConnectedRealms, dungeon dungeons.MythicDungeon) (*MythicLeatherboard, error) {
	authStr := fmt.Sprintf("Bearer %s", token)
	header := req.Header{
		"Authorization": authStr,
	}
	namespaceStr := fmt.Sprintf("dynamic-%s", region)
	param := req.Param{
		"namespace": namespaceStr,
		"locale":    "en_US",
	}
	urlStr := fmt.Sprintf("https://%s.api.blizzard.com/data/wow/connected-realm/%d/mythic-leaderboard/%d/period/%d", region, realms.ID, dungeon.ID, period.ID)
	request, err := req.Get(urlStr, header, param)
	if err != nil {
		return nil, errors.New("leatherboards: could not retrieve mythic leatherboard - " + err.Error())
	}
	var response MythicLeatherboard
	request.ToJSON(&response)
	return &response, nil
}

// GetRaidLeatherboardAlly return raid leatherboard for alliance and specified raid and region
func GetRaidLeatherboardAlly(token string, region string, raid string) (*RaidLeatherboard, error) {
	authStr := fmt.Sprintf("Bearer %s", token)
	header := req.Header{
		"Authorization": authStr,
	}
	namespaceStr := fmt.Sprintf("dynamic-%s", region)
	param := req.Param{
		"namespace": namespaceStr,
		"locale":    "en_US",
	}
	urlStr := fmt.Sprintf("https://%s.api.blizzard.com/data/wow/leaderboard/hall-of-fame/%s/alliance", region, raid)
	request, err := req.Get(urlStr, header, param)
	if err != nil {
		return nil, errors.New("leatherboards: could not retrieve ally raid leatherboard - " + err.Error())
	}
	var response RaidLeatherboard
	request.ToJSON(&response)
	return &response, nil
}

// GetRaidLeatherboardHorde return raid leatherboard for horde and specified raid and region
func GetRaidLeatherboardHorde(token string, region string, raid string) (*RaidLeatherboard, error) {
	authStr := fmt.Sprintf("Bearer %s", token)
	header := req.Header{
		"Authorization": authStr,
	}
	namespaceStr := fmt.Sprintf("dynamic-%s", region)
	param := req.Param{
		"namespace": namespaceStr,
		"locale":    "en_US",
	}
	urlStr := fmt.Sprintf("https://%s.api.blizzard.com/data/wow/leaderboard/hall-of-fame/%s/horde", region, raid)
	request, err := req.Get(urlStr, header, param)
	if err != nil {
		return nil, errors.New("leatherboards: could not retrieve horde raid leatherboard - " + err.Error())
	}
	var response RaidLeatherboard
	request.ToJSON(&response)
	return &response, nil
}

// GetRaidLeatherboard return raid leatherboard for both faction and specified raid and region - ie: nyalotha-the-waking-city
func GetRaidLeatherboard(token string, region string, raid string) (*RaidLeatherboard, error) {
	allyLeatherBoard, err := GetRaidLeatherboardAlly(token, region, raid)
	if err != nil {
		return nil, errors.New("leatherboards: could not retrieve raid leatherboard - " + err.Error())
	}
	hordeLeatherBoard, err := GetRaidLeatherboardHorde(token, region, raid)
	if err != nil {
		return nil, errors.New("leatherboards: could not retrieve raid leatherboard - " + err.Error())
	}
	mainLeatherBoard := *allyLeatherBoard
	for _, entry := range hordeLeatherBoard.Entries {
		mainLeatherBoard.Entries = append(mainLeatherBoard.Entries, entry)
	}
	return &mainLeatherBoard, nil
}

// GetPvpSeasonsIndex return pvp season index for specified region
func GetPvpSeasonsIndex(token string, region string) (*PvpSeasonsIndex, error) {
	authStr := fmt.Sprintf("Bearer %s", token)
	header := req.Header{
		"Authorization": authStr,
	}
	namespaceStr := fmt.Sprintf("dynamic-%s", region)
	param := req.Param{
		"namespace": namespaceStr,
		"locale":    "en_US",
	}
	urlStr := fmt.Sprintf("https://%s.api.blizzard.com/data/wow/pvp-season/index", region)
	request, err := req.Get(urlStr, header, param)
	if err != nil {
		return nil, errors.New("leatherboards: could not retrieve pvp seasons index - " + err.Error())
	}
	var response PvpSeasonsIndex
	request.ToJSON(&response)
	return &response, nil
}

// GetPvpSeason return pvp season for provided pvp seasons index url
func GetPvpSeason(token string, url common.URL) (*PvpSeason, error) {
	authStr := fmt.Sprintf("Bearer %s", token)
	header := req.Header{
		"Authorization": authStr,
	}
	param := req.Param{
		"locale": "en_US",
	}
	request, err := req.Get(url.Href, header, param)
	if err != nil {
		return nil, errors.New("leatherboards: could not retrieve pvp season - " + err.Error())
	}
	var response PvpSeason
	request.ToJSON(&response)
	return &response, nil
}

// GetPvpLeatherboards return pvp leatherboards list for provided pvp season url
func GetPvpLeatherboards(token string, url common.URL) (*PvpLeatherboards, error) {
	authStr := fmt.Sprintf("Bearer %s", token)
	header := req.Header{
		"Authorization": authStr,
	}
	param := req.Param{
		"locale": "en_US",
	}
	request, err := req.Get(url.Href, header, param)
	if err != nil {
		return nil, errors.New("leatherboards: could not retrieve pvp leatherboards list - " + err.Error())
	}
	var response PvpLeatherboards
	request.ToJSON(&response)
	return &response, nil
}

// GetPvpLeatherboard return pvp leatherboard for provided pvp leatherboards list url
func GetPvpLeatherboard(token string, url common.URL) (*PvpLeatherboard, error) {
	authStr := fmt.Sprintf("Bearer %s", token)
	header := req.Header{
		"Authorization": authStr,
	}
	param := req.Param{
		"locale": "en_US",
	}
	request, err := req.Get(url.Href, header, param)
	if err != nil {
		return nil, errors.New("leatherboards: could not retrieve pvp leatherboard - " + err.Error())
	}
	var response PvpLeatherboard
	request.ToJSON(&response)
	return &response, nil
}

// GetMemberSpecialization return member specialization for provided specialization url
func GetMemberSpecialization(token string, url common.URL) (*characters.Specialization, error) {
	authStr := fmt.Sprintf("Bearer %s", token)
	header := req.Header{
		"Authorization": authStr,
	}
	param := req.Param{
		"locale": "en_US",
	}
	request, err := req.Get(url.Href, header, param)
	if err != nil {
		return nil, errors.New("leatherboards: could not retrieve member specialization - " + err.Error())
	}
	var response characters.Specialization
	request.ToJSON(&response)
	return &response, nil
}
