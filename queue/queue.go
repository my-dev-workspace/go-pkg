package queue

import (
	"github.com/jmoiron/sqlx"
	"time"
)

// queue [namespace,name,tag,key,data,lock] processed
// https://github.com/AnalogRepublic/go-mysql-queue/tree/master/msq
// https://github.com/olivere/jobqueue/blob/517172f3c5dfbcf4f548a2631ce4c277f2f54060/mysql/store.go#L727
// https://github.com/mindreframer/golang-webapp-stuff/tree/master/src/github.com/alouca/MongoQueue

// for rows.Next() {
/*
const (
	mysqlSchema = `CREATE TABLE IF NOT EXISTS jobqueue_jobs (
	id varchar(36) primary key,
	topic varchar(255),
	state varchar(30),
	args text,
	priority bigint,
	retry integer,
	max_retry integer,
	correlation_id varchar(255),
	created bigint,
	started bigint,
	completed bigint,
	last_mod bigint,
	index ix_jobs_topic (topic),
	index ix_jobs_state (state),
	index ix_jobs_priority (priority),
	index ix_jobs_correlation_id (correlation_id),
	index ix_jobs_created (created),
	index ix_jobs_started (started),
	index ix_jobs_completed (completed),
	index ix_jobs_last_mod (last_mod));`
)
*/

type QueueConfig struct {
	DB         *sqlx.DB
	MaxRetries int64
	MessageTTL time.Duration
}

type Queue struct {
	QueueConfig
	//listener  map[string]*Listener
}

//type Listener struct {
//	Namespace string
//	QueueConfig
//}

func NewQueue(config QueueConfig) (*Queue, error) {
	queue := &Queue{
		config,
	}

	return queue, nil
}

func (q *Queue) Add(namespace string, payload Payload) (*Event, error) {
	encodedPayload, err := payload.Marshal()

	if err != nil {
		return &Event{}, err
	}

	event := &Event{
		Namespace: namespace,
		Payload:   string(encodedPayload),
		Retries:   1,
		State:     "waiting",
		CreatedAt: time.Now(),
	}
	event.BeforeCreate()

	//_, err = s.db.NamedExec(`INSERT INTO olx_offers (id,url,title,description,params,user,contact,location,photos,category) VALUES (:id,:url,:title,:description,:params,:user,:contact,:location,:photos,:category)`,
	//	map[string]interface{}{
	//		"id": offer.ID,
	//		"url":  offer.URL,
	//		"title":  offer.Title,
	//		"description":  offer.Description,
	//		"params":  string(Params),
	//		"user":  string(User),
	//		"contact":  string(Contact),
	//		"location":  string(Location),
	//		"photos":  string(Photos),
	//		"category":  string(Category),
	//	})
	//if err != nil {
	//	logger.Error("[workerSpider][olx.store.Add] %s", err.Error())
	//}

	return event, nil
}

func (q *Queue) Next(namespace string) (*Event, error) {
	event := &Event{}

	tx, err := q.DB.Begin()
	row := tx.QueryRow("SELECT uid,namespace,payload,retries,state,created_at,updated_at FROM queue WHERE namespace=? AND state=waiting AND created_at <= ? AND retries <= ? LIMIT 1", namespace, time.Now(), q.MaxRetries)
	row.Scan(event.UID, event.Namespace, event.Payload, event.Retries, event.CreatedAt, event.UpdatedAt)

	_, err = tx.Exec("UPDATE queue SET state=processing, started_at=? WHERE uid=?", time.Now(), event.UID)
	if err != nil {
		return event, err
	}

	err = tx.Commit()
	if err != nil {
		return event, err
	}

	return event, nil
}

func (q *Queue) ReQueue(event *Event) {
	now := time.Now()
	pushback := time.Now().Add(time.Millisecond * (time.Duration(event.Retries) * 100))
	retries := event.Retries + 1

	_, _ = q.DB.NamedExec(`UPDATE queue SET state=waiting, created_at=:created_at, updated_at=:updated_at, retries=:retries WHERE uid=:uid`,
		map[string]interface{}{
			"uid":        event.UID,
			"created_at": pushback,
			"updated_at": now,
			"retries":    retries,
		})
}

func (q *Queue) Failed(event *Event) {
	_, _ = q.DB.NamedExec(`UPDATE queue SET state=failed WHERE uid=:uid`,
		map[string]interface{}{
			"uid": event.UID,
		})
}

func (q *Queue) Done(event *Event) {
	_, _ = q.DB.NamedExec(`DELETE FROM queue WHERE uid=:uid`,
		map[string]interface{}{
			"uid": event.UID,
		})
}
