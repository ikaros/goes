package goes

import (
	"time"
)

type NilAggregate struct{}

func (NilAggregate) Apply(event Event) Aggregate {
	panic("trying ti use a NilAggregate")
}

func (NilAggregate) GetID() string {
	panic("trying ti use a NilAggregate")
}

type Aggregate interface {
	Apply(Event) Aggregate
	GetID() string
	//	UpdateVersion() Aggregate
}

type BaseAggregate struct {
	ID        string     `json:"id"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at"`
	Version   uint64     `json:"version"`
}

func (a BaseAggregate) GetID() string {
	return a.ID
}

func (a BaseAggregate) Events() ([]Event, error) {
	events := []EventDB{}
	ret := []Event{}

	DB.Where("aggregate_id = ?", a.ID).Order("timestamp").Find(&events)
	for _, event := range events {
		ev, err := event.Decode()
		if err != nil {
			return []Event{}, err
		}
		ret = append(ret, ev)
	}
	return ret, nil
}
