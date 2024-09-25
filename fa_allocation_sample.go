package ibapi

import (
	"encoding/xml"
)

// ListOf Groups .
type ListOfGroups struct {
	XMLName xml.Name `xml:"ListOfGroups"`
	Groups  []Group
}

// Group .
// DefaultMethod is either one of the IB-computed allocation methods (AvailableEquity, Equal, NetLiq)
// or a user-specified allocation methods (MonetaryAmount, Percent, Ratio, ContractsOrShares).
// AvailableEquity method requires you to specify an order size. This method distributes shares based on the amount of available equity in each account.
// The system calculates ratios based on the Available Equity in each account and allocates shares based on these ratios.
// Equal method requires you to specify an order size. This method distributes shares equally between all accounts in the group.
// NetLiq method requires you to specify an order size. This method distributes shares based on the net liquidation value of each account.
// The system calculates ratios based on the Net Liquidation value in each account and allocates shares based on these ratios.
// MonetaryAmount method calculates the number of units to be allocated based on the monetary value assigned to each account.
// Ratio method calculates the allocation of shares based on the ratios you enter.
// ContractsOrShares method allocates the absolute number of shares you enter to each account listed.
// If you use this method, the order size is calculated by adding together the number of shares allocated to each account in the profile.
type Group struct {
	XMLName        xml.Name `xml:"Group"`
	Name           string   `xml:"name"`
	DefaultMethod  string   `xml:"defaultMethod"`
	ListOfAccounts ListOfAccounts
}

// ListOfAccounts .
type ListOfAccounts struct {
	XMLName  xml.Name `xml:"ListOfAccts"`
	VarName  string   `xml:"varName,attr"`
	Accounts []Account
}

// Account .
type Account struct {
	XMLName xml.Name `xml:"Account"`
	ID      string   `xml:"acct"`
	Amount  string   `xml:"amount,omitempty"`
}

// FAUpdatedGroup returns a list of groups xml sample
func FAUpdatedGroup() string {
	log := ListOfGroups{
		Groups: []Group{
			{
				Name:          "MyTestProfile1",
				DefaultMethod: "ContractsOrShares",
				ListOfAccounts: ListOfAccounts{
					VarName: "list",
					Accounts: []Account{
						{ID: "DU6202167", Amount: "100.0"},
						{ID: "DU6202168", Amount: "200.0"},
					},
				},
			},
			{
				Name:          "MyTestProfile2",
				DefaultMethod: "Ratio",
				ListOfAccounts: ListOfAccounts{
					VarName: "list",
					Accounts: []Account{
						{ID: "DU6202167", Amount: "1.0"},
						{ID: "DU6202168", Amount: "2.0"},
					},
				},
			},
			{
				Name:          "MyTestProfile3",
				DefaultMethod: "Percent",
				ListOfAccounts: ListOfAccounts{
					VarName: "list",
					Accounts: []Account{
						{ID: "DU6202167", Amount: "60.0"},
						{ID: "DU6202168", Amount: "40.0"},
					},
				},
			},
			{
				Name:          "MyTestProfile4",
				DefaultMethod: "MonetaryAmount",
				ListOfAccounts: ListOfAccounts{
					VarName: "list",
					Accounts: []Account{
						{ID: "DU6202167", Amount: "1000.0"},
						{ID: "DU6202168", Amount: "2000.0"},
					},
				},
			},
			{
				Name:          "Group_1",
				DefaultMethod: "NetLiq",
				ListOfAccounts: ListOfAccounts{
					VarName: "list",
					Accounts: []Account{
						{ID: "DU6202167"},
						{ID: "DU6202168"},
					},
				},
			},
			{
				Name:          "MyTestGroup1",
				DefaultMethod: "AvailableEquity",
				ListOfAccounts: ListOfAccounts{
					VarName: "list",
					Accounts: []Account{
						{ID: "DU6202167"},
						{ID: "DU6202168"},
					},
				},
			},
			{
				Name:          "MyTestGroup2",
				DefaultMethod: "NetLiq",
				ListOfAccounts: ListOfAccounts{
					VarName: "list",
					Accounts: []Account{
						{ID: "DU6202167"},
						{ID: "DU6202168"},
					},
				},
			},
			{
				Name:          "MyTestGroup3",
				DefaultMethod: "Equal",
				ListOfAccounts: ListOfAccounts{
					VarName: "list",
					Accounts: []Account{
						{ID: "DU6202167"},
						{ID: "DU6202168"},
					},
				},
			},
			{
				Name:          "Group_2",
				DefaultMethod: "AvailableEquity",
				ListOfAccounts: ListOfAccounts{
					VarName: "list",
					Accounts: []Account{
						{ID: "DU6202167"},
						{ID: "DU6202168"},
					},
				},
			},
		},
	}

	groupXml, _ := xml.MarshalIndent(log, " ", "  ")
	return xml.Header + string(groupXml)
}

