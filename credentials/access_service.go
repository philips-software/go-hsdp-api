package credentials

type AccessService struct {
	client *Client
}

// GetAccessOptions describes the fields on which you can search for policies
type GetAccessOptions struct {
	ProductKey *string `url:"-"`
}

// GetPolicy searches for polices
func (c *AccessService) GetAccess(opt *GetAccessOptions, options ...OptionFunc) ([]*Access, *Response, error) {

	req, err := c.client.NewRequest("GET", "core/credentials/Access", opt, options)
	if err != nil {
		return nil, nil, err
	}
	if opt.ProductKey == nil {
		return nil, nil, ErrMissingProductKey
	}

	req.Header.Set("Api-Version", CredentialsAPIVersion)
	req.Header.Set("X-Product-Key", *opt.ProductKey)

	var accessGetResponse []*Access

	resp, err := c.client.Do(req, &accessGetResponse)
	if err != nil {
		return nil, resp, err
	}
	return accessGetResponse, resp, err
}
