package types

import (
	"database/sql/driver"

	pbentity "awesomeProject/shared/proto/service_name"
)

type FeeFreeStatus int

const (
	FeeFreeStatusUnspecified FeeFreeStatus = iota
	FeeFreeStatusActive
	FeeFreeStatusArchived
	FeeFreeStatusExpireSoon
	FeeFreeStatusCompleted
	FeeFreeStatusExpired
	FeeFreeStatusNoFeeFree
	FeeFreeStatusPendingCompletion
)

const (
	FeeFreeStatusUnspecifiedString       = "unspecified"
	FeeFreeStatusActiveString            = "active"
	FeeFreeStatusArchivedString          = "archived"
	FeeFreeStatusExpireSoonString        = "expire_soon"
	FeeFreeStatusCompletedString         = "completed"
	FeeFreeStatusExpiredString           = "expired"
	FeeFreeStatusNoFeeFreeString         = "no_fee_free"
	FeeFreeStatusPendingCompletionString = "pending_completion"
)

var (
	pbToFeeFreeStatusMap = map[pbentity.FeeFreeStatus]FeeFreeStatus{
		pbentity.FeeFreeStatus_FEE_FREE_STATUS_ACTIVE:             FeeFreeStatusActive,
		pbentity.FeeFreeStatus_FEE_FREE_STATUS_ARCHIVED:           FeeFreeStatusArchived,
		pbentity.FeeFreeStatus_FEE_FREE_STATUS_EXPIRE_SOON:        FeeFreeStatusExpireSoon,
		pbentity.FeeFreeStatus_FEE_FREE_STATUS_COMPLETED:          FeeFreeStatusCompleted,
		pbentity.FeeFreeStatus_FEE_FREE_STATUS_EXPIRED:            FeeFreeStatusExpired,
		pbentity.FeeFreeStatus_FEE_FREE_STATUS_NO_FEE_FREE:        FeeFreeStatusNoFeeFree,
		pbentity.FeeFreeStatus_FEE_FREE_STATUS_PENDING_COMPLETION: FeeFreeStatusPendingCompletion,
	}

	feeFreeStatusToPBMap = map[FeeFreeStatus]pbentity.FeeFreeStatus{
		FeeFreeStatusActive:            pbentity.FeeFreeStatus_FEE_FREE_STATUS_ACTIVE,
		FeeFreeStatusArchived:          pbentity.FeeFreeStatus_FEE_FREE_STATUS_ARCHIVED,
		FeeFreeStatusExpireSoon:        pbentity.FeeFreeStatus_FEE_FREE_STATUS_EXPIRE_SOON,
		FeeFreeStatusCompleted:         pbentity.FeeFreeStatus_FEE_FREE_STATUS_COMPLETED,
		FeeFreeStatusExpired:           pbentity.FeeFreeStatus_FEE_FREE_STATUS_EXPIRED,
		FeeFreeStatusNoFeeFree:         pbentity.FeeFreeStatus_FEE_FREE_STATUS_NO_FEE_FREE,
		FeeFreeStatusPendingCompletion: pbentity.FeeFreeStatus_FEE_FREE_STATUS_PENDING_COMPLETION,
	}

	feeFreeStatusToStringMap = map[FeeFreeStatus]string{
		FeeFreeStatusActive:            FeeFreeStatusActiveString,
		FeeFreeStatusArchived:          FeeFreeStatusArchivedString,
		FeeFreeStatusExpireSoon:        FeeFreeStatusExpireSoonString,
		FeeFreeStatusCompleted:         FeeFreeStatusCompletedString,
		FeeFreeStatusExpired:           FeeFreeStatusExpiredString,
		FeeFreeStatusNoFeeFree:         FeeFreeStatusNoFeeFreeString,
		FeeFreeStatusPendingCompletion: FeeFreeStatusPendingCompletionString,
	}

	stringToFeeFreeStatusMap = map[string]FeeFreeStatus{
		FeeFreeStatusActiveString:            FeeFreeStatusActive,
		FeeFreeStatusArchivedString:          FeeFreeStatusArchived,
		FeeFreeStatusExpireSoonString:        FeeFreeStatusExpireSoon,
		FeeFreeStatusCompletedString:         FeeFreeStatusCompleted,
		FeeFreeStatusExpiredString:           FeeFreeStatusExpired,
		FeeFreeStatusNoFeeFreeString:         FeeFreeStatusNoFeeFree,
		FeeFreeStatusPendingCompletionString: FeeFreeStatusPendingCompletion,
	}
)

func (f *FeeFreeStatus) Scan(value interface{}) error {
	switch data := value.(type) {
	case string:
		*f = FeeFreeStatusFromString(data)
	case []byte:
		*f = FeeFreeStatusFromString(string(data))
	default:
		return errorsx.ErrTypeAssertionToByte
	}

	return nil
}

func (f FeeFreeStatus) Value() (driver.Value, error) {
	return f.String(), nil
}

func FeeFreeStatusFromPB(status pbentity.FeeFreeStatus) FeeFreeStatus {
	if s, exists := pbToFeeFreeStatusMap[status]; exists {
		return s
	}
	return FeeFreeStatusUnspecified
}

func (f FeeFreeStatus) ToPB() pbentity.FeeFreeStatus {
	if pb, exists := feeFreeStatusToPBMap[f]; exists {
		return pb
	}
	return pbentity.FeeFreeStatus_FEE_FREE_STATUS_UNSPECIFIED
}

func (f FeeFreeStatus) String() string {
	if str, exists := feeFreeStatusToStringMap[f]; exists {
		return str
	}
	return FeeFreeStatusUnspecifiedString
}

func FeeFreeStatusFromString(s string) FeeFreeStatus {
	if status, exists := stringToFeeFreeStatusMap[s]; exists {
		return status
	}
	return FeeFreeStatusUnspecified
}
