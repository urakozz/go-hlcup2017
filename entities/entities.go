package entities

import (
	"errors"
	"sync/atomic"
	"time"
)

type ShortVisitContainer struct {
	Visits []*ShortVisit `json:"visits"`
}
type VisitContainer struct {
	Visits []*Visit `json:"visits"`
}
type UserContainer struct {
	Users []*User `json:"users"`
}
type LocationContainer struct {
	Locations []*Location `json:"locations"`
}

// Location
//{
//"distance": 9,
//"city": "Новоомск",
//"place": "Ресторан",
//"id": 1,
//"country": "Венесуэлла"
//},
type Location struct {
	ID       int64   `json:"id"`
	Distance *int64  `json:"distance"`
	City     *string `json:"city"`
	Place    *string `json:"place"`
	Country  *string `json:"country"`
	hasJSON  int32   `json:"-"`
	json     []byte  `json:"-"`
}

func (v *Location) SaveJSON() {
	v.json, _ = v.MarshalJSON()
	atomic.StoreInt32(&v.hasJSON, 1)
}

func (l *Location) Reset() {
	l.ID = 0
	l.Distance, l.City = nil, nil
	l.Place, l.Country = nil, nil
	atomic.StoreInt32(&l.hasJSON, 0)
}

//func (u *Location) Clone() *Location {
//	v := &Location{}
//	*v = *u
//	v.ID = u.ID
//
//	*v.Distance = *u.Distance
//	*v.City = *u.City
//	*v.Place = *u.Place
//	*v.Country = *u.Country
//	return v
//}

func (u *Location) Update(new *Location) {
	atomic.StoreInt32(&u.hasJSON, 0)
	if new.Distance != nil {
		u.Distance = new.Distance
	}
	if new.City != nil {
		u.City = new.City
	}
	if new.Place != nil {
		u.Place = new.Place
	}
	if new.Country != nil {
		u.Country = new.Country
	}
}

func (u *Location) Validate() error {
	if u.Distance == nil {
		return errors.New("dis nil")
	}
	if u.City == nil {
		return errors.New("city nil")
	}
	if len(*u.City) > 50 {
		return errors.New("city too long")
	}
	if u.Place == nil {
		return errors.New("place nil")
	}

	if u.Country == nil {
		return errors.New("country nil")
	}
	if len(*u.Country) > 50 {
		return errors.New("country too long")
	}
	return nil
}

// Users
//{
//"first_name": "Злата",
//"last_name": "Кисатович",
//"birth_date": -627350400,
//"gender": "f",
//"id": 1,
//"email": "coorzaty@me.com"
//},

type User struct {
	ID        int64   `json:"id"`
	Email     *string `json:"email"`
	FirstName *string `json:"first_name"`
	LastName  *string `json:"last_name"`
	Gender    *string `json:"gender"`
	Birthdate *int64  `json:"birth_date,omitempty"`
	hasJSON   int32   `json:"-"`
	json      []byte  `json:"-"`
}

func (v *User) SaveJSON() {
	v.json, _ = v.MarshalJSON()
	atomic.StoreInt32(&v.hasJSON, 1)
}

func (u *User) Reset() {
	u.ID = 0
	u.Email, u.FirstName = nil, nil
	u.LastName, u.Gender = nil, nil
	u.Birthdate = nil
	atomic.StoreInt32(&u.hasJSON, 0)
}

func (u *User) Update(new *User) {
	atomic.StoreInt32(&u.hasJSON, 0)
	if new.Email != nil {
		u.Email = new.Email
	}
	if new.FirstName != nil {
		u.FirstName = new.FirstName
	}
	if new.LastName != nil {
		u.LastName = new.LastName
	}
	if new.Gender != nil {
		u.Gender = new.Gender
	}
	if new.Birthdate != nil {
		u.Birthdate = new.Birthdate
	}
}

func (u *User) Validate() error {
	if u.Email == nil {
		return errors.New("email nil")
	}
	if len(*u.Email) > 100 {
		return errors.New("too long")
	}
	if u.FirstName == nil {
		return errors.New("fn nil")
	}
	if len(*u.FirstName) > 50 {
		return errors.New("too long")
	}
	if u.LastName == nil {
		return errors.New("ln nil")
	}
	if len(*u.LastName) > 50 {
		return errors.New("too long")
	}
	if u.Gender == nil {
		return errors.New("gen nil")
	}
	if *u.Gender != "f" && *u.Gender != "m" {
		return errors.New("wrong gen")
	}
	return nil
}

