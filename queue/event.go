package queue

import (
	"github.com/twinj/uuid"
	"time"
)

type Event struct {
	UID       string    `db:"uid"`
	Namespace string    `db:"namespace"`
	Payload   string    `db:"payload"`
	Retries   int       `db:"retries"`
	State     string    `db:"state"`
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
}

func (e *Event) BeforeCreate() error {
	e.UID = uuid.NewV4().String()
	return nil
}

func (e *Event) GetPayload() (Payload, error) {
	p := Payload{}
	returnPayload, err := p.UnMarshal([]byte(e.Payload))

	if err != nil {
		return Payload{}, err
	}

	return *returnPayload, nil
}
