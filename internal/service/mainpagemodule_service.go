package service

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/redis/go-redis/v9"
	"tilimauth/internal/model"
	"tilimauth/internal/repository"
	"time"
)

type MainPageModuleService struct {
	moduleRepo     *repository.ModuleRepository
	sectionRepo    *repository.SectionRepository
	lessonRepo     *repository.LessonRepository
	completionRepo *repository.LessonCompletionRepository
	redis          *redis.Client
}

func NewMainPageModuleService(
	moduleRepo *repository.ModuleRepository,
	sectionRepo *repository.SectionRepository,
	lessonRepo *repository.LessonRepository,
	completionRepo *repository.LessonCompletionRepository,
	redis *redis.Client,
) *MainPageModuleService {
	return &MainPageModuleService{
		moduleRepo:     moduleRepo,
		sectionRepo:    sectionRepo,
		lessonRepo:     lessonRepo,
		completionRepo: completionRepo,
		redis:          redis,
	}
}

func (s *MainPageModuleService) GetMainPageModuleByID(ctx context.Context, moduleID, userID int64) (*model.Module, error) {
	cacheKey := fmt.Sprintf("module_skeleton:%d", moduleID)
	var module model.Module
	cached, err := s.redis.Get(ctx, cacheKey).Result()

	isCached := false

	if err == nil {
		fmt.Println("CACHED DATA SPOTTED")
		if err := json.Unmarshal([]byte(cached), &module); err == nil {
			isCached = true
		}
	}

	if !isCached {
		modulePtr, err := s.moduleRepo.GetByID(moduleID)
		if err != nil {
			return nil, err
		}
		module = *modulePtr

		sections, err := s.sectionRepo.GetByModuleID(module.ID)
		if err != nil {
			return nil, fmt.Errorf("failed to get sections: %w", err)
		}
		if len(sections) == 0 {
			module.Sections = []model.Section{}
			return &module, nil
		}

		sectionIDs := make([]int64, len(sections))
		for i, sec := range sections {
			sectionIDs[i] = sec.ID
		}

		lessonsBySection, err := s.lessonRepo.GetBySectionIDs(sectionIDs)
		if err != nil {
			return nil, fmt.Errorf("failed to get lessons: %w", err)
		}

		for i := range sections {
			sections[i].Lessons = lessonsBySection[sections[i].ID]
		}
		module.Sections = sections

		raw, _ := json.Marshal(module)
		fmt.Println("saving cache...")
		s.redis.Set(ctx, cacheKey, raw, time.Hour)
	}

	var lessonIDs []int64
	for _, sec := range module.Sections {
		for _, lesson := range sec.Lessons {
			lessonIDs = append(lessonIDs, lesson.ID)
		}
	}

	completedIDs, err := s.completionRepo.GetCompletedLessonIDs(userID, lessonIDs)
	if err != nil {
		return nil, fmt.Errorf("failed to get completed lessons: %w", err)
	}

	completedSet := make(map[int64]struct{}, len(completedIDs))
	for _, id := range completedIDs {
		completedSet[id] = struct{}{}
	}

	for i := range module.Sections {
		unlockedSet := false
		for j := range module.Sections[i].Lessons {
			lesson := &module.Sections[i].Lessons[j]

			if _, ok := completedSet[lesson.ID]; ok {
				lesson.Status = model.StatusCompleted
			} else if !unlockedSet {
				lesson.Status = model.StatusUnlocked
				unlockedSet = true
			} else {
				lesson.Status = model.StatusLocked
			}
		}
	}

	return &module, nil
}
