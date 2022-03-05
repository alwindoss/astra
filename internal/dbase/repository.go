package dbase

import (
	"fmt"

	"go.etcd.io/bbolt"
)

type KeyValuePair struct {
	Key   []byte
	Value []byte
}

type Repository interface {
	CreateBucket(name []byte) error
	FetchBuckets() ([][]byte, error)
	DeleteBucket(name []byte) error
	GetAllData(name []byte) ([]KeyValuePair, error)
	Get(bucket, key []byte) ([]byte, error)
	Set(bucket, key, value []byte) error
}

func NewBoltDBRepository(db *bbolt.DB) Repository {
	return &boltDBRepository{
		DB: db,
		// buckets: make(map[string]*bbolt.Bucket),
	}
}

type boltDBRepository struct {
	DB *bbolt.DB
	// buckets map[string]*bbolt.Bucket
}

func (d *boltDBRepository) CreateBucket(name []byte) error {
	err := d.DB.Update(func(tx *bbolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists(name)
		if err != nil {
			return err
		}
		// d.buckets[string(name)] = b
		return nil
	})
	if err != nil {
		err = fmt.Errorf("unable to create the bucket %s: %w", name, err)
		return err
	}
	return nil
}

// func (d *boltDBRepository) getBucket(name []byte) (*bbolt.Bucket, error) {
// 	b, ok := d.buckets[string(name)]
// 	if !ok {
// 		return nil, fmt.Errorf("no bucket with the name %s", name)
// 	}
// 	return b, nil
// }

func (d *boltDBRepository) DeleteBucket(name []byte) error {
	err := d.DB.Update(func(tx *bbolt.Tx) error {
		err := tx.DeleteBucket(name)
		if err != nil {
			return err
		}
		// delete(d.buckets, string(name))
		return nil
	})
	if err != nil {
		err = fmt.Errorf("unable to delete the bucket %s: %w", name, err)
		return err
	}
	return nil
}
func (d *boltDBRepository) FetchBuckets() ([][]byte, error) {
	buckets := [][]byte{}
	err := d.DB.View(func(tx *bbolt.Tx) error {
		return tx.ForEach(func(name []byte, _ *bbolt.Bucket) error {
			fmt.Println("Bucket Name Fetched: " + string(name))
			buckets = append(buckets, name)
			return nil
		})

	})
	if err != nil {
		err = fmt.Errorf("unable to fetch the buckets: %w", err)
		return nil, err
	}
	return buckets, err
}

func (d *boltDBRepository) Get(bucket, key []byte) ([]byte, error) {
	if bucket == nil {
		return nil, fmt.Errorf("bucket name %s is not valid", bucket)
	}
	var val []byte
	err := d.DB.View(func(tx *bbolt.Tx) error {
		b := tx.Bucket(bucket)
		val = b.Get(key)
		if val == nil {
			return fmt.Errorf("key %s not found", key)
		}

		return nil
	})
	if err != nil {
		err = fmt.Errorf("unable to fetch value for the key %s from bucket %s", key, bucket)
		return nil, err
	}
	return val, nil
}

func (d *boltDBRepository) Set(bucket, key, value []byte) error {
	if bucket == nil {
		return fmt.Errorf("bucket name %s is not valid", bucket)
	}
	err := d.DB.Update(func(tx *bbolt.Tx) error {
		b := tx.Bucket(bucket)
		err := b.Put(key, value)
		if err != nil {
			return fmt.Errorf("error when adding key %s with value %s to the bucket %s", key, value, bucket)
		}
		return nil
	})
	if err != nil {
		err = fmt.Errorf("unable to add value for the key %s with value %s to the bucket %s", key, value, bucket)
		return err
	}
	return nil
}

func (d *boltDBRepository) GetAllData(bucket []byte) ([]KeyValuePair, error) {
	if bucket == nil {
		return nil, fmt.Errorf("bucket name %s is not valid", bucket)
	}
	kvPairs := []KeyValuePair{}
	err := d.DB.View(func(tx *bbolt.Tx) error {
		b := tx.Bucket(bucket)
		b.ForEach(func(k, v []byte) error {
			kvPairs = append(kvPairs, KeyValuePair{
				Key:   k,
				Value: v,
			})
			return nil
		})
		return nil
	})
	if err != nil {
		err = fmt.Errorf("unable to get all data for the bucket %s: %w", bucket, err)
		return nil, err
	}

	return kvPairs, nil
}
