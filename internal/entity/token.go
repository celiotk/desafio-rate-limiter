package entity

type tokenEntity struct {
	limit    int
	reqCount int
}

func NewTokenEntity(limit int, reqCount int) *tokenEntity {
	return &tokenEntity{
		limit:    limit,
		reqCount: reqCount,
	}
}

func (t *tokenEntity) LimitExceeded() bool {
	return t.reqCount > t.limit
}

type TokenSettingsParam struct {
	Token    string
	Limit    int
	Interval int
}

type TokenSettings struct {
	Limit    int
	Interval int
}
