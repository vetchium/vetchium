package main

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/vetchium/vetchium/typespec/common"
	"github.com/vetchium/vetchium/typespec/hub"
)

var countries = []string{
	"IND",
	"USA",
	"PRC",
	"JAP",
	"GBR",
	"CAN",
	"AUS",
	"DEU",
	"FRA",
	"ITA",
}

var cities map[string][]string

func init() {
	cities = make(map[string][]string)
	cities["IND"] = []string{
		"Chennai",
		"Bengaluru",
		"Mumbai",
		"Delhi",
		"Hyderabad",
	}
	cities["USA"] = []string{
		"New York",
		"Los Angeles",
		"Chicago",
		"Houston",
		"Miami",
	}
	cities["PRC"] = []string{
		"Beijing",
		"Shanghai",
		"Guangzhou",
		"Shenzhen",
		"Chengdu",
	}
	cities["JAP"] = []string{"Tokyo", "Osaka", "Nagoya", "Sapporo", "Fukuoka"}
	cities["GBR"] = []string{
		"London",
		"Manchester",
		"Birmingham",
		"Glasgow",
		"Edinburgh",
	}
	cities["CAN"] = []string{
		"Toronto",
		"Vancouver",
		"Montreal",
		"Calgary",
		"Edmonton",
	}
	cities["AUS"] = []string{
		"Sydney",
		"Melbourne",
		"Brisbane",
		"Perth",
		"Adelaide",
	}
	cities["DEU"] = []string{
		"Berlin",
		"Hamburg",
		"Munich",
		"Frankfurt",
		"Hamburg",
	}
	cities["FRA"] = []string{"Paris", "Marseille", "Lyon", "Toulouse", "Nice"}
	cities["ITA"] = []string{"Rome", "Milan", "Naples", "Turin", "Palermo"}
}

type CareerPath struct {
	BroadArea string
	Steps     []string
}

type Job struct {
	Title   string
	Website string
}

// WorkHistoryItem stores detailed information about a job including dates
type WorkHistoryItem struct {
	EmployerID   string
	EmployerName string
	StartDate    time.Time
	EndDate      *time.Time
	JobTitle     string
	Description  string
}

type HubSeedUser struct {
	Name                   string
	Handle                 string
	Email                  string
	Tier                   hub.HubUserTier
	ResidentCountry        string
	ResidentCity           string
	PreferredLanguage      string
	ShortBio               string
	LongBio                string
	ProfilePictureFilename string

	Jobs []Job

	// Track detailed work history for each user
	WorkHistoryItems []WorkHistoryItem

	// TODO: Find out how we can efficiently populate these
	Endorsers             []common.Handle
	ApplyToCompanyDomains []string

	Achievements []hub.AddAchievementRequest

	BroadArea string
}

var firstNames = []string{
	"Liam",
	"Noah",
	"Oliver",
	"Elijah",
	"James",
	"William",
	"Benjamin",
	"Lucas",
	"Henry",
	"Alexander",
	"Aarav",
	"Wei",
	"Hiroshi",
	"Fatima",
	"Amara",
	"Carlos",
	"Sofia",
	"Yusuf",
	"Zhang",
	"Muhammad",
}

var lastNames = []string{
	"Smith",
	"Johnson",
	"Williams",
	"Brown",
	"Jones",
	"Garcia",
	"Miller",
	"Davis",
	"Rodriguez",
	"Martinez",
	"Chen",
	"Kumar",
	"Tanaka",
	"Singh",
	"Ali",
	"Lopez",
	"Gonzalez",
	"Hernandez",
	"Nguyen",
	"Kim",
}

