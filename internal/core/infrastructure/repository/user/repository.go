package user_repository

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/jfelipearaujo-healthmed/review-processor-service/internal/core/domain/entities"
	user_repository_contract "github.com/jfelipearaujo-healthmed/review-processor-service/internal/core/domain/repositories/user"
	"github.com/jfelipearaujo-healthmed/review-processor-service/internal/core/infrastructure/config"
	"github.com/jfelipearaujo-healthmed/review-processor-service/internal/core/infrastructure/shared/token"
)

type repository struct {
	userServiceBaseRoute string
	tokenService         token.TokenService
}

func NewRepository(config *config.Config, tokenService token.TokenService) user_repository_contract.Repository {
	return &repository{
		userServiceBaseRoute: config.ExternalServiceConfig.UserServiceBaseRoute,
		tokenService:         tokenService,
	}
}

func (rp *repository) GetByDoctorID(ctx context.Context, patientID, doctorID uint) (*entities.Doctor, error) {
	route := fmt.Sprintf("%s/users/doctors/%d", rp.userServiceBaseRoute, doctorID)
	req, err := http.NewRequest(http.MethodGet, route, nil)
	if err != nil {
		return nil, err
	}

	token, err := rp.tokenService.CreateJwtToken(patientID)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	var doctor entities.Doctor

	if err := json.NewDecoder(resp.Body).Decode(&doctor); err != nil {
		return nil, err
	}

	return &doctor, nil
}

func (rp *repository) UpdateRating(ctx context.Context, patientID, doctorID uint, rating float64) error {
	body := `{
		"rating": %f,
	}`

	route := fmt.Sprintf("%s/users/doctors/%d/ratings", rp.userServiceBaseRoute, doctorID)
	req, err := http.NewRequest(http.MethodPost, route, bytes.NewBuffer([]byte(fmt.Sprintf(body, rating))))
	if err != nil {
		return err
	}

	token, err := rp.tokenService.CreateJwtToken(patientID)
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	return nil
}
