package types

import (
	"database/sql/driver"
)

type FinanceInvoiceStatus int

const (
	FinanceInvoiceStatusUnspecified FinanceInvoiceStatus = iota
	FinanceInvoiceStatusDraft
	FinanceInvoiceStatusRequested
	FinanceInvoiceStatusInvoiced
	FinanceInvoiceStatusPaid
	FinanceInvoiceStatusPartiallyPaid
	FinanceInvoiceStatusVoided
)

const (
	FinanceInvoiceStatusUnspecifiedString   = "unspecified"
	FinanceInvoiceStatusDraftString         = "draft"
	FinanceInvoiceStatusRequestedString     = "requested"
	FinanceInvoiceStatusInvoicedString      = "invoiced"
	FinanceInvoiceStatusPaidString          = "paid"
	FinanceInvoiceStatusPartiallyPaidString = "partially_paid"
	FinanceInvoiceStatusVoidedString        = "voided"
)

var (
	financeInvoiceStatusToPBMap = map[FinanceInvoiceStatus]int32{
		// FinanceInvoiceStatusDraft is not expected in response
		FinanceInvoiceStatusRequested:     2,
		FinanceInvoiceStatusInvoiced:      3,
		FinanceInvoiceStatusPaid:          4,
		FinanceInvoiceStatusPartiallyPaid: 5,
		FinanceInvoiceStatusVoided:        6,
	}

	financeInvoiceStatusToStringMap = map[FinanceInvoiceStatus]string{
		FinanceInvoiceStatusDraft:         FinanceInvoiceStatusDraftString,
		FinanceInvoiceStatusRequested:     FinanceInvoiceStatusRequestedString,
		FinanceInvoiceStatusInvoiced:      FinanceInvoiceStatusInvoicedString,
		FinanceInvoiceStatusPaid:          FinanceInvoiceStatusPaidString,
		FinanceInvoiceStatusPartiallyPaid: FinanceInvoiceStatusPartiallyPaidString,
		FinanceInvoiceStatusVoided:        FinanceInvoiceStatusVoidedString,
	}

	stringToFinanceInvoiceStatusMap = map[string]FinanceInvoiceStatus{
		FinanceInvoiceStatusDraftString:         FinanceInvoiceStatusDraft,
		FinanceInvoiceStatusRequestedString:     FinanceInvoiceStatusRequested,
		FinanceInvoiceStatusInvoicedString:      FinanceInvoiceStatusInvoiced,
		FinanceInvoiceStatusPaidString:          FinanceInvoiceStatusPaid,
		FinanceInvoiceStatusPartiallyPaidString: FinanceInvoiceStatusPartiallyPaid,
		FinanceInvoiceStatusVoidedString:        FinanceInvoiceStatusVoided,
	}

	pbToFinanceInvoiceStatusMap = map[int32]FinanceInvoiceStatus{
		// FinanceInvoiceStatusDraft is not expected in response
		2: FinanceInvoiceStatusRequested,
		3: FinanceInvoiceStatusInvoiced,
		4: FinanceInvoiceStatusPaid,
		5: FinanceInvoiceStatusPartiallyPaid,
		6: FinanceInvoiceStatusVoided,
	}
)

func (f *FinanceInvoiceStatus) Scan(value interface{}) error {
	switch data := value.(type) {
	case string:
		*f = FinanceInvoiceStatusFromString(data)
	case []byte:
		*f = FinanceInvoiceStatusFromString(string(data))
	default:
		return errorsx.ErrTypeAssertionToByte
	}

	return nil
}

func (f FinanceInvoiceStatus) Value() (driver.Value, error) {
	return f.String(), nil
}

func (f FinanceInvoiceStatus) String() string {
	if str, exists := financeInvoiceStatusToStringMap[f]; exists {
		return str
	}
	return FinanceInvoiceStatusUnspecifiedString
}

func (f FinanceInvoiceStatus) ToPB() int32 {
	if pb, exists := financeInvoiceStatusToPBMap[f]; exists {
		return pb
	}
	return 0
}

func FinanceInvoiceStatusFromString(s string) FinanceInvoiceStatus {
	if status, exists := stringToFinanceInvoiceStatusMap[s]; exists {
		return status
	}
	return FinanceInvoiceStatusUnspecified
}

func FinanceInvoiceStatusFromPB(financeInvoiceStatusType int32) FinanceInvoiceStatus {
	if status, exists := pbToFinanceInvoiceStatusMap[financeInvoiceStatusType]; exists {
		return status
	}
	return FinanceInvoiceStatusUnspecified
}
