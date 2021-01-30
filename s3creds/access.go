package s3creds

type Access struct {
	Allowed struct {
		Resources []string `json:"resources"`
		Actions   []string `json:"actions"`
	} `json:"allowed"`
	Credentials struct {
		AccessKey    string `json:"accessKey"`
		SecretKey    string `json:"secretKey"`
		SessionToken string `json:"sessionToken"`
		Expires      string `json:"expires"`
		Bucket       string `json:"bucket"`
	} `json:"credentials"`
}
