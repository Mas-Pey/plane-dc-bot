package services

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
)

type PlaneService struct {
	BaseURL       string
	ApiKey        string
	WorkspaceSlug string
}

type StatePayload struct {
	StateID string `json:"state"`
}

type IssueResponse struct {
	ID        string `json:"id"`
	ProjectID string `json:"project"`
}

func NewPlaneService() *PlaneService {
	return &PlaneService{
		BaseURL:       "https://plane-hq.runsystemdev.com/api/v1",
		ApiKey:        os.Getenv("PLANE_TOKEN"),
		WorkspaceSlug: "erp-intern",
	}
}

func (s *PlaneService) UpdateStateIssue(identifier string, stateUUID string) error {
	// Find uuid by Identifier ID
	log.Println("[INFO] find id issue by identifier: ", identifier)
	data, err := s.GetIssueData(identifier)
	if err != nil {
		return fmt.Errorf("failed to retrieve issue by identifier: %w", err)
	}

	log.Printf("[INFO] Find ID Succesful! UUID: %s | Project: %s", data.ID, data.ProjectID)

	url := fmt.Sprintf("%s/workspaces/%s/projects/%s/issues/%s/", s.BaseURL, s.WorkspaceSlug, data.ProjectID, data.ID)
	payload := StatePayload{
		StateID: stateUUID,
	}
	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("failed to wrap json data : %w", err)
	}

	// Request PATCH to API Plane
	req, err := http.NewRequest("PATCH", url, bytes.NewBuffer(jsonPayload))
	if err != nil {
		return fmt.Errorf("failed to make request (PATCH): %w", err)
	}
	req.Header.Set("x-api-key", s.ApiKey)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to connect in Server Plane: %w", err)
	}
	defer res.Body.Close()

	bodyBytes, _ := io.ReadAll(res.Body)

	if res.StatusCode < 200 || res.StatusCode >= 300 {
		log.Println("[ERROR] Rejected by Plane API: ", string(bodyBytes))
		return fmt.Errorf("API error code %d: %s", res.StatusCode, string(bodyBytes))
	}

	log.Printf("[SUCCESS] State Issue [%s], (UUID: %s) updated succesfully, new state: %s", identifier, data.ID, stateUUID)

	return nil
}

// helper function for get UUID by sequence_id
func (s *PlaneService) GetIssueData(identifier string) (*IssueResponse, error) {
	url := fmt.Sprintf("%s/workspaces/%s/issues/%s/", s.BaseURL, s.WorkspaceSlug, identifier)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to make request (GET): %w", err)
	}

	req.Header.Set("x-api-key", s.ApiKey)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to connect in Server Plane: %w", err)
	}
	defer res.Body.Close()

	if res.StatusCode != 200 {
		return nil, fmt.Errorf("issue %s not found (status: %d)", identifier, res.StatusCode)
	}

	var result IssueResponse
	if err := json.NewDecoder(res.Body).Decode(&result); err != nil {
		return nil, err
	}

	return &result, nil
}
