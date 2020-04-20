package cmd

import (
	"errors"
	"fmt"
	"log"
	"strings"
	"wowstatistician/auth"
	"wowstatistician/characters"
	"wowstatistician/guilds"
	"wowstatistician/helpers"
	"wowstatistician/helpers/databases"
	"wowstatistician/leatherboards"
	"wowstatistician/realms"
)

// SaveRaidProfiles save player profiles from raid leatherboard to a db
func SaveRaidProfiles(region string, raid string) error {
	token, err := auth.CreateToken()
	if err != nil {
		return errors.New("cmd: could not save raid profiles - " + err.Error())
	}
	db, err := databases.OpenDB("databases/raid")
	if err != nil {
		return errors.New("cmd: could not save raid profiles - " + err.Error())
	}
	defer db.Close()
	fmt.Printf("--- Getting leatherboard for region: %v and raid: %v ---\n", region, raid)
	raidLeatherboard, err := leatherboards.GetRaidLeatherboard(token, region, raid)
	if err != nil {
		return errors.New("cmd: could not save raid profiles - " + err.Error())
	}
	entriesNumber := 0
	// Loop:
	for _, entry := range raidLeatherboard.Entries {
		if entry.Region == region {
			entry.Guild.Slug = guilds.MakeGuildSlug(entry.Guild.Name)
			fmt.Printf("--- Getting roster for guild: %v from: %v ---\n", entry.Guild.Slug, entry.Guild.Realm.Slug)
			guildRoster, err := guilds.GetGuildRoster(token, region, entry.Guild.Realm.Slug, entry.Guild.Slug)
			if err != nil {
				log.Println(err)
				continue
			}
			for _, member := range guildRoster.Members {
				if member.Character.Level == 120 {
					characterProfile, err := characters.GetCharacterProfile(token, region, entry.Guild.Realm.Slug, strings.ToLower(member.Character.Name))
					if err != nil {
						log.Println(err)
						continue
					}
					if helpers.CheckValidProfile(*characterProfile) {
						fmt.Printf("Saving %v as a %v %v with id: %v\n", characterProfile.Name, characterProfile.ActiveSpec.Name, characterProfile.CharacterClass.Name, characterProfile.ID)
						err := databases.WriteProfileToDb(db, *characterProfile)
						if err != nil {
							log.Println(err)
							continue
						}
						entriesNumber++
					}
					// if entriesNumber >= 10 {
					// 	break Loop
					// }
				}
			}
		}
	}
	fmt.Printf("--- Added %v entries ---\n", entriesNumber)
	return nil
}

