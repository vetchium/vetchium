package main

import (
	"github.com/vetchium/vetchium/typespec/common"
	"github.com/vetchium/vetchium/typespec/hub"
)

var broadAreaPostsMap map[string][]hub.AddPostRequest

func init() {
	broadAreaPostsMap = map[string][]hub.AddPostRequest{}

	broadAreaPostsMap["Engineering"] = []hub.AddPostRequest{
		{
			Content: `Just solved a critical production bug by turning it off and on again.

They say "debugging is like being the detective in a crime movie where you are also the murderer."

Today, I cleared my own name.`,
			NewTags: []common.VTagName{
				common.VTagName("Software Engineer"),
				common.VTagName("Blessed"),
			},
		},
		{
			Content: `Woke up at 4AM. Did yoga. Built a microservice. Deployed to production. Got paged. Fixed it. Broke it again. Blamed Jenkins. Meditated. Asked ChatGPT. Took credit.

Remember: It's not about how many bugs you fix, it's about how confidently you explain why they aren't your fault.`,
			NewTags: []common.VTagName{
				common.VTagName("Leadership"),
				common.VTagName("Senior Engineer"),
				common.VTagName("Fake It Till You Deploy It"),
			},
		},
		{
			Content: `Dear Recruiters,
When you say "We're like a startup, but with the stability of a big company,"
Do you mean: chaotic, underpaid, but with more meetings?

Asking for a friend who's crying into his YAML files.`,
			NewTags: []common.VTagName{
				common.VTagName("Life Lessons"),
				common.VTagName("Points to Ponder"),
			},
		},
		{
			Content: `Jenkins failed. Docker misbehaved. Kubernetes restarted everything.
I asked for logs, it gave me 12,000 lines of vibes.

Anyway, I'm a backend engineer now, but spiritually, I'm frontend. I just pretend everything works.`,
			NewTags: []common.VTagName{
				common.VTagName("Vibe Coding"),
			},
		},
		{
			Content: `As a $MINORITY in Tech, My big contributions to Tech are giving keynote talks about improving $MINORITY contributions to Tech, while somehow sidelining the actual $MINORITY workers doing real Technology work. Evangelism beats Engineering.`,
			NewTags: []common.VTagName{
				common.VTagName("DEI"),
				common.VTagName("He/Him/HisMajesty"),
				common.VTagName("TechBrosSuck"),
			},
		},
		{
			Content: "Vibe Coding is meh. Vibe Debugging, that would be a killer.",
			NewTags: []common.VTagName{
				common.VTagName("Vibe Coding"),
				common.VTagName("Artificial Intelligence"),
			},
		},
	}

	broadAreaPostsMap["Finance"] = []hub.AddPostRequest{
		{
			Content: `My portfolio went up 0.01% today. Already ordered the yacht. Financial freedom, baby!`,
			NewTags: []common.VTagName{
				common.VTagName("Stonks"),
				common.VTagName("Investing"),
				common.VTagName("FinanceHumor"),
			},
		},
		{
			Content: `Explained "synergy" using only Excel functions. Client was impressed. I think. They just stared blankly for 5 minutes.`,
			NewTags: []common.VTagName{
				common.VTagName("ExcelWizard"),
				common.VTagName("ConsultingLife"),
				common.VTagName("Finance"),
			},
		},
		{
			Content: `Survived another earnings call where the only thing higher than projections was my blood pressure.`,
			NewTags: []common.VTagName{
				common.VTagName("CorporateLife"),
				common.VTagName("Finance"),
				common.VTagName("StressEating"),
			},
		},
		{
			Content: `"Budgeting" is just telling your money where to go instead of wondering where it went... probably to coffee and emergency snacks.`,
			NewTags: []common.VTagName{
				common.VTagName("PersonalFinance"),
				common.VTagName("Adulting"),
				common.VTagName("Finance"),
			},
		},
		{
			Content: `Attended a webinar on "Blockchain for Beginners." Still a beginner. But now I'm a confused beginner with a certificate.`,
			NewTags: []common.VTagName{
				common.VTagName("Crypto"),
				common.VTagName("FinTech"),
				common.VTagName("ContinuousLearning"),
			},
		},
		{
			Content: `The market is volatile. My mood is volatile. Coincidence? I think not.`,
			NewTags: []common.VTagName{
				common.VTagName("MarketSwings"),
				common.VTagName("EmotionalInvesting"),
				common.VTagName("Finance"),
			},
		},
		{
			Content: `Built a complex financial model. It accurately predicts I need more coffee.`,
			NewTags: []common.VTagName{
				common.VTagName("DataAnalysis"),
				common.VTagName("Finance"),
				common.VTagName("CaffeinePowered"),
			},
		},
		{
			Content: `They say money can't buy happiness, but have you ever seen someone sad on a jet ski made of spreadsheets? Didn't think so.`,
			NewTags: []common.VTagName{
				common.VTagName("WorkLifeBalance"),
				common.VTagName("FinanceHumor"),
				common.VTagName("Goals"),
			},
		},
		{
			Content: `Risk assessment complete: The biggest risk is running out of coffee before the market opens.`,
			NewTags: []common.VTagName{
				common.VTagName("RiskManagement"),
				common.VTagName("FinanceLife"),
				common.VTagName("Priorities"),
			},
		},
		{
			Content: `Asked my AI assistant for stock tips. It suggested investing in "cat memes." Tempting, honestly.`,
			NewTags: []common.VTagName{
				common.VTagName("AI"),
				common.VTagName("Investing"),
				common.VTagName("FutureOfFinance"),
			},
		},
	}

	broadAreaPostsMap["Law"] = []hub.AddPostRequest{
		{
			Content: `Just billed 0.1 hours for thinking about billing. Efficiency.`,
			NewTags: []common.VTagName{
				common.VTagName("LawyerLife"),
				common.VTagName("BillableHours"),
				common.VTagName("LegalHumor"),
			},
		},
		{
			Content: `My brain is 90% case law, 10% remembering where I parked.`,
			NewTags: []common.VTagName{
				common.VTagName("LegalMind"),
				common.VTagName("InformationOverload"),
				common.VTagName("Law"),
			},
		},
		{
			Content: `Drafted a contract so airtight, even light cannot escape it. Client asked if we could "make it less... intense?" No.`,
			NewTags: []common.VTagName{
				common.VTagName("Contracts"),
				common.VTagName("AttentionToDetail"),
				common.VTagName("LawyerProblems"),
			},
		},
		{
			Content: `"It depends" is not a non-committal answer, it's the *most* accurate legal advice you'll ever get.`,
			NewTags: []common.VTagName{
				common.VTagName("LegalAdvice"),
				common.VTagName("Truth"),
				common.VTagName("Law"),
			},
		},
		{
			Content: `Successfully used "hereinafter" in casual conversation. The barista was confused. My work here is done.`,
			NewTags: []common.VTagName{
				common.VTagName("Legalese"),
				common.VTagName("LawyerLife"),
				common.VTagName("SmallVictories"),
			},
		},
		{
			Content: `Attended a 4-hour deposition. Pretty sure I achieved enlightenment somewhere around hour 3. Or maybe it was just the fluorescent lights.`,
			NewTags: []common.VTagName{
				common.VTagName("Litigation"),
				common.VTagName("Endurance"),
				common.VTagName("Law"),
			},
		},
		{
			Content: `Client: "Can we just agree to disagree?" Me: "We can agree that you sign this settlement, or we can disagree... in court."`,
			NewTags: []common.VTagName{
				common.VTagName("Negotiation"),
				common.VTagName("LawyerHumor"),
				common.VTagName("RealityCheck"),
			},
		},
		{
			Content: `Reading legislation is my cardio.`,
			NewTags: []common.VTagName{
				common.VTagName("Dedication"),
				common.VTagName("LawLife"),
				common.VTagName("FitnessGoals"),
			},
		},
		{
			Content: `Explained 'force majeure' to my cat. He seemed unimpressed. Tough crowd.`,
			NewTags: []common.VTagName{
				common.VTagName("WorkFromHome"),
				common.VTagName("LegalTerms"),
				common.VTagName("Law"),
			},
		},
		{
			Content: `Objection! Leading the witness... to the coffee machine. Sustained.`,
			NewTags: []common.VTagName{
				common.VTagName("CourtroomHumor"),
				common.VTagName("LawyerLife"),
				common.VTagName("Caffeine"),
			},
		},
	}

	broadAreaPostsMap["Education"] = []hub.AddPostRequest{
		{
			Content: `Grading papers fueled by coffee and the faint hope that someone actually read the syllabus.`,
			NewTags: []common.VTagName{
				common.VTagName("TeacherLife"),
				common.VTagName("Grading"),
				common.VTagName("EducationHumor"),
			},
		},
		{
			Content: `Spent an hour explaining a concept. Student asks, "Will this be on the test?" Sigh.`,
			NewTags: []common.VTagName{
				common.VTagName("EducatorProblems"),
				common.VTagName("StudentQuestions"),
				common.VTagName("Teaching"),
			},
		},
		{
			Content: `My superpower is deciphering handwriting that looks like abstract art.`,
			NewTags: []common.VTagName{
				common.VTagName("TeacherSkills"),
				common.VTagName("GradingLife"),
				common.VTagName("Education"),
			},
		},
		{
			Content: `"Summer break" is just educator code for "professional development and catching up on sleep... maybe."`,
			NewTags: []common.VTagName{
				common.VTagName("TeacherSummer"),
				common.VTagName("MythVsReality"),
				common.VTagName("Education"),
			},
		},
		{
			Content: `Attended a faculty meeting that could have been an email. Again. Productivity!`,
			NewTags: []common.VTagName{
				common.VTagName("AcademicLife"),
				common.VTagName("Meetings"),
				common.VTagName("Education"),
			},
		},
		{
			Content: `The sound of the bell doesn't dismiss the class. My existential sigh does.`,
			NewTags: []common.VTagName{
				common.VTagName("ClassroomManagement"),
				common.VTagName("TeacherHumor"),
				common.VTagName("EndOfTheDay"),
			},
		},
		{
			Content: `Found a meme in a student's essay. Points for creativity? Or deduction for unprofessionalism? The struggle is real.`,
			NewTags: []common.VTagName{
				common.VTagName("ModernEducation"),
				common.VTagName("GradingDilemmas"),
				common.VTagName("Teaching"),
			},
		},
		{
			Content: `Trying to inspire the next generation while simultaneously trying to remember where I put my coffee mug.`,
			NewTags: []common.VTagName{
				common.VTagName("Multitasking"),
				common.VTagName("EducatorLife"),
				common.VTagName("Priorities"),
			},
		},
		{
			Content: `"Let's circle back to that" - Professor code for "I have no idea, but I'll Google it during the break."`,
			NewTags: []common.VTagName{
				common.VTagName("AcademicJargon"),
				common.VTagName("Honesty"),
				common.VTagName("Education"),
			},
		},
		{
			Content: `That feeling when a student finally understands a difficult concept. It's almost as good as finding a working pen. Almost.`,
			NewTags: []common.VTagName{
				common.VTagName("TeachingWins"),
				common.VTagName("SmallJoys"),
				common.VTagName("Education"),
			},
		},
	}

	broadAreaPostsMap["Marketing"] = []hub.AddPostRequest{
		{
			Content: `Our latest campaign generated massive buzz! Mostly from the server room fan, but still... buzz!`,
			NewTags: []common.VTagName{
				common.VTagName("MarketingLife"),
				common.VTagName("SpinZone"),
				common.VTagName("Metrics"),
			},
		},
		{
			Content: `Content is king. Engagement is queen. The algorithm is the court jester messing everything up.`,
			NewTags: []common.VTagName{
				common.VTagName("SocialMediaMarketing"),
				common.VTagName("AlgorithmBlues"),
				common.VTagName("Marketing"),
			},
		},
		{
			Content: `Added "synergize" and "leverage" to my keyword list. SEO score went up, but my soul died a little.`,
			NewTags: []common.VTagName{
				common.VTagName("Buzzwords"),
				common.VTagName("SEOLife"),
				common.VTagName("MarketingPain"),
			},
		},
		{
			Content: `Client: "Can we make the logo bigger?" Me: *Opens Photoshop, increases size by 1px* "How's this?" Client: "Perfect!"`,
			NewTags: []common.VTagName{
				common.VTagName("ClientFeedback"),
				common.VTagName("DesignLife"),
				common.VTagName("MarketingMagic"),
			},
		},
		{
			Content: `Ran an A/B test. Button A got 10 clicks. Button B got 11 clicks. Conclusion: People randomly click things. Groundbreaking insights here.`,
			NewTags: []common.VTagName{
				common.VTagName("ABTesting"),
				common.VTagName("DataDriven"),
				common.VTagName("MarketingScience"),
			},
		},
		{
			Content: `My marketing funnel is more like a marketing sieve, but we're "optimizing" it.`,
			NewTags: []common.VTagName{
				common.VTagName("LeadGeneration"),
				common.VTagName("MarketingStrategy"),
				common.VTagName("Honesty"),
			},
		},
		{
			Content: `"Going viral" sounds like a disease. And based on some comment sections, it might be.`,
			NewTags: []common.VTagName{
				common.VTagName("ViralMarketing"),
				common.VTagName("SocialMedia"),
				common.VTagName("Observations"),
			},
		},
		{
			Content: `Pretty sure half our website traffic is just me checking if the analytics tag is still working.`,
			NewTags: []common.VTagName{
				common.VTagName("MarketingAnalytics"),
				common.VTagName("DataLife"),
				common.VTagName("Truth"),
			},
		},
		{
			Content: `Just used AI to write ad copy. It suggested: "Buy our stuff. It is... stuff." Nailed it.`,
			NewTags: []common.VTagName{
				common.VTagName("AIinMarketing"),
				common.VTagName("Copywriting"),
				common.VTagName("FutureIsNow"),
			},
		},
		{
			Content: `ROI isn't just Return on Investment. It's also "Really Overworked Individuals."`,
			NewTags: []common.VTagName{
				common.VTagName("MarketingLife"),
				common.VTagName("Workload"),
				common.VTagName("Acronyms"),
			},
		},
	}

	broadAreaPostsMap["Human Resources"] = []hub.AddPostRequest{
		{
			Content: `Onboarding new hires today! Step 1: Explain the coffee machine. Step 2: Everything else.`,
			NewTags: []common.VTagName{
				common.VTagName("Onboarding"),
				common.VTagName("Priorities"),
				common.VTagName("HRLife"),
			},
		},
		{
			Content: `"We're like a family here." Translation: Slightly dysfunctional, passive-aggressive emails are common.`,
			NewTags: []common.VTagName{
				common.VTagName("CompanyCulture"),
				common.VTagName("CorporateSpeak"),
				common.VTagName("HRHumor"),
			},
		},
		{
			Content: `Updated the employee handbook. Added a clause about not microwaving fish in the breakroom. Progress.`,
			NewTags: []common.VTagName{
				common.VTagName("PolicyUpdate"),
				common.VTagName("OfficeLife"),
				common.VTagName("HRWins"),
			},
		},
		{
			Content: `Performance review season: The annual festival of finding creative ways to say "meets expectations."`,
			NewTags: []common.VTagName{
				common.VTagName("PerformanceManagement"),
				common.VTagName("HR"),
				common.VTagName("CorporateRituals"),
			},
		},
		{
			Content: `Interviewing candidates. Asked one "Where do you see yourself in 5 years?" They said "Celebrating the 5-year anniversary of you asking me this question." Hired.`,
			NewTags: []common.VTagName{
				common.VTagName("Recruiting"),
				common.VTagName("Interviews"),
				common.VTagName("HRHumor"),
			},
		},
		{
			Content: `Received an anonymous suggestion for "Mandatory Nap Time." Forwarding to leadership with my strongest recommendation.`,
			NewTags: []common.VTagName{
				common.VTagName("EmployeeWellness"),
				common.VTagName("GoodIdeas"),
				common.VTagName("HR"),
			},
		},
		{
			Content: `My job involves mediating disputes. Today's crisis: Who finished the good biscuits?`,
			NewTags: []common.VTagName{
				common.VTagName("ConflictResolution"),
				common.VTagName("OfficePolitics"),
				common.VTagName("HRLife"),
			},
		},
		{
			Content: `Trying to foster a positive work environment while simultaneously enforcing 37 different compliance regulations. Easy peasy.`,
			NewTags: []common.VTagName{
				common.VTagName("HRChallenges"),
				common.VTagName("Multitasking"),
				common.VTagName("HumanResources"),
			},
		},
		{
			Content: `Another team-building exercise successfully planned. Hope everyone likes lukewarm pizza and awkward icebreakers!`,
			NewTags: []common.VTagName{
				common.VTagName("TeamBuilding"),
				common.VTagName("EmployeeEngagement"),
				common.VTagName("HR"),
			},
		},
		{
			Content: `"Please activate your synergy." - Just trying out phrases for the next all-hands meeting.`,
			NewTags: []common.VTagName{
				common.VTagName("CorporateJargon"),
				common.VTagName("Communication"),
				common.VTagName("HRTesting"),
			},
		},
	}

	broadAreaPostsMap["Medical"] = []hub.AddPostRequest{
		{
			Content: `Survived a 12-hour shift fueled by caffeine, adrenaline, and the sheer terror of reading Dr. Smith's handwriting.`,
			NewTags: []common.VTagName{
				common.VTagName("HospitalLife"),
				common.VTagName("NurseLife"),
				common.VTagName("MedicalHumor"),
			},
		},
		{
			Content: `Patient asked if WebMD was accurate. I prescribed two aspirin and a strong dose of skepticism.`,
			NewTags: []common.VTagName{
				common.VTagName("DoctorProblems"),
				common.VTagName("PatientCare"),
				common.VTagName("MedicalAdvice"),
			},
		},
		{
			Content: `The hospital cafeteria coffee has medicinal properties. Specifically, the property of keeping you awake through sheer bitterness.`,
			NewTags: []common.VTagName{
				common.VTagName("HealthcareWorker"),
				common.VTagName("HospitalFood"),
				common.VTagName("Survival"),
			},
		},
		{
			Content: `"Stat" means "Drop everything and run," usually towards the sound of a beeping machine or an empty coffee pot.`,
			NewTags: []common.VTagName{
				common.VTagName("MedicalJargon"),
				common.VTagName("Urgency"),
				common.VTagName("HospitalLife"),
			},
		},
		{
			Content: `Charting: The art of describing a complex medical event in fewer characters than a tweet, but taking 10 times longer.`,
			NewTags: []common.VTagName{
				common.VTagName("EMR"),
				common.VTagName("Documentation"),
				common.VTagName("MedicalAdmin"),
			},
		},
		{
			Content: `My scrubs have seen things... unspeakable things. Mostly coffee spills, but still.`,
			NewTags: []common.VTagName{
				common.VTagName("NurseProblems"),
				common.VTagName("ScrubsLife"),
				common.VTagName("MedicalHumor"),
			},
		},
		{
			Content: `Attended a medical conference. Learned about groundbreaking new treatments and confirmed the hotel biscuits are still mediocre.`,
			NewTags: []common.VTagName{
				common.VTagName("ContinuingEducation"),
				common.VTagName("MedicalField"),
				common.VTagName("ConferenceLife"),
			},
		},
		{
			Content: `Trying to explain complex medical procedures using simple analogies. "Think of this artery like a really stubborn garden hose..."`,
			NewTags: []common.VTagName{
				common.VTagName("PatientCommunication"),
				common.VTagName("DoctorSkills"),
				common.VTagName("Medical"),
			},
		},
		{
			Content: `The longest distance in the universe isn't between galaxies, it's between the nurses' station and a working pen.`,
			NewTags: []common.VTagName{
				common.VTagName("HospitalHumor"),
				common.VTagName("NurseLife"),
				common.VTagName("Truth"),
			},
		},
		{
			Content: `Night shift brain: Pretty sure I just tried to unlock my car with my stethoscope.`,
			NewTags: []common.VTagName{
				common.VTagName("NightShift"),
				common.VTagName("HealthcareWorker"),
				common.VTagName("SleepDeprived"),
			},
		},
	}

	broadAreaPostsMap["Data Science"] = []hub.AddPostRequest{
		{
			Content: `80% of data science is cleaning data. The other 20% is complaining about cleaning data.`,
			NewTags: []common.VTagName{
				common.VTagName("DataCleaning"),
				common.VTagName("DataScienceLife"),
				common.VTagName("Truth"),
			},
		},
		{
			Content: `My model has an accuracy of 98%! (On the training data. Let's not talk about the test data.)`,
			NewTags: []common.VTagName{
				common.VTagName("MachineLearning"),
				common.VTagName("Overfitting"),
				common.VTagName("DataScienceHumor"),
			},
		},
		{
			Content: `Built a dashboard so complex, even I don't know what it means anymore. But look at the pretty colors!`,
			NewTags: []common.VTagName{
				common.VTagName("DataVisualization"),
				common.VTagName("BI"),
				common.VTagName("DataScienceProblems"),
			},
		},
		{
			Content: `Feature engineering: Also known as "staring at data until it confesses."`,
			NewTags: []common.VTagName{
				common.VTagName("FeatureEngineering"),
				common.VTagName("DataScience"),
				common.VTagName("Process"),
			},
		},
		{
			Content: `Asked ChatGPT to explain my model. It gave a beautiful, confident, and completely wrong explanation. Like a tiny intern.`,
			NewTags: []common.VTagName{
				common.VTagName("AI"),
				common.VTagName("LLM"),
				common.VTagName("DataScience"),
			},
		},
		{
			Content: `The joy of finding a dataset that's already clean and perfectly formatted. ... Just kidding, that doesn't exist.`,
			NewTags: []common.VTagName{
				common.VTagName("DataLife"),
				common.VTagName("DataScientistsDream"),
				common.VTagName("Myth"),
			},
		},
		{
			Content: `Stakeholder: "Can you just sprinkle some AI on this?" Me: *Adds random forest model* "Consider it sprinkled."`,
			NewTags: []common.VTagName{
				common.VTagName("AIHype"),
				common.VTagName("StakeholderManagement"),
				common.VTagName("DataScience"),
			},
		},
		{
			Content: `My code is held together by Stack Overflow answers and sheer willpower.`,
			NewTags: []common.VTagName{
				common.VTagName("CodingLife"),
				common.VTagName("DataScienceReality"),
				common.VTagName("ImposterSyndrome"),
			},
		},
		{
			Content: `Correlation does not imply causation, but it's really good at making graphs look convincing.`,
			NewTags: []common.VTagName{
				common.VTagName("Statistics"),
				common.VTagName("DataAnalysis"),
				common.VTagName("Caution"),
			},
		},
		{
			Content: `Successfully predicted customer churn with my model. It also predicted I'd eat pizza for dinner. 100% accuracy today!`,
			NewTags: []common.VTagName{
				common.VTagName("PredictiveAnalytics"),
				common.VTagName("DataScienceWins"),
				common.VTagName("PersonalInsights"),
			},
		},
	}

	broadAreaPostsMap["Consulting"] = []hub.AddPostRequest{
		{
			Content: `Just created a 100-slide deck explaining a 2-slide concept. Added value. Synergy. Impact.`,
			NewTags: []common.VTagName{
				common.VTagName("ConsultingLife"),
				common.VTagName("PowerPoint"),
				common.VTagName("ValueAdd"),
			},
		},
		{
			Content: `My travel schedule is optimized for maximum airport lounge access and minimal sleep.`,
			NewTags: []common.VTagName{
				common.VTagName("ConsultantTravel"),
				common.VTagName("WorkLifeBalance"),
				common.VTagName("Perks"),
			},
		},
		{
			Content: `Client asked for a "quick analysis." Three weeks and 5 frameworks later... "Here's the quick analysis."`,
			NewTags: []common.VTagName{
				common.VTagName("ClientWork"),
				common.VTagName("Consulting"),
				common.VTagName("ScopeCreep"),
			},
		},
		{
			Content: `"Let's take this offline" = "I don't know the answer, but I'll find someone who does before the next meeting."`,
			NewTags: []common.VTagName{
				common.VTagName("ConsultingSpeak"),
				common.VTagName("MeetingStrategy"),
				common.VTagName("ProblemSolving"),
			},
		},
		{
			Content: `Survived another "blue sky thinking" session. Pretty sure my main contribution was suggesting more coffee.`,
			NewTags: []common.VTagName{
				common.VTagName("Brainstorming"),
				common.VTagName("CorporateJargon"),
				common.VTagName("ConsultingLife"),
			},
		},
		{
			Content: `Expense report submitted. Claimed "strategic alignment fuel" (it was coffee). Wish me luck.`,
			NewTags: []common.VTagName{
				common.VTagName("Expenses"),
				common.VTagName("ConsultingHumor"),
				common.VTagName("Truth"),
			},
		},
		{
			Content: `Just used a 2x2 matrix to decide what to have for lunch. Peak consulting achieved.`,
			NewTags: []common.VTagName{
				common.VTagName("Frameworks"),
				common.VTagName("ConsultingMindset"),
				common.VTagName("Overthinking"),
			},
		},
		{
			Content: `Wearing a suit on a Zoom call from my living room. Peak professionalism or peak absurdity? Yes.`,
			NewTags: []common.VTagName{
				common.VTagName("RemoteWork"),
				common.VTagName("Consulting"),
				common.VTagName("NewNormal"),
			},
		},
		{
			Content: `The client loves the recommendations! (They were mostly the client's ideas repackaged with better graphics).`,
			NewTags: []common.VTagName{
				common.VTagName("ClientManagement"),
				common.VTagName("ConsultingSkills"),
				common.VTagName("WinWin"),
			},
		},
		{
			Content: `My primary skill is looking confident while presenting slides I finished 5 minutes ago.`,
			NewTags: []common.VTagName{
				common.VTagName("PresentationSkills"),
				common.VTagName("ConsultingLife"),
				common.VTagName("ImposterSyndrome"),
			},
		},
	}

	broadAreaPostsMap["Design"] = []hub.AddPostRequest{
		{
			Content: `Client feedback: "Can you make it pop more?" *Increases saturation by 2%* Client: "Perfect!"`,
			NewTags: []common.VTagName{
				common.VTagName("ClientFeedback"),
				common.VTagName("DesignLife"),
				common.VTagName("Magic"),
			},
		},
		{
			Content: `Spent 3 hours debating the perfect shade of grey. My life is a monochrome adventure.`,
			NewTags: []common.VTagName{
				common.VTagName("ColorTheory"),
				common.VTagName("DesignerProblems"),
				common.VTagName("Details"),
			},
		},
		{
			Content: `My therapist told me to embrace whitespace. I told her I do, but stakeholders keep wanting to fill it with more content.`,
			NewTags: []common.VTagName{
				common.VTagName("Whitespace"),
				common.VTagName("DesignPrinciples"),
				common.VTagName("StakeholderManagement"),
			},
		},
		{
			Content: `"Just a small tweak" - famous last words before a complete redesign.`,
			NewTags: []common.VTagName{
				common.VTagName("ScopeCreep"),
				common.VTagName("DesignProcess"),
				common.VTagName("ClientRequests"),
			},
		},
		{
			Content: `Organized my layers in Photoshop. Feeling like I have my life together. It's an illusion, but a well-structured one.`,
			NewTags: []common.VTagName{
				common.VTagName("GraphicDesign"),
				common.VTagName("Organization"),
				common.VTagName("SmallWins"),
			},
		},
		{
			Content: `Judging websites based solely on their font choices. It's not snobbery, it's *professional assessment*.`,
			NewTags: []common.VTagName{
				common.VTagName("Typography"),
				common.VTagName("DesignEye"),
				common.VTagName("Priorities"),
			},
		},
		{
			Content: `User testing: Where you watch people completely miss the giant button you thought was obvious. Humbling.`,
			NewTags: []common.VTagName{
				common.VTagName("UXDesign"),
				common.VTagName("UserTesting"),
				common.VTagName("DesignReality"),
			},
		},
		{
			Content: `My design process involves 10% inspiration, 40% perspiration, and 50% convincing people Comic Sans is not an option.`,
			NewTags: []common.VTagName{
				common.VTagName("DesignProcess"),
				common.VTagName("FontWars"),
				common.VTagName("Advocacy"),
			},
		},
		{
			Content: `Exported final files. Named them Final_Final_ReallyFinal_ThisOne_v3.zip. Seems about right.`,
			NewTags: []common.VTagName{
				common.VTagName("FileNaming"),
				common.VTagName("DesignLife"),
				common.VTagName("Workflow"),
			},
		},
		{
			Content: `That feeling when the design just *clicks*. It's rarer than finding a unicorn riding a pixel-perfect skateboard, but it happens.`,
			NewTags: []common.VTagName{
				common.VTagName("DesignInspiration"),
				common.VTagName("CreativeFlow"),
				common.VTagName("MomentsOfJoy"),
			},
		},
	}

	broadAreaPostsMap["Operations"] = []hub.AddPostRequest{
		{
			Content: `My day involves putting out fires. Sometimes metaphorical, sometimes literal (don't ask about the server room incident).`,
			NewTags: []common.VTagName{
				common.VTagName("OperationsLife"),
				common.VTagName("ProblemSolving"),
				common.VTagName("Firefighter"),
			},
		},
		{
			Content: `Optimized a process today. Saved 3 seconds per transaction. At scale, that's almost enough time for a coffee break!`,
			NewTags: []common.VTagName{
				common.VTagName("ProcessImprovement"),
				common.VTagName("Efficiency"),
				common.VTagName("OpsWins"),
			},
		},
		{
			Content: `Spreadsheets are my love language. Pivot tables are my sonnets.`,
			NewTags: []common.VTagName{
				common.VTagName("ExcelNinja"),
				common.VTagName("DataAnalysis"),
				common.VTagName("OperationsTools"),
			},
		},
		{
			Content: `"Seamless integration" usually means "duct tape and hope."`,
			NewTags: []common.VTagName{
				common.VTagName("SystemsThinking"),
				common.VTagName("OperationsReality"),
				common.VTagName("Honesty"),
			},
		},
		{
			Content: `Supply chain issues again. Pretty sure our shipment is currently vacationing in Bermuda. Can't blame it.`,
			NewTags: []common.VTagName{
				common.VTagName("SupplyChain"),
				common.VTagName("LogisticsLife"),
				common.VTagName("OpsHumor"),
			},
		},
		{
			Content: `Created a flowchart so detailed, it includes branches for existential crises during coffee breaks.`,
			NewTags: []common.VTagName{
				common.VTagName("ProcessMapping"),
				common.VTagName("AttentionToDetail"),
				common.VTagName("Operations"),
			},
		},
		{
			Content: `My job is to ensure things run smoothly. Which mostly means anticipating how they might spectacularly fail.`,
			NewTags: []common.VTagName{
				common.VTagName("RiskManagement"),
				common.VTagName("ContingencyPlanning"),
				common.VTagName("OperationsMindset"),
			},
		},
		{
			Content: `Attended a meeting about efficiency. It ran 30 minutes over schedule. The irony was noted. Silently.`,
			NewTags: []common.VTagName{
				common.VTagName("MeetingCulture"),
				common.VTagName("OperationsLife"),
				common.VTagName("Irony"),
			},
		},
		{
			Content: `The ops team: Unsung heroes keeping the lights on, metaphorically and sometimes literally.`,
			NewTags: []common.VTagName{
				common.VTagName("Teamwork"),
				common.VTagName("BehindTheScenes"),
				common.VTagName("Operations"),
			},
		},
		{
			Content: `Just found the bottleneck. It was me, trying to find the bottleneck. Moving on.`,
			NewTags: []common.VTagName{
				common.VTagName("SelfAwareness"),
				common.VTagName("ProblemSolving"),
				common.VTagName("OpsHumor"),
			},
		},
	}

	broadAreaPostsMap["Sales"] = []hub.AddPostRequest{
		{
			Content: `Hit quota! Time to celebrate by immediately worrying about next month's quota.`,
			NewTags: []common.VTagName{
				common.VTagName("SalesLife"),
				common.VTagName("QuotaLife"),
				common.VTagName("NeverSatisfied"),
			},
		},
		{
			Content: `My CRM is my best friend, my worst enemy, and the only thing that understands my obsession with pipeline velocity.`,
			NewTags: []common.VTagName{
				common.VTagName("CRM"),
				common.VTagName("SalesTools"),
				common.VTagName("RelationshipStatus"),
			},
		},
		{
			Content: `"Just checking in" - Polite sales code for "Please, for the love of commission, sign the contract."`,
			NewTags: []common.VTagName{
				common.VTagName("SalesSpeak"),
				common.VTagName("FollowUp"),
				common.VTagName("SalesHumor"),
			},
		},
		{
			Content: `Closed a deal that's been in the works for months. Feeling like a superhero whose only power is persistent emailing.`,
			NewTags: []common.VTagName{
				common.VTagName("ClosingDeals"),
				common.VTagName("SalesWins"),
				common.VTagName("Persistence"),
			},
		},
		{
			Content: `Survived a cold-calling session. My ears are ringing, my spirit is slightly bruised, but my resilience is +10.`,
			NewTags: []common.VTagName{
				common.VTagName("ColdCalling"),
				common.VTagName("SalesGrind"),
				common.VTagName("Resilience"),
			},
		},
		{
			Content: `Prospect: "We'll review internally and get back to you." Translation: "Prepare for the follow-up abyss."`,
			NewTags: []common.VTagName{
				common.VTagName("SalesProcess"),
				common.VTagName("ObjectionHandling"),
				common.VTagName("RealTalk"),
			},
		},
		{
			Content: `My sales pitch is so smooth, I accidentally sold myself a pen this morning.`,
			NewTags: []common.VTagName{
				common.VTagName("SalesSkills"),
				common.VTagName("PitchPerfect"),
				common.VTagName("AccidentsHappen"),
			},
		},
		{
			Content: `Celebrating the end of the quarter like we just won the lottery. Except the prize is just... the start of the next quarter.`,
			NewTags: []common.VTagName{
				common.VTagName("EndOfQuarter"),
				common.VTagName("SalesCycle"),
				common.VTagName("GroundhogDay"),
			},
		},
		{
			Content: `Always Be Closing... the fridge door, the meeting tab, the deal. It's a lifestyle.`,
			NewTags: []common.VTagName{
				common.VTagName("ABC"),
				common.VTagName("SalesMantra"),
				common.VTagName("Lifestyle"),
			},
		},
		{
			Content: `They say rejection builds character. At this point, my character should be a skyscraper.`,
			NewTags: []common.VTagName{
				common.VTagName("Rejection"),
				common.VTagName("SalesLife"),
				common.VTagName("CharacterBuilding"),
			},
		},
	}

	broadAreaPostsMap["Product Management"] = []hub.AddPostRequest{
		{
			Content: `Roadmap planning: The art of confidently predicting the future while knowing everything will change next week.`,
			NewTags: []common.VTagName{
				common.VTagName("Roadmap"),
				common.VTagName("ProductStrategy"),
				common.VTagName("AgileLife"),
			},
		},
		{
			Content: `My job is 50% saying "no," 40% explaining why, and 10% wondering if I should have said "yes."`,
			NewTags: []common.VTagName{
				common.VTagName("Prioritization"),
				common.VTagName("ProductManager"),
				common.VTagName("DecisionMaking"),
			},
		},
		{
			Content: `User story writing: Translating vague stakeholder wishes into something engineers can actually build, possibly with magic.`,
			NewTags: []common.VTagName{
				common.VTagName("UserStories"),
				common.VTagName("Requirements"),
				common.VTagName("ProductOwner"),
			},
		},
		{
			Content: `"Let's put a pin in that" - Product Manager code for "Good idea, but it's going straight to the backlog abyss."`,
			NewTags: []common.VTagName{
				common.VTagName("ProductSpeak"),
				common.VTagName("BacklogGrooming"),
				common.VTagName("IdeaParkingLot"),
			},
		},
		{
			Content: `Just shipped a new feature! Now accepting bug reports and feature requests for version 2.0.`,
			NewTags: []common.VTagName{
				common.VTagName("ProductLaunch"),
				common.VTagName("ReleaseDay"),
				common.VTagName("NeverDone"),
			},
		},
		{
			Content: `Attended 7 meetings today. Pretty sure I'm now qualified as a professional meeting attendee. Where's my certificate?`,
			NewTags: []common.VTagName{
				common.VTagName("MeetingOverload"),
				common.VTagName("ProductLife"),
				common.VTagName("TimeManagement"),
			},
		},
		{
			Content: `Trying to balance user needs, business goals, and engineering constraints. It's like juggling chainsaws, but with more spreadsheets.`,
			NewTags: []common.VTagName{
				common.VTagName("ProductManagement"),
				common.VTagName("BalancingAct"),
				common.VTagName("StakeholderAlignment"),
			},
		},
		{
			Content: `That feeling when user feedback validates a feature you fought for. Briefly makes the endless meetings worth it. Briefly.`,
			NewTags: []common.VTagName{
				common.VTagName("UserFeedback"),
				common.VTagName("ProductWins"),
				common.VTagName("Validation"),
			},
		},
		{
			Content: `My backlog is longer than a CVS receipt and possibly contains items from the Jurassic period.`,
			NewTags: []common.VTagName{
				common.VTagName("ProductBacklog"),
				common.VTagName("AgileProblems"),
				common.VTagName("NeedsGrooming"),
			},
		},
		{
			Content: `Explaining the product vision with passion, clarity, and a slight tremor of panic about the deadline.`,
			NewTags: []common.VTagName{
				common.VTagName("ProductVision"),
				common.VTagName("Communication"),
				common.VTagName("DeadlinePressure"),
			},
		},
	}

	broadAreaPostsMap["Aerospace"] = []hub.AddPostRequest{
		{
			Content: `Calculated orbital mechanics before my first coffee. Just another Monday.`,
			NewTags: []common.VTagName{
				common.VTagName("RocketScience"),
				common.VTagName("AerospaceEngineering"),
				common.VTagName("MorningRoutine"),
			},
		},
		{
			Content: `My simulation crashed. Again. Either the physics is wrong, or the universe just enjoys messing with me.`,
			NewTags: []common.VTagName{
				common.VTagName("Simulation"),
				common.VTagName("AerospaceProblems"),
				common.VTagName("DebuggingLife"),
			},
		},
		{
			Content: `Building things that fly requires meticulous planning, precise engineering, and ignoring the little voice saying "What if it doesn't?"`,
			NewTags: []common.VTagName{
				common.VTagName("Aerospace"),
				common.VTagName("EngineeringLife"),
				common.VTagName("Mindset"),
			},
		},
		{
			Content: `"It's not rocket science." Oh, wait. Yes, it is. That's why it's taking so long.`,
			NewTags: []common.VTagName{
				common.VTagName("LiteralRocketScience"),
				common.VTagName("AerospaceHumor"),
				common.VTagName("Complexity"),
			},
		},
		{
			Content: `Reviewed launch readiness checklist. Item 347: Ensure coffee supply is adequate. Critical path item.`,
			NewTags: []common.VTagName{
				common.VTagName("LaunchPrep"),
				common.VTagName("Aerospace"),
				common.VTagName("Priorities"),
			},
		},
		{
			Content: `Dealing with tolerances measured in microns. My patience is measured in nanometers today.`,
			NewTags: []common.VTagName{
				common.VTagName("PrecisionEngineering"),
				common.VTagName("AerospaceManufacturing"),
				common.VTagName("Details"),
			},
		},
		{
			Content: `Attended a design review. Used the phrase "thrust-to-weight ratio" five times. Felt powerful.`,
			NewTags: []common.VTagName{
				common.VTagName("AerospaceJargon"),
				common.VTagName("EngineeringMeetings"),
				common.VTagName("PowerPhrases"),
			},
		},
		{
			Content: `My code is designed to handle catastrophic failure scenarios. My brain, less so after a long week.`,
			NewTags: []common.VTagName{
				common.VTagName("SafetyCritical"),
				common.VTagName("AerospaceSoftware"),
				common.VTagName("Burnout"),
			},
		},
		{
			Content: `Explaining aerodynamics using hand gestures. Pretty sure I just invented a new form of interpretive dance.`,
			NewTags: []common.VTagName{
				common.VTagName("Aerodynamics"),
				common.VTagName("CommunicationSkills"),
				common.VTagName("Aerospace"),
			},
		},
		{
			Content: `That moment when the test results match the simulation. Pure, unadulterated, nerdy joy.`,
			NewTags: []common.VTagName{
				common.VTagName("TestingAndValidation"),
				common.VTagName("AerospaceWins"),
				common.VTagName("EngineeringSuccess"),
			},
		},
	}

	broadAreaPostsMap["Automotive"] = []hub.AddPostRequest{
		{
			Content: `Debugged CAN bus issues all morning. I now speak fluent hexadecimal and existential dread.`,
			NewTags: []common.VTagName{
				common.VTagName("AutomotiveTech"),
				common.VTagName("Debugging"),
				common.VTagName("EngineerLife"),
			},
		},
		{
			Content: `Designed a component to withstand extreme temperatures. Tested it by leaving my coffee mug on it. Passed.`,
			NewTags: []common.VTagName{
				common.VTagName("AutomotiveTesting"),
				common.VTagName("EngineeringHumor"),
				common.VTagName("RealWorldTesting"),
			},
		},
		{
			Content: `Talking about torque and horsepower like it's gossip. "Did you hear about the new engine? Scandalous!"`,
			NewTags: []common.VTagName{
				common.VTagName("EngineTalk"),
				common.VTagName("AutomotiveIndustry"),
				common.VTagName("Passion"),
			},
		},
		{
			Content: `"Let's just make this sensor 1mm smaller." - Famous last words before redesigning half the engine bay.`,
			NewTags: []common.VTagName{
				common.VTagName("AutomotiveDesign"),
				common.VTagName("ScopeCreep"),
				common.VTagName("EngineeringProblems"),
			},
		},
		{
			Content: `Attended a meeting on fuel efficiency. Drove my V8 gas guzzler home. Balance.`,
			NewTags: []common.VTagName{
				common.VTagName("AutomotiveLife"),
				common.VTagName("Irony"),
				common.VTagName("Confessions"),
			},
		},
		{
			Content: `My car has more lines of code than the lunar lander. And probably more bugs.`,
			NewTags: []common.VTagName{
				common.VTagName("SoftwareDefinedVehicle"),
				common.VTagName("AutomotiveSoftware"),
				common.VTagName("Complexity"),
			},
		},
		{
			Content: `Ran thermal simulations for the new EV battery pack. Conclusion: It gets hot. Groundbreaking.`,
			NewTags: []common.VTagName{
				common.VTagName("EV"),
				common.VTagName("Simulation"),
				common.VTagName("AutomotiveEngineering"),
			},
		},
		{
			Content: `Trying to reduce vehicle weight. Considered replacing myself with a lighter engineer, but HR advised against it.`,
			NewTags: []common.VTagName{
				common.VTagName("Lightweighting"),
				common.VTagName("AutomotiveHumor"),
				common.VTagName("HRSaidNo"),
			},
		},
		{
			Content: `Crash test ratings are important. My code's crash test rating after pulling an all-nighter? Less stellar.`,
			NewTags: []common.VTagName{
				common.VTagName("AutomotiveSafety"),
				common.VTagName("CodingLife"),
				common.VTagName("Relatable"),
			},
		},
		{
			Content: `That feeling when the prototype finally drives without catching fire. A good day in automotive.`,
			NewTags: []common.VTagName{
				common.VTagName("PrototypeTesting"),
				common.VTagName("AutomotiveWins"),
				common.VTagName("Success"),
			},
		},
	}

	broadAreaPostsMap["Hospitality"] = []hub.AddPostRequest{
		{
			Content: `"The customer is always right," except when they insist their room key opens the minibar for free.`,
			NewTags: []common.VTagName{
				common.VTagName("HospitalityLife"),
				common.VTagName("CustomerService"),
				common.VTagName("HotelProblems"),
			},
		},
		{
			Content: `Survived the check-in rush fueled by complimentary lobby coffee and the ability to smile while crying internally.`,
			NewTags: []common.VTagName{
				common.VTagName("FrontDeskLife"),
				common.VTagName("Hospitality"),
				common.VTagName("SurvivalSkills"),
			},
		},
		{
			Content: `Mastered the art of folding a fitted sheet. Next up: World peace.`,
			NewTags: []common.VTagName{
				common.VTagName("HousekeepingSkills"),
				common.VTagName("HospitalityWins"),
				common.VTagName("LifeGoals"),
			},
		},
		{
			Content: `"Can I speak to the manager?" - Words that strike fear into the heart of every hospitality worker.`,
			NewTags: []common.VTagName{
				common.VTagName("ManagementPlease"),
				common.VTagName("HospitalityHumor"),
				common.VTagName("CodeRed"),
			},
		},
		{
			Content: `Dealing with bizarre guest requests. "Can you arrange for a unicorn to deliver my room service?" Let me check on that for you...`,
			NewTags: []common.VTagName{
				common.VTagName("GuestRequests"),
				common.VTagName("ConciergeLife"),
				common.VTagName("AnythingIsPossible"),
			},
		},
		{
			Content: `My steps counter goes crazy during a shift. Pretty sure I walk a marathon around this hotel daily.`,
			NewTags: []common.VTagName{
				common.VTagName("HospitalityFit"),
				common.VTagName("OnYourFeet"),
				common.VTagName("WorkLife"),
			},
		},
		{
			Content: `That moment a guest leaves a genuinely nice review. Restores my faith in humanity (for about 5 minutes).`,
			NewTags: []common.VTagName{
				common.VTagName("PositiveReviews"),
				common.VTagName("HospitalityJoy"),
				common.VTagName("Motivation"),
			},
		},
		{
			Content: `Explaining the difference between 'ocean view' and 'ocean glimpse' requires diplomatic skills worthy of the UN.`,
			NewTags: []common.VTagName{
				common.VTagName("Reservations"),
				common.VTagName("SettingExpectations"),
				common.VTagName("Hospitality"),
			},
		},
		{
			Content: `The night audit: Where time, logic, and basic math skills go on a little vacation.`,
			NewTags: []common.VTagName{
				common.VTagName("NightShift"),
				common.VTagName("HotelOperations"),
				common.VTagName("HospitalityLife"),
			},
		},
		{
			Content: `"Service with a smile" - even when the coffee machine exploded and there's a conga line forming at reception.`,
			NewTags: []common.VTagName{
				common.VTagName("Professionalism"),
				common.VTagName("GraceUnderPressure"),
				common.VTagName("HospitalityStrong"),
			},
		},
	}

	broadAreaPostsMap["Retail"] = []hub.AddPostRequest{
		{
			Content: `Just perfectly folded a mountain of sweaters. It will remain perfect for approximately 7 seconds.`,
			NewTags: []common.VTagName{
				common.VTagName("RetailLife"),
				common.VTagName("VisualMerchandising"),
				common.VTagName("FutileEfforts"),
			},
		},
		{
			Content: `Customer: "Do you work here?" Me: *Wearing branded uniform, name tag, standing behind register* "Just visiting."`,
			NewTags: []common.VTagName{
				common.VTagName("RetailProblems"),
				common.VTagName("CustomerQuestions"),
				common.VTagName("Sarcasm"),
			},
		},
		{
			Content: `Survived the weekend rush. My feet hate me, but my ability to upsell socks is stronger than ever.`,
			NewTags: []common.VTagName{
				common.VTagName("RetailWorker"),
				common.VTagName("SalesFloor"),
				common.VTagName("SmallWins"),
			},
		},
		{
			Content: `"The item scanned at the wrong price? Let me just manually override reality for you."`,
			NewTags: []common.VTagName{
				common.VTagName("CashierLife"),
				common.VTagName("RetailHumor"),
				common.VTagName("MagicWand"),
			},
		},
		{
			Content: `Working inventory day. Pretty sure I've seen boxes older than me in that stockroom.`,
			NewTags: []common.VTagName{
				common.VTagName("StockroomAdventures"),
				common.VTagName("RetailOperations"),
				common.VTagName("Archaeology"),
			},
		},
		{
			Content: `Hearing holiday music in October triggers my retail PTSD. Fa-la-la-la-NO.`,
			NewTags: []common.VTagName{
				common.VTagName("HolidaySeason"),
				common.VTagName("RetailTrauma"),
				common.VTagName("TooSoon"),
			},
		},
		{
			Content: `That feeling when a customer genuinely thanks you for your help. It's like finding a $20 bill in an old coat.`,
			NewTags: []common.VTagName{
				common.VTagName("CustomerServiceWin"),
				common.VTagName("RetailJoy"),
				common.VTagName("RareOccasion"),
			},
		},
		{
			Content: `Explaining the return policy for the 50th time today. My patience is also non-refundable after 30 days.`,
			NewTags: []common.VTagName{
				common.VTagName("ReturnPolicy"),
				common.VTagName("RetailLife"),
				common.VTagName("TestingLimits"),
			},
		},
		{
			Content: `Closing shift: The magical time when everything needs to be cleaned, restocked, and faced, usually by one person. Me.`,
			NewTags: []common.VTagName{
				common.VTagName("ClosingTime"),
				common.VTagName("RetailGrind"),
				common.VTagName("MultitaskingHero"),
			},
		},
		{
			Content: `"Can I get a discount?" - The official soundtrack of my retail career.`,
			NewTags: []common.VTagName{
				common.VTagName("DiscountHunters"),
				common.VTagName("RetailLife"),
				common.VTagName("ThemeSong"),
			},
		},
	}

	broadAreaPostsMap["Pharmaceuticals"] = []hub.AddPostRequest{
		{
			Content: `Spent the day pipetting tiny amounts of liquid. Felt like a giant playing with a very expensive, sterile dollhouse.`,
			NewTags: []common.VTagName{
				common.VTagName("LabLife"),
				common.VTagName("Research"),
				common.VTagName("PharmaHumor"),
			},
		},
		{
			Content: `Reading clinical trial data. Side effects may include drowsiness, dizziness, and questioning all your life choices.`,
			NewTags: []common.VTagName{
				common.VTagName("ClinicalTrials"),
				common.VTagName("DataAnalysis"),
				common.VTagName("Pharma"),
			},
		},
		{
			Content: `Drug naming convention meeting. Rejected "MiracleCureXtreme." Suggested "SlightlyHelpfulMaybe." We'll workshop it.`,
			NewTags: []common.VTagName{
				common.VTagName("DrugDevelopment"),
				common.VTagName("Marketing"),
				common.VTagName("PharmaLife"),
			},
		},
		{
			Content: `"This formulation requires precise temperature control." *Looks nervously at the office thermostat 전쟁*`,
			NewTags: []common.VTagName{
				common.VTagName("Manufacturing"),
				common.VTagName("QualityControl"),
				common.VTagName("PharmaProblems"),
			},
		},
		{
			Content: `Attended a regulatory affairs seminar. Learned 100 new ways to fill out forms incorrectly. Progress!`,
			NewTags: []common.VTagName{
				common.VTagName("RegulatoryAffairs"),
				common.VTagName("Compliance"),
				common.VTagName("Pharma"),
			},
		},
		{
			Content: `My experiment failed. Again. Time to invoke the scientific method: Step 1, cry. Step 2, coffee. Step 3, try again.`,
			NewTags: []common.VTagName{
				common.VTagName("ResearchLife"),
				common.VTagName("Setbacks"),
				common.VTagName("Persistence"),
			},
		},
		{
			Content: `Explaining pharmacokinetics using breakfast analogies. "This drug is like slow-release oatmeal..."`,
			NewTags: []common.VTagName{
				common.VTagName("ScienceCommunication"),
				common.VTagName("Pharma"),
				common.VTagName("Analogies"),
			},
		},
		{
			Content: `That moment your Western Blot actually works. You feel like a wizard who just summoned a faint, blurry band.`,
			NewTags: []common.VTagName{
				common.VTagName("LabWork"),
				common.VTagName("Biotech"),
				common.VTagName("SmallVictories"),
			},
		},
		{
			Content: `Trying to synthesize a new compound. Currently synthesizing new levels of frustration.`,
			NewTags: []common.VTagName{
				common.VTagName("Chemistry"),
				common.VTagName("DrugDiscovery"),
				common.VTagName("PharmaResearch"),
			},
		},
		{
			Content: `Celebrating a successful drug approval like we just cured everything. (Spoiler: We didn't, but let us have this moment).`,
			NewTags: []common.VTagName{
				common.VTagName("DrugApproval"),
				common.VTagName("PharmaWins"),
				common.VTagName("Milestone"),
			},
		},
	}

	broadAreaPostsMap["Construction"] = []hub.AddPostRequest{
		{
			Content: `Project site visit today. Confirmed mud exists and hard hats mess up your hair. Crucial findings.`,
			NewTags: []common.VTagName{
				common.VTagName("SiteLife"),
				common.VTagName("Construction"),
				common.VTagName("FieldWork"),
			},
		},
		{
			Content: `Reading blueprints that look like abstract spaghetti monsters. Pretty sure this line means "wall," maybe?`,
			NewTags: []common.VTagName{
				common.VTagName("Blueprints"),
				common.VTagName("ConstructionManagement"),
				common.VTagName("Interpretation"),
			},
		},
		{
			Content: `Meeting about budget overruns. Suggested replacing solid gold fixtures with slightly less solid gold. We'll see.`,
			NewTags: []common.VTagName{
				common.VTagName("Budgeting"),
				common.VTagName("ProjectManagement"),
				common.VTagName("ConstructionHumor"),
			},
		},
		{
			Content: `"We need to value engineer this." Translation: "How can we make this cheaper without it collapsing immediately?"`,
			NewTags: []common.VTagName{
				common.VTagName("ValueEngineering"),
				common.VTagName("ConstructionSpeak"),
				common.VTagName("CostCutting"),
			},
		},
		{
			Content: `Dealing with unexpected delays. Today's culprit: A flock of pigeons decided the scaffolding was prime real estate.`,
			NewTags: []common.VTagName{
				common.VTagName("ProjectDelays"),
				common.VTagName("ConstructionLife"),
				common.VTagName("NatureWins"),
			},
		},
		{
			Content: `Safety briefing: "Don't stand under the thing being lifted." Profound stuff.`,
			NewTags: []common.VTagName{
				common.VTagName("ConstructionSafety"),
				common.VTagName("CommonSense"),
				common.VTagName("ToolboxTalk"),
			},
		},
		{
			Content: `That satisfying feeling when the concrete pour goes smoothly. It's the little things... and the giant spinning truck.`,
			NewTags: []common.VTagName{
				common.VTagName("ConcreteLife"),
				common.VTagName("ConstructionWins"),
				common.VTagName("Milestone"),
			},
		},
		{
			Content: `Coordinating subcontractors is like herding cats. Except the cats have power tools and opinions on rebar spacing.`,
			NewTags: []common.VTagName{
				common.VTagName("SubcontractorManagement"),
				common.VTagName("Construction"),
				common.VTagName("Coordination"),
			},
		},
		{
			Content: `Trying to explain project timelines to clients. "Yes, the building will magically appear on Tuesday, weather permitting."`,
			NewTags: []common.VTagName{
				common.VTagName("ClientCommunication"),
				common.VTagName("ProjectScheduling"),
				common.VTagName("Realism"),
			},
		},
		{
			Content: `End of the day. Covered in dust, slightly deafened, but the structure is still standing. Success!`,
			NewTags: []common.VTagName{
				common.VTagName("ConstructionWorker"),
				common.VTagName("HardWork"),
				common.VTagName("JobDone"),
			},
		},
	}

	broadAreaPostsMap["Real Estate"] = []hub.AddPostRequest{
		{
			Content: `Just hosted an open house. Served artisanal cheese. Pretty sure people came for the cheese, not the house. Still counts?`,
			NewTags: []common.VTagName{
				common.VTagName("OpenHouse"),
				common.VTagName("RealEstateLife"),
				common.VTagName("MarketingStrategy"),
			},
		},
		{
			Content: `"Cozy" = Small. "Charming" = Old. "Needs TLC" = Bring a bulldozer. Mastering the art of real estate euphemisms.`,
			NewTags: []common.VTagName{
				common.VTagName("RealEstateJargon"),
				common.VTagName("ListingDescription"),
				common.VTagName("TruthInAdvertising"),
			},
		},
		{
			Content: `Negotiating offers like a high-stakes poker game, except with more paperwork and less cool sunglasses.`,
			NewTags: []common.VTagName{
				common.VTagName("Negotiation"),
				common.VTagName("RealEstateAgent"),
				common.VTagName("DealMaking"),
			},
		},
		{
			Content: `Showing houses all day. My car now permanently smells like air freshener and desperation.`,
			NewTags: []common.VTagName{
				common.VTagName("RealtorLife"),
				common.VTagName("ShowingHomes"),
				common.VTagName("OccupationalHazards"),
			},
		},
		{
			Content: `Client: "I want a 5-bedroom house with a pool, downtown, under $100k." Me: "Have you considered Mars?"`,
			NewTags: []common.VTagName{
				common.VTagName("ClientExpectations"),
				common.VTagName("RealEstateHumor"),
				common.VTagName("RealityCheck"),
			},
		},
		{
			Content: `The thrill of getting a signed contract! Almost makes up for the 50 unanswered calls that preceded it. Almost.`,
			NewTags: []common.VTagName{
				common.VTagName("ClosingDeals"),
				common.VTagName("RealEstateWins"),
				common.VTagName("Persistence"),
			},
		},
		{
			Content: `Market analysis: Prices are up! Prices are down! Prices are sideways! Basically, nobody knows, but buy now!`,
			NewTags: []common.VTagName{
				common.VTagName("MarketTrends"),
				common.VTagName("RealEstateAdvice"),
				common.VTagName("CrystalBall"),
			},
		},
		{
			Content: `Staging a house: The art of making it look like nobody actually lives there, which is ironically what buyers want.`,
			NewTags: []common.VTagName{
				common.VTagName("HomeStaging"),
				common.VTagName("RealEstateMarketing"),
				common.VTagName("Psychology"),
			},
		},
		{
			Content: `My commission check is playing hide-and-seek. It's very good at hiding.`,
			NewTags: []common.VTagName{
				common.VTagName("RealEstateIncome"),
				common.VTagName("AgentProblems"),
				common.VTagName("WaitingGame"),
			},
		},
		{
			Content: `"Location, Location, Location!" Also important: "Paperwork, Paperwork, Paperwork!"`,
			NewTags: []common.VTagName{
				common.VTagName("RealEstateMantra"),
				common.VTagName("AdminLife"),
				common.VTagName("Truth"),
			},
		},
	}

	broadAreaPostsMap["Entertainment"] = []hub.AddPostRequest{
		{
			Content: `On set today. 90% waiting, 10% frantic activity. Showbiz, baby!`,
			NewTags: []common.VTagName{
				common.VTagName("SetLife"),
				common.VTagName("Production"),
				common.VTagName("EntertainmentIndustry"),
			},
		},
		{
			Content: `Just read a script where the main character's motivation is "revenge... for his parking spot?" Bold choice.`,
			NewTags: []common.VTagName{
				common.VTagName("Screenwriting"),
				common.VTagName("EntertainmentHumor"),
				common.VTagName("CreativeDecisions"),
			},
		},
		{
			Content: `Budget meeting. Suggested cutting the craft services budget. Almost got fired. Lesson learned.`,
			NewTags: []common.VTagName{
				common.VTagName("FilmFinance"),
				common.VTagName("ProductionManagement"),
				common.VTagName("Priorities"),
			},
		},
		{
			Content: `"We'll fix it in post." - The magical phrase that solves all production problems (until post-production).`,
			NewTags: []common.VTagName{
				common.VTagName("PostProduction"),
				common.VTagName("FilmMaking"),
				common.VTagName("EntertainmentSpeak"),
			},
		},
		{
			Content: `Casting call today. Saw 50 people convincingly pretend to be a talking squirrel. This industry is wild.`,
			NewTags: []common.VTagName{
				common.VTagName("Casting"),
				common.VTagName("Acting"),
				common.VTagName("EntertainmentLife"),
			},
		},
		{
			Content: `Trying to secure distribution. It's easier to get a meeting with Bigfoot.`,
			NewTags: []common.VTagName{
				common.VTagName("FilmDistribution"),
				common.VTagName("EntertainmentBusiness"),
				common.VTagName("Challenges"),
			},
		},
		{
			Content: `That feeling when the audience actually laughs at the joke you wrote. Pure, unadulterated validation.`,
			NewTags: []common.VTagName{
				common.VTagName("Writing"),
				common.VTagName("Comedy"),
				common.VTagName("EntertainmentWins"),
			},
		},
		{
			Content: `Dealing with talent agents. Requires the patience of a saint and the negotiating skills of a warlord.`,
			NewTags: []common.VTagName{
				common.VTagName("TalentManagement"),
				common.VTagName("EntertainmentIndustry"),
				common.VTagName("AgentLife"),
			},
		},
		{
			Content: `Wrap party! Celebrating the end of sleep deprivation and the beginning of worrying about the reviews.`,
			NewTags: []common.VTagName{
				common.VTagName("WrapParty"),
				common.VTagName("ProductionLife"),
				common.VTagName("CycleOfWorry"),
			},
		},
		{
			Content: `"It's got heart." Entertainment code for "The plot makes no sense, but maybe you'll cry?"`,
			NewTags: []common.VTagName{
				common.VTagName("FilmCritique"),
				common.VTagName("EntertainmentJargon"),
				common.VTagName("Spin"),
			},
		},
	}

	broadAreaPostsMap["Media"] = []hub.AddPostRequest{
		{
			Content: `Deadline looming. Fueled by coffee, adrenaline, and the fear of the editor's red pen.`,
			NewTags: []common.VTagName{
				common.VTagName("Journalism"),
				common.VTagName("DeadlineDriven"),
				common.VTagName("MediaLife"),
			},
		},
		{
			Content: `Chasing down a source who insists on speaking only in riddles. Investigative journalism or LARPing? Hard to tell.`,
			NewTags: []common.VTagName{
				common.VTagName("Reporting"),
				common.VTagName("Sources"),
				common.VTagName("MediaProblems"),
			},
		},
		{
			Content: `Fact-checking an article. Discovered the "expert" quoted based their entire argument on a tweet they misread. Sigh.`,
			NewTags: []common.VTagName{
				common.VTagName("FactChecking"),
				common.VTagName("MediaEthics"),
				common.VTagName("InformationAge"),
			},
		},
		{
			Content: `"We need more clicks!" - The battle cry of modern media. Let's add a listicle about cats!`,
			NewTags: []common.VTagName{
				common.VTagName("DigitalMedia"),
				common.VTagName("Clickbait"),
				common.VTagName("MediaStrategy"),
			},
		},
		{
			Content: `Attended a press conference. Got a free pen and vague non-answers. Success?`,
			NewTags: []common.VTagName{
				common.VTagName("PressConference"),
				common.VTagName("MediaEvents"),
				common.VTagName("Perks"),
			},
		},
		{
			Content: `Trying to explain complex global events in 800 words. It's like summarizing War and Peace on a cocktail napkin.`,
			NewTags: []common.VTagName{
				common.VTagName("Writing"),
				common.VTagName("MediaChallenges"),
				common.VTagName("Brevity"),
			},
		},
		{
			Content: `That feeling when your story gets picked up by major outlets. Briefly forget you're paid in exposure and coffee vouchers.`,
			NewTags: []common.VTagName{
				common.VTagName("MediaWins"),
				common.VTagName("JournalismSuccess"),
				common.VTagName("Motivation"),
			},
		},
		{
			Content: `Dealing with angry commenters who clearly only read the headline. My block button is getting a workout.`,
			NewTags: []common.VTagName{
				common.VTagName("OnlineComments"),
				common.VTagName("MediaLife"),
				common.VTagName("DigitalAge"),
			},
		},
		{
			Content: `The news never sleeps. Unfortunately, journalists do. Occasionally.`,
			NewTags: []common.VTagName{
				common.VTagName("NewsCycle"),
				common.VTagName("MediaGrind"),
				common.VTagName("SleepDeprived"),
			},
		},
		{
			Content: `"Off the record..." - Famous last words before someone tells you the juiciest story you can't publish.`,
			NewTags: []common.VTagName{
				common.VTagName("JournalismEthics"),
				common.VTagName("Sources"),
				common.VTagName("MediaDilemmas"),
			},
		},
	}

	broadAreaPostsMap["Telecommunications"] = []hub.AddPostRequest{
		{
			Content: `Traced a network outage to a squirrel chewing on a fiber optic cable. Never underestimate nature's chaos monkeys.`,
			NewTags: []common.VTagName{
				common.VTagName("NetworkOperations"),
				common.VTagName("TelecomLife"),
				common.VTagName("SquirrelIncident"),
			},
		},
		{
			Content: `Configuring routers all day. Pretty sure I can now communicate directly with machines via blinking lights.`,
			NewTags: []common.VTagName{
				common.VTagName("NetworkEngineer"),
				common.VTagName("TelecomTech"),
				common.VTagName("NewSkills"),
			},
		},
		{
			Content: `Explaining bandwidth limitations to customers. "No sir, you can't download the entire internet in 5 seconds."`,
			NewTags: []common.VTagName{
				common.VTagName("CustomerSupport"),
				common.VTagName("Telecom"),
				common.VTagName("SettingExpectations"),
			},
		},
		{
			Content: `"Five nines uptime" is the goal. Reality involves hoping the duct tape holds during peak hours.`,
			NewTags: []common.VTagName{
				common.VTagName("NetworkReliability"),
				common.VTagName("TelecomHumor"),
				common.VTagName("OpsLife"),
			},
		},
		{
			Content: `Attended a meeting about 5G deployment. Mostly understood the acronyms. Progress.`,
			NewTags: []common.VTagName{
				common.VTagName("5G"),
				common.VTagName("TelecomIndustry"),
				common.VTagName("AcronymSoup"),
			},
		},
		{
			Content: `My troubleshooting process: 1. Reboot it. 2. Check the cables. 3. Blame the user. 4. Panic. 5. Coffee. 6. Actually fix it.`,
			NewTags: []common.VTagName{
				common.VTagName("Troubleshooting"),
				common.VTagName("TelecomSupport"),
				common.VTagName("Process"),
			},
		},
		{
			Content: `That satisfying feeling when the signal bars go from one to full. It's like watching a tiny miracle unfold.`,
			NewTags: []common.VTagName{
				common.VTagName("NetworkPerformance"),
				common.VTagName("TelecomWins"),
				common.VTagName("SmallJoys"),
			},
		},
		{
			Content: `Dealing with legacy systems held together by hope and undocumented Perl scripts.`,
			NewTags: []common.VTagName{
				common.VTagName("LegacyTech"),
				common.VTagName("TelecomChallenges"),
				common.VTagName("Archaeology"),
			},
		},
		{
			Content: `Climbing a cell tower. Great views, slightly terrifying. Worth it for the 'gram? Debatable.`,
			NewTags: []common.VTagName{
				common.VTagName("FieldTechnician"),
				common.VTagName("TelecomLife"),
				common.VTagName("Heights"),
			},
		},
		{
			Content: `"The network is slow." - The universal complaint that could mean anything from sunspots to someone microwaving a burrito.`,
			NewTags: []common.VTagName{
				common.VTagName("NetworkIssues"),
				common.VTagName("TelecomMysteries"),
				common.VTagName("TroubleshootingLife"),
			},
		},
	}

	broadAreaPostsMap["Renewable Energy"] = []hub.AddPostRequest{
		{
			Content: `Calculating solar panel efficiency. Mostly involves praying for sunny days and minimal bird poop.`,
			NewTags: []common.VTagName{
				common.VTagName("SolarEnergy"),
				common.VTagName("Renewables"),
				common.VTagName("Reality"),
			},
		},
		{
			Content: `Designing a wind turbine blade. It needs to be strong, efficient, and hopefully not scare the local cows.`,
			NewTags: []common.VTagName{
				common.VTagName("WindEnergy"),
				common.VTagName("EngineeringDesign"),
				common.VTagName("RenewableTech"),
			},
		},
		{
			Content: `Site assessment for a new solar farm. Discovered the optimal location is currently occupied by very stubborn goats. Negotiations pending.`,
			NewTags: []common.VTagName{
				common.VTagName("SiteAssessment"),
				common.VTagName("RenewableEnergyLife"),
				common.VTagName("GoatDiplomacy"),
			},
		},
		{
			Content: `"Grid integration" sounds simple. In reality, it's like teaching calculus to a toaster.`,
			NewTags: []common.VTagName{
				common.VTagName("GridTechnology"),
				common.VTagName("RenewablesIntegration"),
				common.VTagName("Challenges"),
			},
		},
		{
			Content: `Attended a conference on battery storage. Learned that the future is bright, rechargeable, and slightly explosive if mishandled.`,
			NewTags: []common.VTagName{
				common.VTagName("EnergyStorage"),
				common.VTagName("RenewableFuture"),
				common.VTagName("SafetyFirst"),
			},
		},
		{
			Content: `Trying to explain renewable energy credits. Pretty sure I confused myself halfway through.`,
			NewTags: []common.VTagName{
				common.VTagName("EnergyPolicy"),
				common.VTagName("Renewables"),
				common.VTagName("Complexity"),
			},
		},
		{
			Content: `That feeling when the turbines start spinning and the power meter goes up. Saving the planet, one rotation at a time!`,
			NewTags: []common.VTagName{
				common.VTagName("WindPower"),
				common.VTagName("RenewableWins"),
				common.VTagName("MakingADifference"),
			},
		},
		{
			Content: `Dealing with intermittent power generation. The sun sets, the wind stops, my anxiety spikes. Normal Tuesday.`,
			NewTags: []common.VTagName{
				common.VTagName("Intermittency"),
				common.VTagName("RenewableChallenges"),
				common.VTagName("GridManagement"),
			},
		},
		{
			Content: `My job involves harnessing the power of nature. Which mostly means dealing with weather delays and unexpected wildlife encounters.`,
			NewTags: []common.VTagName{
				common.VTagName("RenewableEnergyOps"),
				common.VTagName("Nature"),
				common.VTagName("FieldWork"),
			},
		},
		{
			Content: `"Carbon neutral" is the goal. My coffee consumption? Less so. Baby steps.`,
			NewTags: []common.VTagName{
				common.VTagName("Sustainability"),
				common.VTagName("RenewableEnergyHumor"),
				common.VTagName("PersonalGoals"),
			},
		},
	}

}
