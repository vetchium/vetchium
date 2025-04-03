package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"

	"github.com/fatih/color"
	"github.com/psankar/vetchi/typespec/common"
	"github.com/psankar/vetchi/typespec/hub"
)

// AchievementTemplate represents a possible achievement for a career path
type AchievementTemplate struct {
	Type        common.AchievementType
	Title       string
	Description string
	URL         string
}

// Map of career paths to possible achievements
var careerAchievements = map[string][]AchievementTemplate{
	"Engineering": {
		{
			Type:        common.Patent,
			Title:       "Distributed System Architecture",
			Description: "Patent for innovative distributed system design",
			URL:         "https://patents.example.com/distributed-systems",
		},
		{
			Type:        common.Patent,
			Title:       "Cloud Computing Optimization",
			Description: "Patent for cloud resource optimization algorithm",
			URL:         "https://patents.example.com/cloud-optimization",
		},
		{
			Type:        common.Patent,
			Title:       "Quantum Computing Interface",
			Description: "Patent for quantum computing data interface",
			URL:         "https://patents.example.com/quantum-interface",
		},
		{
			Type:        common.Patent,
			Title:       "Blockchain Consensus Protocol",
			Description: "Patent for novel blockchain consensus mechanism",
			URL:         "https://patents.example.com/blockchain-consensus",
		},
		{
			Type:        common.Publication,
			Title:       "Modern Microservices Architecture",
			Description: "Research paper on scalable microservices design",
			URL:         "https://journals.example.com/microservices",
		},
		{
			Type:        common.Publication,
			Title:       "Serverless Computing at Scale",
			Description: "Research on large-scale serverless architectures",
			URL:         "https://journals.example.com/serverless",
		},
		{
			Type:        common.Publication,
			Title:       "Edge Computing Optimization",
			Description: "Research on optimizing edge computing networks",
			URL:         "https://journals.example.com/edge-computing",
		},
		{
			Type:        common.Certification,
			Title:       "AWS Solutions Architect Professional",
			Description: "Advanced certification for AWS architecture",
			URL:         "https://aws.example.com/certification",
		},
		{
			Type:        common.Certification,
			Title:       "Google Cloud Professional Architect",
			Description: "Professional certification for GCP architecture",
			URL:         "https://google.example.com/certification",
		},
		{
			Type:        common.Certification,
			Title:       "Azure Solutions Architect Expert",
			Description: "Expert level certification for Azure architecture",
			URL:         "https://azure.example.com/certification",
		},
		{
			Type:        common.Certification,
			Title:       "Kubernetes Administrator (CKA)",
			Description: "Advanced certification for Kubernetes administration",
			URL:         "https://kubernetes.example.com/certification",
		},
	},
	"Data Science": {
		{
			Type:        common.Patent,
			Title:       "Machine Learning Algorithm",
			Description: "Patent for novel ML prediction system",
			URL:         "https://patents.example.com/ml-algorithm",
		},
		{
			Type:        common.Patent,
			Title:       "Neural Network Architecture",
			Description: "Patent for innovative neural network design",
			URL:         "https://patents.example.com/neural-network",
		},
		{
			Type:        common.Patent,
			Title:       "Automated Feature Engineering",
			Description: "Patent for automated feature discovery system",
			URL:         "https://patents.example.com/feature-engineering",
		},
		{
			Type:        common.Publication,
			Title:       "Deep Learning in Computer Vision",
			Description: "Research on advanced CV techniques",
			URL:         "https://journals.example.com/deep-learning",
		},
		{
			Type:        common.Publication,
			Title:       "Natural Language Processing Advances",
			Description: "Research on modern NLP architectures",
			URL:         "https://journals.example.com/nlp-advances",
		},
		{
			Type:        common.Publication,
			Title:       "Reinforcement Learning in Robotics",
			Description: "Research on RL applications in robotics",
			URL:         "https://journals.example.com/rl-robotics",
		},
		{
			Type:        common.Certification,
			Title:       "TensorFlow Developer Certificate",
			Description: "Professional certification in TensorFlow",
			URL:         "https://tensorflow.example.com/certification",
		},
		{
			Type:        common.Certification,
			Title:       "PyTorch Advanced Developer",
			Description: "Advanced certification in PyTorch development",
			URL:         "https://pytorch.example.com/certification",
		},
		{
			Type:        common.Certification,
			Title:       "Deep Learning Specialization",
			Description: "Comprehensive certification in deep learning",
			URL:         "https://deeplearning.example.com/certification",
		},
	},
	"Finance": {
		{
			Type:        common.Patent,
			Title:       "Algorithmic Trading System",
			Description: "Patent for automated trading platform",
			URL:         "https://patents.example.com/algo-trading",
		},
		{
			Type:        common.Patent,
			Title:       "Fraud Detection Algorithm",
			Description: "Patent for real-time fraud detection",
			URL:         "https://patents.example.com/fraud-detection",
		},
		{
			Type:        common.Patent,
			Title:       "Cryptocurrency Trading Protocol",
			Description: "Patent for secure crypto trading system",
			URL:         "https://patents.example.com/crypto-trading",
		},
		{
			Type:        common.Publication,
			Title:       "Risk Management in Banking",
			Description: "Research on modern risk assessment",
			URL:         "https://journals.example.com/risk-management",
		},
		{
			Type:        common.Publication,
			Title:       "Quantitative Investment Strategies",
			Description: "Research on advanced quant strategies",
			URL:         "https://journals.example.com/quant-strategies",
		},
		{
			Type:        common.Publication,
			Title:       "ESG Investment Analysis",
			Description: "Research on ESG investment metrics",
			URL:         "https://journals.example.com/esg-analysis",
		},
		{
			Type:        common.Certification,
			Title:       "Chartered Financial Analyst (CFA)",
			Description: "Professional certification in financial analysis",
			URL:         "https://cfa.example.com/certification",
		},
		{
			Type:        common.Certification,
			Title:       "Financial Risk Manager (FRM)",
			Description: "Professional certification in risk management",
			URL:         "https://frm.example.com/certification",
		},
		{
			Type:        common.Certification,
			Title:       "Certified Financial Planner (CFP)",
			Description: "Professional certification in financial planning",
			URL:         "https://cfp.example.com/certification",
		},
	},
	"Medical": {
		{
			Type:        common.Patent,
			Title:       "Medical Diagnostic System",
			Description: "Patent for AI-based diagnosis",
			URL:         "https://patents.example.com/medical-diagnosis",
		},
		{
			Type:        common.Patent,
			Title:       "Drug Discovery Platform",
			Description: "Patent for ML-driven drug discovery",
			URL:         "https://patents.example.com/drug-discovery",
		},
		{
			Type:        common.Patent,
			Title:       "Remote Patient Monitoring",
			Description: "Patent for IoT-based patient monitoring",
			URL:         "https://patents.example.com/patient-monitoring",
		},
		{
			Type:        common.Publication,
			Title:       "Advances in Telemedicine",
			Description: "Research on remote healthcare delivery",
			URL:         "https://journals.example.com/telemedicine",
		},
		{
			Type:        common.Publication,
			Title:       "Personalized Medicine Approaches",
			Description: "Research on genomic medicine",
			URL:         "https://journals.example.com/personalized-medicine",
		},
		{
			Type:        common.Publication,
			Title:       "AI in Healthcare",
			Description: "Research on AI applications in healthcare",
			URL:         "https://journals.example.com/ai-healthcare",
		},
		{
			Type:        common.Certification,
			Title:       "Medical Device Safety",
			Description: "Certification in medical device standards",
			URL:         "https://medical.example.com/certification",
		},
		{
			Type:        common.Certification,
			Title:       "Healthcare Data Security",
			Description: "Certification in healthcare data protection",
			URL:         "https://healthsecurity.example.com/certification",
		},
		{
			Type:        common.Certification,
			Title:       "Clinical Research Management",
			Description: "Certification in clinical trial management",
			URL:         "https://clinical.example.com/certification",
		},
	},
	"Law": {
		{
			Type:        common.Publication,
			Title:       "Digital Privacy Laws",
			Description: "Research on modern privacy regulations",
			URL:         "https://journals.example.com/privacy-law",
		},
		{
			Type:        common.Publication,
			Title:       "Blockchain Legal Framework",
			Description: "Research on cryptocurrency regulations",
			URL:         "https://journals.example.com/blockchain-law",
		},
		{
			Type:        common.Publication,
			Title:       "AI and Legal Liability",
			Description: "Research on AI system liability",
			URL:         "https://journals.example.com/ai-liability",
		},
		{
			Type:        common.Publication,
			Title:       "International IP Protection",
			Description: "Research on global IP rights",
			URL:         "https://journals.example.com/ip-protection",
		},
		{
			Type:        common.Certification,
			Title:       "International Business Law",
			Description: "Certification in global business regulations",
			URL:         "https://law.example.com/certification",
		},
		{
			Type:        common.Certification,
			Title:       "Technology Law Specialist",
			Description: "Certification in technology law",
			URL:         "https://techlaw.example.com/certification",
		},
		{
			Type:        common.Certification,
			Title:       "Data Privacy Professional",
			Description: "Certification in privacy law compliance",
			URL:         "https://privacy.example.com/certification",
		},
	},
}

