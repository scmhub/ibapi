package ibapi

import "fmt"

// TagValue maps a tag to a value. Both of them are strings.
// They are used in a slice to convey extra info with the requests.
type TagValue struct {
	Tag   string
	Value string
}

func NewTagValue() TagValue {
	return TagValue{}
}

func (tv TagValue) String() string {
	return fmt.Sprintf("%s=%s;", tv.Tag, tv.Value)
}
