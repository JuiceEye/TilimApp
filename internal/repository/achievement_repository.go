package repository

import (
	"database/sql"
	"fmt"
	"tilimauth/internal/model"
	"time"
)

type AchievementRepository struct {
	db *sql.DB
}

func NewAchievementRepository(db *sql.DB) *AchievementRepository {
	return &AchievementRepository{
		db: db,
	}
}

// GetAllAchievements retrieves all achievements from the database
func (r *AchievementRepository) GetAllAchievements() ([]model.Achievement, error) {
	query := `SELECT id, code, name, description, xp_reward, created_at FROM app.achievements`
	
	rows, err := r.db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("failed to get achievements: %w", err)
	}
	defer rows.Close()
	
	var achievements []model.Achievement
	for rows.Next() {
		var achievement model.Achievement
		err := rows.Scan(
			&achievement.ID,
			&achievement.Code,
			&achievement.Name,
			&achievement.Description,
			&achievement.XPReward,
			&achievement.CreatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan achievement: %w", err)
		}
		achievements = append(achievements, achievement)
	}
	
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating achievements: %w", err)
	}
	
	return achievements, nil
}

// GetAchievementByCode retrieves an achievement by its code
func (r *AchievementRepository) GetAchievementByCode(code string) (*model.Achievement, error) {
	query := `SELECT id, code, name, description, xp_reward, created_at FROM app.achievements WHERE code = $1`
	
	var achievement model.Achievement
	err := r.db.QueryRow(query, code).Scan(
		&achievement.ID,
		&achievement.Code,
		&achievement.Name,
		&achievement.Description,
		&achievement.XPReward,
		&achievement.CreatedAt,
	)
	
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrNotFound
		}
		return nil, fmt.Errorf("failed to get achievement by code: %w", err)
	}
	
	return &achievement, nil
}

// GetUserAchievements retrieves all achievements earned by a user
func (r *AchievementRepository) GetUserAchievements(userID int64) ([]model.UserAchievement, error) {
	query := `
		SELECT ua.id, ua.user_id, ua.achievement_id, ua.earned_at 
		FROM app.user_achievements ua
		WHERE ua.user_id = $1
	`
	
	rows, err := r.db.Query(query, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user achievements: %w", err)
	}
	defer rows.Close()
	
	var userAchievements []model.UserAchievement
	for rows.Next() {
		var ua model.UserAchievement
		err := rows.Scan(
			&ua.ID,
			&ua.UserID,
			&ua.AchievementID,
			&ua.EarnedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan user achievement: %w", err)
		}
		userAchievements = append(userAchievements, ua)
	}
	
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating user achievements: %w", err)
	}
	
	return userAchievements, nil
}

// HasUserEarnedAchievement checks if a user has already earned a specific achievement
func (r *AchievementRepository) HasUserEarnedAchievement(userID int64, achievementID int64) (bool, error) {
	query := `SELECT 1 FROM app.user_achievements WHERE user_id = $1 AND achievement_id = $2`
	
	var dummy int
	err := r.db.QueryRow(query, userID, achievementID).Scan(&dummy)
	
	if err != nil {
		if err == sql.ErrNoRows {
			return false, nil
		}
		return false, fmt.Errorf("failed to check if user earned achievement: %w", err)
	}
	
	return true, nil
}

// GrantAchievementTx grants an achievement to a user within a transaction
func (r *AchievementRepository) GrantAchievementTx(tx *sql.Tx, userID int64, achievementID int64) error {
	query := `
		INSERT INTO app.user_achievements (user_id, achievement_id, earned_at)
		VALUES ($1, $2, $3)
		ON CONFLICT (user_id, achievement_id) DO NOTHING
	`
	
	_, err := tx.Exec(query, userID, achievementID, time.Now().UTC())
	if err != nil {
		return fmt.Errorf("failed to grant achievement: %w", err)
	}
	
	return nil
}

// GetLessonCompletionsCountForDay returns the number of lessons completed by a user on a specific day
func (r *AchievementRepository) GetLessonCompletionsCountForDay(userID int64, date time.Time) (int, error) {
	startOfDay := time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, time.UTC)
	endOfDay := startOfDay.Add(24 * time.Hour)
	
	query := `
		SELECT COUNT(*) FROM app.lesson_completions
		WHERE user_id = $1 AND date_completed >= $2 AND date_completed < $3
	`
	
	var count int
	err := r.db.QueryRow(query, userID, startOfDay, endOfDay).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("failed to get lesson completions count: %w", err)
	}
	
	return count, nil
}

// GetTotalLessonCompletionsCount returns the total number of lessons completed by a user
func (r *AchievementRepository) GetTotalLessonCompletionsCount(userID int64) (int, error) {
	query := `SELECT COUNT(*) FROM app.lesson_completions WHERE user_id = $1`
	
	var count int
	err := r.db.QueryRow(query, userID).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("failed to get total lesson completions count: %w", err)
	}
	
	return count, nil
}