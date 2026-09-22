package criteria

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

type fakeCriterionRepository struct {
	getAllFunc    func() ([]Criterion, error)
	getActiveFunc func() ([]Criterion, error)
	getByIDFunc   func(id int) (*Criterion, error)
	createFunc    func(criterion Criterion) (*Criterion, error)
	updateFunc    func(id int, criterion Criterion) (*Criterion, error)
	deleteFunc    func(id int) error
}

func (f *fakeCriterionRepository) GetAll() ([]Criterion, error) {
	if f.getAllFunc != nil {
		return f.getAllFunc()
	}

	return nil, nil
}

func (f *fakeCriterionRepository) GetActive() ([]Criterion, error) {
	if f.getActiveFunc != nil {
		return f.getActiveFunc()
	}

	return nil, nil
}

func (f *fakeCriterionRepository) GetByID(id int) (*Criterion, error) {
	if f.getByIDFunc != nil {
		return f.getByIDFunc(id)
	}

	return nil, nil
}

func (f *fakeCriterionRepository) Create(
	criterion Criterion,
) (*Criterion, error) {
	if f.createFunc != nil {
		return f.createFunc(criterion)
	}

	return nil, nil
}

func (f *fakeCriterionRepository) Update(
	id int,
	criterion Criterion,
) (*Criterion, error) {
	if f.updateFunc != nil {
		return f.updateFunc(id, criterion)
	}

	return nil, nil
}

func (f *fakeCriterionRepository) Delete(id int) error {
	if f.deleteFunc != nil {
		return f.deleteFunc(id)
	}

	return nil
}

// --------------------------------------------------
// GET ALL
// --------------------------------------------------

func TestService_GetAll_Success(t *testing.T) {
	repository := &fakeCriterionRepository{
		getAllFunc: func() ([]Criterion, error) {
			return []Criterion{
				{
					ID:          1,
					Name:        "Clarté",
					Description: "Clarté des explications",
					MaxScore:    5,
					Active:      true,
				},
				{
					ID:          2,
					Name:        "Ponctualité",
					Description: "Respect des horaires",
					MaxScore:    5,
					Active:      true,
				},
			}, nil
		},
	}

	service := NewService(repository)

	result, err := service.GetAll()

	assert.NoError(t, err)
	assert.Len(t, result, 2)
	assert.Equal(t, "Clarté", result[0].Name)
	assert.Equal(t, "Ponctualité", result[1].Name)
}

// --------------------------------------------------
// GET ACTIVE
// --------------------------------------------------

func TestService_GetActive_Success(t *testing.T) {
	repository := &fakeCriterionRepository{
		getActiveFunc: func() ([]Criterion, error) {
			return []Criterion{
				{
					ID:       1,
					Name:     "Clarté",
					MaxScore: 5,
					Active:   true,
				},
			}, nil
		},
	}

	service := NewService(repository)

	result, err := service.GetActive()

	assert.NoError(t, err)
	assert.Len(t, result, 1)
	assert.True(t, result[0].Active)
}

func TestService_GetActive_Error(t *testing.T) {
	repository := &fakeCriterionRepository{
		getActiveFunc: func() ([]Criterion, error) {
			return nil, errors.New("database error")
		},
	}

	service := NewService(repository)

	result, err := service.GetActive()

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.EqualError(t, err, "database error")
}

// --------------------------------------------------
// GET BY ID
// --------------------------------------------------

func TestService_GetByID_InvalidID(t *testing.T) {
	repository := &fakeCriterionRepository{}

	service := NewService(repository)

	result, err := service.GetByID(0)

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.EqualError(t, err, "invalid criterion ID")
}

func TestService_GetByID_Success(t *testing.T) {
	repository := &fakeCriterionRepository{
		getByIDFunc: func(id int) (*Criterion, error) {
			assert.Equal(t, 1, id)

			return &Criterion{
				ID:       1,
				Name:     "Clarté",
				MaxScore: 5,
				Active:   true,
			}, nil
		},
	}

	service := NewService(repository)

	result, err := service.GetByID(1)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, 1, result.ID)
	assert.Equal(t, "Clarté", result.Name)
}

func TestService_GetByID_NotFound(t *testing.T) {
	repository := &fakeCriterionRepository{
		getByIDFunc: func(id int) (*Criterion, error) {
			return nil, ErrCriterionNotFound
		},
	}

	service := NewService(repository)

	result, err := service.GetByID(999)

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.ErrorIs(t, err, ErrCriterionNotFound)
}

// --------------------------------------------------
// CREATE
// --------------------------------------------------

func TestService_Create_Success(t *testing.T) {
	repository := &fakeCriterionRepository{
		createFunc: func(criterion Criterion) (*Criterion, error) {
			assert.Equal(t, "Clarté des explications", criterion.Name)
			assert.Equal(t, "Évalue la clarté du professeur", criterion.Description)
			assert.Equal(t, 5, criterion.MaxScore)
			assert.True(t, criterion.Active)

			criterion.ID = 1

			return &criterion, nil
		},
	}

	service := NewService(repository)

	request := CreateCriterionRequest{
		Name:        "  Clarté des explications  ",
		Description: "  Évalue la clarté du professeur  ",
		MaxScore:    5,
	}

	result, err := service.Create(request)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, 1, result.ID)
	assert.Equal(t, "Clarté des explications", result.Name)
	assert.Equal(t, "Évalue la clarté du professeur", result.Description)
	assert.Equal(t, 5, result.MaxScore)
	assert.True(t, result.Active)
}

