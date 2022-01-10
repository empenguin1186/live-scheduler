package domain

import (
	"fmt"
	"time"
)

type LiveDescService interface {
	GetByDate(date *time.Time) (*LiveModel, error)
}

type LiveDescServiceImpl struct {
	liveRepository       LiveRepository
	bandRepository       BandRepository
	bandMemberRepository BandMemberRepository
}

func NewLiveDescServiceImpl(liveRepository LiveRepository, bandRepository BandRepository, bandMemberRepository BandMemberRepository) *LiveDescServiceImpl {
	return &LiveDescServiceImpl{liveRepository: liveRepository, bandRepository: bandRepository, bandMemberRepository: bandMemberRepository}
}

func (i *LiveDescServiceImpl) GetByDate(date *time.Time) (*LiveModel, error) {
	lives, err := i.liveRepository.FindByPeriod(date, date)
	if err != nil {
		return nil, err
	}
	if len(lives) == 0 {
		return nil, fmt.Errorf("live not found")
	}
	live := lives[0]

	bands, err := i.bandRepository.FindByLiveId(live.Id)
	if err != nil {
		return nil, err
	}

	var bandModels []*BandModel
	for _, band := range bands {
		players, err := i.bandMemberRepository.FindByLiveIdAndTurn(band.LiveId, band.Turn)
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
