package internal

type Credentials struct {
	BHUrl      string `config:"bloodhound.url"`
	BHTokenID  string `config:"bloodhound.tokenID"`
	BHTokenKey string `config:"bloodhound.tokenKey"`
}
