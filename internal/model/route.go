package model

type Route struct {
	Subdomain string `bson:"subdomain"`
	TargetURL string `bson:"target_url"`
	IsActive  bool   `bson:"is_active"`
}
