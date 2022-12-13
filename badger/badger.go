package badger

import (
	"encoding/json"
	"fmt"
	"github.com/dgraph-io/badger/v3"
	"scorepower4cours/api"
)

func IterateRecords(db *badger.DB, callback func(record api.GameRecord, i int)) error {
	return db.View(func(txn *badger.Txn) error {
		it := txn.NewIterator(badger.DefaultIteratorOptions)
		defer it.Close()
		var i int
		for it.Rewind(); it.Valid(); it.Next() {
			item := it.Item()
			err := item.Value(func(val []byte) error {
				var record api.GameRecord
				if err := json.Unmarshal(val, &record); err != nil {
					return err
				}
				fmt.Println(string(val))
				callback(record, i)
				return nil
			})
			if err != nil {
				return err
			}
			i++
		}
		return nil
	})
}
