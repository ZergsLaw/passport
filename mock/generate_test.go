package mock

//go:generate mockgen -source=../core.go -destination=mock.go -package mock -mock_names OauthClient=OauthClient
