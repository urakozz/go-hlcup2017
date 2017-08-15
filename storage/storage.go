package storage

import (
	"errors"
	"sync"

	"sort"

	"time"

	"../entities"
)

type Container struct {
	sync.RWMutex
	Opts

	userStorage []*entities.User
	//userMaxId   int64

	locationStorage []*entities.Location
	//locationNextId  int64

	visitStorage []*entities.Visit
	//visitNextId  int64

	userToVisits     [][]int64
	locationToVisits [][]int64
}

var ErrNotFound = errors.New("no such entity")
var ErrBadRequest = errors.New("bad request")

type Opts struct {
}

func NewStorage(o Opts) *Container {
	return &Container{
		Opts:             o,
		userStorage:      make([]*entities.User, 1000000),
		locationStorage:  make([]*entities.Location, 1000000),
		visitStorage:     make([]*entities.Visit, 1000000),
		userToVisits:     make([][]int64, 1000000),
		locationToVisits: make([][]int64, 1000000),
	}
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
		if l := c.userToVisits[*visit.UserID]; l == nil {
			c.userToVisits[*visit.UserID] = []int64{}
		}
		c.userToVisits[*visit.UserID] = append(c.userToVisits[*visit.UserID], visit.ID)

		// to Location
		if l := c.locationToVisits[*visit.LocationID]; l == nil {
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
	c.growUser(u.ID)
	c.growUserToVisits(u.ID)
	c.growLocationToVisits(u.ID)
	c.userStorage[u.ID] = u
	c.Unlock()
}
func (c *Container) growUserToVisits(n int64) {
	if n >= int64(len(c.userToVisits)) {
		tmp := make([][]int64, n*3/2)
		copy(tmp, c.userToVisits)
		//c.userToVisits = tmp
	}
}
func (c *Container) growLocationToVisits(n int64) {
	if n >= int64(len(c.locationToVisits)) {
		tmp := make([][]int64, n*3/2)
		copy(tmp, c.locationToVisits)
		//c.locationToVisits = tmp
	}
}
func (c *Container) growUser(n int64) {
	if n >= int64(len(c.userStorage)) {
		tmp := make([]*entities.User, n*3/2)
		copy(tmp, c.userStorage)
		//c.userStorage = tmp
	}
}

func (c *Container) LoadUsers(vs []*entities.User) {
	c.Lock()
	for _, v := range vs {
		c.growUser(int64(v.ID))
		c.userStorage[v.ID] = v
	}
	c.Unlock()
}

func (c *Container) UpdateUser(u *entities.User) error {
	c.Lock()
	defer c.Unlock()
	if u.ID >= int64(len(c.userStorage)) {
		return ErrNotFound
	}
	if v := c.userStorage[u.ID]; v != nil {
		v.Update(u)
		return nil
	}
	return ErrNotFound
}
func (c *Container) GetUser(ID int64) (*entities.User, error) {
	c.RLock()
	defer c.RUnlock()
	if ID >= int64(len(c.userStorage)) {
		return nil, ErrNotFound
	}
	if u := c.userStorage[ID]; u != nil {
		return u, nil
	}
	return nil, ErrNotFound
}

// ----

func (c *Container) NewLocation(v *entities.Location) {
	//id := atomic.AddInt32(&c.userMaxId, 1)
	//u.ID = id
	c.Lock()
	c.growLocation(int64(v.ID))
	c.locationStorage[v.ID] = v
	c.Unlock()
}
func (c *Container) growLocation(n int64) {
	if n >= int64(len(c.locationStorage)) {
		tmp := make([]*entities.Location, n*3/2)
		copy(tmp, c.locationStorage)
		//c.locationStorage = tmp
	}
}

func (c *Container) LoadLocations(vs []*entities.Location) {
	c.Lock()
	for _, v := range vs {
		c.growLocation(int64(v.ID))
		c.locationStorage[v.ID] = v
	}
	c.Unlock()
}

func (c *Container) UpdateLocation(u *entities.Location) error {
	c.Lock()
	defer c.Unlock()
	if u.ID >= int64(len(c.locationStorage)) {
		return ErrNotFound
	}
	if v := c.locationStorage[u.ID]; v != nil {
		v.Update(u)
		return nil
	}
	return ErrNotFound
}
func (c *Container) GetLocation(ID int64) (*entities.Location, error) {
	c.RLock()
	defer c.RUnlock()
	if ID >= int64(len(c.locationStorage)) {
		return nil, ErrNotFound
	}

	if u := c.locationStorage[ID]; u != nil {
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

	if *v.UserID >= int64(len(c.userStorage)) || c.userStorage[*v.UserID] == nil {
		return ErrBadRequest
	}

	if *v.LocationID >= int64(len(c.locationStorage)) || c.locationStorage[*v.LocationID] == nil {
		return ErrBadRequest
	}

	c.growVisit(int64(v.ID))
	v.User = c.userStorage[*v.UserID]
	v.Location = c.locationStorage[*v.LocationID]
	c.visitStorage[v.ID] = v
	c.userToVisits[*v.UserID] = append(c.userToVisits[*v.UserID], v.ID)
	c.locationToVisits[*v.LocationID] = append(c.locationToVisits[*v.LocationID], v.ID)
	return nil
}
func (c *Container) growVisit(n int64) {
	if n >= int64(len(c.visitStorage)) {
		tmp := make([]*entities.Visit, n*3/2)
		copy(tmp, c.visitStorage)
		//c.visitStorage = tmp
	}
}

func (c *Container) LoadVisits(vs []*entities.Visit) {
	c.Lock()
	for _, v := range vs {
		c.growVisit(int64(v.ID))
		c.visitStorage[v.ID] = v
	}
	c.Unlock()
}

func (c *Container) UpdateVisit(u *entities.Visit) error {
	c.Lock()
	defer c.Unlock()
	if int64(len(c.visitStorage))-1 < u.ID {
		return ErrNotFound
	}
	visit := c.visitStorage[u.ID]
	if visit == nil {
		return ErrNotFound
	}
	if u.UserID != nil {
		if *u.UserID >= int64(len(c.userStorage)) || c.userStorage[*u.UserID] == nil {
			return ErrBadRequest
		}
	}
	if u.LocationID != nil {
		if *u.LocationID >= int64(len(c.locationStorage)) || c.locationStorage[*u.LocationID] == nil {
			return ErrBadRequest
		}
	}

	diff := visit.Update(u)
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
	if ID >= int64(len(c.visitStorage)) {
		return nil, ErrNotFound
	}
	if u := c.visitStorage[ID]; u != nil {
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

func (c *Container) GetUserVisitsFiltered(ID int64, opts GetUserVisitsOpts) (res []*entities.ShortVisit) {
	var list shortVisitList
	for _, v := range c.GetUserVisits(ID) {
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
		sv := &entities.ShortVisit{
			Mark:      *v.Mark,
			Place:     *v.Location.Place,
			VisitedAt: *v.VisitedAt,
		}
		list = append(list, sv)
	}
	sort.Sort(list)
	return []*entities.ShortVisit(list)
}

func (c *Container) GetUserVisits(ID int64) (res []*entities.Visit) {
	c.RLock()
	defer c.RUnlock()
	if ID >= int64(len(c.userToVisits)) {
		return nil
	}
	if visits := c.userToVisits[ID]; visits != nil {
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

func (c *Container) GetLocationVisitsFiltered(ID int64, opts GetLocationVisitsOpts) (res []*entities.Visit) {
	for _, v := range c.GetLocationVisits(ID) {
		//j, _ := json.Marshal(v)
		//log.Println(string(j))
		if opts.FromDate != nil && *v.VisitedAt <= *opts.FromDate {
			continue
		}
		if opts.ToDate != nil && *v.VisitedAt >= *opts.ToDate {
			continue
		}

		age := int64(time.Since(time.Unix(int64(*v.User.Birthdate), 0)).Hours() / 24 / 365)
		if opts.FromAge != nil && age <= *opts.FromAge {
			continue
		}
		if opts.ToAge != nil && age >= *opts.ToAge {
			continue
		}
		if opts.Gender != nil && *v.User.Gender != *opts.Gender {
			continue
		}

		res = append(res, v)
	}
	return
}

func (c *Container) GetLocationVisits(ID int64) (res []*entities.Visit) {
	c.RLock()
	defer c.RUnlock()
	if ID >= int64(len(c.locationToVisits)) {
		return nil
	}
	if visits := c.locationToVisits[ID]; visits != nil {
		for _, vid := range visits {
			res = append(res, c.visitStorage[vid])
		}
		return
	}
	return nil
}