var careerPaths = []CareerPath{
	{
		"Engineering",
		[]string{
			"Software Engineer",
			"Senior Software Engineer",
			"Staff Engineer",
			"Principal Engineer",
			"Senior Principal Engineer",
			"Distinguished Engineer",
			"Fellow",
			"Senior Fellow",
		},
	},
	{
		"Engineering",
		[]string{
			"Software Engineer",
			"Senior Software Engineer",
			"Engineering Manager",
			"Senior Engineering Manager",
			"Director of Engineering",
			"VP of Engineering",
			"CTO",
		},
	},
	{
		"Finance",
		[]string{
			"Financial Analyst",
			"Investment Associate",
			"Investment Manager",
			"Senior Investment Manager",
			"Managing Director",
			"Chief Investment Officer",
		},
	},
	{
		"Finance",
		[]string{
			"Risk Analyst",
			"Risk Manager",
			"Senior Risk Manager",
			"Director of Risk",
			"Chief Risk Officer",
		},
	},
	{
		"Finance",
		[]string{
			"Investment Banking Analyst",
			"Associate Investment Banker",
			"Senior Investment Banker",
			"Director of Investment Banking",
			"Managing Director",
		},
	},
	{
		"Law",
		[]string{
			"Paralegal",
			"Associate Lawyer",
			"Senior Lawyer",
			"Partner",
			"Managing Partner",
		},
	},
	{
		"Law",
		[]string{
			"Legal Assistant",
			"Legal Counsel",
			"Senior Legal Advisor",
			"Director of Legal Affairs",
			"Chief Legal Officer",
		},
	},
	{
		"Education",
		[]string{
			"Teaching Assistant",
			"Lecturer",
			"Assistant Professor",
			"Professor",
			"Dean",
			"Vice Chancellor",
		},
	},
	{
		"Education",
		[]string{
			"Curriculum Developer",
			"Education Consultant",
			"Director of Curriculum",
			"VP of Academic Affairs",
		},
	},
	{
		"Marketing",
		[]string{
			"Marketing Associate",
			"Marketing Manager",
			"Senior Marketing Manager",
			"Director of Marketing",
			"Chief Marketing Officer",
		},
	},
	{
		"Marketing",
		[]string{
			"Brand Manager",
			"Senior Brand Manager",
			"Director of Brand Strategy",
			"VP of Branding",
		},
	},
	{
		"Human Resources",
		[]string{
			"HR Associate",
			"HR Manager",
			"Senior HR Manager",
			"Director of HR",
			"VP of HR",
			"Chief Human Resources Officer",
		},
	},
	{
		"Medical",
		[]string{
			"Medical Intern",
			"Resident Doctor",
			"Attending Physician",
			"Chief Surgeon",
			"Medical Director",
		},
	},
	{
		"Medical",
		[]string{
			"Healthcare Administrator",
			"Senior Administrator",
			"Director of Healthcare",
			"VP of Healthcare",
		},
	},
	{
		"Data Science",
		[]string{
			"Data Analyst",
			"Senior Data Analyst",
			"Data Scientist",
			"Senior Data Scientist",
			"Head of Data Science",
			"Chief Data Officer",
		},
	},
	{
		"Consulting",
		[]string{
			"Consultant",
			"Senior Consultant",
			"Managing Consultant",
			"Partner",
			"Managing Partner",
		},
	},
	{
		"Design",
		[]string{
			"UX Designer",
			"Senior UX Designer",
			"Design Manager",
			"Director of UX",
			"VP of Design",
		},
	},
	{
		"Design",
		[]string{
			"Graphic Designer",
			"Senior Graphic Designer",
			"Art Director",
			"Creative Director",
			"Chief Design Officer",
		},
	},
	{
		"Operations",
		[]string{
			"Operations Analyst",
			"Operations Manager",
			"Senior Operations Manager",
			"Director of Operations",
			"VP of Operations",
			"COO",
		},
	},
	{
		"Sales",
		[]string{
			"Sales Associate",
			"Sales Manager",
			"Senior Sales Manager",
			"Director of Sales",
			"VP of Sales",
			"Chief Revenue Officer",
		},
	},
	{
		"Product Management",
		[]string{
			"Product Analyst",
			"Associate Product Manager",
			"Product Manager",
			"Senior Product Manager",
			"Director of Product",
			"VP of Product",
			"Chief Product Officer",
		},
	},
	{
		"Aerospace",
		[]string{
			"Aerospace Engineer",
			"Senior Aerospace Engineer",
			"Lead Engineer",
			"Chief Engineer",
			"Director of Aerospace Engineering",
		},
	},
	{
		"Aerospace",
		[]string{
			"Aeronautical Engineer",
			"Flight Systems Engineer",
			"Senior Flight Systems Engineer",
			"Spacecraft Engineer",
			"Chief Spacecraft Engineer",
		},
	},
	{
		"Automotive",
		[]string{
			"Automotive Engineer",
			"Senior Automotive Engineer",
			"Powertrain Engineer",
			"Lead Powertrain Engineer",
			"Chief Vehicle Engineer",
		},
	},
	{
		"Hospitality",
		[]string{
			"Front Desk Associate",
			"Concierge",
			"Hotel Manager",
			"Director of Guest Services",
			"General Manager",
		},
	},
	{
		"Hospitality",
		[]string{
			"Restaurant Server",
			"Restaurant Manager",
			"Food and Beverage Director",
			"Executive Chef",
			"Director of Hospitality",
		},
	},
	{
		"Retail",
		[]string{
			"Retail Associate",
			"Department Manager",
			"Store Manager",
			"Regional Manager",
			"Director of Retail Operations",
		},
	},
	{
		"Pharmaceuticals",
		[]string{
			"Research Associate",
			"Research Scientist",
			"Principal Scientist",
			"Director of Research",
			"Chief Scientific Officer",
		},
	},
	{
		"Pharmaceuticals",
		[]string{
			"Clinical Trial Coordinator",
			"Clinical Research Manager",
			"Director of Clinical Development",
			"VP of Drug Development",
		},
	},
	{
		"Construction",
		[]string{
			"Construction Worker",
			"Project Coordinator",
			"Project Manager",
			"Construction Manager",
			"Director of Construction",
		},
	},
	{
		"Real Estate",
		[]string{
			"Real Estate Agent",
			"Senior Real Estate Agent",
			"Broker",
			"Managing Broker",
			"Real Estate Director",
		},
	},
	{
		"Entertainment",
		[]string{
			"Production Assistant",
			"Associate Producer",
			"Producer",
			"Executive Producer",
			"Studio Executive",
		},
	},
	{
		"Media",
		[]string{
			"Reporter",
			"Senior Reporter",
			"Editor",
			"Managing Editor",
			"Editor-in-Chief",
		},
	},
	{
		"Telecommunications",
		[]string{
			"Network Engineer",
			"Senior Network Engineer",
			"Network Architect",
			"Director of Network Operations",
			"Chief Technology Officer",
		},
	},
	{
		"Renewable Energy",
		[]string{
			"Energy Analyst",
			"Renewable Energy Engineer",
			"Senior Energy Engineer",
			"Director of Energy Systems",
			"Chief Sustainability Officer",
		},
	},
}

