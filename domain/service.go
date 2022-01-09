package domain

import (
	"time"
)

type LiveService interface {
	GetByDate(date *time.Time) (*LiveModel, error)
	Register(live *Live) error
	Update(live *Live) error
	Delete(live *Live) error
}

type LiveServiceImpl struct {
	liveRepository       LiveRepository
	bandRepository       BandRepository
	bandMemberRepository BandMemberRepository
}

func NewLiveServiceImpl(liveRepository LiveRepository, bandRepository BandRepository, bandMemberRepository BandMemberRepository) *LiveServiceImpl {
	return &LiveServiceImpl{liveRepository: liveRepository, bandRepository: bandRepository, bandMemberRepository: bandMemberRepository}
}

func (s *LiveServiceImpl) GetByDate(date *time.Time) (*LiveModel, error) {
	live, err := s.liveRepository.FindByDate(date)
	if err != nil {
		return nil, err
	}

	bands, err := s.bandRepository.FindByLiveId(live.Id)
	if err != nil {
		return nil, err
	}

	var bandModels []*BandModel
	for _, band := range bands {
		players, err := s.bandMemberRepository.FindByLiveIdAndTurn(band.LiveId, band.Turn)
		if err != nil {
			return nil, err
		}
		bandModels = append(bandModels, &BandModel{Name: band.Name, LiveId: band.LiveId, Turn: band.Turn, Player: players})
	}

	return &LiveModel{
		Id:             live.Id,
		Name:           live.Name,
		Location:       live.Location,
		Date:           live.Date,
		PerformanceFee: live.PerformanceFee,
		EquipmentCost:  live.EquipmentCost,
		Band:           bandModels,
	}, nil
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
