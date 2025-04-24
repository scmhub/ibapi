package ibapi

import "fmt"

// FamilyCode .
type FamilyCode struct {
	AccountID     string
	FamilyCodeStr string
}

func NewFamilyCode() FamilyCode {
	return FamilyCode{}
}

func (f FamilyCode) String() string {
	return fmt.Sprintf("AccountId: %s, FamilyCodeStr: %s", f.AccountID, f.FamilyCodeStr)
}
