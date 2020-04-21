package databases

import (
	"errors"
	"log"
	"strconv"
	"time"
	"wowstatistician/characters"
	"wowstatistician/helpers"
	"wowstatistician/models"

	"github.com/dgraph-io/badger/v2"
)

// WriteProfileToDb write a character profile to a db provided db pointer
func WriteProfileToDb(db *badger.DB, characterProfile characters.CharacterProfile) error {
	ID := strconv.Itoa(characterProfile.ID)
	data, err := helpers.EncodeProfile(characterProfile)
	if err != nil {
		return errors.New("databases: could not write profile to db - " + err.Error())
	}
	err = db.Update(func(txn *badger.Txn) error {
		err := txn.Set([]byte(ID), data)
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return errors.New("maidatabasesn: could not write profile to db - " + err.Error())
	}
	return nil
}

// WriteStatsToDb write stats struct to the stats db provided a dbname to compute stats against
func WriteStatsToDb(stats models.Stats, dbname string) error {
	db, err := OpenDB("databases/stats")
	if err != nil {
		return errors.New("databases: could not write stats for db: " + dbname + " - " + err.Error())
	}
	defer db.Close()
	stats.Source = dbname
	stats.SyncDate = time.Now().Format("01-02-2006")
	data, err := helpers.EncodeStats(stats)
	if err != nil {
		return errors.New("databases: could not write stats for db: " + dbname + " - " + err.Error())
	}
	err = db.Update(func(txn *badger.Txn) error {
		err := txn.Set([]byte(dbname), data)
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return errors.New("databases: could not write stats for db: " + dbname + " - " + err.Error())
	}
	return nil
}

// ReadProfileFromDb read a character profile from a db provided db pointer
func ReadProfileFromDb(db *badger.DB, ID int) (*characters.CharacterProfile, error) {
	var data []byte
	strID := strconv.Itoa(ID)
	err := db.View(func(txn *badger.Txn) error {
		item, err := txn.Get([]byte(strID))
		if err != nil {
			return err
		}
		data, err = item.ValueCopy(nil)
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return nil, errors.New("databases: could not read profile from db - " + err.Error())
	}
	characterProfile, err := helpers.DecodeProfile(data)
	if err != nil {
		return nil, errors.New("databases: could not read profile from db - " + err.Error())
	}
	return characterProfile, nil
}

// ReadStatsDb read stats struct from the stats db provided a dbname
func ReadStatsDb(dbname string) (*models.Stats, error) {
	db, err := OpenDB("databases/stats")
	if err != nil {
		return nil, errors.New("databases: could not read stats db - " + err.Error())
	}
	defer db.Close()
	var data []byte
	err = db.View(func(txn *badger.Txn) error {
		item, err := txn.Get([]byte(dbname))
		if err != nil {
			return err
		}
		data, err = item.ValueCopy(nil)
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return nil, errors.New("databases: could not read stats db - " + err.Error())
	}
	stats, err := helpers.DecodeStats(data)
	if err != nil {
		return nil, errors.New("databases: could not read stats db - " + err.Error())
	}
	return stats, nil
}

func GetStatsFromDb(db *badger.DB, dbname string) (*models.Stats, error) {
	var data []byte
	err := db.View(func(txn *badger.Txn) error {
		item, err := txn.Get([]byte(dbname))
		if err != nil {
			return err
		}
		data, err = item.ValueCopy(nil)
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return nil, errors.New("databases: could not read stats from db - " + err.Error())
	}
	stats, err := helpers.DecodeStats(data)
	if err != nil {
		return nil, errors.New("databases: could not read stats from db - " + err.Error())
	}
	return stats, nil
}

// GenerateStatistics generate stats for a db provided a db pointer
func GenerateStatistics(db *badger.DB) (*models.Stats, error) {
	stats := &models.Stats{}
	err := db.View(func(tnx *badger.Txn) error {
		options := badger.DefaultIteratorOptions
		iterator := tnx.NewIterator(options)
		defer iterator.Close()
		for iterator.Rewind(); iterator.Valid(); iterator.Next() {
			item := iterator.Item()
			data, err := item.ValueCopy(nil)
			if err != nil {
				log.Println(err)
				continue
			}
			characterProfile, err := helpers.DecodeProfile(data)
			if err != nil {
				log.Println(err)
				continue
			}
			distrib := stats.FindDistribution(characterProfile.CharacterClass.Name)
			if distrib == nil {
				distrib = &models.Distribution{
					Class: characterProfile.CharacterClass.Name,
					Total: 1,
				}
				stats.Distributions = append(stats.Distributions, distrib)
			} else {
				distrib.Total++
			}
			spec := distrib.FindSpec(characterProfile.ActiveSpec.Name)
			if spec == nil {
				spec = &models.Spec{
					Spec:  characterProfile.ActiveSpec.Name,
					Count: 1,
				}
				distrib.Specs = append(distrib.Specs, spec)
			} else {
				spec.Count++
			}
			stats.Overall++
		}
		return nil
	})
	if err != nil {
		return nil, errors.New("databases: could not generate stats from db - " + err.Error())
	}
	return stats, nil
}

// WriteStatsForDb compute and write stats for provided path db
func WriteStatsForDb(dbname string) error {
	db, err := OpenDB("databases/" + dbname)
	if err != nil {
		return errors.New("databases: could not save stats for db " + dbname + " - " + err.Error())
	}
	defer db.Close()
	stats, err := GenerateStatistics(db)
	if err != nil {
		return errors.New("databases: could not save stats for db " + dbname + " - " + err.Error())
	}
	err = WriteStatsToDb(*stats, dbname)
	if err != nil {
		return errors.New("databases: could not save stats for db " + dbname + " - " + err.Error())
	}
	return nil
}

// OpenDB open a db at provided path and return a db pointer
func OpenDB(path string) (*badger.DB, error) {
	options := badger.DefaultOptions(path)
	options.Logger = nil
	options.Truncate = true
	db, err := badger.Open(options)
	if err != nil {
		return nil, errors.New("databases: could not open db - " + err.Error())
	}
	return db, nil
}
