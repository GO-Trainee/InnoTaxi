package types

import (
	"database/sql/driver"

	pbentity "awesomeProject/shared/proto/service_name"
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
	financeInvoiceStatusToPBMap = map[FinanceInvoiceStatus]pbentity.Status{
		FinanceInvoiceStatusRequested:     pbentity.Status_INVOICE_STATUS_REQUESTED,
		FinanceInvoiceStatusInvoiced:      pbentity.Status_INVOICE_STATUS_INVOICED,
		FinanceInvoiceStatusPaid:          pbentity.Status_INVOICE_STATUS_PAID,
		FinanceInvoiceStatusPartiallyPaid: pbentity.Status_INVOICE_STATUS_PARTIALLY_PAID,
		FinanceInvoiceStatusVoided:        pbentity.Status_INVOICE_STATUS_VOIDED,
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

	pbToFinanceInvoiceStatusMap = map[pbentity.Status]FinanceInvoiceStatus{
		pbentity.Status_INVOICE_STATUS_REQUESTED:      FinanceInvoiceStatusRequested,
		pbentity.Status_INVOICE_STATUS_INVOICED:       FinanceInvoiceStatusInvoiced,
		pbentity.Status_INVOICE_STATUS_PAID:           FinanceInvoiceStatusPaid,
		pbentity.Status_INVOICE_STATUS_PARTIALLY_PAID: FinanceInvoiceStatusPartiallyPaid,
		pbentity.Status_INVOICE_STATUS_VOIDED:         FinanceInvoiceStatusVoided,
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

func (f FinanceInvoiceStatus) ToPB() pbentity.Status {
	if pb, exists := financeInvoiceStatusToPBMap[f]; exists {
		return pb
	}
	return pbentity.Status_INVOICE_STATUS_UNSPECIFIED
}

func FinanceInvoiceStatusFromString(s string) FinanceInvoiceStatus {
	if status, exists := stringToFinanceInvoiceStatusMap[s]; exists {
		return status
	}
	return FinanceInvoiceStatusUnspecified
}

func FinanceInvoiceStatusFromPB(status pbentity.Status) FinanceInvoiceStatus {
	if s, exists := pbToFinanceInvoiceStatusMap[status]; exists {
		return s
	}
	return FinanceInvoiceStatusUnspecified
}
