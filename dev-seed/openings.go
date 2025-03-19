package main

import (
	"log"

	"github.com/fatih/color"
	"github.com/psankar/vetchi/typespec/common"
	"github.com/psankar/vetchi/typespec/employer"
)

// Track openings per company using company domain as key, openingID as value
var companyOpenings = make(map[string][]string)

// Track active openings per company using company domain as key, openingID as value
var activeOpenings = make(map[string][]string)

func createOpening(token string, req employer.CreateOpeningRequest) string {
	var resp employer.CreateOpeningResponse
	makeRequest("POST", "/employer/create-opening", token, req, &resp)
	return resp.OpeningID
}

func changeOpeningState(
	token string,
	openingID string,
	fromState, toState common.OpeningState,
) {
	req := employer.ChangeOpeningStateRequest{
		OpeningID: openingID,
		FromState: fromState,
		ToState:   toState,
	}
	makeRequest("POST", "/employer/change-opening-state", token, req, nil)
}

func createOpenings() {
	// Get tokens from the global map
	gryffindorVal, ok := employerSessionTokens.Load("admin@gryffindor.example")
	if !ok {
		log.Fatal("failed to get gryffindor token")
	}
	gryffindorToken := gryffindorVal.(string)

	hufflepuffVal, ok := employerSessionTokens.Load("admin@hufflepuff.example")
	if !ok {
		log.Fatal("failed to get hufflepuff token")
	}
	hufflepuffToken := hufflepuffVal.(string)

	ravenclawVal, ok := employerSessionTokens.Load("admin@ravenclaw.example")
	if !ok {
		log.Fatal("failed to get ravenclaw token")
	}
	ravenclawToken := ravenclawVal.(string)

	slytherinVal, ok := employerSessionTokens.Load("admin@slytherin.example")
	if !ok {
		log.Fatal("failed to get slytherin token")
	}
	slytherinToken := slytherinVal.(string)

	openings := []struct {
		domain string
		token  string
		req    employer.CreateOpeningRequest
	}{
		// Gryffindor openings
		{
			domain: "gryffindor.example",
			token:  gryffindorToken,
			req: employer.CreateOpeningRequest{
				Title:             "Senior Backend Engineer",
				Positions:         2,
				JD:                `Senior Backend Engineer with 10+ years experience in designing and implementing high-performance, scalable systems. Must have expert-level knowledge in at least one backend language such as Go, Java, or Python. Experience with distributed systems, microservices architecture, and cloud platforms (AWS, GCP, or Azure). Strong understanding of database technologies (SQL and NoSQL), messaging systems, and caching solutions. Proficiency in container technologies (Docker, Kubernetes) and CI/CD pipelines. Demonstrated experience in system design, performance optimization, and handling high-traffic applications. Must be able to lead technical teams and mentor junior engineers.`,
				Recruiter:         "hermione@gryffindor.example",
				HiringManager:     "harry@gryffindor.example",
				CostCenterName:    "UK Operations",
				OpeningType:       common.FullTimeOpening,
				YoeMin:            5,
				YoeMax:            10,
				MinEducationLevel: common.BachelorEducation,
				LocationTitles:    []string{"Chennai Diagon"},
				NewTags:           []string{"DevOps"},
			},
		},
		{
			domain: "gryffindor.example",
			token:  gryffindorToken,
			req: employer.CreateOpeningRequest{
				Title:             "Product Manager",
				Positions:         1,
				JD:                `Experienced Product Manager needed to drive product strategy and execution for our Irish expansion. You will work cross-functionally with engineering, design, marketing, and business teams to define and deliver innovative products that delight our customers. The ideal candidate has 8-15 years of experience leading product development for consumer-facing technology products, with demonstrated success in product launches and growth. Strong analytical skills and data-driven decision making are essential. Must excel at customer discovery, market analysis, roadmap development, and agile methodologies. MBA or equivalent experience preferred. Experience with international markets, particularly in Europe, is a significant plus. Must be an outstanding communicator with the ability to influence stakeholders at all levels and translate complex technical concepts for diverse audiences.`,
				Recruiter:         "hermione@gryffindor.example",
				HiringManager:     "ron@gryffindor.example",
				CostCenterName:    "Ireland Division",
				OpeningType:       common.FullTimeOpening,
				YoeMin:            8,
				YoeMax:            15,
				MinEducationLevel: common.MasterEducation,
				LocationTitles:    []string{"Chennai Diagon"},
				NewTags:           []string{"DevOps"},
			},
		},
		{
			domain: "gryffindor.example",
			token:  gryffindorToken,
			req: employer.CreateOpeningRequest{
				Title:             "DevOps Engineer",
				Positions:         2,
				JD:                `DevOps Engineer needed to design, implement, and maintain our global cloud infrastructure. You will be responsible for building and operating scalable, highly available systems while optimizing performance and cost. Requires 3-8 years of hands-on experience with AWS/GCP/Azure, infrastructure as code (Terraform, CloudFormation), container orchestration (Kubernetes, ECS), and CI/CD pipelines (Jenkins, GitHub Actions, CircleCI). Strong programming skills in Python, Go, or similar languages for automation. Experience with monitoring tools (Prometheus, Grafana), log management systems (ELK stack), and security best practices in cloud environments. Must be proficient in Linux systems administration and networking concepts. The ideal candidate has experience supporting microservices architectures in production environments and possesses excellent troubleshooting skills. Must be willing to participate in on-call rotations and thrive in a fast-paced, global team environment.`,
				Recruiter:         "ron@gryffindor.example",
				HiringManager:     "harry@gryffindor.example",
				CostCenterName:    "APAC Operations",
				OpeningType:       common.FullTimeOpening,
				YoeMin:            3,
				YoeMax:            8,
				MinEducationLevel: common.BachelorEducation,
				LocationTitles:    []string{"Chennai Diagon"},
				NewTags:           []string{"DevOps"},
			},
		},
		{
			domain: "gryffindor.example",
			token:  gryffindorToken,
			req: employer.CreateOpeningRequest{
				Title:             "Frontend Developer",
				Positions:         3,
				JD:                `Frontend Developer needed for our Canadian office to build exceptional user interfaces for our web applications. You'll work with our product and design teams to implement responsive, accessible, and high-performance user experiences. Required qualifications include 2-6 years of professional experience with React and modern JavaScript (ES6+). Strong understanding of state management (Redux, Context API), CSS preprocessors (SASS/LESS), and component-based architecture. Experience with TypeScript, responsive design, cross-browser compatibility, and optimization techniques for web performance. Familiarity with testing frameworks (Jest, React Testing Library) and CI/CD workflows. Bonus skills include experience with Next.js, GraphQL, WebSockets, and animation libraries. You should have a keen eye for detail, understand web accessibility standards (WCAG), and be passionate about creating intuitive user experiences. Must be comfortable in an agile environment and have excellent communication skills.`,
				Recruiter:         "hermione@gryffindor.example",
				HiringManager:     "ron@gryffindor.example",
				CostCenterName:    "Canada Business",
				OpeningType:       common.FullTimeOpening,
				YoeMin:            2,
				YoeMax:            6,
				MinEducationLevel: common.BachelorEducation,
				LocationTitles:    []string{"Chennai Diagon"},
				NewTags:           []string{"DevOps"},
			},
		},
		{
			domain: "gryffindor.example",
			token:  gryffindorToken,
			req: employer.CreateOpeningRequest{
				Title:             "Marketing Lead",
				Positions:         1,
				JD:                `Marketing Lead needed to develop and execute comprehensive global marketing strategies that drive business growth and brand awareness. The ideal candidate has 10-15 years of experience in technology marketing with a proven track record of successful global campaigns. You will oversee market research, competitive analysis, campaign development, and performance measurement across all channels (digital, social, events, PR). Must have expertise in digital marketing, content strategy, marketing automation, and analytics tools. Experience managing substantial marketing budgets and cross-functional teams is essential. Strong leadership skills with the ability to hire, develop, and inspire a diverse marketing team. MBA or equivalent experience preferred. International marketing experience is required, with demonstrated success in multiple geographic markets. Must be data-driven with excellent verbal and written communication skills, and have the ability to translate complex technical concepts for various audiences.`,
				Recruiter:         "ron@gryffindor.example",
				HiringManager:     "harry@gryffindor.example",
				CostCenterName:    "Global Marketing",
				OpeningType:       common.FullTimeOpening,
				YoeMin:            10,
				YoeMax:            15,
				MinEducationLevel: common.MasterEducation,
				LocationTitles:    []string{"Chennai Diagon"},
				NewTags:           []string{"DevOps"},
			},
		},

		// Hufflepuff openings
		{
			domain: "hufflepuff.example",
			token:  hufflepuffToken,
			req: employer.CreateOpeningRequest{
				Title:             "Software Architect",
				Positions:         1,
				JD:                `Software Architect needed for our Benelux operations to design and oversee the implementation of complex software solutions. You will be responsible for making high-level design choices, technical standards, and defining coding standards. Requires 8-15 years of software development experience with at least 5 years in architectural roles. Expert knowledge of software architecture patterns, distributed systems design, and microservices. Proficiency in multiple programming languages and platforms. Experience with cloud-native architectures (AWS/GCP/Azure), containerization, and serverless computing. Strong understanding of security practices, performance optimization, and scalability considerations. The ideal candidate has led architectural transformations and can balance technical requirements with business goals. Must possess excellent communication skills to articulate complex architectural concepts to both technical and non-technical stakeholders. Experience with international development teams and Benelux region market requirements is highly desirable.`,
				Recruiter:         "cedric@hufflepuff.example",
				HiringManager:     "newt@hufflepuff.example",
				CostCenterName:    "Benelux Operations",
				OpeningType:       common.FullTimeOpening,
				YoeMin:            8,
				YoeMax:            15,
				MinEducationLevel: common.MasterEducation,
				LocationTitles:    []string{"Chennai Cauldron"},
				NewTags:           []string{"DevOps"},
			},
		},
		{
			domain: "hufflepuff.example",
			token:  hufflepuffToken,
			req: employer.CreateOpeningRequest{
				Title:             "Data Scientist",
				Positions:         2,
				JD:                `Data Scientist needed for our Nordic team to extract valuable insights from large datasets and develop machine learning models to solve complex business problems. Requires 3-8 years of experience in data science or related field. Advanced degree in Computer Science, Statistics, Mathematics, or related field required. Expert proficiency in Python and R for data analysis and modeling. Experience with machine learning frameworks (TensorFlow, PyTorch, scikit-learn) and deep learning techniques. Strong background in statistical analysis, experimental design, and causal inference. Expertise in data visualization tools (Tableau, PowerBI) and SQL for data querying. Experience with big data technologies (Spark, Hadoop) and cloud-based data platforms. Must possess excellent communication skills to translate complex analyses into actionable business insights. The ideal candidate has experience working in cross-functional teams and delivering data-driven solutions in production environments. Knowledge of Nordic markets and business trends is a plus.`,
				Recruiter:         "nymphadora@hufflepuff.example",
				HiringManager:     "cedric@hufflepuff.example",
				CostCenterName:    "Nordic Division",
				OpeningType:       common.FullTimeOpening,
				YoeMin:            3,
				YoeMax:            8,
				MinEducationLevel: common.MasterEducation,
				LocationTitles:    []string{"Chennai Cauldron"},
				NewTags:           []string{"DevOps"},
			},
		},
		{
			domain: "hufflepuff.example",
			token:  hufflepuffToken,
			req: employer.CreateOpeningRequest{
				Title:             "Marketing Manager",
				Positions:         1,
				JD:                `Marketing Manager needed for EU operations to develop and execute strategic marketing plans that drive business growth across European markets. Requires 5-10 years of experience in marketing with a focus on B2B or B2C technology products. You will oversee market research, campaign development, budget management, and performance analysis. Must have expertise in digital marketing channels, content strategy, marketing automation tools, and marketing analytics. Experience managing integrated marketing campaigns across multiple European countries is essential. Proficiency with CRM systems, SEO/SEM strategies, and social media marketing. Strong understanding of the European regulatory environment for marketing and advertising. The ideal candidate has demonstrated success in building brand awareness and generating leads in competitive markets. Must possess excellent project management skills, multilingual capabilities (English plus at least one other European language), and the ability to work effectively with cross-functional teams. MBA or equivalent experience is preferred.`,
				Recruiter:         "newt@hufflepuff.example",
				HiringManager:     "nymphadora@hufflepuff.example",
				CostCenterName:    "EU Marketing",
				OpeningType:       common.FullTimeOpening,
				YoeMin:            5,
				YoeMax:            10,
				MinEducationLevel: common.BachelorEducation,
				LocationTitles:    []string{"Chennai Cauldron"},
				NewTags:           []string{"DevOps"},
			},
		},
		{
			domain: "hufflepuff.example",
			token:  hufflepuffToken,
			req: employer.CreateOpeningRequest{
				Title:             "Research Engineer",
				Positions:         3,
				JD:                `Research Engineer needed for our R&D team to explore, design, and prototype cutting-edge technology solutions. Requires 4-12 years of experience in research and development in computer science, engineering, or related fields. PhD or equivalent research experience preferred. Expert-level programming skills in multiple languages. Strong background in one or more specialized areas such as computer vision, natural language processing, robotics, augmented reality, distributed systems, or quantum computing. Experience translating research concepts into practical prototypes and production systems. Publication history in top-tier conferences or journals is a plus. Must be comfortable working in ambiguous problem spaces and developing novel approaches to unsolved challenges. The ideal candidate is intellectually curious, self-directed, and passionate about advancing the state of the art. Excellent collaboration skills for working with interdisciplinary teams and communicating complex technical concepts clearly. Patent experience is a plus.`,
				Recruiter:         "cedric@hufflepuff.example",
				HiringManager:     "newt@hufflepuff.example",
				CostCenterName:    "R&D Labs",
				OpeningType:       common.FullTimeOpening,
				YoeMin:            4,
				YoeMax:            12,
				MinEducationLevel: common.DoctorateEducation,
				LocationTitles:    []string{"Chennai Cauldron"},
				NewTags:           []string{"DevOps"},
			},
		},
		{
			domain: "hufflepuff.example",
			token:  hufflepuffToken,
			req: employer.CreateOpeningRequest{
				Title:             "Operations Manager",
				Positions:         1,
				JD:                `Operations Manager needed for EU administration to oversee day-to-day operational activities and drive continuous improvement. Requires 6-12 years of operations management experience, with at least 3 years in technology companies or multinational environments. You will be responsible for developing operational strategies, optimizing workflows, managing resources, and ensuring compliance with EU regulations. Strong background in process improvement methodologies (Lean, Six Sigma) and project management. Experience with ERP systems, data analysis, and operational reporting. Demonstrated ability to lead cross-functional teams and manage vendor relationships. The ideal candidate has expertise in supply chain management, facilities operations, and business continuity planning. Must possess excellent problem-solving skills and the ability to make data-driven decisions. Knowledge of EU labor laws, data protection regulations (GDPR), and business practices is essential. Multilingual capabilities (English plus at least one other European language) are highly desirable.`,
				Recruiter:         "nymphadora@hufflepuff.example",
				HiringManager:     "cedric@hufflepuff.example",
				CostCenterName:    "EU Admin",
				OpeningType:       common.FullTimeOpening,
				YoeMin:            6,
				YoeMax:            12,
				MinEducationLevel: common.BachelorEducation,
				LocationTitles:    []string{"Chennai Cauldron"},
				NewTags:           []string{"DevOps"},
			},
		},

		// Ravenclaw openings
		{
			domain: "ravenclaw.example",
			token:  ravenclawToken,
			req: employer.CreateOpeningRequest{
				Title:             "Technical Lead",
				Positions:         1,
				JD:                `Technical Lead needed for our APAC headquarters to provide technical leadership and direction for development teams. Requires 8-15 years of software development experience with a proven track record of leading engineering teams. Expert-level proficiency in multiple programming languages and technology stacks. Strong architectural skills with experience designing scalable, resilient systems. Deep understanding of software development methodologies, CI/CD practices, and quality assurance processes. Experience mentoring junior developers and conducting technical interviews. The ideal candidate has worked on complex projects across multiple domains and technologies. Must possess excellent problem-solving abilities and stay current with emerging technologies and industry trends. Strong communication skills for collaborating with product managers, designers, and business stakeholders. Experience working in APAC markets and understanding regional technology adoption patterns is highly valuable. Must be comfortable working across multiple time zones and with distributed teams.`,
				Recruiter:         "luna@ravenclaw.example",
				HiringManager:     "filius@ravenclaw.example",
				CostCenterName:    "APAC Headquarters",
				OpeningType:       common.FullTimeOpening,
				YoeMin:            8,
				YoeMax:            15,
				MinEducationLevel: common.MasterEducation,
				LocationTitles:    []string{"Chennai Flourish"},
				NewTags:           []string{"DevOps"},
			},
		},
		{
			domain: "ravenclaw.example",
			token:  ravenclawToken,
			req: employer.CreateOpeningRequest{
				Title:             "Mobile Developer",
				Positions:         2,
				JD:                `Mobile Developer needed for our Japan team to design and build advanced applications for iOS and Android platforms. Requires 3-8 years of mobile development experience with proficiency in Swift/Objective-C for iOS or Kotlin/Java for Android. Experience with cross-platform frameworks (React Native, Flutter) is a plus. Strong understanding of mobile UI/UX principles, performance optimization, and memory management. Experience with RESTful APIs, local data storage, and offline functionality. Proficiency in building responsive layouts, animations, and custom UI components. Knowledge of mobile security best practices, app lifecycle management, and push notifications. Experience with mobile testing frameworks and CI/CD pipelines for mobile apps. The ideal candidate has published apps on the App Store or Google Play with high user ratings. Must stay current with platform updates and best practices. Understanding of Japanese market preferences and mobile usage patterns is highly desirable. Ability to work in a global team environment and communicate effectively with designers and backend developers.`,
				Recruiter:         "cho@ravenclaw.example",
				HiringManager:     "luna@ravenclaw.example",
				CostCenterName:    "Japan Operations",
				OpeningType:       common.FullTimeOpening,
				YoeMin:            3,
				YoeMax:            8,
				MinEducationLevel: common.BachelorEducation,
				LocationTitles:    []string{"Chennai Flourish"},
				NewTags:           []string{"DevOps"},
			},
		},
		{
			domain: "ravenclaw.example",
			token:  ravenclawToken,
			req: employer.CreateOpeningRequest{
				Title:             "QA Engineer",
				Positions:         2,
				JD:                `QA Engineer needed for our Korean team to ensure the quality and reliability of our products. Requires 2-6 years of experience in software quality assurance with a strong focus on automation testing. Proficiency in test automation frameworks (Selenium, Cypress, Appium) and programming languages (Python, Java, JavaScript). Experience with API testing tools (Postman, RestAssured) and performance testing (JMeter, Gatling). Strong understanding of testing methodologies, test planning, and defect management processes. Experience with continuous integration tools and implementing testing in CI/CD pipelines. Knowledge of database testing, security testing, and accessibility testing is a plus. The ideal candidate has experience testing complex applications across web and mobile platforms. Must possess excellent analytical and problem-solving skills, with strong attention to detail. Ability to communicate effectively with developers to resolve issues. Understanding of Korean language and local market requirements is advantageous. Must be comfortable working in an agile environment with rapidly changing priorities.`,
				Recruiter:         "filius@ravenclaw.example",
				HiringManager:     "cho@ravenclaw.example",
				CostCenterName:    "Korea Division",
				OpeningType:       common.FullTimeOpening,
				YoeMin:            2,
				YoeMax:            6,
				MinEducationLevel: common.BachelorEducation,
				LocationTitles:    []string{"Chennai Flourish"},
				NewTags:           []string{"DevOps"},
			},
		},
		{
			domain: "ravenclaw.example",
			token:  ravenclawToken,
			req: employer.CreateOpeningRequest{
				Title:             "Full Stack Developer",
				Positions:         3,
				JD:                `Full Stack Developer needed for our India operations to design and implement end-to-end solutions across frontend and backend systems. Requires 4-10 years of experience in full stack development with expertise in modern JavaScript frameworks (React, Angular, or Vue) and backend technologies (Node.js, Python, Java, or Go). Strong proficiency in database design and optimization (SQL and NoSQL). Experience with RESTful APIs, GraphQL, and microservices architecture. Solid understanding of web fundamentals (HTML5, CSS3, JavaScript ES6+) and responsive design principles. Proficiency with version control systems, testing methodologies, and CI/CD pipelines. Knowledge of cloud platforms (AWS, GCP, Azure) and containerization (Docker, Kubernetes). The ideal candidate has experience building scalable, production-grade applications and working in agile environments. Must possess excellent problem-solving skills and the ability to learn new technologies quickly. Understanding of Indian tech ecosystem and local market requirements is a plus. Strong teamwork and communication skills are essential.`,
				Recruiter:         "luna@ravenclaw.example",
				HiringManager:     "filius@ravenclaw.example",
				CostCenterName:    "India Operations",
				OpeningType:       common.FullTimeOpening,
				YoeMin:            4,
				YoeMax:            10,
				MinEducationLevel: common.BachelorEducation,
				LocationTitles:    []string{"Chennai Flourish"},
				NewTags:           []string{"DevOps"},
			},
		},
		{
			domain: "ravenclaw.example",
			token:  ravenclawToken,
			req: employer.CreateOpeningRequest{
				Title:             "Solutions Architect",
				Positions:         1,
				JD:                `Solutions Architect needed for Middle East expansion to design comprehensive technical solutions that address complex business challenges. Requires 10-15 years of experience in technology consulting, systems integration, or solution architecture. Expert-level knowledge of enterprise architecture frameworks, cloud platforms (AWS, GCP, Azure), and integration patterns. Experience designing solutions involving multiple technologies, platforms, and vendors. Strong understanding of security, compliance, and regulatory requirements in the Middle East region. Proficiency in creating architecture diagrams, technical specifications, and implementation roadmaps. Experience with large-scale data management, analytics solutions, and modernization of legacy systems. The ideal candidate has led digital transformation initiatives and can effectively translate business requirements into technical solutions. Must possess exceptional communication skills to present complex technical concepts to business stakeholders. Ability to build relationships with clients and partners at senior levels. Experience working in the Middle East market and understanding regional business practices is highly valuable. Some travel within the region required.`,
				Recruiter:         "cho@ravenclaw.example",
				HiringManager:     "luna@ravenclaw.example",
				CostCenterName:    "Middle East Division",
				OpeningType:       common.FullTimeOpening,
				YoeMin:            10,
				YoeMax:            15,
				MinEducationLevel: common.MasterEducation,
				LocationTitles:    []string{"Chennai Flourish"},
				NewTags:           []string{"DevOps"},
			},
		},

		// Slytherin openings
		{
			domain: "slytherin.example",
			token:  slytherinToken,
			req: employer.CreateOpeningRequest{
				Title:             "Engineering Manager",
				Positions:         1,
				JD:                `Engineering Manager needed for DACH region to lead and grow engineering teams while delivering high-quality software products. Requires 8-15 years of software engineering experience with at least 5 years in management roles. Strong technical background with hands-on experience in software development and system design. Exceptional leadership skills with proven ability to recruit, develop, and retain top engineering talent. Experience managing multiple engineering teams and coordinating complex projects. Strong understanding of agile methodologies, engineering best practices, and technical debt management. Proficiency in resource planning, performance management, and cross-functional collaboration. The ideal candidate has experience scaling engineering organizations during periods of rapid growth. Must possess excellent communication skills and the ability to balance technical excellence with business priorities. Experience working in German-speaking markets (Germany, Austria, Switzerland) and familiarity with regional business practices is highly valuable. Fluency in English required, German language skills strongly preferred.`,
				Recruiter:         "draco@slytherin.example",
				HiringManager:     "severus@slytherin.example",
				CostCenterName:    "DACH Operations",
				OpeningType:       common.FullTimeOpening,
				YoeMin:            8,
				YoeMax:            15,
				MinEducationLevel: common.MasterEducation,
				LocationTitles:    []string{"Chennai Ollivanders"},
				NewTags:           []string{"DevOps"},
			},
		},
		{
			domain: "slytherin.example",
			token:  slytherinToken,
			req: employer.CreateOpeningRequest{
				Title:             "Security Engineer",
				Positions:         2,
				JD:                `Security Engineer needed for our French team to design, implement, and maintain security measures that protect our systems and data. Requires 5-10 years of experience in information security with deep expertise in application security, network security, and cloud security. Strong understanding of security principles, threat modeling, and risk assessment methodologies. Experience with security tools and technologies (SIEM, WAF, IDS/IPS, vulnerability scanners). Proficiency in security automation and implementing security in CI/CD pipelines. Knowledge of secure coding practices, penetration testing, and incident response. Certifications such as CISSP, CEH, or OSCP are highly desired. Experience with compliance frameworks (ISO 27001, SOC 2, GDPR) and implementing security controls. The ideal candidate has experience securing complex, distributed systems and cloud environments. Must possess excellent problem-solving skills and the ability to balance security requirements with business needs. Understanding of the European security landscape and French data protection regulations is essential. Fluency in English required, French language skills strongly preferred.`,
				Recruiter:         "severus@slytherin.example",
				HiringManager:     "horace@slytherin.example",
				CostCenterName:    "France Division",
				OpeningType:       common.FullTimeOpening,
				YoeMin:            5,
				YoeMax:            10,
				MinEducationLevel: common.BachelorEducation,
				LocationTitles:    []string{"Chennai Ollivanders"},
				NewTags:           []string{"DevOps"},
			},
		},
		{
			domain: "slytherin.example",
			token:  slytherinToken,
			req: employer.CreateOpeningRequest{
				Title:             "Cloud Engineer",
				Positions:         2,
				JD:                `Cloud Engineer needed for Southern Europe operations to design, implement, and manage our cloud infrastructure. Requires 3-8 years of experience with major cloud platforms (AWS, GCP, Azure) and infrastructure as code tools (Terraform, CloudFormation, Pulumi). Strong understanding of cloud architecture principles, networking, and security best practices. Experience with container orchestration platforms (Kubernetes, ECS) and serverless architectures. Proficiency in scripting and automation (Python, Bash, PowerShell). Knowledge of monitoring, logging, and observability solutions for cloud environments. Experience optimizing cloud costs and implementing FinOps practices. Certification in one or more cloud platforms is highly desired. The ideal candidate has experience designing and operating large-scale production environments in the cloud. Must possess excellent problem-solving skills and the ability to work in fast-paced environments. Understanding of Southern European business requirements and regional cloud infrastructure considerations is a plus. Must be comfortable participating in on-call rotations and handling production incidents efficiently.`,
				Recruiter:         "horace@slytherin.example",
				HiringManager:     "draco@slytherin.example",
				CostCenterName:    "Southern Europe",
				OpeningType:       common.FullTimeOpening,
				YoeMin:            3,
				YoeMax:            8,
				MinEducationLevel: common.BachelorEducation,
				LocationTitles:    []string{"Chennai Ollivanders"},
				NewTags:           []string{"DevOps"},
			},
		},
		{
			domain: "slytherin.example",
			token:  slytherinToken,
			req: employer.CreateOpeningRequest{
				Title:             "Project Manager",
				Positions:         1,
				JD:                `Project Manager needed for EU special projects to plan, execute, and deliver complex technology initiatives across European markets. Requires 6-12 years of project management experience in technology or similar fast-paced environments. PMP or Prince2 certification preferred. Strong expertise in project management methodologies (Agile, Waterfall, Hybrid) and project management tools. Demonstrated ability to manage multiple stakeholders, resolve conflicts, and drive projects to successful completion. Experience managing budgets, resources, timelines, and scope for large-scale projects. Strong risk management skills with the ability to identify, assess, and mitigate project risks. Experience working with distributed teams across multiple countries and time zones. The ideal candidate has managed cross-functional projects involving engineering, product, design, and business teams. Must possess excellent communication, negotiation, and leadership skills. Understanding of EU business practices, regional differences, and regulatory requirements is essential. Fluency in English required, proficiency in one or more additional European languages is a significant advantage.`,
				Recruiter:         "draco@slytherin.example",
				HiringManager:     "severus@slytherin.example",
				CostCenterName:    "EU Projects",
				OpeningType:       common.FullTimeOpening,
				YoeMin:            6,
				YoeMax:            12,
				MinEducationLevel: common.BachelorEducation,
				LocationTitles:    []string{"Chennai Ollivanders"},
				NewTags:           []string{"DevOps"},
			},
		},
		{
			domain: "slytherin.example",
			token:  slytherinToken,
			req: employer.CreateOpeningRequest{
				Title:             "Research Scientist",
				Positions:         2,
				JD:                `Research Scientist needed for our continental R&D team to advance the state of the art in key technology areas. Requires 5-12 years of research experience in computer science, engineering, or related fields with a PhD or equivalent research experience. Expertise in one or more specialized areas such as machine learning, artificial intelligence, computer vision, natural language processing, or computational biology. Strong publication record in top-tier conferences or journals. Experience translating research into practical applications and working with engineering teams to implement research breakthroughs. Proficiency in programming languages commonly used in research (Python, C++, R). Experience with research tools, frameworks, and methodologies specific to your area of expertise. The ideal candidate has a proven track record of innovative research with real-world impact. Must possess excellent analytical thinking, problem-solving skills, and the ability to work in ambiguous problem spaces. Strong communication skills for presenting research findings to technical and non-technical audiences. Experience collaborating with academic institutions and participating in research communities across Europe is highly valuable.`,
				Recruiter:         "severus@slytherin.example",
				HiringManager:     "horace@slytherin.example",
				CostCenterName:    "Continental R&D",
				OpeningType:       common.FullTimeOpening,
				YoeMin:            5,
				YoeMax:            12,
				MinEducationLevel: common.DoctorateEducation,
				LocationTitles:    []string{"Chennai Ollivanders"},
				NewTags:           []string{"DevOps"},
			},
		},
	}

	for _, opening := range openings {
		openingID := createOpening(opening.token, opening.req)
		color.Green("Created opening %s for %s", openingID, opening.domain)
		// Track openings by domain
		companyOpenings[opening.domain] = append(
			companyOpenings[opening.domain],
			openingID,
		)
	}

	// Publish first two openings for each company
	for domain, openings := range companyOpenings {
		employerTokenRaw, ok := employerSessionTokens.Load("admin@" + domain)
		if !ok {
			log.Fatalf("failed to get employer token for %s", domain)
		}
		employerToken, ok := employerTokenRaw.(string)
		if !ok {
			log.Fatalf("failed to cast employer token for %s", domain)
		}

		for i := 0; i < 2 && i < len(openings); i++ {
			changeOpeningState(
				employerToken,
				openings[i],
				common.DraftOpening,
				common.ActiveOpening,
			)
			color.Green("Published opening %s for %s", openings[i], domain)
			activeOpenings[domain] = append(activeOpenings[domain], openings[i])
		}

	}
}
