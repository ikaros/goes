package goes

import (
	"encoding/json"
	"reflect"
	"strconv"
	"time"

	"github.com/jinzhu/gorm/dialects/postgres"
	uuid "github.com/satori/go.uuid"
)

var eventRegistry = map[string]reflect.Type{}

type EventInterface interface {
	AggregateType() string
	Action() string
	Version() uint64
}

type Event struct {
	ID            string      `json:"id" gorm:"type:uuid;primary_key"`
	Timestamp     time.Time   `json:"timestamp"`
	AggregateID   string      `json:"aggregate_id"`
	AggregateType string      `json:"aggregate_type"`
	Action        string      `json:"action"`
	Version       uint64      `json:"version"`
	Type          string      `json:"type"`
	Data          interface{} `json:"data"`
	Metadata      Metadata    `json:"metadata"`
}

type Metadata = map[string]interface{}

type EventDB struct {
	ID            string    `json:"id"`
	Timestamp     time.Time `json:"timestamp"`
	AggregateID   string    `json:"aggregate_id"`
	AggregateType string    `json:"aggregate_type"`
	Action        string    `json:"action"`
	Version       uint64    `json:"version"`
	Type          string    `json:"type"`

	RawData     postgres.Jsonb `json:"-" gorm:"type:jsonb;column:data"`
	RawMetadata postgres.Jsonb `json:"-" gorm:"type:jsonb;column:metadata"`
}

func buildBaseEvent(evi EventInterface, metadata Metadata, aggregateID string) Event {
	event := Event{}
	uuidV4, _ := uuid.NewV4()

	if metadata == nil {
		metadata = Metadata{}
	}

	event.ID = uuidV4.String()
	event.Timestamp = time.Now()
	event.AggregateID = aggregateID
	event.AggregateType = evi.AggregateType()
	event.Action = evi.Action()
	event.Type = evi.AggregateType() + "." + evi.Action()
	event.Metadata = metadata
	event.Version = evi.Version()

	return event
}

func RegisterEvents(events ...EventInterface) {

	for _, event := range events {
		eventType := event.AggregateType() +
			"." + event.Action() +
			"." + strconv.FormatUint(event.Version(), 10)

		eventRegistry[eventType] = reflect.TypeOf(event)
	}
}

func (event Event) Encode() (EventDB, error) {
	ret := EventDB{}
	var err error

	ret.ID = event.ID
	ret.Timestamp = event.Timestamp
	ret.AggregateID = event.AggregateID
	ret.AggregateType = event.AggregateType
	ret.Action = event.Action
	ret.Type = event.Type
	ret.Version = event.Version

	ret.RawMetadata.RawMessage, err = json.Marshal(event.Metadata)
	if err != nil {
		return EventDB{}, err
	}

	ret.RawData.RawMessage, err = json.Marshal(event.Data)
	if err != nil {
		return EventDB{}, err
	}

	return ret, nil
}

func (event EventDB) Decode() (Event, error) {
	// deserialize json
	var err error
	ret := Event{}

	// reflexion magic
	eventType := event.AggregateType +
		"." + event.Action +
		"." + strconv.FormatUint(event.Version, 10)
	dataPointer := reflect.New(eventRegistry[eventType])
	dataValue := dataPointer.Elem()
	var data map[string]interface{}

	err = json.Unmarshal(event.RawData.RawMessage, &data)
	if err != nil {
		return Event{}, err
	}

	n := dataValue.NumField()
	for i := 0; i < n; i += 1 {
		field := dataValue.Type().Field(i)
		jsonName := field.Tag.Get("json")
		//fmt.Println(jsonName)
		if jsonName == "" {
			jsonName = field.Name
		}
		val := dataValue.FieldByName(field.Name)
		val.Set(reflect.ValueOf(data[jsonName]))
	}

	ret.ID = event.ID
	ret.Timestamp = event.Timestamp
	ret.AggregateID = event.AggregateID
	ret.AggregateType = event.AggregateType
	ret.Action = event.Action
	ret.Type = event.Type
	ret.Version = event.Version

	dataInterface := dataValue.Interface()
	ret.Data = dataInterface

	err = json.Unmarshal(event.RawMetadata.RawMessage, &ret.Metadata)
	if err != nil {
		return Event{}, err
	}

	return ret, nil
}

func (EventDB) TableName() string {
	return "events"
}
