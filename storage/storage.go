package storage

import (
	"errors"
	"sort"
	"sync"
	"time"

	"github.com/urakozz/highloadcamp/entities"
)

type Container struct {
	sync.RWMutex
	Opts

	now time.Time

	userStorage     map[int64]*entities.User
	locationStorage map[int64]*entities.Location
	visitStorage    map[int64]*entities.Visit

	userList     []*entities.User
	locationList []*entities.Location
	visitList    []*entities.Visit

	userToVisits     map[int64][]int64
	locationToVisits map[int64][]int64
	//sortedL2V        int64
}

var ErrNotFound = errors.New("no such entity")
var ErrBadRequest = errors.New("bad request")

type Opts struct {
}

func NewStorage(o Opts) *Container {
	return &Container{
		Opts:             o,
		now:              time.Now(),
		userStorage:      make(map[int64]*entities.User, 1000000),
		locationStorage:  make(map[int64]*entities.Location, 1000000),
		visitStorage:     make(map[int64]*entities.Visit, 1000000),
		userToVisits:     make(map[int64][]int64, 1000000),
		locationToVisits: make(map[int64][]int64, 1000000),
	}
}
func (c *Container) SetNow(t time.Time) {
	c.now = t
}
func (c *Container) ProcessLoad() {
	c.Lock()
	for _, visit := range c.visitStorage {
		if visit == nil {
			continue
		}
		visit.User = c.userStorage[*visit.UserID]
		visit.Location = c.locationStorage[*visit.LocationID]

		// to User
		if l, ok := c.userToVisits[*visit.UserID]; ok && l == nil {
			c.userToVisits[*visit.UserID] = []int64{}
		}
		c.userToVisits[*visit.UserID] = append(c.userToVisits[*visit.UserID], visit.ID)

		// to Location
		if l, ok := c.locationToVisits[*visit.LocationID]; ok && l == nil {
			c.locationToVisits[*visit.LocationID] = []int64{}
		}
		c.locationToVisits[*visit.LocationID] = append(c.locationToVisits[*visit.LocationID], visit.ID)
	}
	c.Unlock()
}

func (c *Container) NewUser(u *entities.User) {
	//id := atomic.AddInt32(&c.userMaxId, 1)
	//u.ID = id
	c.Lock()
	c.userStorage[u.ID] = u
	c.userToVisits[u.ID] = []int64{}
	c.Unlock()
}

func (c *Container) LoadUsers(vs []*entities.User) {
	c.Lock()
	for _, v := range vs {
		c.userStorage[v.ID] = v
		v.SaveJSON()
	}
	c.Unlock()
}

func (c *Container) UpdateUser(u *entities.User) error {
	c.Lock()
	defer c.Unlock()
	if v, ok := c.userStorage[u.ID]; ok && v != nil {
		v.Update(u)
		return nil
	}
	return ErrNotFound
}
func (c *Container) GetUser(ID int64) (*entities.User, error) {
	c.RLock()
	defer c.RUnlock()
	if u, ok := c.userStorage[ID]; ok && u != nil {
		return u, nil
	}
	return nil, ErrNotFound
}

// ----

func (c *Container) NewLocation(v *entities.Location) {
	c.Lock()
	c.locationStorage[v.ID] = v
	c.locationToVisits[v.ID] = []int64{}
	c.Unlock()
}

func (c *Container) LoadLocations(vs []*entities.Location) {
	c.Lock()
	for _, v := range vs {
		c.locationStorage[v.ID] = v
		v.SaveJSON()
	}
	c.Unlock()
}

func (c *Container) UpdateLocation(u *entities.Location) error {
	c.Lock()
	defer c.Unlock()
	if v, ok := c.locationStorage[u.ID]; ok && v != nil {
		v.Update(u)
		return nil
	}
	return ErrNotFound
}
func (c *Container) GetLocation(ID int64) (*entities.Location, error) {
	c.RLock()
	defer c.RUnlock()

	if u, ok := c.locationStorage[ID]; ok && u != nil {
		return u, nil
	}
	return nil, ErrNotFound
}

