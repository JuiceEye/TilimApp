package service

import (
	"fmt"
	"tilimauth/internal/model"
	"tilimauth/internal/repository"
)

type MainPageModuleService struct {
	moduleRepo     *repository.ModuleRepository
	sectionRepo    *repository.SectionRepository
	lessonRepo     *repository.LessonRepository
	completionRepo *repository.LessonCompletionRepository
}

func NewMainPageModuleService(
	moduleRepo *repository.ModuleRepository,
	sectionRepo *repository.SectionRepository,
	lessonRepo *repository.LessonRepository,
	completionRepo *repository.LessonCompletionRepository,
) *MainPageModuleService {
	return &MainPageModuleService{
		moduleRepo:     moduleRepo,
		sectionRepo:    sectionRepo,
		lessonRepo:     lessonRepo,
		completionRepo: completionRepo,
	}
}

func (s *MainPageModuleService) GetMainPageModuleByID(moduleID, userID int64) (*model.Module, error) {
	module, err := s.moduleRepo.GetByID(moduleID)
	if err != nil {
		return nil, err
	}

	sections, err := s.sectionRepo.GetByModuleID(module.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to get sections: %w", err)
	}

	if len(sections) == 0 {
		module.Sections = []model.Section{}
		return module, nil
	}

	sectionIDs := make([]int64, len(sections))
	for i, section := range sections {
		sectionIDs[i] = section.ID
	}

	lessonsBySection, err := s.lessonRepo.GetBySectionIDs(sectionIDs)
	if err != nil {
		return nil, fmt.Errorf("failed to get lessons: %w", err)
	}

	var lessonIDs []int64
	for _, lessons := range lessonsBySection {
		for _, lesson := range lessons {
			lessonIDs = append(lessonIDs, lesson.ID)
		}
	}

	completedLessonIDs, err := s.completionRepo.GetCompletedLessonIDs(userID, lessonIDs)

	if err != nil {
		return nil, fmt.Errorf("failed to get completed lessons: %w", err)
	}
	completedSet := make(map[int64]struct{}, len(completedLessonIDs))
	for _, id := range completedLessonIDs {
		completedSet[id] = struct{}{}
	}

	for i := range sections {
		sections[i].Lessons = lessonsBySection[sections[i].ID]

		unlockedSet := false

		for j := range sections[i].Lessons {
			lesson := &sections[i].Lessons[j]
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

	module.Sections = sections
	return module, nil
}
