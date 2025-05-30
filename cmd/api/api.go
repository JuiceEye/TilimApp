package api

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"tilimauth/internal/handler"
	"tilimauth/internal/middleware"
	"tilimauth/internal/repository"
	"tilimauth/internal/service"
)

type Server struct {
	address string
	db      *sql.DB
}

func NewServer(address string, db *sql.DB) *Server {
	return &Server{
		address: address,
		db:      db,
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

	profileService := service.NewProfileService(userRepo, userProgressRepo)
	profileHandler := handler.NewProfileHandler(profileService)
	profileHandler.RegisterRoutes(protectedRouter)

	moduleRepo := repository.NewModuleRepo(s.db)
	sectionRepo := repository.NewSectionRepo(s.db)
	lessonRepo := repository.NewLessonRepo(s.db)
	lessonCompletionRepo := repository.NewLessonCompletionRepo(s.db)
	moduleMainPageService := service.NewMainPageModuleService(moduleRepo, sectionRepo, lessonRepo, lessonCompletionRepo)
	moduleMainPageHandler := handler.NewMainPageModuleHandler(moduleMainPageService)
	moduleMainPageHandler.RegisterRoutes(protectedRouter)

	answerRepo := repository.NewAnswerRepo(s.db)
	exerciseRepo := repository.NewExerciseRepo(s.db)
	lessonService := service.NewLessonService(lessonRepo, exerciseRepo, answerRepo)
	lessonCompletionService := service.NewLessonCompletionService(lessonRepo, lessonCompletionRepo, userRepo)
	lessonHandler := handler.NewLessonHandler(lessonService, lessonCompletionService)
	lessonHandler.RegisterRoutes(protectedRouter)

	leaderboardsService := service.NewLeaderboardsService(userRepo)
	leaderboardsHandler := handler.NewLeaderboardsHandler(leaderboardsService)
	leaderboardsHandler.RegisterRoutes(protectedRouter)

	deleteUserDlyaFrontov(publicRouter, s.db)

	mainRouter := http.NewServeMux()

	standardChain := middleware.CreateChain(
		middleware.Logging,
		middleware.SetCorsPolicy,
	)

	protectedChain := middleware.CreateChain(
		middleware.Auth,
	)

	mainRouter.Handle("/auth/", publicRouter)

	// idk chatgpt said I need to do this instead of Handle(), no clue what the difference is
	mainRouter.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		protectedChain(protectedRouter).ServeHTTP(w, r)
	})

	log.Printf("[INFO] Starting server on %s...", s.address)
	fmt.Println("***************************************************************************************************************************************")

	return http.ListenAndServe(s.address, standardChain(mainRouter))
}

// delete later
func deleteUserDlyaFrontov(router *http.ServeMux, db *sql.DB) {
	router.HandleFunc("DELETE /users/{user_id}", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		id, _ := strconv.Atoi(r.PathValue("user_id"))
		query := `DELETE FROM auth.users WHERE id = $1 RETURNING id`
		var deletedUserID int
		err := db.QueryRow(query, id).Scan(&deletedUserID)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			err := json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
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
