package storage

import (
	"testing"

	"time"

	"github.com/urakozz/highloadcamp/entities"
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

func TestContainer_GetLocationVisitsFilteredAvg(t *testing.T) {
	c := NewStorage(Opts{})
	now := time.Date(2000, time.January, 2, 0, 0, 0, 0, time.UTC)
	c.SetNow(now)
	c.LoadLocations([]*entities.Location{{
		ID: 1,
	}})
	users := make([]*entities.User, 0, 100)
	visits := make([]*entities.Visit, 0, 100)
	for j := 0; j < 100; j++ {
		j64 := int64(j)
		t := time.Date(1900+j, time.January, 1, 0, 0, 0, 0, time.UTC).Unix()
		users = append(users, &entities.User{
			ID:        j64,
			Birthdate: &t,
		})
		i1 := int64(1)
		mk := uint8(j / 10)
		visitedAt := now.Unix()
		visits = append(visits, &entities.Visit{
			ID:         j64,
			UserID:     &j64,
			LocationID: &i1,
			Mark:       &mk,
			VisitedAt:  &visitedAt,
		})
	}
	c.LoadUsers(users)
	c.LoadVisits(visits)
	c.ProcessLoad()

	toAge := int64(2)
	v := c.GetLocationVisitsFilteredAvg(1, GetLocationVisitsOpts{ToAge: &toAge})
	if v != 9 {
		t.Error(v, "is not 9")
	}
	toAge = int64(12)
	v = c.GetLocationVisitsFilteredAvg(1, GetLocationVisitsOpts{ToAge: &toAge})
	if v != 8.909090909090908 {
		t.Error(v, "v == 8.909090909090908")
	}
	toAge = int64(12)
	fromAge := int64(9)
	v = c.GetLocationVisitsFilteredAvg(1, GetLocationVisitsOpts{ToAge: &toAge, FromAge: &fromAge})
	if v != 8.5 {
		t.Error(v, "v == 8.909090909090908")
	}
	now = time.Date(2000, time.January, 1, 0, 0, 0, 0, time.UTC)
	c.SetNow(now)
	toAge = int64(2)
	v = c.GetLocationVisitsFilteredAvg(1, GetLocationVisitsOpts{ToAge: &toAge})
	if v != 9 {
		t.Error(v, "is not 9")
	}

}
