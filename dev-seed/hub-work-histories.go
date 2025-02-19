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

var jobs = []struct {
	Title       string
	Description string
}{
	{
		Title:       "Account Manager",
		Description: "Managed key client relationships, developed strategic account plans, and ensured customer satisfaction while meeting revenue targets.",
	},
	{
		Title:       "Accountant",
		Description: "Maintained financial records, prepared reports and tax returns, analyzed business operations, and ensured regulatory compliance.",
	},
	{
		Title:       "Administrative Assistant",
		Description: "Provided comprehensive administrative support including scheduling, document management, and office coordination.",
	},
	{
		Title:       "Business Development Representative",
		Description: "Generated new business opportunities through prospecting, cold calling, and relationship building with potential clients.",
	},
	{
		Title:       "Business Operations Manager",
		Description: "Oversaw daily operations, optimized business processes, and implemented strategic initiatives to improve efficiency.",
	},
	{
		Title:       "Chief Executive Officer",
		Description: "Led overall company strategy, growth, and operations while building strong relationships with board members and stakeholders.",
	},
	{
		Title:       "Chief Financial Officer",
		Description: "Directed financial planning, risk management, and accounting practices while ensuring long-term financial health of the organization.",
	},
	{
		Title:       "Chief Marketing Officer",
		Description: "Developed and executed marketing strategies to drive growth, brand awareness, and market positioning.",
	},
	{
		Title:       "Chief Operating Officer",
		Description: "Managed day-to-day operations, optimized organizational processes, and implemented strategic business initiatives.",
	},
	{
		Title:       "Chief Technology Officer",
		Description: "Led technology strategy, innovation, and digital transformation while managing technical teams and infrastructure.",
	},
	{
		Title:       "Cloud Architect",
		Description: "Designed and implemented scalable cloud infrastructure solutions while ensuring security, reliability, and cost-effectiveness.",
	},
	{
		Title:       "Cloud Engineer",
		Description: "Built and maintained cloud-based systems, implemented automation, and optimized cloud resource utilization.",
	},
	{
		Title:       "Data Scientist",
		Description: "Analyzed complex data sets, developed predictive models, and provided data-driven insights to guide business decisions.",
	},
	{
		Title:       "Database Administrator",
		Description: "Managed database systems, ensured data integrity, and optimized database performance while maintaining security.",
	},
	{
		Title:       "DevOps Engineer",
		Description: "Implemented CI/CD pipelines, automated deployment processes, and maintained infrastructure as code.",
	},
	{
		Title:       "Network Engineer",
		Description: "Designed and maintained network infrastructure, implemented security measures, and resolved connectivity issues.",
	},
	{
		Title:       "Product Manager",
		Description: "Led product strategy, gathered requirements, and coordinated with engineering teams to deliver successful products.",
	},
	{
		Title:       "Security Engineer",
		Description: "Implemented security measures, conducted vulnerability assessments, and ensured compliance with security standards.",
	},
	{
		Title:       "Site Reliability Engineer",
		Description: "Ensured system reliability, implemented monitoring solutions, and automated operational tasks.",
	},
	{
		Title:       "Software Engineer",
		Description: "Developed and maintained software applications, collaborated with cross-functional teams, and implemented technical solutions.",
	},
	{
		Title:       "System Administrator",
		Description: "Managed IT infrastructure, maintained system security, and provided technical support to ensure smooth operations.",
	},
	{
		Title:       "System Engineer",
		Description: "Designed and implemented system architecture, optimized performance, and resolved complex technical issues.",
	},
	{
		Title:       "UI/UX Designer",
		Description: "Created user-centered designs, developed wireframes and prototypes, and improved user experience through iterative design.",
	},
	{
		Title:       "VP of Engineering",
		Description: "Led engineering teams, defined technical strategy, and ensured delivery of high-quality software products.",
	},
	{
		Title:       "VP of Product",
		Description: "Directed product strategy, roadmap development, and cross-functional execution of product initiatives.",
	},
	{
		Title:       "VP of Sales",
		Description: "Led sales strategy, managed sales teams, and drove revenue growth through new business development.",
	},
	{
		Title:       "VP of Marketing",
		Description: "Directed marketing initiatives, brand strategy, and demand generation programs to drive business growth.",
	},
	{
		Title:       "VP of Finance",
		Description: "Managed financial planning, reporting, and analysis while providing strategic financial guidance to leadership.",
	},
	{
		Title:       "VP of HR",
		Description: "Led talent acquisition, employee development, and HR operations while fostering a positive company culture.",
	},
	{
		Title:       "VP of Legal",
		Description: "Managed legal affairs, regulatory compliance, and risk management while providing strategic legal counsel.",
	},
	{
		Title:       "VP of Customer Success",
		Description: "Led customer success strategy, retention initiatives, and team development to ensure customer satisfaction and growth.",
	},
	{
		Title:       "VP of Customer Support",
		Description: "Directed support operations, service quality improvements, and team development to enhance customer experience.",
	},
}

func createWorkHistories() {
	for _, user := range hubUsers {
		// 90% chance of having a work history
		if rand.Float32() < 0.9 {
			// Get the auth token from the session map
			tokenI, ok := hubSessionTokens.Load(user.Email)
			if !ok {
				log.Fatalf("no auth token found for %s", user.Email)
			}
			authToken := tokenI.(string)

			err := createWorkHistory(authToken)
			if err != nil {
				log.Fatalf(
					"error creating work history for %s: %v",
					user.Email,
					err,
				)
			}
			color.Magenta("created work history for %s", user.Email)
		}
	}
}

func createWorkHistory(authToken string) error {
	numWorkHistories := rand.Intn(3) + 1

	// 1 year ago is the start date of the first work history
	startDate := time.Now().AddDate(-1, 0, 0)
	prevStartDate := startDate

	// Fill work history in reverse chronological order
	for i := 0; i < numWorkHistories; i++ {
		var endDatePtr *string
		if i == 0 {
			// First work history is current job
			endDatePtr = nil
		} else {
			// Assuming this job ended 30 days before the next job
			endDate := prevStartDate.AddDate(0, 0, -30)
			endDateStr := endDate.Format("2006-01-02")
			endDatePtr = &endDateStr

			startDate = prevStartDate.AddDate(-1, -2, -3)
			prevStartDate = startDate
		}

		err := createWorkHistoryItem(authToken, startDate, endDatePtr)
		if err != nil {
			return err
		}
	}
	return nil
}

func createWorkHistoryItem(
	authToken string,
	startDate time.Time,
	endDatePtr *string,
) error {
	job := jobs[rand.Intn(len(jobs))]
	employerDomain := employerDomains[rand.Intn(len(employerDomains))]

	var addWorkHistoryRequest = hub.AddWorkHistoryRequest{
		EmployerDomain: employerDomain,
		Title:          job.Title,
		StartDate:      startDate.Format("2006-01-02"),
		EndDate:        endDatePtr,
		Description:    &job.Description,
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
