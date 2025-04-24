package ibapi

import "fmt"

// DepthMktDataDescription .
type DepthMktDataDescription struct {
	Exchange        string
	SecType         string
	ListingExch     string
	ServiceDataType string
	AggGroup        int64
}

func NewDepthMktDataDescription() DepthMktDataDescription {
	dmdd := DepthMktDataDescription{}
	dmdd.AggGroup = UNSET_INT
	return dmdd
}

// DepthMktDataDescription .
func (d DepthMktDataDescription) String() string {
	return fmt.Sprintf("Exchange: %s, SecType: %s, ListingExchange: %s, ServiceDataType: %s, AggGroup: %s",
		d.Exchange, d.SecType, d.ListingExch, d.ServiceDataType, IntMaxString(d.AggGroup))
}