var employers = []struct {
	Name               string
	Website            string
	HiringInBroadAreas []string
}{
	{
		"Google",
		"google.example",
		[]string{"Engineering", "Data Science", "Product Management"},
	},
	{
		"Microsoft",
		"microsoft.example",
		[]string{"Engineering", "Product Management", "Sales"},
	},
	{
		"Goldman Sachs",
		"goldmansachs.example",
		[]string{"Finance", "Consulting", "Data Science"},
	},
	{
		"JP Morgan",
		"jpmorgan.example",
		[]string{"Finance", "Data Science", "Operations"},
	},
	{
		"Harvard University",
		"harvard.example",
		[]string{"Education", "Research", "Medical"},
	},
	{
		"Stanford University",
		"stanford.example",
		[]string{"Education", "Engineering", "Medical"},
	},
	{
		"Pfizer",
		"pfizer.example",
		[]string{"Medical", "Pharmaceuticals", "Research"},
	},
	{
		"Mayo Clinic",
		"mayoclinic.example",
		[]string{"Medical", "Healthcare", "Research"},
	},
	{
		"McKinsey & Company",
		"mckinsey.example",
		[]string{"Consulting", "Finance", "Data Science"},
	},
	{
		"Boston Consulting Group",
		"bcg.example",
		[]string{"Consulting", "Data Science", "Operations"},
	},
	{
		"Baker & McKenzie",
		"bakermckenzie.example",
		[]string{"Law", "Consulting", "Finance"},
	},
	{
		"Skadden, Arps, Slate, Meagher & Flom",
		"skadden.example",
		[]string{"Law", "Finance", "Consulting"},
	},
	{"Ogilvy", "ogilvy.example", []string{"Marketing", "Design", "Sales"}},
	{
		"Publicis Group",
		"publicis.example",
		[]string{"Marketing", "Design", "Media"},
	},
	{
		"Adecco",
		"adecco.example",
		[]string{"Human Resources", "Consulting", "Sales"},
	},
	{
		"Randstad",
		"randstad.example",
		[]string{"Human Resources", "Operations", "Consulting"},
	},
	{
		"IBM Watson",
		"ibmwatson.example",
		[]string{"Data Science", "Engineering", "Consulting"},
	},
	{
		"Palantir Technologies",
		"palantir.example",
		[]string{"Data Science", "Engineering", "Consulting"},
	},
	{
		"IDEO",
		"ideo.example",
		[]string{"Design", "Consulting", "Product Management"},
	},
	{
		"Pentagram",
		"pentagram.example",
		[]string{"Design", "Marketing", "Media"},
	},
	{"Maersk", "maersk.example", []string{"Operations", "Logistics", "Sales"}},
	{"FedEx", "fedex.example", []string{"Operations", "Logistics", "Sales"}},
	{
		"Salesforce",
		"salesforce.example",
		[]string{"Sales", "Engineering", "Product Management"},
	},
	{
		"Oracle",
		"oracle.example",
		[]string{"Sales", "Engineering", "Consulting"},
	},
	{
		"Atlassian",
		"atlassian.example",
		[]string{"Product Management", "Engineering", "Sales"},
	},
	{
		"Slack",
		"slack.example",
		[]string{"Product Management", "Engineering", "Sales"},
	},
	{
		"Boeing",
		"boeing.example",
		[]string{"Aerospace", "Engineering", "Operations"},
	},
	{
		"Airbus",
		"airbus.example",
		[]string{"Aerospace", "Engineering", "Operations"},
	},
	{
		"Tesla",
		"tesla.example",
		[]string{"Automotive", "Engineering", "Product Management"},
	},
	{
		"Toyota",
		"toyota.example",
		[]string{"Automotive", "Engineering", "Operations"},
	},
	{
		"Marriott International",
		"marriott.example",
		[]string{"Hospitality", "Operations", "Sales"},
	},
	{
		"Hilton Worldwide",
		"hilton.example",
		[]string{"Hospitality", "Operations", "Sales"},
	},
	{"Walmart", "walmart.example", []string{"Retail", "Operations", "Sales"}},
	{"Target", "target.example", []string{"Retail", "Operations", "Marketing"}},
	{
		"Novartis",
		"novartis.example",
		[]string{"Pharmaceuticals", "Medical", "Research"},
	},
	{
		"Merck",
		"merck.example",
		[]string{"Pharmaceuticals", "Medical", "Research"},
	},
	{
		"Bechtel",
		"bechtel.example",
		[]string{"Construction", "Engineering", "Operations"},
	},
	{
		"AECOM",
		"aecom.example",
		[]string{"Construction", "Engineering", "Design"},
	},
	{
		"Coldwell Banker",
		"coldwellbanker.example",
		[]string{"Real Estate", "Sales", "Marketing"},
	},
	{"RE/MAX", "remax.example", []string{"Real Estate", "Sales", "Marketing"}},
	{
		"Warner Bros.",
		"warnerbros.example",
		[]string{"Entertainment", "Media", "Marketing"},
	},
	{
		"Universal Pictures",
		"universal.example",
		[]string{"Entertainment", "Media", "Marketing"},
	},
	{
		"The New York Times",
		"nyt.example",
		[]string{"Media", "Marketing", "Sales"},
	},
	{"BBC", "bbc.example", []string{"Media", "Entertainment", "Marketing"}},
	{
		"Verizon",
		"verizon.example",
		[]string{"Telecommunications", "Engineering", "Sales"},
	},
	{
		"AT&T",
		"att.example",
		[]string{"Telecommunications", "Engineering", "Sales"},
	},
	{
		"SunPower",
		"sunpower.example",
		[]string{"Renewable Energy", "Engineering", "Sales"},
	},
	{
		"Vestas",
		"vestas.example",
		[]string{"Renewable Energy", "Engineering", "Operations"},
	},
	{
		"Deloitte",
		"deloitte.example",
		[]string{"Consulting", "Finance", "Operations"},
	},
	{"PwC", "pwc.example", []string{"Consulting", "Finance", "Data Science"}},
	{
		"Spotify",
		"spotify.example",
		[]string{"Engineering", "Product Management", "Media"},
	},
	{
		"Netflix",
		"netflix.example",
		[]string{"Entertainment", "Engineering", "Data Science"},
	},
	{
		"Shopify",
		"shopify.example",
		[]string{"Engineering", "Product Management", "Sales"},
	},
	{
		"Uber",
		"uber.example",
		[]string{"Engineering", "Operations", "Data Science"},
	},
	{"HSBC", "hsbc.example", []string{"Finance", "Data Science", "Sales"}},
	{
		"Barclays",
		"barclays.example",
		[]string{"Finance", "Technology", "Consulting"},
	},
	{
		"Philips",
		"philips.example",
		[]string{"Medical", "Engineering", "Design"},
	},
	{
		"GE Healthcare",
		"gehealthcare.example",
		[]string{"Medical", "Engineering", "Sales"},
	},
}

