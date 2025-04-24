package ibapi

import "fmt"

// NewsProvider .
type NewsProvider struct {
	Code string
	Name string
}

func NewNewsProvider() NewsProvider {
	return NewsProvider{}
}

func (np NewsProvider) String() string {
	return fmt.Sprintf("Code: %s, Name: %s", np.Code, np.Name)
}
