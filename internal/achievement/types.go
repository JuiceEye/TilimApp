package achievement

// EventType represents the type of event that can trigger an achievement
type EventType string

const (
	EventLessonCompleted EventType = "lesson_completed"
	EventStreakUpdated   EventType = "streak_updated"
	EventXPGained        EventType = "xp_gained"
)

// EventContext contains information about the event that triggered an achievement check
type EventContext struct {
	UserID    int64
	EventType EventType
	Payload   map[string]interface{}
}

// Achievement defines the interface for all achievements
type Achievement interface {
	// ID returns the unique identifier for this achievement
	ID() string

	// Trigger returns the event type that should trigger this achievement check
	Trigger() EventType

	// Check determines if the achievement conditions are met
	Check(ctx EventContext) (bool, error)

	// Grant awards the achievement to the user
	Grant(ctx EventContext) (awardableXP int, err error)
}