// createAchievements generates and creates achievements for all hub users
func createAchievements() {
	for i := 0; i < len(hubUsers); i++ {
		user := hubUsers[i]
		// Get the auth token from the session map
		tokenI, ok := hubSessionTokens.Load(user.Email)
		if !ok {
			log.Fatalf("no auth token found for %s", user.Email)
		}
		authToken := tokenI.(string)

		// Find career path for the user based on their jobs
		var careerPath string
		for _, path := range careerPaths {
			for _, job := range user.Jobs {
				if contains(path.Steps, job.Title) {
					careerPath = path.Tag
					break
				}
			}
			if careerPath != "" {
				break
			}
		}

		if careerPath == "" {
			continue // Skip if no matching career path found
		}

		// Get possible achievements for this career
		achievements := careerAchievements[careerPath]
		if len(achievements) == 0 {
			continue
		}

		// Select achievements based on career level
		selectedAchievements := selectAchievements(achievements, len(user.Jobs))

		// Create each achievement
		for _, achievement := range selectedAchievements {
			// First save to the user struct
			achievementReq := hub.AddAchievementRequest{
				Type:        achievement.Type,
				Title:       achievement.Title,
				Description: &achievement.Description,
				URL:         &achievement.URL,
			}
			user.Achievements = append(user.Achievements, achievementReq)

			// Then create via API
			err := createAchievement(authToken, achievement)
			if err != nil {
				log.Printf(
					"Failed to create achievement for %s: %v",
					user.Email,
					err,
				)
				continue
			}
			color.Magenta(
				"created %s achievement for %s",
				achievement.Type,
				user.Email,
			)
		}
	}
}