// ----

func (c *Container) NewVisit(v *entities.Visit) error {
	//id := atomic.AddInt32(&c.userMaxId, 1)
	//u.ID = id
	c.Lock()
	defer c.Unlock()

	if _, ok := c.userStorage[*v.UserID]; !ok {
		return ErrBadRequest
	}

	if _, ok := c.locationStorage[*v.LocationID]; !ok {
		return ErrBadRequest
	}

	v.User = c.userStorage[*v.UserID]
	v.Location = c.locationStorage[*v.LocationID]
	c.visitStorage[v.ID] = v
	c.userToVisits[*v.UserID] = append(c.userToVisits[*v.UserID], v.ID)
	c.locationToVisits[*v.LocationID] = append(c.locationToVisits[*v.LocationID], v.ID)
	return nil
}

func (c *Container) LoadVisits(vs []*entities.Visit) {
	c.Lock()
	for _, v := range vs {
		c.visitStorage[v.ID] = v
		v.SaveJSON()
	}
	c.Unlock()
}

func (c *Container) UpdateVisit(u *entities.Visit) error {
	c.Lock()
	defer c.Unlock()
	visit, ok := c.visitStorage[u.ID]
	if !ok || visit == nil {
		return ErrNotFound
	}
	if u.UserID != nil {
		if _, ok := c.userStorage[*u.UserID]; !ok {
			return ErrBadRequest
		}
	}
	if u.LocationID != nil {
		if _, ok := c.locationStorage[*u.LocationID]; !ok {
			return ErrBadRequest
		}
	}

	diff := visit.Update(u)
	visit.User = c.userStorage[*visit.UserID]
	visit.Location = c.locationStorage[*visit.LocationID]

	// process UserChange
	if diff.UserID.HasDiff && diff.UserID.Old != diff.UserID.New {
		// remove from old
		if tmp := c.userToVisits[diff.UserID.Old]; len(tmp) != 0 {
			var updatedOld []int64
			for _, v := range tmp {
				if v != u.ID {
					updatedOld = append(updatedOld, v)
				}
			}
			c.userToVisits[diff.UserID.Old] = updatedOld
		}
		// append to new
		c.userToVisits[diff.UserID.New] = append(c.userToVisits[diff.UserID.New], u.ID)
	}
	// process Location Change
	//log.Println("start updating location", u.ID, diff.LocationID.Old, diff.LocationID.New)
	if diff.LocationID.HasDiff && diff.LocationID.Old != diff.LocationID.New {
		if tmp := c.locationToVisits[diff.LocationID.Old]; len(tmp) != 0 {
			var updatedOld []int64
			for _, v := range tmp {
				if v != u.ID {
					updatedOld = append(updatedOld, v)
				}
			}
			c.locationToVisits[diff.LocationID.Old] = updatedOld
		}
		c.locationToVisits[diff.LocationID.New] = append(c.locationToVisits[diff.LocationID.New], u.ID)
	}
	return nil
}
func (c *Container) GetVisit(ID int64) (*entities.Visit, error) {
	c.RLock()
	defer c.RUnlock()
	if u, ok := c.visitStorage[ID]; ok && u != nil {
		return u, nil
	}
	return nil, ErrNotFound
}
func (c *Container) ListVisits() []*entities.Visit {
	c.RLock()
	defer c.RUnlock()
	tmp := make([]*entities.Visit, 0, len(c.visitStorage)/4)
	for _, v := range c.visitStorage {
		if v != nil {
			tmp = append(tmp, v)
		}
	}
	return tmp
}

type GetUserVisitsOpts struct {
	FromDate   *int64
	ToDate     *int64
	Country    *string
	ToDistance *int64
}

type shortVisitList []*entities.ShortVisit

func (l shortVisitList) Len() int {
	return len(l)
}
func (l shortVisitList) Less(i, j int) bool {
	return l[i].VisitedAt < l[j].VisitedAt
}
func (l shortVisitList) Swap(i, j int) {
	l[i], l[j] = l[j], l[i]
}

