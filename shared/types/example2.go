package types

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
	pbToFeeFreeStatusMap = map[int32]FeeFreeStatus{
		1: FeeFreeStatusActive,
		2: FeeFreeStatusArchived,
		3: FeeFreeStatusExpireSoon,
		4: FeeFreeStatusCompleted,
		5: FeeFreeStatusExpired,
		6: FeeFreeStatusNoFeeFree,
		7: FeeFreeStatusPendingCompletion,
	}

	feeFreeStatusToPBMap = map[FeeFreeStatus]int32{
		FeeFreeStatusActive:            1,
		FeeFreeStatusArchived:          2,
		FeeFreeStatusExpireSoon:        3,
		FeeFreeStatusCompleted:         4,
		FeeFreeStatusExpired:           5,
		FeeFreeStatusNoFeeFree:         6,
		FeeFreeStatusPendingCompletion: 7,
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

func FeeFreeStatusFromPB(pbStatus int32) FeeFreeStatus {
	if status, exists := pbToFeeFreeStatusMap[pbStatus]; exists {
		return status
	}
	return FeeFreeStatusUnspecified
}

func (f FeeFreeStatus) ToPB() int32 {
	if pb, exists := feeFreeStatusToPBMap[f]; exists {
		return pb
	}
	return 0
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
