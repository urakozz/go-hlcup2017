package storage

import (
	"testing"

	"../entities"
)

func TestContainer_LoadUsers(t *testing.T) {
	s := NewStorage(Opts{})
	bulksize := int64(10000)
	bulkwrites := int64(100)
	for i := int64(0); i < bulkwrites; i++ {
		list := make([]*entities.User, 0, bulksize)
		for j := bulksize * i; j < bulksize*(i+1); j++ {
			list = append(list, &entities.User{ID: int64(j)})
		}
		s.LoadUsers(list)
	}

	for i := int64(0); i < bulksize*bulkwrites; i++ {
		u, err := s.GetUser(i)
		if err != nil {
			t.Error("unable to get", i, err)
		}
		if u == nil {
			t.Error("nil response", i)
		}
		if i != u.ID {
			t.Error("collision? saved", i, "got", u.ID)
		}
	}
}
