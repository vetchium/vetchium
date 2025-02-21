package db

// ColleagueState represents the state of a colleague connection
type ColleagueState string

const (
	// ColleaguePending represents a pending colleague connection request
	ColleaguePending ColleagueState = "COLLEAGUING_PENDING"

	// ColleagueAccepted represents an accepted colleague connection
	ColleagueAccepted ColleagueState = "COLLEAGUING_ACCEPTED"

	// ColleagueRejected represents a rejected colleague connection request
	ColleagueRejected ColleagueState = "COLLEAGUING_REJECTED"

	// ColleagueUnlinked represents a previously accepted connection that was unlinked
	ColleagueUnlinked ColleagueState = "COLLEAGUING_UNLINKED"
)
