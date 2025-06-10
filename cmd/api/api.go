package api

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/redis/go-redis/v9"
	"log"
	"net/http"
	"strconv"
	"tilimauth/internal/achievement"
	"tilimauth/internal/handler"
	"tilimauth/internal/middleware"
	"tilimauth/internal/repository"
	"tilimauth/internal/service"
)

type Server struct {
	address string
	db      *sql.DB
	redis   *redis.Client
}

func NewServer(address string, db *sql.DB, redis *redis.Client) *Server {
	return &Server{
		address: address,
		db:      db,
		redis:   redis,
	}
}

func (s *Server) Run() error {
	publicRouter := http.NewServeMux()
	protectedRouter := http.NewServeMux()

	userRepo := repository.NewUserRepo(s.db)
	userProgressRepo := repository.NewUserProgressRepo(s.db)

	userService := service.NewAuthService(userRepo, userProgressRepo)
	authHandler := handler.NewAuthHandler(userService)
	authHandler.RegisterRoutes(publicRouter)

	subRepo := repository.NewSubscriptionRepo(s.db)
	dailyTaskRepo := repository.NewDailyTaskRepository(s.db)
	achievementRepo := repository.NewAchievementRepository(s.db)

	// Initialize achievement service
	achievementService := achievement.NewAchievementService(achievementRepo, userRepo)

	// Register default achievements
	achievement.CreateDefaultAchievements(achievementService, achievementRepo, userRepo)

	dailyTaskService := service.NewDailyTaskService(dailyTaskRepo)
	profileService := service.NewProfileService(userRepo, userProgressRepo, subRepo, achievementService)
	profileHandler := handler.NewProfileHandler(profileService)
	profileHandler.RegisterRoutes(protectedRouter)

	dailyTaskHandler := handler.NewDailyTaskHandler(dailyTaskService)
	dailyTaskHandler.RegisterRoutes(protectedRouter)

	moduleRepo := repository.NewModuleRepo(s.db)
	sectionRepo := repository.NewSectionRepo(s.db)
	lessonRepo := repository.NewLessonRepo(s.db)
	lessonCompletionRepo := repository.NewLessonCompletionRepo(s.db)
	moduleMainPageService := service.NewMainPageModuleService(moduleRepo, sectionRepo, lessonRepo, lessonCompletionRepo, s.redis)
	moduleMainPageHandler := handler.NewMainPageModuleHandler(moduleMainPageService)
	moduleMainPageHandler.RegisterRoutes(protectedRouter)

	answerRepo := repository.NewAnswerRepo(s.db)
	exerciseRepo := repository.NewExerciseRepo(s.db)
	lessonService := service.NewLessonService(lessonRepo, exerciseRepo, answerRepo)
	lessonCompletionService := service.NewLessonCompletionService(lessonRepo, lessonCompletionRepo, userRepo, profileService, dailyTaskService, achievementService)
	lessonHandler := handler.NewLessonHandler(lessonService, lessonCompletionService, s.redis)
	lessonHandler.RegisterRoutes(protectedRouter)

	leaderboardsService := service.NewLeaderboardsService(userRepo)
	leaderboardsHandler := handler.NewLeaderboardsHandler(leaderboardsService)
	leaderboardsHandler.RegisterRoutes(protectedRouter)

	subService := service.NewSubscriptionService(subRepo, userRepo)
	subHandler := handler.NewSubscriptionHandler(subService)
	subHandler.RegisterRoutes(protectedRouter)

	deleteUserDlyaFrontov(publicRouter, s.db)

	mainRouter := http.NewServeMux()

	protectedChain := middleware.CreateChain(
		middleware.Auth,
	)

	standardChain := middleware.CreateChain(
		middleware.Logging,
		middleware.SetCorsPolicy,
	)

	// idk chatgpt said I need to do this instead of Handle(), no clue what the difference is
	mainRouter.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		protectedChain(protectedRouter).ServeHTTP(w, r)
	})

	mainRouter.Handle("/auth/", http.StripPrefix("/auth", publicRouter))

	log.Printf("[INFO] Starting server on %s...", s.address)
	fmt.Println("***********************************************************************************************************************************")

	return http.ListenAndServe(s.address, standardChain(mainRouter))
}

// delete later
func deleteUserDlyaFrontov(router *http.ServeMux, db *sql.DB) {
	router.HandleFunc("DELETE /users/{user_id}", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		userID, _ := strconv.Atoi(r.PathValue("user_id"))
		query := `DELETE FROM app.lesson_completions WHERE user_id = $1`
		_, err := db.Exec(query, userID)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			err = json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
			if err != nil {
				return
			}
			return
		}
		query = `DELETE FROM auth.users WHERE id = $1 RETURNING id`
		var deletedUserID int
		err = db.QueryRow(query, userID).Scan(&deletedUserID)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			err = json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
			if err != nil {
				return
			}
			return
		}

		w.WriteHeader(http.StatusOK)
		err = json.NewEncoder(w).Encode(map[string]int{"user_id": deletedUserID})
		if err != nil {
			return
		}
	})
}