// selectAchievements chooses appropriate achievements based on career level
// The number and type of achievements are selected based on the user's experience:
// - Junior (1-2 jobs): 2-3 achievements, mostly certifications
// - Mid-level (3-4 jobs): 3-4 achievements, mix of certifications and publications
// - Senior (5+ jobs): 4-5 achievements, including patents if available
func selectAchievements(
	available []AchievementTemplate,
	jobCount int,
) []AchievementTemplate {
	var selected []AchievementTemplate
	var numAchievements int

	// Determine number of achievements based on experience
	switch {
	case jobCount <= 2:
		numAchievements = rand.Intn(2) + 2 // 2-3 achievements
	case jobCount <= 4:
		numAchievements = rand.Intn(2) + 3 // 3-4 achievements
	default:
		numAchievements = rand.Intn(2) + 4 // 4-5 achievements
	}

	// Ensure we don't exceed available achievements
	if numAchievements > len(available) {
		numAchievements = len(available)
	}

	// Group achievements by type
	certifications := filterByType(available, common.Certification)
	publications := filterByType(available, common.Publication)
	patents := filterByType(available, common.Patent)

	// Select achievements based on experience level
	switch {
	case jobCount <= 2:
		// Junior: Focus on certifications
		selected = append(selected, selectRandom(certifications, 2)...)
		if len(selected) < numAchievements {
			selected = append(selected, selectRandom(publications, 1)...)
		}
	case jobCount <= 4:
		// Mid-level: Mix of certifications and publications
		selected = append(selected, selectRandom(certifications, 2)...)
		selected = append(selected, selectRandom(publications, 1)...)
		if len(patents) > 0 && len(selected) < numAchievements {
			selected = append(selected, selectRandom(patents, 1)...)
		}
	default:
		// Senior: Include patents if available
		if len(patents) > 0 {
			selected = append(selected, selectRandom(patents, 2)...)
		}
		selected = append(selected, selectRandom(publications, 1)...)
		selected = append(selected, selectRandom(certifications, 1)...)
	}

	// If we still need more achievements, add random ones
	remaining := available
	for _, s := range selected {
		remaining = removeAchievement(remaining, s)
	}
	if len(selected) < numAchievements && len(remaining) > 0 {
		selected = append(
			selected,
			selectRandom(remaining, numAchievements-len(selected))...)
	}

	return selected
}