// SaveMythicProfiles save player profiles from mythic leatherboard to a db
func SaveMythicProfiles(region string) error {
	token, err := auth.CreateToken()
	if err != nil {
		return errors.New("cmd: could not save mythic profiles - " + err.Error())
	}
	db, err := databases.OpenDB("databases/mythic")
	if err != nil {
		return errors.New("cmd: could not save mythic profiles - " + err.Error())
	}
	defer db.Close()
	fmt.Printf("--- Getting connected realms index for region: %v ---\n", region)
	connectedRealmsIndex, err := realms.GetConnectedRealmsIndex(token, region)
	if err != nil {
		return errors.New("cmd: could not save mythic profiles - " + err.Error())
	}
	entriesNumber := 0
	// Loop:
	for _, url := range connectedRealmsIndex.ConnectedRealms {
		connectedRealms, err := realms.GetConnectedRealms(token, url)
		if err != nil {
			log.Println(err)
			continue
		}
		if connectedRealms.ID != 0 {
			fmt.Printf("------ Getting leatherboards for connected realms: %v ------\n", connectedRealms.ID)
			realmsMythicLeatherboards, err := leatherboards.GetRealmsMythicLeatherboards(token, connectedRealms.MythicLeaderboards)
			if err != nil {
				log.Println(err)
				continue
			}
			for _, boards := range realmsMythicLeatherboards.CurrentLeaderboards {
				leatherboard, err := leatherboards.GetMythicLeatherboard(token, boards.Key)
				if err != nil {
					log.Println(err)
					continue
				}
				if leatherboard.Name != "" {
					fmt.Printf("--- Getting details for leatherboard: %v ---\n", leatherboard.Name)
					for _, group := range leatherboard.LeadingGroups {
						for _, member := range group.Members {
							var characterProfile characters.CharacterProfile
							characterProfile.Name = member.Profile.Name
							characterProfile.ID = member.Profile.ID
							characterProfile.Realm = member.Profile.Realm
							characterProfile.Level = member.Profile.Level
							characterProfile.Race = member.Profile.PlayableRace
							characterProfile.Faction = member.Faction
							activeSpec, err := leatherboards.GetMemberSpecialization(token, member.Specialization.Key)
							if err != nil {
								log.Println(err)
								continue
							}
							characterProfile.ActiveSpec = *activeSpec
							characterProfile.CharacterClass = characterProfile.ActiveSpec.PlayableClass
							if helpers.CheckValidProfile(characterProfile) {
								fmt.Printf("Saving %v as a %v %v with id: %v\n", characterProfile.Name, characterProfile.ActiveSpec.Name, characterProfile.CharacterClass.Name, characterProfile.ID)
								err := databases.WriteProfileToDb(db, characterProfile)
								if err != nil {
									log.Println(err)
									continue
								}
								entriesNumber++
							}
							// if entriesNumber >= 10 {
							// 	break Loop
							// }
						}
					}
				}
			}
		}
	}
	fmt.Printf("--- Added %v entries ---\n", entriesNumber)
	return nil
}

// SaveArenaProfiles save player profiles from arena leatherboard to a db
func SaveArenaProfiles(region string) error {
	token, err := auth.CreateToken()
	if err != nil {
		return errors.New("cmd: could not save arena profiles - " + err.Error())
	}
	db, err := databases.OpenDB("databases/arena")
	if err != nil {
		return errors.New("cmd: could not save arena profiles - " + err.Error())
	}
	defer db.Close()
	fmt.Printf("--- Getting pvp season index for region: %v ---\n", region)
	pvpSeasonsIndex, err := leatherboards.GetPvpSeasonsIndex(token, region)
	if err != nil {
		return errors.New("cmd: could not save arena profiles - " + err.Error())
	}
	fmt.Printf("--- Getting current pvp season---\n")
	pvpSeason, err := leatherboards.GetPvpSeason(token, pvpSeasonsIndex.CurrentSeason.Key)
	if err != nil {
		return errors.New("cmd: could not save arena profiles - " + err.Error())
	}
	fmt.Printf("--- Getting pvp leatherboards---\n")
	pvpLeatherboards, err := leatherboards.GetPvpLeatherboards(token, pvpSeason.Leaderboards)
	if err != nil {
		return errors.New("cmd: could not save arena profiles - " + err.Error())
	}
	entriesNumber := 0
	for _, leatherboard := range pvpLeatherboards.Leaderboards {
		if leatherboard.Name == "2v2" {
			fmt.Printf("--- Getting 2v2 data---\n")
			vTwoLeatherboard, err := leatherboards.GetPvpLeatherboard(token, leatherboard.Key)
			if err != nil {
				log.Println(err)
			} else {
				// Loop:
				for _, entry := range vTwoLeatherboard.Entries {
					characterProfile, err := characters.GetCharacterProfile(token, region, entry.Character.Realm.Slug, strings.ToLower(entry.Character.Name))
					if err != nil {
						log.Println(err)
						continue
					}
					if helpers.CheckValidProfile(*characterProfile) {
						fmt.Printf("Saving %v as a %v %v with id: %v\n", characterProfile.Name, characterProfile.ActiveSpec.Name, characterProfile.CharacterClass.Name, characterProfile.ID)
						err := databases.WriteProfileToDb(db, *characterProfile)
						if err != nil {
							log.Println(err)
							continue
						}
						entriesNumber++
					}
					// if entriesNumber >= 10 {
					// 	break Loop
					// }
				}
			}
		}
		if leatherboard.Name == "3v3" {
			fmt.Printf("--- Getting 3v3 data---\n")
			vThreeLeatherboard, err := leatherboards.GetPvpLeatherboard(token, leatherboard.Key)
			if err != nil {
				log.Println(err)
			} else {
				// Loop:
				for _, entry := range vThreeLeatherboard.Entries {
					characterProfile, err := characters.GetCharacterProfile(token, region, entry.Character.Realm.Slug, strings.ToLower(entry.Character.Name))
					if err != nil {
						log.Println(err)
						continue
					}
					if helpers.CheckValidProfile(*characterProfile) {
						fmt.Printf("Saving %v as a %v %v with id: %v\n", characterProfile.Name, characterProfile.ActiveSpec.Name, characterProfile.CharacterClass.Name, characterProfile.ID)
						err := databases.WriteProfileToDb(db, *characterProfile)
						if err != nil {
							log.Println(err)
							continue
						}
						entriesNumber++
					}
					// if entriesNumber >= 10 {
					// 	break Loop
					// }
				}
			}
		}
	}
	fmt.Printf("--- Added %v entries ---\n", entriesNumber)
	return nil
}

