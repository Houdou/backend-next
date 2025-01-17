package types

import "gopkg.in/guregu/null.v3"

type ArkDrop struct {
	DropType string `json:"dropType" validate:"required,oneof=REGULAR_DROP NORMAL_DROP SPECIAL_DROP EXTRA_DROP FURNITURE"`
	ItemID   string `json:"itemId" validate:"required,printascii"`
	Quantity int    `json:"quantity" validate:"required,lte=1000"`
}

type Drop struct {
	DropType string
	ItemID   int
	Quantity int
}

type SingleReportRequest struct {
	FragmentStageID
	FragmentReportCommon

	Drops     []ArkDrop `json:"drops" validate:"dive"`
	PenguinID string    `json:"-"`
}

type SingleReportRecallRequest struct {
	ReportHash string `json:"reportHash" validate:"required,printascii"`
}

type BatchReportDrop struct {
	FragmentStageID

	Drops    []ArkDrop             `json:"drops" validate:"dive"`
	Metadata ReportRequestMetadata `json:"metadata" validate:"dive"`
}

type ReportRequestMetadata struct {
	Fingerprint  string      `json:"fingerprint" validate:"lte=128"`
	MD5          null.String `json:"md5" validate:"lte=32"`
	FileName     string      `json:"fileName" validate:"lte=512"`
	LastModified int         `json:"lastModified"`
}

type BatchReportRequest struct {
	FragmentReportCommon

	BatchDrops []BatchReportDrop `json:"batchDrops" validate:"dive"`
}

type BatchReportError struct {
	Index  int    `json:"index"`
	Reason string `json:"reason,omitempty"`
}
