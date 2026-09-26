package eligibility

import (
	"errors"
)

type EligibilityRepository interface {
	Create(eligibility *Eligibility) error
	GetByStudentID(studentID int) (*Eligibility, error)
	Update(eligibility *Eligibility) error
	Delete(studentID int) error
}

type Service struct {
	repository EligibilityRepository
}

func NewService(repository EligibilityRepository) *Service {
	return &Service{
		repository: repository,
	}
}

// Create crée une nouvelle éligibilité.
func (s *Service) Create(eligibility *Eligibility) (*Eligibility, error) {

	if eligibility == nil {
		return nil, errors.New("eligibility is required")
	}

	if eligibility.StudentID <= 0 {
		return nil, errors.New("invalid student ID")
	}

	// Vérifier si une éligibilité existe déjà.
	existing, err := s.repository.GetByStudentID(eligibility.StudentID)

	if err == nil && existing != nil {
		return nil, errors.New(
			"eligibility already exists for this student",
		)
	}

	if err != nil && !errors.Is(err, ErrEligibilityNotFound) {
		return nil, err
	}

	// Calcul automatique.
	s.calculateEligibility(eligibility)

	// Enregistrement.
	if err := s.repository.Create(eligibility); err != nil {
		return nil, err
	}

	return eligibility, nil
}

// Check vérifie l'éligibilité d'un étudiant.
func (s *Service) Check(studentID int) (*Eligibility, error) {

	if studentID <= 0 {
		return nil, errors.New("invalid student ID")
	}

	eligibility, err := s.repository.GetByStudentID(studentID)

	if err != nil {
		return nil, err
	}

	// Recalculer pour garantir la cohérence.
	s.calculateEligibility(eligibility)

	return eligibility, nil
}

// Update modifie une éligibilité existante.
func (s *Service) Update(eligibility *Eligibility) (*Eligibility, error) {

	if eligibility == nil {
		return nil, errors.New("eligibility is required")
	}

	if eligibility.StudentID <= 0 {
		return nil, errors.New("invalid student ID")
	}

	// Vérifier que l'éligibilité existe.
	existing, err := s.repository.GetByStudentID(
		eligibility.StudentID,
	)

	if err != nil {
		return nil, err
	}

	// Conserver l'ID existant.
	eligibility.ID = existing.ID

	// Recalcul automatique.
	s.calculateEligibility(eligibility)

	// Mise à jour.
	if err := s.repository.Update(eligibility); err != nil {
		return nil, err
	}

	return eligibility, nil
}

// Delete supprime l'éligibilité d'un étudiant.
func (s *Service) Delete(studentID int) error {

	if studentID <= 0 {
		return errors.New("invalid student ID")
	}

	return s.repository.Delete(studentID)
}

// calculateEligibility applique les règles d'éligibilité.
func (s *Service) calculateEligibility(
	eligibility *Eligibility,
) {

	var reasons []string

	if !eligibility.Enrollment {
		reasons = append(
			reasons,
			"ENROLLMENT_INVALID",
		)
	}

	if !eligibility.AcademicFees {
		reasons = append(
			reasons,
			"ACADEMIC_FEES_NOT_PAID",
		)
	}

	if !eligibility.LaboratoryFees {
		reasons = append(
			reasons,
			"LABORATORY_FEES_NOT_PAID",
		)
	}

	if !eligibility.AccessFees {
		reasons = append(
			reasons,
			"ACCESS_FEES_NOT_PAID",
		)
	}

	eligibility.Reasons = reasons

	eligibility.Eligible = len(reasons) == 0
}