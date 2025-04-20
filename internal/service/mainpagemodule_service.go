package service

import (
	"fmt"
	"tilimauth/internal/model"
	"tilimauth/internal/repository"
)

type MainPageModuleService struct {
	moduleRepo  *repository.ModuleRepository
	sectionRepo *repository.SectionRepository
	lessonRepo  *repository.LessonRepository
}

func NewMainPageModuleService(
	moduleRepo *repository.ModuleRepository,
	sectionRepo *repository.SectionRepository,
	lessonRepo *repository.LessonRepository,
) *MainPageModuleService {
	return &MainPageModuleService{
		moduleRepo:  moduleRepo,
		sectionRepo: sectionRepo,
		lessonRepo:  lessonRepo,
	}
}

func (s *MainPageModuleService) GetMainPageModuleByID(ModuleID int64) (*model.Module, error) {
	module, err := s.moduleRepo.GetByID(ModuleID)
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

	for i := range sections {
		sections[i].Lessons = lessonsBySection[sections[i].ID]
		if sections[i].Lessons == nil {
			sections[i].Lessons = []model.Lesson{} // Ensure empty array rather than nil
		}
	}

	module.Sections = sections
	return module, nil
}
