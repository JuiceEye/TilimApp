package achievement

type EventType string

const (
	EventLessonCompleted EventType = "lesson_completed"
	EventStreakUpdated   EventType = "streak_updated"
	EventXPGained        EventType = "xp_gained"
)

type EventContext struct {
	UserID    int64
	EventType EventType
	Payload   map[string]interface{}
}

type Achievement interface {
	ID() string

	Trigger() EventType

	Check(ctx EventContext) (bool, error)

	Grant(ctx EventContext) (awardableXP int, err error)
}
