package professors

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewService(t *testing.T) {
	repo := &Repository{}
	service := NewService(repo)

	assert.NotNil(t, service)
	assert.NotNil(t, service.repository)
}

func TestService_Create_RequiredFields(t *testing.T) {
	repo := &Repository{}
	service := NewService(repo)

	tests := []struct {
		name    string
		request CreateProfessorRequest
		error   string
	}{
		{
			name: "missing matricule",
			request: CreateProfessorRequest{
				FirstName:  "John",
				LastName:   "Doe",
				Department: "CS",
			},
			error: "professor matricule is required",
		},
		{
			name: "missing first name",
			request: CreateProfessorRequest{
				Matricule:  "12345",
				LastName:   "Doe",
				Department: "CS",
			},
			error: "first name is required",
		},
		{
			name: "missing last name",
			request: CreateProfessorRequest{
				Matricule:  "12345",
				FirstName:  "John",
				Department: "CS",
			},
			error: "last name is required",
		},
		{
			name: "missing department",
			request: CreateProfessorRequest{
				Matricule: "12345",
				FirstName: "John",
				LastName:  "Doe",
			},
			error: "department is required",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := service.Create(tt.request)

			assert.Error(t, err)
			assert.Nil(t, result)
			assert.Equal(t, tt.error, err.Error())
		})
	}
}

func TestService_Update_InvalidID(t *testing.T) {
	repo := &Repository{}
	service := NewService(repo)

	request := UpdateProfessorRequest{
		Matricule:  "12345",
		FirstName:  "John",
		LastName:   "Doe",
		Department: "CS",
	}

	result, err := service.Update(0, request)

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Equal(t, "invalid professor ID", err.Error())
}

func TestService_Delete_InvalidID(t *testing.T) {
	repo := &Repository{}
	service := NewService(repo)

	err := service.Delete(0)

	assert.Error(t, err)
	assert.Equal(t, "invalid professor ID", err.Error())
}

func TestService_GetByID_InvalidID(t *testing.T) {
	repo := &Repository{}
	service := NewService(repo)

	result, err := service.GetByID(0)

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Equal(t, "invalid professor ID", err.Error())
}