/*
<?xml version="1.0" encoding="UTF-8"?>
 <ListOfGroups>
   <Group>
     <name>MyTestProfile1</name>
     <defaultMethod>ContractsOrShares</defaultMethod>
     <ListOfAccts varName="list">
       <Account>
         <acct>DU6202167</acct>
         <amount>100.0</amount>
       </Account>
       <Account>
         <acct>DU6202168</acct>
         <amount>200.0</amount>
       </Account>
     </ListOfAccts>
   </Group>
   <Group>
     <name>MyTestProfile2</name>
     <defaultMethod>Ratio</defaultMethod>
     <ListOfAccts varName="list">
       <Account>
         <acct>DU6202167</acct>
         <amount>1.0</amount>
       </Account>
       <Account>
         <acct>DU6202168</acct>
         <amount>2.0</amount>
       </Account>
     </ListOfAccts>
   </Group>
   <Group>
     <name>MyTestProfile3</name>
     <defaultMethod>Percent</defaultMethod>
     <ListOfAccts varName="list">
       <Account>
         <acct>DU6202167</acct>
         <amount>60.0</amount>
       </Account>
       <Account>
         <acct>DU6202168</acct>
         <amount>40.0</amount>
       </Account>
     </ListOfAccts>
   </Group>
   <Group>
     <name>MyTestProfile4</name>
     <defaultMethod>MonetaryAmount</defaultMethod>
     <ListOfAccts varName="list">
       <Account>
         <acct>DU6202167</acct>
         <amount>1000.0</amount>
       </Account>
       <Account>
         <acct>DU6202168</acct>
         <amount>2000.0</amount>
       </Account>
     </ListOfAccts>
   </Group>
   <Group>
     <name>Group_1</name>
     <defaultMethod>NetLiq</defaultMethod>
     <ListOfAccts varName="list">
       <Account>
         <acct>DU6202167</acct>
       </Account>
       <Account>
         <acct>DU6202168</acct>
       </Account>
     </ListOfAccts>
   </Group>
   <Group>
     <name>MyTestGroup1</name>
     <defaultMethod>AvailableEquity</defaultMethod>
     <ListOfAccts varName="list">
       <Account>
         <acct>DU6202167</acct>
       </Account>
       <Account>
         <acct>DU6202168</acct>
       </Account>
     </ListOfAccts>
   </Group>
   <Group>
     <name>MyTestGroup2</nameout>
     <defaultMethod>NetLiq</defaultMethod>
     <ListOfAccts varName="list">
       <Account>
         <acct>DU6202167</acct>
       </Account>
       <Account>
         <acct>DU6202168</acct>
       </Account>
     </ListOfAccts>
   </Group>
   <Group>
     <name>MyTestGroup3</name>
     <defaultMethod>Equal</defaultMethod>
     <ListOfAccts varName="list">
       <Account>
         <acct>DU6202167</acct>
       </Account>
       <Account>
         <acct>DU6202168</acct>
       </Account>
     </ListOfAccts>
   </Group>
   <Group>
     <name>Group_2</name>
     <defaultMethod>AvailableEquity</defaultMethod>
     <ListOfAccts varName="list">
       <Account>
         <acct>DU6202167</acct>
       </Account>
       <Account>
         <acct>DU6202168</acct>
       </Account>
     </ListOfAccts>
   </Group>
 </ListOfGroups>
*/
