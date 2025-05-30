package main

import (
	"bytes"
	"encoding/json"
	"log"
	"math/rand"
	"net/http"
	"sync"

	"github.com/fatih/color"
	"github.com/vetchium/vetchium/typespec/common"
	"github.com/vetchium/vetchium/typespec/employer"
)

func writeEmployerPosts() {
	// List of admin users for each employer
	employerAdmins := []struct {
		email  string
		domain string
	}{
		{
			email:  "admin@gryffindor.example",
			domain: "gryffindor.example",
		},
		{
			email:  "admin@hufflepuff.example",
			domain: "hufflepuff.example",
		},
		{
			email:  "admin@ravenclaw.example",
			domain: "ravenclaw.example",
		},
		{
			email:  "admin@slytherin.example",
			domain: "slytherin.example",
		},
	}

	// Sample employer post content organized by categories
	employerPostContent := map[string][]struct {
		content string
		tags    []common.VTagID
	}{
		"company-updates": {
			{
				content: "We're excited to announce our Q4 results! Our team has achieved remarkable growth and we're expanding to new markets. Thanks to all our dedicated employees for making this possible.",
				tags:    []common.VTagID{"leadership"},
			},
			{
				content: "New office opening next month! We're expanding our presence and creating more opportunities for talented individuals to join our team.",
				tags:    []common.VTagID{"leadership"},
			},
			{
				content: "Proud to announce our partnership with leading technology companies. This collaboration will drive innovation and bring cutting-edge solutions to our clients.",
				tags:    []common.VTagID{"leadership"},
			},
		},
		"hiring": {
			{
				content: "We're hiring! Looking for passionate software engineers to join our growing team. Competitive salary, excellent benefits, and opportunity to work on exciting projects.",
				tags: []common.VTagID{
					"human-resources",
					"software-engineering",
				},
			},
			{
				content: "Join our marketing team! We're seeking creative minds to help us tell our story and connect with customers in meaningful ways.",
				tags:    []common.VTagID{"human-resources", "marketing"},
			},
			{
				content: "Open positions in our finance department. Looking for detail-oriented professionals who want to be part of our financial success story.",
				tags:    []common.VTagID{"human-resources", "finance"},
			},
		},
		"culture": {
			{
				content: "Team building day was a huge success! Our employees enjoyed outdoor activities, team challenges, and great food. Building strong relationships is key to our success.",
				tags:    []common.VTagID{"leadership"},
			},
			{
				content: "Celebrating our employees' achievements this month. Recognition and appreciation are fundamental values in our workplace culture.",
				tags:    []common.VTagID{"leadership"},
			},
			{
				content: "Lunch and learn session today featuring industry experts. We believe in continuous learning and professional development for all our team members.",
				tags:    []common.VTagID{"leadership", "entrepreneurship"},
			},
		},
		"product": {
			{
				content: "Introducing our latest product features! We've listened to customer feedback and implemented improvements that enhance user experience and functionality.",
				tags:    []common.VTagID{"product-management"},
			},
			{
				content: "Behind the scenes: How our engineering team builds reliable, scalable solutions. Innovation and quality are at the heart of everything we do.",
				tags:    []common.VTagID{"product-management"},
			},
			{
				content: "Customer success story: How our solutions helped a client achieve 300% growth in efficiency. Real results that make a difference.",
				tags:    []common.VTagID{"product-management"},
			},
		},
	}

	var wg sync.WaitGroup
	for i, admin := range employerAdmins {
		if i%2 == 0 {
			color.Yellow("employer thread %d waiting", i)
			wg.Wait()
			color.Yellow("employer thread %d resumes", i)
		}
		wg.Add(1)

		// Parallelism to have posts from different employers at similar times
		go func(admin struct {
			email  string
			domain string
		}, i int) {
			defer wg.Done()

			tokenI, ok := employerSessionTokens.Load(admin.email)
			if !ok {
				log.Fatalf("no auth token found for employer %s", admin.email)
			}
			authToken := tokenI.(string)

			// Each employer creates 3-7 posts
			numPosts := 3 + rand.Intn(5) // 3-7 posts
			var posts []employer.AddEmployerPostRequest

			// Select random posts from different categories
			categories := []string{
				"company-updates",
				"hiring",
				"culture",
				"product",
			}
			for j := 0; j < numPosts; j++ {
				category := categories[rand.Intn(len(categories))]
				categoryPosts := employerPostContent[category]
				selectedPost := categoryPosts[rand.Intn(len(categoryPosts))]

				posts = append(posts, employer.AddEmployerPostRequest{
					Content: selectedPost.content,
					TagIDs:  selectedPost.tags,
				})
			}

			for k, post := range posts {
				color.Yellow("employer thread %d creating post %d", i, k)
				createEmployerPost(post, authToken)
			}

			color.Green(
				"Employer %s created %d posts",
				admin.domain,
				len(posts),
			)
		}(admin, i)
	}

	wg.Wait()
}

func createEmployerPost(
	post employer.AddEmployerPostRequest,
	authToken string,
) {
	var body bytes.Buffer
	err := json.NewEncoder(&body).Encode(post)
	if err != nil {
		log.Fatalf("failed to encode employer post: %v", err)
	}

	req, err := http.NewRequest(
		http.MethodPost,
		serverURL+"/employer/add-post",
		&body,
	)
	if err != nil {
		log.Fatalf("failed to create employer post request: %v", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+authToken)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Fatalf("failed to send employer post request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := json.Marshal(post)
		log.Fatalf(
			"failed to create employer post: %v, request: %s",
			resp.Status,
			string(bodyBytes),
		)
	}
}
