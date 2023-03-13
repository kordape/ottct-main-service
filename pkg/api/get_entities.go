package api

type GetEntitiesResponse []Entity

type Entity struct {
	Id               string `json:"id"`
	TwitterAccountId string `json:"twitterAccountId"`
	Name             string `json:"name"`
}
