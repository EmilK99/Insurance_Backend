package store

import (
	"context"
	"github.com/jackc/pgx/v4/pgxpool"
	log "github.com/sirupsen/logrus"
	"time"
)

type Event struct {
	ID       uint
	Name     string
	FlightId string
}

type Listeners map[string]ListenFunc

func (s Scheduler) AddListener(event string, listenFunc ListenFunc) {
	s.listeners[event] = listenFunc
}

type ListenFunc func(string)

type Scheduler struct {
	pool      *pgxpool.Pool
	listeners Listeners
}

func NewScheduler(pool *pgxpool.Pool, listeners Listeners) Scheduler {
	return Scheduler{
		pool:      pool,
		listeners: listeners,
	}
}

func (s Scheduler) Schedule(event, flightId string, runAt time.Time) {
	log.Print("🚀 Scheduling event ", event, " to run at ", runAt)
	_, err := s.pool.Exec(context.Background(), `INSERT INTO "public"."flights" ("name", "flight_id", "runAt") VALUES ($1, $2, $3)`, event, flightId, runAt)
	if err != nil {
		log.Print("schedule insert error: ", err)
	}
}

func (s Scheduler) checkDueEvents() []Event {
	events := []Event{}
	rows, err := s.pool.Query(context.Background(), `SELECT "id", "name", "flight_id" FROM "public"."flights" WHERE "runAt" < $1`, time.Now())
	if err != nil {
		log.Print("💀 error: ", err)
		return nil
	}
	for rows.Next() {
		evt := Event{}
		rows.Scan(&evt.ID, &evt.Name, &evt.FlightId)
		events = append(events, evt)
	}
	return events
}

func (s Scheduler) callListeners(event Event) {
	eventFn, ok := s.listeners[event.Name]
	if ok {
		go eventFn(event.FlightId)
		_, err := s.pool.Exec(context.Background(), `DELETE FROM "public"."flights" WHERE "id" = $1`, event.ID)
		if err != nil {
			log.Print("💀 error: ", err)
		}
	} else {
		log.Print("💀 error: couldn't find event listeners attached to ", event.Name)
	}

}

func (s Scheduler) CheckEventsInInterval(ctx context.Context, duration time.Duration) {
	ticker := time.NewTicker(duration)
	go func() {
		for {
			select {
			case <-ctx.Done():
				ticker.Stop()
				return
			case <-ticker.C:
				log.Println("⏰ Ticks Received...")
				events := s.checkDueEvents()
				for _, e := range events {
					s.callListeners(e)
				}
			}

		}
	}()
}

func (s *Scheduler) StartScheduler(ctx context.Context, interval time.Duration) (chan<- struct{}, <-chan struct{}) {
	if interval <= 0 {
		interval = 60
	}
	quit, done := make(chan struct{}), make(chan struct{})
	go s.CheckEventsInInterval(ctx, interval)
	return quit, done
}
