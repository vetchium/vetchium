package main

import "github.com/fatih/color"

func createApplications() {
	for _, user := range hubUsers {
		createApplication(user)
	}
}

func createApplication(user HubUser) {
	for _, company := range user.PreferredCompanyDomains {
		createApplicationForCompany(user, company)
	}
}

func createApplicationForCompany(user HubUser, company string) {
	color.Green("Creating application for %q for %q", user.Name, company)
}
