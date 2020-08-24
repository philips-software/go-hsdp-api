package has

type ImagesService struct {
	orgID  string
	client *Client
}

type Image struct {
	ID          string   `json:"id"`
	Name        string   `json:"name"`
	Description string   `json:"description"`
	Regions     []string `json:"regions"`
}

type imageResult struct {
	Images []Image `json:"images"`
}

// GetImages retrieves images in HAS
func (c *ImagesService) GetImages(options ...OptionFunc) (*[]Image, *Response, error) {
	req, err := c.client.NewHASRequest("GET", "has/image", nil, options)
	if err != nil {
		return nil, nil, err
	}
	req.Header.Set("organizationId", c.orgID)
	req.Header.Set("Api-Version", HASAPIVersion)

	var ir imageResult

	resp, err := c.client.Do(req, &ir)
	if err != nil {
		return nil, resp, err
	}
	return &ir.Images, resp, nil
}
