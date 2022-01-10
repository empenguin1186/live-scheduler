package domain

import (
	"time"
)

type LiveService interface {
	GetByPeriod(start *time.Time, end *time.Time) ([]*Live, error)
	Register(live *Live) error
	Update(live *Live) error
	Delete(live *Live) error
}

type LiveServiceImpl struct {
	liveRepository LiveRepository
}

func NewLiveServiceImpl(liveRepository LiveRepository) *LiveServiceImpl {
	return &LiveServiceImpl{liveRepository: liveRepository}
}

func (s *LiveServiceImpl) GetByPeriod(start *time.Time, end *time.Time) ([]*Live, error) {
	live, err := s.liveRepository.FindByPeriod(start, end)
	if err != nil {
		return nil, err
	}
	return live, nil
}

func (s *LiveServiceImpl) Register(live *Live) error {
	err := s.liveRepository.Create(live)
	return verifyAndGetError(err)
}

func (s *LiveServiceImpl) Update(live *Live) error {
	err := s.liveRepository.Update(live)
	return verifyAndGetError(err)
}

func (s *LiveServiceImpl) Delete(live *Live) error {
	err := s.liveRepository.Delete(live)
	return verifyAndGetError(err)
}

func verifyAndGetError(err error) error {
	if err != nil {
		return err
	}
	return nil
}
