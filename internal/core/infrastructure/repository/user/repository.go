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
)

type repository struct {
	UserServiceBaseRoute string
}

func NewRepository(config *config.Config) user_repository_contract.Repository {
	return &repository{
		UserServiceBaseRoute: config.ExternalServiceConfig.UserServiceBaseRoute,
	}
}

func (rp *repository) GetByDoctorID(ctx context.Context, doctorID uint) (*entities.Doctor, error) {
	route := fmt.Sprintf("%s/users/doctors/%d", rp.UserServiceBaseRoute, doctorID)
	req, err := http.NewRequest(http.MethodGet, route, nil)
	if err != nil {
		return nil, err
	}

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

func (rp *repository) UpdateRating(ctx context.Context, doctorID uint, rating float64) error {
	body := `{
		"rating": %f,
	}`

	route := fmt.Sprintf("%s/users/doctors/%d/ratings", rp.UserServiceBaseRoute, doctorID)
	req, err := http.NewRequest(http.MethodPost, route, bytes.NewBuffer([]byte(fmt.Sprintf(body, rating))))
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	return nil
}
