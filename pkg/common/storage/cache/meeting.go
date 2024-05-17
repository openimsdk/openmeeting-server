package cache

type Meeting interface {
	Meta
	NewCache() Meeting
}
