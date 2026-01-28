package domain

import "time"

// AuditEventType limits what we can log to the DB
type AuditEventType string

const (
	AuditUserLogin           AuditEventType = "user_login"
	AuditHMRCConnected       AuditEventType = "hmrc_connected"
	AuditUploadIngested      AuditEventType = "upload_ingested"
	AuditSubmissionSuccess   AuditEventType = "submission_success"
	AuditPurgeOnSubmission   AuditEventType = "purge_on_submission"
	AuditPurgeAbandonedDraft AuditEventType = "purge_abandoned_draft"
)

// PurgeEventPayload is the JSON structure for the 'details' column.
// It proves WHAT was deleted without saving the sensitive data itself.
type PurgeEventPayload struct {
	PolicyVersion string    `json:"policy_version"`
	Stream        string    `json:"stream"` // "trade" or "property"
	Reason        string    `json:"reason"` // "submission_success" or "ttl_expired"
	RowsDeleted   int64     `json:"rows_deleted"`
	Timestamp     time.Time `json:"timestamp"`
	SubmissionID  string    `json:"submission_id,omitempty"`
}
