package domain

type BandService interface {
	GetByLiveId(id int) ([]*Band, error)
	Register(band *Band) error
	Update(id int, turn int, band *Band) error
	Delete(id int, turn int) error
}

type BandServiceImpl struct {
	bandRepository BandRepository
}

func NewBandServiceImpl(bandRepository BandRepository) *BandServiceImpl {
	return &BandServiceImpl{bandRepository: bandRepository}
}

func (b *BandServiceImpl) GetByLiveId(id int) ([]*Band, error) {
	return b.bandRepository.FindByLiveId(id)
}

func (b *BandServiceImpl) Register(band *Band) error {
	return b.bandRepository.Create(band)
}

func (b *BandServiceImpl) Update(id int, turn int, band *Band) error {
	return b.bandRepository.Update(id, turn, band)
}

func (b *BandServiceImpl) Delete(id int, turn int) error {
	return b.bandRepository.Delete(id, turn)
}
