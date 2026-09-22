package results

import (
	"errors"
	"sort"
	"strings"

	"github.com/Mulamba-Mukoma-Olivier/professor-evaluation-api/internal/evaluations"
)

var (
	ErrInvalidProfessorID = errors.New("invalid professor ID")
	ErrInvalidCourseID    = errors.New("invalid course ID")
	ErrAcademicYearRequired = errors.New("academic year is required")
	ErrPeriodRequired       = errors.New("period is required")
	ErrNoEvaluations        = errors.New("no evaluations found")
	ErrNoAnswers            = errors.New("evaluations contain no answers")
)

type ResultsRepository interface {
	GetByProfessor(
		professorID int,
		courseID int,
		academicYear string,
		period string,
	) ([]evaluations.Evaluation, error)
}

type Service struct {
	repository ResultsRepository
}

func NewService(repository ResultsRepository) *Service {
	return &Service{
		repository: repository,
	}
}

func (s *Service) GetProfessorResult(
	professorID int,
	courseID int,
	academicYear string,
	period string,
) (*ProfessorResult, error) {

	if professorID <= 0 {
		return nil, ErrInvalidProfessorID
	}

	if courseID <= 0 {
		return nil, ErrInvalidCourseID
	}

	academicYear = strings.TrimSpace(academicYear)

	if academicYear == "" {
		return nil, ErrAcademicYearRequired
	}

	period = strings.TrimSpace(period)

	if period == "" {
		return nil, ErrPeriodRequired
	}

	evaluationList, err := s.repository.GetByProfessor(
		professorID,
		courseID,
		academicYear,
		period,
	)

	if err != nil {
		return nil, err
	}

	if len(evaluationList) == 0 {
		return nil, ErrNoEvaluations
	}

	// criterionID -> total des scores
	totals := make(map[int]int)

	// criterionID -> nombre de réponses
	counts := make(map[int]int)

	globalTotal := 0
	globalCount := 0

	for _, evaluation := range evaluationList {
		for _, answer := range evaluation.Answers {

			// Ignore les réponses avec un critère invalide.
			if answer.CriterionID <= 0 {
				continue
			}

			totals[answer.CriterionID] += answer.Score
			counts[answer.CriterionID]++

			globalTotal += answer.Score
			globalCount++
		}
	}

	if globalCount == 0 {
		return nil, ErrNoAnswers
	}

	criteria := make([]CriterionResult, 0, len(totals))

	for criterionID, total := range totals {
		count := counts[criterionID]

		if count == 0 {
			continue
		}

		average := float64(total) / float64(count)

		criteria = append(criteria, CriterionResult{
			CriterionID: criterionID,
			Average:     average,
			Responses:   count,
		})
	}

	// Les maps Go n'ont pas d'ordre garanti.
	// On trie donc les critères par ID.
	sort.Slice(criteria, func(i, j int) bool {
		return criteria[i].CriterionID < criteria[j].CriterionID
	})

	globalAverage := float64(globalTotal) /
		float64(globalCount)

	return &ProfessorResult{
		ProfessorID:   professorID,
		CourseID:      courseID,
		AcademicYear:  academicYear,
		Period:        period,
		TotalReviews:  len(evaluationList),
		GlobalAverage: globalAverage,
		Criteria:      criteria,
	}, nil
}
