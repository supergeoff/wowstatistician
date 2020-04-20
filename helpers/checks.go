package helpers

import "wowstatistician/characters"

// CheckValidProfile check a profile and return true if it contains all the variable requiered for stats
func CheckValidProfile(characterProfile characters.CharacterProfile) bool {
	if characterProfile.Name != "" && characterProfile.ActiveSpec.Name != "" && characterProfile.CharacterClass.Name != "" && characterProfile.ID != 0 {
		return true
	}
	return false
}
