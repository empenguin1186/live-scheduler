package domain

type LiveDescService interface {
	GetById(id int) (*LiveModel, error)
}

type LiveDescServiceImpl struct {
	liveRepository       LiveRepository
	bandRepository       BandRepository
	bandMemberRepository BandMemberRepository
}

func NewLiveDescServiceImpl(liveRepository LiveRepository, bandRepository BandRepository, bandMemberRepository BandMemberRepository) *LiveDescServiceImpl {
	return &LiveDescServiceImpl{liveRepository: liveRepository, bandRepository: bandRepository, bandMemberRepository: bandMemberRepository}
}

func (i *LiveDescServiceImpl) GetById(id int) (*LiveModel, error) {
	live, err := i.liveRepository.FindById(id)
	if err != nil {
		return nil, err
	}

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