func generateHubSeedUsers(num int) []HubSeedUser {
	var hubSeedUsers []HubSeedUser
	for i := 0; i < num; i++ {
		name := fmt.Sprintf(
			"%s %s",
			firstNames[rand.Intn(len(firstNames))],
			lastNames[rand.Intn(len(lastNames))],
		)

		tier := hub.FreeHubUserTier
		if rand.Float32() < 0.5 {
			tier = hub.PaidHubUserTier
		}

		country := countries[rand.Intn(len(countries))]
		city := cities[country][rand.Intn(len(cities[country]))]
		pic := fmt.Sprintf("avatar%d.jpg", rand.Intn(17)+1)

		career := careerPaths[rand.Intn(len(careerPaths))]
		levels := rand.Intn(len(career.Steps)) + 1

		jobs := make([]Job, levels)

		for j := 0; j < levels; j++ {
			// TODO: possibleEmployers should be calculated in init
			var possibleEmployers []struct {
				Name    string
				Website string
			}
			for _, employer := range employers {
				for _, k := range employer.HiringInBroadAreas {
					if k == career.BroadArea {
						possibleEmployers = append(possibleEmployers, struct {
							Name    string
							Website string
						}{employer.Name, employer.Website})
						break
					}
				}
			}
			if len(possibleEmployers) == 0 {
				// Fallback in case no employer matches this career tag
				possibleEmployers = append(possibleEmployers, struct {
					Name    string
					Website string
				}{"Generic Company", "generic.example"})
			}
			employer := possibleEmployers[rand.Intn(len(possibleEmployers))]
			job := Job{
				Title:   career.Steps[j],
				Website: employer.Website,
			}
			jobs[j] = job
		}

		longBioOpts := jobBioMap[jobs[len(jobs)-1].Title]
		longBio := longBioOpts[rand.Intn(len(longBioOpts))]

		hubUser := HubSeedUser{
			Name:                   name,
			Handle:                 fmt.Sprintf("user%d", i),
			Email:                  fmt.Sprintf("user%d@example.com", i),
			Tier:                   tier,
			ResidentCountry:        country,
			ResidentCity:           city,
			PreferredLanguage:      "en",
			ShortBio:               jobs[len(jobs)-1].Title,
			LongBio:                longBio,
			ProfilePictureFilename: pic,
			Jobs:                   jobs,
			BroadArea:              career.BroadArea,
		}

		// Generate work history details for this user
		trackWorkHistory(&hubUser)

		hubSeedUsers = append(hubSeedUsers, hubUser)
	}
	return hubSeedUsers
}

