package engines

// EngineType represents a type of statuspage.
type EngineType int

const (
	// Unknown is a default value for EngineType.
	Unknown EngineType = iota
	// StatusPageIO represents statuspage.io engine.
	StatusPageIO
	// StatusIO represents status.io engine.
	StatusIO
)