// filterByType returns achievements of a specific type
func filterByType(
	achievements []AchievementTemplate,
	achievementType common.AchievementType,
) []AchievementTemplate {
	var filtered []AchievementTemplate
	for _, a := range achievements {
		if a.Type == achievementType {
			filtered = append(filtered, a)
		}
	}
	return filtered
}

// selectRandom returns up to n random items from the slice
func selectRandom(items []AchievementTemplate, n int) []AchievementTemplate {
	if n > len(items) {
		n = len(items)
	}
	if n == 0 {
		return nil
	}

	// Create a copy to avoid modifying the original
	temp := make([]AchievementTemplate, len(items))
	copy(temp, items)

	// Fisher-Yates shuffle
	for i := len(temp) - 1; i > 0; i-- {
		j := rand.Intn(i + 1)
		temp[i], temp[j] = temp[j], temp[i]
	}

	return temp[:n]
}

// removeAchievement removes an achievement from a slice
func removeAchievement(
	achievements []AchievementTemplate,
	achievement AchievementTemplate,
) []AchievementTemplate {
	var result []AchievementTemplate
	for _, a := range achievements {
		if a.Title != achievement.Title {
			result = append(result, a)
		}
	}
	return result
}

func createAchievement(
	authToken string,
	achievement AchievementTemplate,
) error {
	request := hub.AddAchievementRequest{
		Type:        achievement.Type,
		Title:       achievement.Title,
		Description: &achievement.Description,
		URL:         &achievement.URL,
	}

	body, err := json.Marshal(request)
	if err != nil {
		return fmt.Errorf("failed to marshal request: %v", err)
	}

	req, err := http.NewRequest(
		http.MethodPost,
		serverURL+"/hub/add-achievement",
		bytes.NewBuffer(body),
	)
	if err != nil {
		return fmt.Errorf("failed to create request: %v", err)
	}

	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", authToken))
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("request failed with status: %s", resp.Status)
	}

	return nil
}

// Helper function to check if a slice contains a string
func contains(slice []string, str string) bool {
	for _, s := range slice {
		if s == str {
			return true
		}
	}
	return false
}
