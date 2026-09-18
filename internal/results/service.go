package results

import (
	"errors"
	"sort"
)

type Service struct {
	repository *Repository
}

func NewService(repository *Repository) *Service {
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
		return nil, errors.New("invalid professor ID")
	}

	if courseID <= 0 {
		return nil, errors.New("invalid course ID")
	}

	if academicYear == "" {
		return nil, errors.New("academic year is required")
	}

	if period == "" {
		return nil, errors.New("period is required")
	}

	evaluations, err := s.repository.GetByProfessor(
		professorID,
		courseID,
		academicYear,
		period,
	)

	if err != nil {
		return nil, err
	}

	if len(evaluations) == 0 {
		return nil, errors.New("no evaluations found")
	}

	// criterionID -> total score
	totals := make(map[int]int)

	// criterionID -> number of responses
	counts := make(map[int]int)

	globalTotal := 0
	globalCount := 0

	for _, evaluation := range evaluations {

		for _, answer := range evaluation.Answers {

			totals[answer.CriterionID] += answer.Score
			counts[answer.CriterionID]++

			globalTotal += answer.Score
			globalCount++
		}
	}

	if globalCount == 0 {
		return nil, errors.New("evaluations contain no answers")
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

	// Les maps Go ne garantissent pas l'ordre.
	// On trie les critères par ID pour obtenir
	// une réponse API stable.
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
		TotalReviews:  len(evaluations),
		GlobalAverage: globalAverage,
		Criteria:      criteria,
	}, nil
}