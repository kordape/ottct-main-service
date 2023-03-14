package api

type EntitiesResponse []Entity

type Entity struct {
	Id          string `json:"id"`
	TwitterId   string `json:"twitterId"`
	DisplayName string `json:"displayName"`
}