func TestService_Create_EmptyName(t *testing.T) {
	repository := &fakeCriterionRepository{}

	service := NewService(repository)

	request := CreateCriterionRequest{
		Name:     "",
		MaxScore: 5,
	}

	result, err := service.Create(request)

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.EqualError(t, err, "criterion name is required")
}

func TestService_Create_OnlySpacesName(t *testing.T) {
	repository := &fakeCriterionRepository{}

	service := NewService(repository)

	request := CreateCriterionRequest{
		Name:     "   ",
		MaxScore: 5,
	}

	result, err := service.Create(request)

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.EqualError(t, err, "criterion name is required")
}

func TestService_Create_InvalidMaxScore(t *testing.T) {
	repository := &fakeCriterionRepository{}

	service := NewService(repository)

	request := CreateCriterionRequest{
		Name:     "Clarté",
		MaxScore: 0,
	}

	result, err := service.Create(request)

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.EqualError(t, err, "max score must be greater than zero")
}

func TestService_Create_RepositoryError(t *testing.T) {
	repository := &fakeCriterionRepository{
		createFunc: func(criterion Criterion) (*Criterion, error) {
			return nil, errors.New("database error")
		},
	}

	service := NewService(repository)

	request := CreateCriterionRequest{
		Name:     "Clarté",
		MaxScore: 5,
	}

	result, err := service.Create(request)

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.EqualError(t, err, "database error")
}

// --------------------------------------------------
// UPDATE
// --------------------------------------------------

func TestService_Update_InvalidID(t *testing.T) {
	repository := &fakeCriterionRepository{}

	service := NewService(repository)

	request := UpdateCriterionRequest{
		Name:     "Clarté",
		MaxScore: 5,
	}

	result, err := service.Update(0, request)

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.EqualError(t, err, "invalid criterion ID")
}

func TestService_Update_Success(t *testing.T) {
	repository := &fakeCriterionRepository{
		updateFunc: func(id int, criterion Criterion) (*Criterion, error) {
			assert.Equal(t, 1, id)
			assert.Equal(t, "Clarté avancée", criterion.Name)
			assert.Equal(t, "Nouvelle description", criterion.Description)
			assert.Equal(t, 10, criterion.MaxScore)
			assert.False(t, criterion.Active)

			criterion.ID = id

			return &criterion, nil
		},
	}

	service := NewService(repository)

	request := UpdateCriterionRequest{
		Name:        "  Clarté avancée  ",
		Description: "  Nouvelle description  ",
		MaxScore:    10,
		Active:      false,
	}

	result, err := service.Update(1, request)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, 1, result.ID)
	assert.Equal(t, "Clarté avancée", result.Name)
	assert.Equal(t, 10, result.MaxScore)
	assert.False(t, result.Active)
}

func TestService_Update_EmptyName(t *testing.T) {
	repository := &fakeCriterionRepository{}

	service := NewService(repository)

	request := UpdateCriterionRequest{
		Name:     "",
		MaxScore: 5,
	}

	result, err := service.Update(1, request)

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.EqualError(t, err, "criterion name is required")
}

func TestService_Update_OnlySpacesName(t *testing.T) {
	repository := &fakeCriterionRepository{}

	service := NewService(repository)

	request := UpdateCriterionRequest{
		Name:     "   ",
		MaxScore: 5,
	}

	result, err := service.Update(1, request)

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.EqualError(t, err, "criterion name is required")
}

func TestService_Update_InvalidMaxScore(t *testing.T) {
	repository := &fakeCriterionRepository{}

	service := NewService(repository)

	request := UpdateCriterionRequest{
		Name:     "Clarté",
		MaxScore: 0,
	}

	result, err := service.Update(1, request)

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.EqualError(t, err, "max score must be greater than zero")
}

func TestService_Update_NotFound(t *testing.T) {
	repository := &fakeCriterionRepository{
		updateFunc: func(
			id int,
			criterion Criterion,
		) (*Criterion, error) {
			return nil, ErrCriterionNotFound
		},
	}

	service := NewService(repository)

	request := UpdateCriterionRequest{
		Name:     "Clarté",
		MaxScore: 5,
	}

	result, err := service.Update(999, request)

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.ErrorIs(t, err, ErrCriterionNotFound)
}

// --------------------------------------------------
// DELETE
// --------------------------------------------------

func TestService_Delete_InvalidID(t *testing.T) {
	repository := &fakeCriterionRepository{}

	service := NewService(repository)

	err := service.Delete(0)

	assert.Error(t, err)
	assert.EqualError(t, err, "invalid criterion ID")
}

func TestService_Delete_Success(t *testing.T) {
	repository := &fakeCriterionRepository{
		deleteFunc: func(id int) error {
			assert.Equal(t, 1, id)
			return nil
		},
	}

	service := NewService(repository)

	err := service.Delete(1)

	assert.NoError(t, err)
}

func TestService_Delete_NotFound(t *testing.T) {
	repository := &fakeCriterionRepository{
		deleteFunc: func(id int) error {
			return ErrCriterionNotFound
		},
	}

	service := NewService(repository)

	err := service.Delete(999)

	assert.Error(t, err)
	assert.ErrorIs(t, err, ErrCriterionNotFound)
}

func TestService_Delete_RepositoryError(t *testing.T) {
	repository := &fakeCriterionRepository{
		deleteFunc: func(id int) error {
			return errors.New("database error")
		},
	}

	service := NewService(repository)

	err := service.Delete(1)

	assert.Error(t, err)
	assert.EqualError(t, err, "database error")
}
