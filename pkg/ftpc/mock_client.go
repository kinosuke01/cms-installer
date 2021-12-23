package ftpc

import "errors"

type MockClient struct {
	ErrorMessage string
	FilePath     string
	Content      string
}

func (c *MockClient) Upload(filePath string, pContent *string) error {
	c.FilePath = filePath
	c.Content = *pContent

	if c.ErrorMessage != "" {
		return errors.New(c.ErrorMessage)
	} else {
		return nil
	}
}

func (c *MockClient) Delete(filePath string, onlyExists bool) error {
	c.FilePath = filePath

	if c.ErrorMessage != "" {
		return errors.New(c.ErrorMessage)
	} else {
		return nil
	}
}
