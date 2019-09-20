package boltdb

import (
	"strconv"
	"time"

	"github.com/boltdb/bolt"
)

//Save Save a Bucket, Key/Value
func Save(bucket string, key string, value string) {

	db, _ := bolt.Open("my.db", 0600, nil)

	defer db.Close()

	db.Update(func(tx *bolt.Tx) error {

		b, _ := tx.CreateBucketIfNotExists([]byte(bucket))

		err := b.Put([]byte(key), []byte(value))

		return err

	})

}

//Get Get a Bucket, Key/Value
func Get(bucket string, key string) string {

	var result = ""

	db, _ := bolt.Open("my.db", 0600, nil)

	defer db.Close()

	db.View(func(tx *bolt.Tx) error {

		b := tx.Bucket([]byte(bucket))

		if b != nil {

			v := b.Get([]byte(key))

			result = string(v)

		}

		return nil
	})

	return result

}

//GetTimestamp Get the current timestamp in unix format.
func GetTimestamp() string {

	var s string

	timestamp := time.Now().Unix()

	s = strconv.FormatInt(timestamp, 10) // use base 10 for sanity purpose

	return s

}
