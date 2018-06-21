package goes

func Call(command Command, aggregate Aggregate, metadata Metadata) (Aggregate, Event, error) {
	var err error

	tx := DB.Begin()

	// if aggregate instance exists, ensure to lock the row before processing the command
	if aggregate.GetID() != "" {
		tx.Set("gorm:query_option", "FOR UPDATE").First(aggregate)
	}

	err = command.Validate(aggregate)
	if err != nil {
		tx.Rollback()
		return NilAggregate{}, Event{}, err
	}

	data, err := command.BuildEvent()
	if err != nil {
		tx.Rollback()
		return NilAggregate{}, Event{}, err
	}

	event := buildBaseEvent(data.(EventInterface), metadata, aggregate.GetID())
	event.Data = data
	aggregate = aggregate.Apply(event)

	// in Case of Create event
	event.AggregateID = aggregate.GetID()

	err = tx.Save(aggregate).Error
	if err != nil {
		tx.Rollback()
		return NilAggregate{}, Event{}, err
	}

	eventDBToSave, err := event.Encode()
	if err != nil {
		tx.Rollback()
		return NilAggregate{}, Event{}, err
	}

	err = tx.Create(&eventDBToSave).Error
	if err != nil {
		tx.Rollback()
		return NilAggregate{}, Event{}, err
	}

	Dispatch(event)
	//p dispatch

	err = tx.Commit().Error
	if err != nil {
		tx.Rollback()
		return NilAggregate{}, Event{}, err
	}

	return aggregate, event, nil
}

type Command interface {
	BuildEvent() (interface{}, error)
	Validate(interface{}) error
}
