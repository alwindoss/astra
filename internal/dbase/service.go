package dbase

import (
	"encoding/json"
	"fmt"
)

type Service interface {
	CreateBucket(name string) error
	DeleteBucket(name string) error
	GetBuckets() ([]string, error)
	Get(bucket, key string) (interface{}, error)
	Set(bucket, key string, value interface{}) error
}

func NewService(repo Repository) Service {
	return &boltDBService{
		Repo: repo,
	}
}

type boltDBService struct {
	Repo Repository
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
