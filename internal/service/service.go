package service

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/alwindoss/astra/internal/dbase"
)

type KeyValuePair struct {
	Key   string
	Value string
}

type Service interface {
	CreateBucket(name string) error
	DeleteBucket(name string) error
	GetBuckets() ([]string, error)
	GetAllData(bucket string) (interface{}, error)
	Get(bucket, key string) (interface{}, error)
	Set(bucket, key string, value interface{}) error
}

func NewService(repo dbase.Repository) Service {
	return &boltDBService{
		Repo: repo,
	}
}

type boltDBService struct {
	Repo dbase.Repository
}

func (s *boltDBService) CreateBucket(name string) error {
	return s.Repo.CreateBucket([]byte(name))
}

func (s *boltDBService) DeleteBucket(name string) error {
	return s.Repo.DeleteBucket([]byte(name))
}

func (s *boltDBService) GetBuckets() ([]string, error) {
	bucketBytes, err := s.Repo.FetchBuckets()
	if err != nil {
		return nil, err
	}
	fmt.Printf("Length of buckets: %d\n", len(bucketBytes))
	buckets := []string{}
	for _, name := range bucketBytes {
		buckets = append(buckets, string(name))
	}
	return buckets, nil
}

func (s *boltDBService) Get(bucket, key string) (interface{}, error) {
	value, err := s.Repo.Get([]byte(bucket), []byte(key))
	if err != nil {
		return nil, err
	}
	var data interface{}
	json.Unmarshal(value, &data)
	return &data, nil
}

func (s *boltDBService) Set(bucket, key string, value interface{}) error {
	data, err := json.Marshal(value)
	if err != nil {
		err = fmt.Errorf("unable to marshal the value: %w", err)
		return err
	}
	return s.Repo.Set([]byte(bucket), []byte(key), data)
}

func (s *boltDBService) GetAllData(bucket string) (interface{}, error) {
	dbKVPairs, err := s.Repo.GetAllData([]byte(bucket))
	if err != nil {
		return nil, err
	}
	svcKVPairs := []KeyValuePair{}
	for _, dbKVPair := range dbKVPairs {
		var valStr string
		log.Printf("Alwin: %s", valStr)
		json.Unmarshal(dbKVPair.Value, &valStr)
		svcKVPair := KeyValuePair{
			Key:   string(dbKVPair.Key),
			Value: valStr,
		}
		svcKVPairs = append(svcKVPairs, svcKVPair)

	}
	return svcKVPairs, nil
}
