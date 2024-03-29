package containerscan

const (
	//defines Relevancy as enum-like
	Unknown   = "Unknown"
	Relevant  = "Relevant"
	Irelevant = "Irelevant"
	NoSP      = "No signature profile to compare"

	//Clair Severities
	UnknownSeverity    = "Unknown"
	NegligibleSeverity = "Negligible"
	LowSeverity        = "Low"
	MediumSeverity     = "Medium"
	HighSeverity       = "High"
	CriticalSeverity   = "Critical"

	ContainerScanRedisPrefix = "_containerscan"

	UnknownScore    = 1
	NegligibleScore = 100
	LowScore        = 200
	MediumScore     = 300
	HighScore       = 400
	CriticalScore   = 500
)

var KnownSeverities = map[string]bool{
	UnknownSeverity:    true,
	NegligibleSeverity: true,
	LowSeverity:        true,
	MediumSeverity:     true,
	HighSeverity:       true,
	CriticalSeverity:   true,
}

var knowScores = map[int]string{
	UnknownScore:    UnknownSeverity,
	NegligibleScore: NegligibleSeverity,
	LowScore:        LowSeverity,
	MediumScore:     MediumSeverity,
	HighScore:       HighSeverity,
	CriticalScore:   CriticalSeverity,
}

func SeverityScoreToString(score int) string {
	return knowScores[score]
}

func CalculateFixed(Fixes []FixedIn) int {
	for _, fix := range Fixes {
		if fix.Version != "None" && fix.Version != "" {
			return 1
		}
	}
	return 0
}
