package dbase

import (
	"bytes"
	"path/filepath"
	"testing"

	"go.etcd.io/bbolt"
)

func TestRepository(t *testing.T) {
	tempDir := t.TempDir()
	dbLoc := filepath.Join(tempDir, "test.db")
	db, err := bbolt.Open(dbLoc, 0600, nil)
	if err != nil {
		t.Logf("expected error to be nil but was %v", err)
		t.FailNow()
	}
	r := NewBoltDBRepository(db)
	bucketName := []byte("mybucket")
	if err = r.CreateBucket(bucketName); err != nil {
		t.Logf("expected error to be nil but was %v", err)
		t.FailNow()
	}
	if err = r.Set(bucketName, []byte("key1"), []byte("value1")); err != nil {
		t.Logf("expected error to be nil but was %v", err)
		t.FailNow()
	}
	val, err := r.Get(bucketName, []byte("key1"))
	if err != nil {
		t.Logf("expected error to be nil but was %v", err)
		t.FailNow()
	}
	if !bytes.Equal(val, []byte("value1")) {
		t.Logf("expected value to be %s but was %s", []byte("value1"), val)
		t.FailNow()
	}
}
