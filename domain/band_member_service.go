package domain

type BandMemberService interface {
	Register(bandMember *BandMember) error
	GetByLiveIdAndTurn(id int, turn int) ([]*Player, error)
	Update(bandMember *BandMember, id int, turn int) error
	Delete(bandMember *BandMember) error
}

type BandMemberServiceImpl struct {
	bandMemberRepository BandMemberRepository
}

func NewBandMemberServiceImpl(bandMemberRepository BandMemberRepository) *BandMemberServiceImpl {
	return &BandMemberServiceImpl{
		bandMemberRepository: bandMemberRepository,
	}
}

func (b *BandMemberServiceImpl) Register(bandMember *BandMember) error {
	return b.bandMemberRepository.Create(bandMember)
}

func (b *BandMemberServiceImpl) GetByLiveIdAndTurn(id int, turn int) ([]*Player, error) {
	return b.bandMemberRepository.FindByLiveIdAndTurn(id, turn)
}

func (b *BandMemberServiceImpl) Update(bandMember *BandMember, id int, turn int) error {
	return b.bandMemberRepository.Update(bandMember, id, turn)
}

func (b *BandMemberServiceImpl) Delete(bandMember *BandMember) error {
	return b.bandMemberRepository.Delete(bandMember)
}
