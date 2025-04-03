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
	"Renewable Energy": {
		{
			Type:        common.Patent,
			Title:       "Solar Energy Optimization",
			Description: "Patent for advanced solar panel efficiency system",
			URL:         "https://patents.example.com/solar-optimization",
		},
		{
			Type:        common.Patent,
			Title:       "Wind Turbine Design",
			Description: "Patent for innovative wind turbine architecture",
			URL:         "https://patents.example.com/wind-turbine",
		},
		{
			Type:        common.Publication,
			Title:       "Grid-Scale Energy Storage",
			Description: "Research on advanced battery technologies",
			URL:         "https://journals.example.com/energy-storage",
		},
		{
			Type:        common.Publication,
			Title:       "Smart Grid Integration",
			Description: "Research on renewable energy grid integration",
			URL:         "https://journals.example.com/smart-grid",
		},
		{
			Type:        common.Certification,
			Title:       "Renewable Energy Professional",
			Description: "Certification in renewable energy systems",
			URL:         "https://renewable.example.com/certification",
		},
		{
			Type:        common.Certification,
			Title:       "Solar System Designer",
			Description: "Professional certification in solar system design",
			URL:         "https://solar.example.com/certification",
		},
	},
	"Telecommunications": {
		{
			Type:        common.Patent,
			Title:       "5G Network Architecture",
			Description: "Patent for advanced 5G network design",
			URL:         "https://patents.example.com/5g-network",
		},
		{
			Type:        common.Patent,
			Title:       "Network Security Protocol",
			Description: "Patent for secure communication system",
			URL:         "https://patents.example.com/network-security",
		},
		{
			Type:        common.Publication,
			Title:       "6G Technology Research",
			Description: "Research on next-generation wireless",
			URL:         "https://journals.example.com/6g-research",
		},
		{
			Type:        common.Publication,
			Title:       "IoT Network Optimization",
			Description: "Research on IoT communication protocols",
			URL:         "https://journals.example.com/iot-network",
		},
		{
			Type:        common.Certification,
			Title:       "Cisco Network Professional",
			Description: "Advanced certification in network architecture",
			URL:         "https://cisco.example.com/certification",
		},
		{
			Type:        common.Certification,
			Title:       "Wireless Systems Expert",
			Description: "Certification in wireless communications",
			URL:         "https://wireless.example.com/certification",
		},
	},
	"Aerospace": {
		{
			Type:        common.Patent,
			Title:       "Aircraft Propulsion System",
			Description: "Patent for efficient propulsion technology",
			URL:         "https://patents.example.com/propulsion",
		},
		{
			Type:        common.Patent,
			Title:       "Satellite Communication System",
			Description: "Patent for advanced satellite communications",
			URL:         "https://patents.example.com/satellite-comm",
		},
		{
			Type:        common.Publication,
			Title:       "Advanced Materials in Aviation",
			Description: "Research on aerospace materials",
			URL:         "https://journals.example.com/aero-materials",
		},
		{
			Type:        common.Publication,
			Title:       "Space Navigation Systems",
			Description: "Research on spacecraft navigation",
			URL:         "https://journals.example.com/space-nav",
		},
		{
			Type:        common.Certification,
			Title:       "Aerospace Systems Engineer",
			Description: "Certification in aerospace systems",
			URL:         "https://aerospace.example.com/certification",
		},
		{
			Type:        common.Certification,
			Title:       "Aviation Safety Expert",
			Description: "Certification in aviation safety",
			URL:         "https://aviation.example.com/certification",
		},
	},
	"Human Resources": {
		{
			Type:        common.Publication,
			Title:       "Modern Workforce Management",
			Description: "Research on employee engagement strategies",
			URL:         "https://journals.example.com/workforce",
		},
		{
			Type:        common.Publication,
			Title:       "Digital HR Transformation",
			Description: "Research on HR technology integration",
			URL:         "https://journals.example.com/digital-hr",
		},
		{
			Type:        common.Publication,
			Title:       "Remote Work Best Practices",
			Description: "Research on distributed team management",
			URL:         "https://journals.example.com/remote-work",
		},
		{
			Type:        common.Certification,
			Title:       "Senior HR Professional",
			Description: "Advanced certification in HR management",
			URL:         "https://hrci.example.com/certification",
		},
		{
			Type:        common.Certification,
			Title:       "Talent Management Specialist",
			Description: "Certification in talent acquisition",
			URL:         "https://shrm.example.com/certification",
		},
		{
			Type:        common.Certification,
			Title:       "Compensation and Benefits Expert",
			Description: "Certification in compensation management",
			URL:         "https://compensation.example.com/certification",
		},
	},
	"Design": {
		{
			Type:        common.Patent,
			Title:       "UI/UX Interaction System",
			Description: "Patent for innovative user interface design",
			URL:         "https://patents.example.com/ui-system",
		},
		{
			Type:        common.Publication,
			Title:       "Design Systems at Scale",
			Description: "Research on enterprise design systems",
			URL:         "https://journals.example.com/design-systems",
		},
		{
			Type:        common.Publication,
			Title:       "Accessible Design Patterns",
			Description: "Research on inclusive design principles",
			URL:         "https://journals.example.com/accessible-design",
		},
		{
			Type:        common.Publication,
			Title:       "Future of Digital Product Design",
			Description: "Research on emerging design trends",
			URL:         "https://journals.example.com/future-design",
		},
		{
			Type:        common.Certification,
			Title:       "UX Research Professional",
			Description: "Certification in user research",
			URL:         "https://uxr.example.com/certification",
		},
		{
			Type:        common.Certification,
			Title:       "Product Design Expert",
			Description: "Certification in product design",
			URL:         "https://design.example.com/certification",
		},
	},
	"Marketing": {
		{
			Type:        common.Publication,
			Title:       "Digital Marketing Innovation",
			Description: "Research on emerging digital marketing trends",
			URL:         "https://journals.example.com/digital-marketing",
		},
		{
			Type:        common.Publication,
			Title:       "Social Media Strategy",
			Description: "Research on social media marketing effectiveness",
			URL:         "https://journals.example.com/social-media",
		},
		{
			Type:        common.Publication,
			Title:       "Marketing Analytics",
			Description: "Research on data-driven marketing approaches",
			URL:         "https://journals.example.com/marketing-analytics",
		},
		{
			Type:        common.Certification,
			Title:       "Digital Marketing Professional",
			Description: "Advanced certification in digital marketing",
			URL:         "https://marketing.example.com/certification",
		},
		{
			Type:        common.Certification,
			Title:       "Marketing Analytics Expert",
			Description: "Certification in marketing data analysis",
			URL:         "https://analytics.example.com/certification",
		},
		{
			Type:        common.Certification,
			Title:       "Brand Strategy Professional",
			Description: "Certification in brand management",
			URL:         "https://brand.example.com/certification",
		},
	},
	"Construction": {
		{
			Type:        common.Patent,
			Title:       "Construction Automation System",
			Description: "Patent for automated construction processes",
			URL:         "https://patents.example.com/construction-automation",
		},
		{
			Type:        common.Patent,
			Title:       "Smart Building Technology",
			Description: "Patent for intelligent building systems",
			URL:         "https://patents.example.com/smart-building",
		},
		{
			Type:        common.Publication,
			Title:       "Sustainable Construction Methods",
			Description: "Research on eco-friendly building practices",
			URL:         "https://journals.example.com/sustainable-construction",
		},
		{
			Type:        common.Publication,
			Title:       "Construction Project Management",
			Description: "Research on efficient project delivery",
			URL:         "https://journals.example.com/construction-management",
		},
		{
			Type:        common.Certification,
			Title:       "Construction Management Professional",
			Description: "Advanced certification in construction management",
			URL:         "https://construction.example.com/certification",
		},
		{
			Type:        common.Certification,
			Title:       "Sustainable Building Expert",
			Description: "Certification in green building practices",
			URL:         "https://green.example.com/certification",
		},
	},
	"Real Estate": {
		{
			Type:        common.Publication,
			Title:       "Real Estate Market Analysis",
			Description: "Research on market trends and forecasting",
			URL:         "https://journals.example.com/real-estate-market",
		},
		{
			Type:        common.Publication,
			Title:       "Property Technology Innovation",
			Description: "Research on PropTech advancements",
			URL:         "https://journals.example.com/proptech",
		},
		{
			Type:        common.Publication,
			Title:       "Commercial Real Estate Strategies",
			Description: "Research on commercial property management",
			URL:         "https://journals.example.com/commercial-real-estate",
		},
		{
			Type:        common.Certification,
			Title:       "Real Estate Investment Analyst",
			Description: "Certification in real estate investment",
			URL:         "https://realestate.example.com/certification",
		},
		{
			Type:        common.Certification,
			Title:       "Property Management Professional",
			Description: "Certification in property management",
			URL:         "https://property.example.com/certification",
		},
		{
			Type:        common.Certification,
			Title:       "Commercial Real Estate Expert",
			Description: "Certification in commercial real estate",
			URL:         "https://commercial.example.com/certification",
		},
	},
	"Hospitality": {
		{
			Type:        common.Publication,
			Title:       "Hospitality Service Innovation",
			Description: "Research on service excellence",
			URL:         "https://journals.example.com/hospitality-innovation",
		},
		{
			Type:        common.Publication,
			Title:       "Hotel Management Strategies",
			Description: "Research on hotel operations",
			URL:         "https://journals.example.com/hotel-management",
		},
		{
			Type:        common.Publication,
			Title:       "Guest Experience Design",
			Description: "Research on customer experience optimization",
			URL:         "https://journals.example.com/guest-experience",
		},
		{
			Type:        common.Certification,
			Title:       "Hospitality Management Professional",
			Description: "Advanced certification in hospitality management",
			URL:         "https://hospitality.example.com/certification",
		},
		{
			Type:        common.Certification,
			Title:       "Hotel Operations Expert",
			Description: "Certification in hotel operations",
			URL:         "https://hotel.example.com/certification",
		},
		{
			Type:        common.Certification,
			Title:       "Food and Beverage Director",
			Description: "Certification in F&B management",
			URL:         "https://foodbeverage.example.com/certification",
		},
	},
	"Automotive": {
		{
			Type:        common.Patent,
			Title:       "Electric Vehicle Technology",
			Description: "Patent for EV powertrain system",
			URL:         "https://patents.example.com/ev-technology",
		},
		{
			Type:        common.Patent,
			Title:       "Autonomous Driving System",
			Description: "Patent for self-driving technology",
			URL:         "https://patents.example.com/autonomous-driving",
		},
		{
			Type:        common.Publication,
			Title:       "Vehicle Safety Systems",
			Description: "Research on automotive safety",
			URL:         "https://journals.example.com/vehicle-safety",
		},
		{
			Type:        common.Publication,
			Title:       "Future of Mobility",
			Description: "Research on transportation trends",
			URL:         "https://journals.example.com/future-mobility",
		},
		{
			Type:        common.Certification,
			Title:       "Automotive Engineering Professional",
			Description: "Advanced certification in automotive engineering",
			URL:         "https://automotive.example.com/certification",
		},
		{
			Type:        common.Certification,
			Title:       "Electric Vehicle Systems Expert",
			Description: "Certification in EV systems",
			URL:         "https://ev.example.com/certification",
		},
	},
	"Product Management": {
		{
			Type:        common.Publication,
			Title:       "Product Strategy Innovation",
			Description: "Research on product development methodologies",
			URL:         "https://journals.example.com/product-strategy",
		},
		{
			Type:        common.Publication,
			Title:       "User-Centered Design",
			Description: "Research on product design principles",
			URL:         "https://journals.example.com/user-centered-design",
		},
		{
			Type:        common.Publication,
			Title:       "Product Analytics",
			Description: "Research on product metrics and KPIs",
			URL:         "https://journals.example.com/product-analytics",
		},
		{
			Type:        common.Certification,
			Title:       "Product Management Professional",
			Description: "Advanced certification in product management",
			URL:         "https://product.example.com/certification",
		},
		{
			Type:        common.Certification,
			Title:       "Agile Product Owner",
			Description: "Certification in agile product ownership",
			URL:         "https://agile.example.com/certification",
		},
		{
			Type:        common.Certification,
			Title:       "Product Analytics Expert",
			Description: "Certification in product analytics",
			URL:         "https://productanalytics.example.com/certification",
		},
	},
	"Education": {
		{
			Type:        common.Publication,
			Title:       "Educational Technology Innovation",
			Description: "Research on EdTech implementation",
			URL:         "https://journals.example.com/edtech",
		},
		{
			Type:        common.Publication,
			Title:       "Learning Analytics",
			Description: "Research on educational data analysis",
			URL:         "https://journals.example.com/learning-analytics",
		},
		{
			Type:        common.Publication,
			Title:       "Curriculum Development",
			Description: "Research on modern curriculum design",
			URL:         "https://journals.example.com/curriculum",
		},
		{
			Type:        common.Certification,
			Title:       "Educational Leadership",
			Description: "Advanced certification in education management",
			URL:         "https://education.example.com/certification",
		},
		{
			Type:        common.Certification,
			Title:       "Instructional Design Expert",
			Description: "Certification in course design",
			URL:         "https://instruction.example.com/certification",
		},
		{
			Type:        common.Certification,
			Title:       "Digital Learning Specialist",
			Description: "Certification in online education",
			URL:         "https://elearning.example.com/certification",
		},
	},
	"Operations": {
		{
			Type:        common.Patent,
			Title:       "Supply Chain Optimization",
			Description: "Patent for logistics optimization system",
			URL:         "https://patents.example.com/supply-chain",
		},
		{
			Type:        common.Publication,
			Title:       "Operations Analytics",
			Description: "Research on operational efficiency",
			URL:         "https://journals.example.com/operations-analytics",
		},
		{
			Type:        common.Publication,
			Title:       "Process Automation",
			Description: "Research on automated operations",
			URL:         "https://journals.example.com/process-automation",
		},
		{
			Type:        common.Certification,
			Title:       "Operations Management Professional",
			Description: "Advanced certification in operations",
			URL:         "https://operations.example.com/certification",
		},
		{
			Type:        common.Certification,
			Title:       "Supply Chain Expert",
			Description: "Certification in supply chain management",
			URL:         "https://supplychain.example.com/certification",
		},
		{
			Type:        common.Certification,
			Title:       "Lean Six Sigma Black Belt",
			Description: "Certification in process improvement",
			URL:         "https://sixsigma.example.com/certification",
		},
	},
	"Sales": {
		{
			Type:        common.Publication,
			Title:       "Sales Strategy Innovation",
			Description: "Research on modern sales methodologies",
			URL:         "https://journals.example.com/sales-strategy",
		},
		{
			Type:        common.Publication,
			Title:       "Digital Sales Transformation",
			Description: "Research on digital sales processes",
			URL:         "https://journals.example.com/digital-sales",
		},
		{
			Type:        common.Publication,
			Title:       "Sales Analytics",
			Description: "Research on sales performance metrics",
			URL:         "https://journals.example.com/sales-analytics",
		},
		{
			Type:        common.Certification,
			Title:       "Sales Management Professional",
			Description: "Advanced certification in sales management",
			URL:         "https://sales.example.com/certification",
		},
		{
			Type:        common.Certification,
			Title:       "Enterprise Sales Expert",
			Description: "Certification in enterprise sales",
			URL:         "https://enterprise.example.com/certification",
		},
		{
			Type:        common.Certification,
			Title:       "Sales Operations Specialist",
			Description: "Certification in sales operations",
			URL:         "https://salesops.example.com/certification",
		},
	},
	"Consulting": {
		{
			Type:        common.Publication,
			Title:       "Management Consulting Practices",
			Description: "Research on consulting methodologies",
			URL:         "https://journals.example.com/consulting-practices",
		},
		{
			Type:        common.Publication,
			Title:       "Digital Transformation Strategy",
			Description: "Research on business transformation",
			URL:         "https://journals.example.com/digital-transformation",
		},
		{
			Type:        common.Publication,
			Title:       "Business Analytics",
			Description: "Research on business intelligence",
			URL:         "https://journals.example.com/business-analytics",
		},
		{
			Type:        common.Certification,
			Title:       "Management Consulting Professional",
			Description: "Advanced certification in consulting",
			URL:         "https://consulting.example.com/certification",
		},
		{
			Type:        common.Certification,
			Title:       "Strategy Consulting Expert",
			Description: "Certification in strategic consulting",
			URL:         "https://strategy.example.com/certification",
		},
		{
			Type:        common.Certification,
			Title:       "Business Transformation Specialist",
			Description: "Certification in change management",
			URL:         "https://transformation.example.com/certification",
		},
	},
	"Media": {
		{
			Type:        common.Publication,
			Title:       "Digital Media Innovation",
			Description: "Research on media transformation",
			URL:         "https://journals.example.com/digital-media",
		},
		{
			Type:        common.Publication,
			Title:       "Content Strategy",
			Description: "Research on content development",
			URL:         "https://journals.example.com/content-strategy",
		},
		{
			Type:        common.Publication,
			Title:       "Media Analytics",
			Description: "Research on media metrics",
			URL:         "https://journals.example.com/media-analytics",
		},
		{
			Type:        common.Certification,
			Title:       "Digital Media Professional",
			Description: "Advanced certification in digital media",
			URL:         "https://media.example.com/certification",
		},
		{
			Type:        common.Certification,
			Title:       "Content Strategy Expert",
			Description: "Certification in content strategy",
			URL:         "https://content.example.com/certification",
		},
		{
			Type:        common.Certification,
			Title:       "Media Management Specialist",
			Description: "Certification in media management",
			URL:         "https://mediamanagement.example.com/certification",
		},
	},
	"Entertainment": {
		{
			Type:        common.Patent,
			Title:       "Interactive Entertainment System",
			Description: "Patent for interactive media platform",
			URL:         "https://patents.example.com/interactive-entertainment",
		},
		{
			Type:        common.Publication,
			Title:       "Digital Entertainment Trends",
			Description: "Research on entertainment industry",
			URL:         "https://journals.example.com/entertainment-trends",
		},
		{
			Type:        common.Publication,
			Title:       "Content Production Innovation",
			Description: "Research on production techniques",
			URL:         "https://journals.example.com/content-production",
		},
		{
			Type:        common.Certification,
			Title:       "Entertainment Production Professional",
			Description: "Advanced certification in production",
			URL:         "https://entertainment.example.com/certification",
		},
		{
			Type:        common.Certification,
			Title:       "Digital Media Producer",
			Description: "Certification in digital production",
			URL:         "https://producer.example.com/certification",
		},
		{
			Type:        common.Certification,
			Title:       "Content Development Expert",
			Description: "Certification in content development",
			URL:         "https://contentdev.example.com/certification",
		},
	},
	"Retail": {
		{
			Type:        common.Patent,
			Title:       "Retail Analytics System",
			Description: "Patent for retail intelligence platform",
			URL:         "https://patents.example.com/retail-analytics",
		},
		{
			Type:        common.Publication,
			Title:       "Digital Retail Innovation",
			Description: "Research on retail transformation",
			URL:         "https://journals.example.com/digital-retail",
		},
		{
			Type:        common.Publication,
			Title:       "Customer Experience Design",
			Description: "Research on retail experience",
			URL:         "https://journals.example.com/retail-experience",
		},
		{
			Type:        common.Certification,
			Title:       "Retail Management Professional",
			Description: "Advanced certification in retail",
			URL:         "https://retail.example.com/certification",
		},
		{
			Type:        common.Certification,
			Title:       "Digital Retail Expert",
			Description: "Certification in e-commerce",
			URL:         "https://ecommerce.example.com/certification",
		},
		{
			Type:        common.Certification,
			Title:       "Retail Operations Specialist",
			Description: "Certification in store operations",
			URL:         "https://retailops.example.com/certification",
		},
	},
	"Pharmaceuticals": {
		{
			Type:        common.Patent,
			Title:       "Drug Delivery System",
			Description: "Patent for pharmaceutical delivery",
			URL:         "https://patents.example.com/drug-delivery",
		},
		{
			Type:        common.Patent,
			Title:       "Pharmaceutical Formulation",
			Description: "Patent for drug formulation",
			URL:         "https://patents.example.com/drug-formulation",
		},
		{
			Type:        common.Publication,
			Title:       "Clinical Research Innovation",
			Description: "Research on clinical trials",
			URL:         "https://journals.example.com/clinical-research",
		},
		{
			Type:        common.Certification,
			Title:       "Pharmaceutical Research Professional",
			Description: "Advanced certification in pharma research",
			URL:         "https://pharma.example.com/certification",
		},
		{
			Type:        common.Certification,
			Title:       "Clinical Trial Manager",
			Description: "Certification in trial management",
			URL:         "https://clinical.example.com/certification",
		},
		{
			Type:        common.Certification,
			Title:       "Drug Safety Expert",
			Description: "Certification in pharmacovigilance",
			URL:         "https://drugsafety.example.com/certification",
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
			color.Magenta("no career path found for %s", user.Email)
			continue // Skip if no matching career path found
		}

		// Get possible achievements for this career
		achievements := careerAchievements[careerPath]
		if len(achievements) == 0 {
			color.Magenta(
				"no achievements found for %s %s",
				user.Email,
				careerPath,
			)
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
