package domain

import (
	"time"
)

type LiveService struct {
	liveRepository       LiveRepository
	bandRepository       BandRepository
	bandMemberRepository BandMemberRepository
}

func NewLiveService(liveRepository LiveRepository, bandRepository BandRepository, bandMemberRepository BandMemberRepository) LiveService {
	return LiveService{liveRepository: liveRepository, bandRepository: bandRepository, bandMemberRepository: bandMemberRepository}
}

func (s LiveService) GetDate(date time.Time) LiveModel {
	live := s.liveRepository.FindByDate(date)
	bands := s.bandRepository.FindByLiveId(live.Id)
	var bandModels []BandModel
	for _, band := range bands {
		members := s.bandMemberRepository.FindByLiveIdAndTurn(band.LiveId, band.Turn)
		bandModels = append(bandModels, BandModel{Name: band.Name, LiveId: band.LiveId, Turn: band.Turn, Member: members})
	}
	return LiveModel{
		Id:             live.Id,
		Name:           live.Name,
		Location:       live.Location,
		Date:           live.Date,
		PerformanceFee: live.PerformanceFee,
		EquipmentCost:  live.EquipmentCost,
		Band:           bandModels,
	}
}