// SaveRbgProfiles save player profiles from rbg leatherboard to a db
func SaveRbgProfiles(region string) error {
	token, err := auth.CreateToken()
	if err != nil {
		return errors.New("cmd: could not save rbg profiles - " + err.Error())
	}
	db, err := databases.OpenDB("databases/rbg")
	if err != nil {
		return errors.New("cmd: could not save rbg profiles - " + err.Error())
	}
	defer db.Close()
	fmt.Printf("--- Getting pvp season index for region: %v ---\n", region)
	pvpSeasonsIndex, err := leatherboards.GetPvpSeasonsIndex(token, region)
	if err != nil {
		return errors.New("cmd: could not save rbg profiles - " + err.Error())
	}
	fmt.Printf("--- Getting current pvp season---\n")
	pvpSeason, err := leatherboards.GetPvpSeason(token, pvpSeasonsIndex.CurrentSeason.Key)
	if err != nil {
		return errors.New("cmd: could not save rbg profiles - " + err.Error())
	}
	fmt.Printf("--- Getting pvp leatherboards---\n")
	pvpLeatherboards, err := leatherboards.GetPvpLeatherboards(token, pvpSeason.Leaderboards)
	if err != nil {
		return errors.New("cmd: could not save rbg profiles - " + err.Error())
	}
	entriesNumber := 0
	for _, leatherboard := range pvpLeatherboards.Leaderboards {
		if leatherboard.Name == "rbg" {
			fmt.Printf("--- Getting rbg data---\n")
			vTwoLeatherboard, err := leatherboards.GetPvpLeatherboard(token, leatherboard.Key)
			if err != nil {
				log.Println(err)
			} else {
				// Loop:
				for _, entry := range vTwoLeatherboard.Entries {
					characterProfile, err := characters.GetCharacterProfile(token, region, entry.Character.Realm.Slug, strings.ToLower(entry.Character.Name))
					if err != nil {
						log.Println(err)
						continue
					}
					if helpers.CheckValidProfile(*characterProfile) {
						fmt.Printf("Saving %v as a %v %v with id: %v\n", characterProfile.Name, characterProfile.ActiveSpec.Name, characterProfile.CharacterClass.Name, characterProfile.ID)
						err := databases.WriteProfileToDb(db, *characterProfile)
						if err != nil {
							log.Println(err)
							continue
						}
						entriesNumber++
					}
					// if entriesNumber >= 10 {
					// 	break Loop
					// }
				}
			}
		}
	}
	fmt.Printf("--- Added %v entries ---\n", entriesNumber)
	return nil
}
