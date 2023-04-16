package types

// EngineType represents a type of statuspage.
type EngineType int

const (
	// UnknownType is a default value for EngineType.
	UnknownType EngineType = iota
	// StatusPageIOType represents statuspage.io engine.
	StatusPageIOType
	// StatusIOType represents status.io engine.
	StatusIOType
)

// Status represents a status of statuspage.
type Status int

const (
	// UnknownStatus is a default value for Status.
	UnknownStatus Status = iota
	// OperationalStatus represents operational status.
	OperationalStatus
	// PlannedMaintenanceStatus represents planned maintenance status.
	PlannedMaintenanceStatus
	// DegradedPerformanceStatus represents degraded performance status.
	DegradedPerformanceStatus
	// PartialOutageStatus represents partial outage status.
	PartialOutageStatus
	// MajorOutageStatus represents major outage status.
	MajorOutageStatus
	// SecurityIssueStatus represents security issue status.
	SecurityIssueStatus
)
