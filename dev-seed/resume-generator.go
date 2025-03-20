package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/fatih/color"
	"github.com/jung-kurt/gofpdf"
)

var pdfTempDir string

// initResumePDFDirectory creates a temporary directory for PDF resumes
func initResumePDFDirectory() {
	tempDir := filepath.Join(
		os.TempDir(),
		fmt.Sprintf("vetchi-resumes-%d", time.Now().Unix()),
	)
	if err := os.MkdirAll(tempDir, 0755); err != nil {
		log.Fatalf("Failed to create PDF directory: %v", err)
	}
	pdfTempDir = tempDir
	color.Green("Created temporary PDF directory: %s", pdfTempDir)
}

// generateResumePDF creates a PDF resume for the given user
func generateResumePDF(user HubSeedUser) (string, []byte, error) {
	pdf := gofpdf.New("P", "mm", "A4", "")

	// Set default margins
	pdf.SetMargins(20, 20, 20)
	pdf.SetAutoPageBreak(true, 20)

	// Add a new page
	pdf.AddPage()

	// Set font
	pdf.SetFont("Arial", "B", 18)

	// Add user name
	pdf.Cell(0, 12, user.Name)
	pdf.Ln(8)

	// Add current position
	pdf.SetFont("Arial", "I", 12)
	pdf.Cell(0, 8, user.ShortBio)
	pdf.Ln(12)

	// Add contact info section
	pdf.SetFont("Arial", "", 10)
	pdf.Cell(40, 6, fmt.Sprintf("Email: %s", user.Email))
	pdf.Ln(6)
	pdf.Cell(
		40,
		6,
		fmt.Sprintf(
			"Location: %s, %s",
			user.ResidentCity,
			user.ResidentCountry,
		),
	)
	pdf.Ln(15)

	// Add professional summary
	pdf.SetFont("Arial", "B", 14)
	pdf.Cell(0, 10, "Professional Summary")
	pdf.Ln(8)

	pdf.SetFont("Arial", "", 10)

	// Format the bio text to fit within the page width
	pdf.MultiCell(0, 6, user.LongBio, "", "", false)
	pdf.Ln(8)

	// Add work experience section
	pdf.SetFont("Arial", "B", 14)
	pdf.Cell(0, 10, "Work Experience")
	pdf.Ln(10)

	// Add each job
	for i, workItem := range user.WorkHistoryItems {
		// Add employer name
		pdf.SetFont("Arial", "B", 12)
		pdf.Cell(0, 8, workItem.EmployerName)
		pdf.Ln(8)

		// Add job title
		pdf.SetFont("Arial", "I", 11)
		pdf.Cell(0, 6, workItem.JobTitle)
		pdf.Ln(6)

		// Add work period
		pdf.SetFont("Arial", "", 10)
		startDateStr := workItem.StartDate.Format("Jan 2006")
		endDateStr := "Present"
		if workItem.EndDate != nil {
			endDateStr = workItem.EndDate.Format("Jan 2006")
		}
		pdf.Cell(0, 6, fmt.Sprintf("%s - %s", startDateStr, endDateStr))
		pdf.Ln(8)

		// Add job description if available
		if workItem.Description != "" {
			pdf.MultiCell(0, 6, workItem.Description, "", "", false)
			pdf.Ln(6)
		}

		// Add space between jobs except for the last one
		if i < len(user.WorkHistoryItems)-1 {
			pdf.Ln(4)
		}
	}

	// Add skills section
	pdf.Ln(10)
	pdf.SetFont("Arial", "B", 14)
	pdf.Cell(0, 10, "Skills")
	pdf.Ln(8)

	// Generate some skills based on the career path
	skills := generateSkills(user.Jobs[len(user.Jobs)-1].Title)
	pdf.SetFont("Arial", "", 10)
	pdf.MultiCell(0, 6, strings.Join(skills, ", "), "", "", false)

	// Save PDF to temp directory
	filename := filepath.Join(
		pdfTempDir,
		fmt.Sprintf("%s-resume.pdf", user.Handle),
	)
	err := pdf.OutputFileAndClose(filename)
	if err != nil {
		return "", nil, err
	}

	// Read the file back to return the bytes
	pdfBytes, err := ioutil.ReadFile(filename)
	if err != nil {
		return filename, nil, err
	}

	return filename, pdfBytes, nil
}

// generateSkills creates a list of skills based on job title
func generateSkills(jobTitle string) []string {
	// Base skills everyone should have
	baseSkills := []string{
		"Communication",
		"Problem Solving",
		"Teamwork",
		"Time Management",
	}

	// Job-specific skills
	jobSkills := map[string][]string{
		"Software Engineer": {
			"Python", "JavaScript", "Go", "Java", "React", "Node.js",
			"SQL", "Docker", "Kubernetes", "AWS",
		},
		"Data Scientist": {
			"Python", "R", "SQL", "Machine Learning", "Data Visualization",
			"Statistical Analysis", "TensorFlow", "PyTorch",
		},
		"Product Manager": {
			"Product Strategy", "User Research", "Agile", "Scrum",
			"Roadmap Planning", "Stakeholder Management", "Market Analysis",
		},
		"Marketing Manager": {
			"Digital Marketing", "Social Media", "Content Strategy",
			"SEO", "Analytics", "Campaign Management",
		},
		"Financial Analyst": {
			"Financial Modeling", "Excel", "Bloomberg Terminal",
			"Valuation", "Financial Reporting", "Risk Analysis",
		},
		"HR Manager": {
			"Recruitment", "Employee Relations", "Performance Management",
			"Compensation", "Labor Laws", "Training & Development",
		},
	}

	// Add job-specific skills if available
	specificSkills := []string{}
	for key, skills := range jobSkills {
		if strings.Contains(jobTitle, key) {
			specificSkills = skills
			break
		}
	}

	// If no specific skills found, add some generic professional skills
	if len(specificSkills) == 0 {
		specificSkills = []string{
			"Project Management",
			"Leadership",
			"Strategic Planning",
			"Analytical Thinking",
			"Negotiation",
		}
	}

	// Combine base skills and job-specific skills
	allSkills := append(baseSkills, specificSkills...)

	// Randomly select a subset of skills (between 6-10)
	numSkills := rand.Intn(5) + 6
	if numSkills > len(allSkills) {
		numSkills = len(allSkills)
	}

	// Shuffle and take the first numSkills
	rand.Shuffle(len(allSkills), func(i, j int) {
		allSkills[i], allSkills[j] = allSkills[j], allSkills[i]
	})

	return allSkills[:numSkills]
}

// generateResumesForAllUsers generates PDF resumes for all hub users
func generateResumesForAllUsers() {
	color.Cyan("Generating PDF resumes for all users...")

	// Make sure PDF directory is initialized
	if pdfTempDir == "" {
		initResumePDFDirectory()
	}

	generatedCount := 0

	// Generate resume for each user
	for _, user := range hubUsers {
		_, _, err := generateResumePDF(user)
		if err != nil {
			log.Printf("Failed to generate resume for %s: %v", user.Email, err)
			continue
		}

		generatedCount++
		if generatedCount%10 == 0 {
			color.Magenta("Generated %d resumes...", generatedCount)
		}
	}

	color.Green(
		"Successfully generated %d resume PDFs in: %s",
		generatedCount,
		pdfTempDir,
	)
}
