package repository

import (
	"live-scheduler/domain/model"
	"time"
)

type LiveRepository interface {
	FindByDate(time time.Time) model.Live
}
