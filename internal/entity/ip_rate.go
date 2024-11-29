package entity

const (
	TTL_NOT_SET = -1
)

type ipEntity struct {
	limit    int
	reqCount int
}

func NewIpEntity(limit int, reqCount int) *ipEntity {
	return &ipEntity{
		limit:    limit,
		reqCount: reqCount,
	}
}

func (i *ipEntity) LimitExceeded() bool {
	return i.reqCount > i.limit
}