// VISITS
//{
//"user": 42,
//"location": 13,
//"visited_at": 1123175509,
//"id": 1,
//"mark": 4
//},

type Visit struct {
	ID         int64     `json:"id"`
	LocationID *int64    `json:"location"`
	UserID     *int64    `json:"user"`
	VisitedAt  *int64    `json:"visited_at"`
	Mark       *uint8    `json:"mark"`
	User       *User     `json:"-"`
	Location   *Location `json:"-"`
	hasJSON    int32     `json:"-"`
	json       []byte    `json:"-"`
}

func (v *Visit) SaveJSON() {
	v.json, _ = v.MarshalJSON()
	atomic.StoreInt32(&v.hasJSON, 1)
}

func (v *Visit) Reset() {
	v.ID = 0
	v.LocationID, v.UserID = nil, nil
	v.VisitedAt, v.Mark = nil, nil
	v.User, v.Location = nil, nil
	atomic.StoreInt32(&v.hasJSON, 0)
}

//func (u *Visit) Clone() *Visit {
//	newv := &Visit{}
//	*newv = *u
//	*newv.LocationID = *u.LocationID
//	*newv.UserID = *u.UserID
//	*newv.VisitedAt = *u.VisitedAt
//	*newv.Mark = *u.Mark
//	newv.User = u.User.Clone()
//	newv.Location = u.Location.Clone()
//	return newv
//}

func (u *Visit) Validate() error {
	if u.LocationID == nil {
		return errors.New("locid nil")
	}
	if u.UserID == nil {
		return errors.New("userid nil")
	}
	if u.VisitedAt == nil { // TODO add min max
		return errors.New("vis is nil")
	}
	_ = time.January
	//if time.Unix(*u.VisitedAt, 0) < time.Date(2000, time.January, 1, 0, 0, 0, 0, nil)
	if u.Mark == nil {
		return errors.New("mark at nil")
	}
	if *u.Mark > 5 {
		return errors.New("mark too high")
	}
	return nil
}

type ShortVisit struct {
	Mark      uint8  `json:"mark"`
	VisitedAt int64  `json:"visited_at"`
	Place     string `json:"place"`
}

func (sv *ShortVisit) Reset() {
	sv.Mark = 0
	sv.VisitedAt = 0
	sv.Place = ""
}

//func (v *ShortVisit) MarshalJSON() (b []byte, err error) {
//	//"mark": 3,
//	//	"visited_at": 1196539893,
//	//	"place": "Ратуша"
//	buf := bytes.NewBufferString(`{"mark":`)
//	buf.WriteString(strconv.Itoa(int(v.Mark)))
//	buf.WriteString(`,"visited_at":`)
//	buf.WriteString(strconv.Itoa(int(v.VisitedAt)))
//	buf.WriteString(`,"place":"`)
//	buf.WriteString(v.Place)
//	buf.WriteString(`"}`)
//	return buf.Bytes(), nil
//}

type VisitDiff struct {
	LocationID struct {
		HasDiff bool
		Old     int64
		New     int64
	}
	UserID struct {
		HasDiff bool
		Old     int64
		New     int64
	}
	//VisitedAt struct{
	//	Old int64
	//	New int64
	//}
	//Mark struct{
	//	Old uint8
	//	New uint8
	//}
}

func (u *Visit) Update(new *Visit) (diff *VisitDiff) {
	atomic.StoreInt32(&u.hasJSON, 0)
	diff = &VisitDiff{}
	if new.LocationID != nil {
		diff.LocationID.HasDiff = true
		diff.LocationID.Old, diff.LocationID.New = *u.LocationID, *new.LocationID
		*u.LocationID = *new.LocationID
	}
	if new.UserID != nil {
		diff.UserID.HasDiff = true
		diff.UserID.Old, diff.UserID.New = *u.UserID, *new.UserID
		*u.UserID = *new.UserID
	}
	if new.VisitedAt != nil {
		//diff.VisitedAt.Old, diff.VisitedAt.New = *u.VisitedAt, *new.VisitedAt
		*u.VisitedAt = *new.VisitedAt
	}
	if new.Mark != nil {
		//diff.Mark.Old, diff.Mark.New = *u.Mark, *new.Mark
		*u.Mark = *new.Mark
	}
	return diff
}