// trackWorkHistory generates detailed work history items with dates for a user
func trackWorkHistory(user *HubSeedUser) {
	var prevStartDate time.Time

	// Initialize the work history items slice
	user.WorkHistoryItems = make([]WorkHistoryItem, 0, len(user.Jobs))

	for i := len(user.Jobs) - 1; i >= 0; i-- {
		job := user.Jobs[i]

		var startDateRaw time.Time
		var endDatePtr *time.Time

		if i == len(user.Jobs)-1 {
			// Last job is current job
			endDatePtr = nil
			randYears := rand.Intn(7) + 1
			startDateRaw = time.Now().AddDate(-randYears, 0, 0)
			prevStartDate = startDateRaw
		} else {
			// Assuming a 30-90 day gap exists between jobs
			gapDays := rand.Intn(60) + 30
			endDate := prevStartDate.AddDate(0, 0, -gapDays)
			endDateCopy := endDate // Create a copy to avoid pointer issues
			endDatePtr = &endDateCopy

			numberOfYears := rand.Intn(7) + 1
			startDateRaw = endDate.AddDate(-numberOfYears, 0, 0)
			prevStartDate = startDateRaw
		}

		// Find employer name from website
		var employerName string
		for _, employer := range employers {
			if employer.Website == job.Website {
				employerName = employer.Name
				break
			}
		}

		// Create and add the work history item
		workItem := WorkHistoryItem{
			EmployerID:   job.Website,
			EmployerName: employerName,
			StartDate:    startDateRaw,
			EndDate:      endDatePtr,
			JobTitle:     job.Title,
			Description:  job.Title, // Using job title as default description
		}

		user.WorkHistoryItems = append(user.WorkHistoryItems, workItem)
	}
}
