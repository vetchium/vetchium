package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"time"

	"github.com/fatih/color"
	"github.com/psankar/vetchi/typespec/hub"
)

// HistoryItem represents a work history item with employer and time info
type HistoryItem struct {
	EmployerDomain string
	Title          string
	StartDate      time.Time
	EndDate        *time.Time
}

// WorkHistoryMap stores work histories organized by hub user email
var WorkHistoryMap = make(map[string][]HistoryItem)

func createWorkHistories() {
	for _, user := range hubUsers {
		tokenI, ok := hubSessionTokens.Load(user.Email)
		if !ok {
			log.Fatalf("no auth token found for %s", user.Email)
		}
		authToken := tokenI.(string)

		workHistoryItems, err := createWorkHistory(authToken, user.Jobs)
		if err != nil {
			log.Fatalf(
				"error creating work history for %s: %v",
				user.Email,
				err,
			)
		}

		// Store work history in the map
		WorkHistoryMap[user.Email] = workHistoryItems

		color.Magenta("created work history for %s", user.Email)
	}
}

func createWorkHistory(
	authToken string,
	jobs []Job,
) ([]HistoryItem, error) {
	var prevStartDate time.Time
	var workHistoryItems []HistoryItem

	for i := len(jobs) - 1; i >= 0; i-- {
		job := jobs[i]

		var startDateRaw time.Time
		var startDate string
		var endDatePtr *string
		var endDateRaw *time.Time

		if i == len(jobs)-1 {
			// Last job is current job
			endDatePtr = nil
			randYears := rand.Intn(7) + 1
			startDateRaw = time.Now().AddDate(-randYears, 0, 0)
			startDate = startDateRaw.Format("2006-01-02")
			prevStartDate = startDateRaw
		} else {
			// Assuming that a 30 day gap exists between jobs
			gapDays := rand.Intn(30) + 30
			endDate := prevStartDate.AddDate(0, 0, -gapDays)
			endDateStr := endDate.Format("2006-01-02")
			endDatePtr = &endDateStr
			endDateRaw = &endDate

			numberOfYears := rand.Intn(7) + 1
			startDateRaw = endDate.AddDate(-numberOfYears, 0, 0)
			startDate = startDateRaw.Format("2006-01-02")

			prevStartDate = startDateRaw
		}

		err := createWorkHistoryItem(
			authToken,
			startDate,
			endDatePtr,
			job,
		)
		if err != nil {
			return nil, err
		}

		// Add to our local tracking
		workHistoryItems = append(workHistoryItems, HistoryItem{
			EmployerDomain: job.Website,
			Title:          job.Title,
			StartDate:      startDateRaw,
			EndDate:        endDateRaw,
		})
	}
	return workHistoryItems, nil
}

func createWorkHistoryItem(
	authToken string,
	startDate string,
	endDatePtr *string,
	job Job,
) error {
	var addWorkHistoryRequest = hub.AddWorkHistoryRequest{
		EmployerDomain: job.Website,
		Title:          job.Title,
		StartDate:      startDate,
		EndDate:        endDatePtr,
		Description:    &job.Title,
	}

	body, err := json.Marshal(addWorkHistoryRequest)
	if err != nil {
		return err
	}

	req, err := http.NewRequest(
		http.MethodPost,
		serverURL+"/hub/add-work-history",
		bytes.NewBuffer(body),
	)
	if err != nil {
		return err
	}

	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", authToken))

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to create work history: %s", resp.Status)
	}

	return nil
}
