package dolores

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/psankar/vetchi/typespec/common"
	"github.com/psankar/vetchi/typespec/employer"
	"github.com/psankar/vetchi/typespec/hub"
)

var _ = FDescribe("Interviews", Ordered, func() {
	var db *pgxpool.Pool
	var adminToken, recruiterToken, hiringManagerToken string
	var interviewer1Token, interviewer2Token, interviewer3Token string
	var hubUserToken string

	BeforeAll(func() {
		db = setupTestDB()
		seedDatabase(db, "0014-interviews-up.pgsql")

		var wg sync.WaitGroup
		tokens := map[string]*string{
			"admin@0014-interview.example":          &adminToken,
			"recruiter@0014-interview.example":      &recruiterToken,
			"hiring-manager@0014-interview.example": &hiringManagerToken,
			"interviewer1@0014-interview.example":   &interviewer1Token,
			"interviewer2@0014-interview.example":   &interviewer2Token,
			"interviewer3@0014-interview.example":   &interviewer3Token,
		}

		for email, token := range tokens {
			wg.Add(1)
			employerSigninAsync(
				"0014-interview.example",
				email,
				"NewPassword123$",
				token,
				&wg,
			)
		}
		wg.Wait()

		// Sign in hub user
		hubUserToken = hubSignin("interview@hub.example", "NewPassword123$")
	})

	AfterAll(func() {
		seedDatabase(db, "0014-interviews-down.pgsql")
		db.Close()
	})

	Context("Interview Management", func() {
		It("should handle full interview workflow", func() {
			// 1. Add Interview
			addInterviewReq := employer.AddInterviewRequest{
				CandidacyID:   "candidacy-001",
				StartTime:     time.Now().Add(24 * time.Hour),
				EndTime:       time.Now().Add(25 * time.Hour),
				InterviewType: employer.InPersonInterviewType,
				Description:   "Technical Interview Round",
				Interviewers: []common.EmailAddress{
					"interviewer1@0014-interview.example",
					"interviewer2@0014-interview.example",
				},
			}

			reqBody, err := json.Marshal(addInterviewReq)
			Expect(err).ShouldNot(HaveOccurred())

			req, err := http.NewRequest(
				"POST",
				serverURL+"/employer/add-interview",
				bytes.NewBuffer(reqBody),
			)
			Expect(err).ShouldNot(HaveOccurred())
			req.Header.Set("Authorization", "Bearer "+recruiterToken)

			resp, err := http.DefaultClient.Do(req)
			Expect(err).ShouldNot(HaveOccurred())
			Expect(resp.StatusCode).Should(Equal(http.StatusOK))

			// 2. Add another interviewer
			addInterviewerReq := employer.AddInterviewerRequest{
				InterviewID:  "interview-001",
				OrgUserEmail: "interviewer3@0014-interview.example",
			}

			reqBody, err = json.Marshal(addInterviewerReq)
			Expect(err).ShouldNot(HaveOccurred())

			req, err = http.NewRequest(
				"POST",
				serverURL+"/employer/add-interviewer",
				bytes.NewBuffer(reqBody),
			)
			Expect(err).ShouldNot(HaveOccurred())
			req.Header.Set("Authorization", "Bearer "+recruiterToken)

			resp, err = http.DefaultClient.Do(req)
			Expect(err).ShouldNot(HaveOccurred())
			Expect(resp.StatusCode).Should(Equal(http.StatusOK))

			// 3. Remove an interviewer
			removeInterviewerReq := employer.RemoveInterviewerRequest{
				InterviewID:  "interview-001",
				OrgUserEmail: "interviewer2@0014-interview.example",
			}

			reqBody, err = json.Marshal(removeInterviewerReq)
			Expect(err).ShouldNot(HaveOccurred())

			req, err = http.NewRequest(
				"POST",
				serverURL+"/employer/remove-interviewer",
				bytes.NewBuffer(reqBody),
			)
			Expect(err).ShouldNot(HaveOccurred())
			req.Header.Set("Authorization", "Bearer "+recruiterToken)

			resp, err = http.DefaultClient.Do(req)
			Expect(err).ShouldNot(HaveOccurred())
			Expect(resp.StatusCode).Should(Equal(http.StatusOK))

			// 4. Get interviews by opening
			getInterviewsByOpeningReq := employer.GetInterviewsByOpeningRequest{
				OpeningID: "2024-Mar-15-001",
				States: []common.InterviewState{
					common.ScheduledInterviewState,
				},
				Limit: 10,
			}

			reqBody, err = json.Marshal(getInterviewsByOpeningReq)
			Expect(err).ShouldNot(HaveOccurred())

			req, err = http.NewRequest(
				"POST",
				serverURL+"/employer/get-interviews-by-opening",
				bytes.NewBuffer(reqBody),
			)
			Expect(err).ShouldNot(HaveOccurred())
			req.Header.Set("Authorization", "Bearer "+recruiterToken)

			resp, err = http.DefaultClient.Do(req)
			Expect(err).ShouldNot(HaveOccurred())
			Expect(resp.StatusCode).Should(Equal(http.StatusOK))

			// 5. Get interviews by candidacy
			getInterviewsByCandidacyReq := employer.GetInterviewsByCandidacyRequest{
				CandidacyID: "candidacy-001",
				States: []common.InterviewState{
					common.ScheduledInterviewState,
				},
			}

			reqBody, err = json.Marshal(getInterviewsByCandidacyReq)
			Expect(err).ShouldNot(HaveOccurred())

			req, err = http.NewRequest(
				"POST",
				serverURL+"/employer/get-interviews-by-candidacy",
				bytes.NewBuffer(reqBody),
			)
			Expect(err).ShouldNot(HaveOccurred())
			req.Header.Set("Authorization", "Bearer "+recruiterToken)

			resp, err = http.DefaultClient.Do(req)
			Expect(err).ShouldNot(HaveOccurred())
			Expect(resp.StatusCode).Should(Equal(http.StatusOK))

			// 6. Submit assessment
			assessment := employer.Assessment{
				InterviewID:         "interview-001",
				Decision:            common.StrongYesInterviewersDecision,
				Positives:           "Strong technical skills",
				Negatives:           "Could improve communication",
				OverallAssessment:   "Good candidate, recommended for hire",
				FeedbackToCandidate: "Thank you for interviewing with us",
			}

			reqBody, err = json.Marshal(assessment)
			Expect(err).ShouldNot(HaveOccurred())

			req, err = http.NewRequest(
				"POST",
				serverURL+"/employer/put-assessment",
				bytes.NewBuffer(reqBody),
			)
			Expect(err).ShouldNot(HaveOccurred())
			req.Header.Set("Authorization", "Bearer "+interviewer1Token)

			resp, err = http.DefaultClient.Do(req)
			Expect(err).ShouldNot(HaveOccurred())
			Expect(resp.StatusCode).Should(Equal(http.StatusOK))

			// 7. Get assessment
			getAssessmentReq := employer.GetAssessmentRequest{
				InterviewID: "interview-001",
			}

			reqBody, err = json.Marshal(getAssessmentReq)
			Expect(err).ShouldNot(HaveOccurred())

			req, err = http.NewRequest(
				"POST",
				serverURL+"/employer/get-assessment",
				bytes.NewBuffer(reqBody),
			)
			Expect(err).ShouldNot(HaveOccurred())
			req.Header.Set("Authorization", "Bearer "+interviewer1Token)

			resp, err = http.DefaultClient.Do(req)
			Expect(err).ShouldNot(HaveOccurred())
			Expect(resp.StatusCode).Should(Equal(http.StatusOK))

			var gotAssessment employer.Assessment
			err = json.NewDecoder(resp.Body).Decode(&gotAssessment)
			Expect(err).ShouldNot(HaveOccurred())
			Expect(gotAssessment.Decision).Should(Equal(assessment.Decision))
			Expect(gotAssessment.Positives).Should(Equal(assessment.Positives))
			Expect(gotAssessment.Negatives).Should(Equal(assessment.Negatives))

			// 8. Hub user RSVP
			rsvpReq := hub.HubRSVPInterviewRequest{
				InterviewID: "interview-001",
				RSVP:        common.YesRSVP,
			}

			reqBody, err = json.Marshal(rsvpReq)
			Expect(err).ShouldNot(HaveOccurred())

			req, err = http.NewRequest(
				"POST",
				serverURL+"/hub/rsvp-interview",
				bytes.NewBuffer(reqBody),
			)
			Expect(err).ShouldNot(HaveOccurred())
			req.Header.Set("Authorization", "Bearer "+hubUserToken)

			resp, err = http.DefaultClient.Do(req)
			Expect(err).ShouldNot(HaveOccurred())
			Expect(resp.StatusCode).Should(Equal(http.StatusOK))
		})

		It("should handle interview state transitions", func() {
			// Create a new interview
			addInterviewReq := employer.AddInterviewRequest{
				CandidacyID:   "candidacy-001",
				StartTime:     time.Now().Add(1 * time.Hour),
				EndTime:       time.Now().Add(2 * time.Hour),
				InterviewType: "TECHNICAL",
				Description:   "Technical Interview Round",
				Interviewers: []common.EmailAddress{
					"interviewer1@0014-interview.example",
				},
			}

			reqBody, err := json.Marshal(addInterviewReq)
			Expect(err).ShouldNot(HaveOccurred())

			req, err := http.NewRequest(
				"POST",
				serverURL+"/employer/add-interview",
				bytes.NewBuffer(reqBody),
			)
			Expect(err).ShouldNot(HaveOccurred())
			req.Header.Set("Authorization", "Bearer "+recruiterToken)

			resp, err := http.DefaultClient.Do(req)
			Expect(err).ShouldNot(HaveOccurred())
			Expect(resp.StatusCode).Should(Equal(http.StatusOK))

			// Try to submit assessment before interview completion (should fail)
			assessment := employer.Assessment{
				InterviewID:         "interview-002", // New interview ID
				Decision:            common.StrongYesInterviewersDecision,
				Positives:           "Strong technical skills",
				Negatives:           "Could improve communication",
				OverallAssessment:   "Good candidate, recommended for hire",
				FeedbackToCandidate: "Thank you for interviewing with us",
			}

			reqBody, err = json.Marshal(assessment)
			Expect(err).ShouldNot(HaveOccurred())

			req, err = http.NewRequest(
				"POST",
				serverURL+"/employer/put-assessment",
				bytes.NewBuffer(reqBody),
			)
			Expect(err).ShouldNot(HaveOccurred())
			req.Header.Set("Authorization", "Bearer "+interviewer1Token)

			resp, err = http.DefaultClient.Do(req)
			Expect(err).ShouldNot(HaveOccurred())
			Expect(
				resp.StatusCode,
			).Should(Equal(http.StatusUnprocessableEntity))
		})

		It("should handle pagination in get-interviews endpoints", func() {
			// Create multiple interviews
			for i := 0; i < 5; i++ {
				addInterviewReq := employer.AddInterviewRequest{
					CandidacyID: "candidacy-001",
					StartTime: time.Now().
						Add(time.Duration(24+i) * time.Hour),
					EndTime: time.Now().
						Add(time.Duration(25+i) * time.Hour),
					InterviewType: "TECHNICAL",
					Description: fmt.Sprintf(
						"Technical Interview Round %d",
						i+1,
					),
					Interviewers: []common.EmailAddress{
						"interviewer1@0014-interview.example",
					},
				}

				reqBody, err := json.Marshal(addInterviewReq)
				Expect(err).ShouldNot(HaveOccurred())

				req, err := http.NewRequest(
					"POST",
					serverURL+"/employer/add-interview",
					bytes.NewBuffer(reqBody),
				)
				Expect(err).ShouldNot(HaveOccurred())
				req.Header.Set("Authorization", "Bearer "+recruiterToken)

				resp, err := http.DefaultClient.Do(req)
				Expect(err).ShouldNot(HaveOccurred())
				Expect(resp.StatusCode).Should(Equal(http.StatusOK))
			}

			// Test pagination with limit
			getInterviewsReq := employer.GetInterviewsByOpeningRequest{
				OpeningID: "2024-Mar-15-001",
				States: []common.InterviewState{
					common.ScheduledInterviewState,
				},
				Limit: 2,
			}

			reqBody, err := json.Marshal(getInterviewsReq)
			Expect(err).ShouldNot(HaveOccurred())

			req, err := http.NewRequest(
				"POST",
				serverURL+"/employer/get-interviews-by-opening",
				bytes.NewBuffer(reqBody),
			)
			Expect(err).ShouldNot(HaveOccurred())
			req.Header.Set("Authorization", "Bearer "+recruiterToken)

			resp, err := http.DefaultClient.Do(req)
			Expect(err).ShouldNot(HaveOccurred())
			Expect(resp.StatusCode).Should(Equal(http.StatusOK))

			var interviews []employer.Interview
			err = json.NewDecoder(resp.Body).Decode(&interviews)
			Expect(err).ShouldNot(HaveOccurred())
			Expect(len(interviews)).Should(Equal(2))
			Expect(interviews[0].InterviewID).ShouldNot(BeEmpty())
			Expect(interviews[1].InterviewID).ShouldNot(BeEmpty())
		})

		It("should validate interview types", func() {
			invalidTypeReq := employer.AddInterviewRequest{
				CandidacyID:   "candidacy-001",
				StartTime:     time.Now().Add(24 * time.Hour),
				EndTime:       time.Now().Add(25 * time.Hour),
				InterviewType: "INVALID_TYPE",
				Description:   "Invalid Interview Type",
				Interviewers: []common.EmailAddress{
					"interviewer1@0014-interview.example",
				},
			}

			reqBody, err := json.Marshal(invalidTypeReq)
			Expect(err).ShouldNot(HaveOccurred())

			req, err := http.NewRequest(
				"POST",
				serverURL+"/employer/add-interview",
				bytes.NewBuffer(reqBody),
			)
			Expect(err).ShouldNot(HaveOccurred())
			req.Header.Set("Authorization", "Bearer "+recruiterToken)

			resp, err := http.DefaultClient.Do(req)
			Expect(err).ShouldNot(HaveOccurred())
			Expect(resp.StatusCode).Should(Equal(http.StatusBadRequest))
		})

		// Add error case tests here for each endpoint
		// For example:
		// - Invalid interview times
		// - Non-existent interviewers
		// - Invalid states
		// - Unauthorized access
		// - etc.
	})
})
