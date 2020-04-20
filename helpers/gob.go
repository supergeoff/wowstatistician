package helpers

import (
	"bytes"
	"encoding/gob"
	"errors"
	"wowstatistician/characters"
	"wowstatistician/models"
)

// EncodeProfile encode a character profile to a byte slice
func EncodeProfile(characterProfile characters.CharacterProfile) ([]byte, error) {
	var buffer bytes.Buffer
	encoder := gob.NewEncoder(&buffer)
	err := encoder.Encode(characterProfile)
	if err != nil {
		return buffer.Bytes(), errors.New("gob: could not encode profile - " + err.Error())
	}
	return buffer.Bytes(), nil
}

// EncodeStats encode a stats map to a byte slice
func EncodeStats(stats models.Stats) ([]byte, error) {
	var buffer bytes.Buffer
	encoder := gob.NewEncoder(&buffer)
	err := encoder.Encode(stats)
	if err != nil {
		return buffer.Bytes(), errors.New("gob: could not encode stats - " + err.Error())
	}
	return buffer.Bytes(), nil
}

// DecodeProfile decode a byte slice to a character profile
func DecodeProfile(data []byte) (*characters.CharacterProfile, error) {
	var characterProfile characters.CharacterProfile
	buffer := bytes.NewBuffer(data)
	decoder := gob.NewDecoder(buffer)
	err := decoder.Decode(&characterProfile)
	if err != nil {
		return nil, errors.New("gob: could not decode profile - " + err.Error())
	}
	return &characterProfile, nil
}

// DecodeProfile decode a byte slice to a stats map
func DecodeStats(data []byte) (*models.Stats, error) {
	var stats models.Stats
	buffer := bytes.NewBuffer(data)
	decoder := gob.NewDecoder(buffer)
	err := decoder.Decode(&stats)
	if err != nil {
		return nil, errors.New("gob: could not decode stats - " + err.Error())
	}
	return &stats, nil
}
