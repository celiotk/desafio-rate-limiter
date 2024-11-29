package database

type IpRateStorage struct {
	storage RateLimiterStorage
}

func NewIpRateStorage(storage RateLimiterStorage) *IpRateStorage {
	return &IpRateStorage{storage: storage}
}

func (i *IpRateStorage) Increment(key string) (int, error) {
	return i.storage.Increment("ip:count:" + key)
}

func (i *IpRateStorage) SetExpiration(key string, ttl int) error {
	return i.storage.SetExpiration("ip:count:"+key, ttl)
}

func (i *IpRateStorage) GetTTL(key string) (int, error) {
	return i.storage.GetTTL("ip:count:" + key)
}

func (i *IpRateStorage) IsBlocked(key string) (bool, error) {
	return i.storage.Exists("ip:block:" + key)
}

func (i *IpRateStorage) Block(key string, duration int) error {
	return i.storage.Block("ip:block:"+key, duration)
}
