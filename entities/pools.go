package entities

import (
	"sync"
)

var DefaultUserPool = &UserPool{}
var DefaultLocationPool = &LocationPool{}
var DefaultVisitPool = &VisitPool{}
var DefaultShortVisitPool = &ShortVisitPool{}

type UserPool struct {
	pool sync.Pool
}

func (p *UserPool) Get() *User {
	v := p.pool.Get()
	if v != nil {
		return v.(*User)
	}
	return &User{}
}
func (p *UserPool) Put(u *User) {
	u.Reset()
	p.pool.Put(u)
}

type VisitPool struct {
	pool sync.Pool
}

func (p *VisitPool) Get() *Visit {
	v := p.pool.Get()
	if v != nil {
		return v.(*Visit)
	}
	return &Visit{}
}
func (p *VisitPool) Put(u *Visit) {
	u.Reset()
	p.pool.Put(u)
}

type LocationPool struct {
	pool sync.Pool
}

func (p *LocationPool) Get() *Location {
	v := p.pool.Get()
	if v != nil {
		return v.(*Location)
	}
	return &Location{}
}
func (p *LocationPool) Put(u *Location) {
	u.Reset()
	p.pool.Put(u)
}

type ShortVisitPool struct {
	pool sync.Pool
}

func (p *ShortVisitPool) Get() *ShortVisit {
	v := p.pool.Get()
	if v != nil {
		return v.(*ShortVisit)
	}
	return &ShortVisit{}
}
func (p *ShortVisitPool) Put(u *ShortVisit) {
	u.Reset()
	p.pool.Put(u)
}
