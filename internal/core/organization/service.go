package organization

type Service struct {
	org       orgRepo
	member    memberRepo
	tombstone tombstoneRepo
	tx        transaction
}

type ServiceDeps struct {
	Org       orgRepo
	Member    memberRepo
	Tombstone tombstoneRepo
	Tx        transaction
}

func NewService(deps ServiceDeps) *Service {
	return &Service{
		org:       deps.Org,
		member:    deps.Member,
		tombstone: deps.Tombstone,
		tx:        deps.Tx,
	}
}