type visitList []*entities.Visit

func (l visitList) Len() int {
	return len(l)
}
func (l visitList) Less(i, j int) bool {
	return *l[i].VisitedAt < *l[j].VisitedAt
}
func (l visitList) Swap(i, j int) {
	l[i], l[j] = l[j], l[i]
}
func (l visitList) Ids() []int64 {
	ids := make([]int64, len(l))
	for i, v := range l {
		ids[i] = v.ID
	}
	return ids
}

func (c *Container) GetUserVisitsFiltered(ID int64, opts GetUserVisitsOpts) *entities.ShortVisitContainer {
	list := shortVisitList{}
	for _, v := range c.getUserVisits(ID) {
		if opts.FromDate != nil && *v.VisitedAt <= *opts.FromDate {
			continue
		}
		if opts.ToDate != nil && *v.VisitedAt >= *opts.ToDate {
			continue
		}
		if opts.Country != nil && *opts.Country != *v.Location.Country {
			continue
		}
		if opts.ToDistance != nil && *v.Location.Distance >= *opts.ToDistance {
			continue
		}
		sv := entities.DefaultShortVisitPool.Get()
		sv.Mark = *v.Mark
		sv.Place = *v.Location.Place
		sv.VisitedAt = *v.VisitedAt

		list = append(list, sv)
	}
	// TODO save presorted list
	sort.Sort(list)
	return &entities.ShortVisitContainer{Visits: []*entities.ShortVisit(list)}
}

func (c *Container) getUserVisits(ID int64) (res []*entities.Visit) {
	c.RLock()
	defer c.RUnlock()

	if visits, ok := c.userToVisits[ID]; ok && visits != nil {
		for _, vid := range visits {
			res = append(res, c.visitStorage[vid])
		}
		return
	}
	return nil
}

type GetLocationVisitsOpts struct {
	FromDate *int64
	ToDate   *int64
	FromAge  *int64
	ToAge    *int64
	Gender   *string
}

func (c *Container) GetLocationVisitsFilteredAvg(ID int64, opts GetLocationVisitsOpts) float64 {
	var sum int64
	var i int64
	fromBdTs := time.Date(1800, c.now.Month(), c.now.Day(), 0, 0, 0, 0, time.UTC).Unix()
	toBdTs := c.now.Unix()
	if opts.FromAge != nil {
		toBdTs = time.Date(c.now.Year()-int(*opts.FromAge), c.now.Month(), c.now.Day(), c.now.Hour(), c.now.Minute(), c.now.Second(), 0, time.UTC).Unix()
	}
	if opts.ToAge != nil {
		fromBdTs = time.Date(c.now.Year()-int(*opts.ToAge), c.now.Month(), c.now.Day(), c.now.Hour(), c.now.Minute(), c.now.Second(), 0, time.UTC).Unix()
	}

	for _, v := range c.getLocationVisits(ID) {
		//j, _ := json.Marshal(v)
		//log.Println(string(j))
		if opts.FromDate != nil && *v.VisitedAt <= *opts.FromDate {
			continue
		}
		if opts.ToDate != nil && *v.VisitedAt >= *opts.ToDate {
			continue
		}
		if opts.Gender != nil && *v.User.Gender != *opts.Gender {
			continue
		}
		if fromBdTs < int64(*v.User.Birthdate) && int64(*v.User.Birthdate) < toBdTs {
			sum += int64(*v.Mark)
			i++
		}
	}
	if i == 0 {
		return 0
	}
	return float64(sum) / float64(i)
}

func (c *Container) getLocationVisits(ID int64) (res []*entities.Visit) {
	c.RLock()
	defer c.RUnlock()
	if visits, ok := c.locationToVisits[ID]; ok && visits != nil {
		for _, vid := range visits {
			res = append(res, c.visitStorage[vid])
		}
		return
	}
	return nil
}
