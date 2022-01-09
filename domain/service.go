package domain

import (
	"time"
)

type LiveService interface {
	GetByDate(date *time.Time) *LiveModel
}

type LiveServiceImpl struct {
	liveRepository       LiveRepository
	bandRepository       BandRepository
	bandMemberRepository BandMemberRepository
}

func NewLiveServiceImpl(liveRepository LiveRepository, bandRepository BandRepository, bandMemberRepository BandMemberRepository) *LiveServiceImpl {
	return &LiveServiceImpl{liveRepository: liveRepository, bandRepository: bandRepository, bandMemberRepository: bandMemberRepository}
}

func (s *LiveServiceImpl) GetByDate(date *time.Time) *LiveModel {
	live, _ := s.liveRepository.FindByDate(date)
	bands, _ := s.bandRepository.FindByLiveId(live.Id)
	var bandModels []*BandModel
	for _, band := range bands {
		players, _ := s.bandMemberRepository.FindByLiveIdAndTurn(band.LiveId, band.Turn)
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
	}
}
