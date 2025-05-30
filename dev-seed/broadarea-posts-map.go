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
			TagIDs: []common.VTagID{
				common.VTagID("software-engineering"),
			},
		},
		{
			Content: `Woke up at 4AM. Did yoga. Built a microservice. Deployed to production. Got paged. Fixed it. Broke it again. Blamed Jenkins. Meditated. Asked ChatGPT. Took credit.

Remember: It's not about how many bugs you fix, it's about how confidently you explain why they aren't your fault.`,
			TagIDs: []common.VTagID{
				common.VTagID("leadership"),
			},
		},
		{
			Content: `Dear Recruiters,
When you say "We're like a startup, but with the stability of a big company,"
Do you mean: chaotic, underpaid, but with more meetings?

Asking for a friend who's crying into his YAML files.`,
			TagIDs: []common.VTagID{},
		},
		{
			Content: `Jenkins failed. Docker misbehaved. Kubernetes restarted everything.
I asked for logs, it gave me 12,000 lines of vibes.

Anyway, I'm a backend engineer now, but spiritually, I'm frontend. I just pretend everything works.`,
			TagIDs: []common.VTagID{},
		},
		{
			Content: `As a $MINORITY in Tech, My big contributions to Tech are giving keynote talks about improving $MINORITY contributions to Tech, while somehow sidelining the actual $MINORITY workers doing real Technology work. Evangelism beats Engineering.`,
			TagIDs: []common.VTagID{
				common.VTagID("diversity-and-inclusion"),
			},
		},
		{
			Content: "Vibe Coding is meh. Vibe Debugging, that would be a killer.",
			TagIDs: []common.VTagID{
				common.VTagID("artificial-intelligence"),
			},
		},
	}

	broadAreaPostsMap["Finance"] = []hub.AddPostRequest{
		{
			Content: `My portfolio went up 0.01% today. Already ordered the yacht. Financial freedom, baby!`,
			TagIDs: []common.VTagID{
				common.VTagID("finance"),
			},
		},
		{
			Content: `Explained "synergy" using only Excel functions. Client was impressed. I think. They just stared blankly for 5 minutes.`,
			TagIDs: []common.VTagID{
				common.VTagID("finance"),
			},
		},
		{
			Content: `Survived another earnings call where the only thing higher than projections was my blood pressure.`,
			TagIDs: []common.VTagID{
				common.VTagID("finance"),
			},
		},
		{
			Content: `"Budgeting" is just telling your money where to go instead of wondering where it went... probably to coffee and emergency snacks.`,
			TagIDs: []common.VTagID{
				common.VTagID("finance"),
			},
		},
		{
			Content: `Attended a webinar on "Blockchain for Beginners." Still a beginner. But now I'm a confused beginner with a certificate.`,
			TagIDs: []common.VTagID{
				common.VTagID("cryptocurrency"),
				common.VTagID("fintech"),
			},
		},
		{
			Content: `The market is volatile. My mood is volatile. Coincidence? I think not.`,
			TagIDs: []common.VTagID{
				common.VTagID("finance"),
			},
		},
		{
			Content: `Built a complex financial model. It accurately predicts I need more coffee.`,
			TagIDs: []common.VTagID{
				common.VTagID("finance"),
			},
		},
		{
			Content: `They say money can't buy happiness, but have you ever seen someone sad on a jet ski made of spreadsheets? Didn't think so.`,
			TagIDs: []common.VTagID{
				common.VTagID("work-life-balance"),
				common.VTagID("finance"),
			},
		},
		{
			Content: `Risk assessment complete: The biggest risk is running out of coffee before the market opens.`,
			TagIDs: []common.VTagID{
				common.VTagID("finance"),
			},
		},
		{
			Content: `Asked my AI assistant for stock tips. It suggested investing in "cat memes." Tempting, honestly.`,
			TagIDs: []common.VTagID{
				common.VTagID("artificial-intelligence"),
				common.VTagID("finance"),
			},
		},
	}

	broadAreaPostsMap["Law"] = []hub.AddPostRequest{
		{
			Content: `Just billed 0.1 hours for thinking about billing. Efficiency.`,
			TagIDs:  []common.VTagID{},
		},
		{
			Content: `My brain is 90% case law, 10% remembering where I parked.`,
			TagIDs:  []common.VTagID{},
		},
		{
			Content: `Drafted a contract so airtight, even light cannot escape it. Client asked if we could "make it less... intense?" No.`,
			TagIDs:  []common.VTagID{},
		},
		{
			Content: `"It depends" is not a non-committal answer, it's the *most* accurate legal advice you'll ever get.`,
			TagIDs:  []common.VTagID{},
		},
		{
			Content: `Successfully used "hereinafter" in casual conversation. The barista was confused. My work here is done.`,
			TagIDs:  []common.VTagID{},
		},
		{
			Content: `Attended a 4-hour deposition. Pretty sure I achieved enlightenment somewhere around hour 3. Or maybe it was just the fluorescent lights.`,
			TagIDs:  []common.VTagID{},
		},
		{
			Content: `Client: "Can we just agree to disagree?" Me: "We can agree that you sign this settlement, or we can disagree... in court."`,
			TagIDs:  []common.VTagID{},
		},
		{
			Content: `Reading legislation is my cardio.`,
			TagIDs:  []common.VTagID{},
		},
		{
			Content: `Explained 'force majeure' to my cat. He seemed unimpressed. Tough crowd.`,
			TagIDs: []common.VTagID{
				common.VTagID("remote-work"),
			},
		},
		{
			Content: `Objection! Leading the witness... to the coffee machine. Sustained.`,
			TagIDs:  []common.VTagID{},
		},
	}

	broadAreaPostsMap["Education"] = []hub.AddPostRequest{
		{
			Content: `Grading papers fueled by coffee and the faint hope that someone actually read the syllabus.`,
			TagIDs: []common.VTagID{
				common.VTagID("education"),
			},
		},
		{
			Content: `Spent an hour explaining a concept. Student asks, "Will this be on the test?" Sigh.`,
			TagIDs: []common.VTagID{
				common.VTagID("education"),
			},
		},
		{
			Content: `My superpower is deciphering handwriting that looks like abstract art.`,
			TagIDs: []common.VTagID{
				common.VTagID("education"),
			},
		},
		{
			Content: `"Summer break" is just educator code for "professional development and catching up on sleep... maybe."`,
			TagIDs: []common.VTagID{
				common.VTagID("education"),
			},
		},
		{
			Content: `Attended a faculty meeting that could have been an email. Again. Productivity!`,
			TagIDs: []common.VTagID{
				common.VTagID("education"),
				common.VTagID("productivity"),
			},
		},
		{
			Content: `The sound of the bell doesn't dismiss the class. My existential sigh does.`,
			TagIDs:  []common.VTagID{},
		},
		{
			Content: `Found a meme in a student's essay. Points for creativity? Or deduction for unprofessionalism? The struggle is real.`,
			TagIDs: []common.VTagID{
				common.VTagID("education"),
				common.VTagID("creativity"),
			},
		},
		{
			Content: `Trying to inspire the next generation while simultaneously trying to remember where I put my coffee mug.`,
			TagIDs: []common.VTagID{
				common.VTagID("education"),
			},
		},
		{
			Content: `"Let's circle back to that" - Professor code for "I have no idea, but I'll Google it during the break."`,
			TagIDs: []common.VTagID{
				common.VTagID("education"),
			},
		},
		{
			Content: `That feeling when a student finally understands a difficult concept. It's almost as good as finding a working pen. Almost.`,
			TagIDs: []common.VTagID{
				common.VTagID("education"),
			},
		},
	}

	broadAreaPostsMap["Marketing"] = []hub.AddPostRequest{
		{
			Content: `Our latest campaign generated massive buzz! Mostly from the server room fan, but still... buzz!`,
			TagIDs: []common.VTagID{
				common.VTagID("marketing"),
			},
		},
		{
			Content: `Content is king. Engagement is queen. The algorithm is the court jester messing everything up.`,
			TagIDs: []common.VTagID{
				common.VTagID("social-media-marketing"),
				common.VTagID("marketing"),
			},
		},
		{
			Content: `Added "synergize" and "leverage" to my keyword list. SEO score went up, but my soul died a little.`,
			TagIDs: []common.VTagID{
				common.VTagID("seo"),
				common.VTagID("marketing"),
			},
		},
		{
			Content: `Client: "Can we make the logo bigger?" Me: *Opens Photoshop, increases size by 1px* "How's this?" Client: "Perfect!"`,
			TagIDs: []common.VTagID{
				common.VTagID("marketing"),
				common.VTagID("graphic-design"),
			},
		},
		{
			Content: `Ran an A/B test. Button A got 10 clicks. Button B got 11 clicks. Conclusion: People randomly click things. Groundbreaking insights here.`,
			TagIDs: []common.VTagID{
				common.VTagID("marketing"),
			},
		},
		{
			Content: `My marketing funnel is more like a marketing sieve, but we're "optimizing" it.`,
			TagIDs: []common.VTagID{
				common.VTagID("marketing"),
				common.VTagID("strategy"),
			},
		},
		{
			Content: `"Going viral" sounds like a disease. And based on some comment sections, it might be.`,
			TagIDs: []common.VTagID{
				common.VTagID("social-media"),
				common.VTagID("marketing"),
			},
		},
		{
			Content: `Pretty sure half our website traffic is just me checking if the analytics tag is still working.`,
			TagIDs: []common.VTagID{
				common.VTagID("marketing"),
			},
		},
		{
			Content: `Just used AI to write ad copy. It suggested: "Buy our stuff. It is... stuff." Nailed it.`,
			TagIDs: []common.VTagID{
				common.VTagID("artificial-intelligence"),
				common.VTagID("content-creation"),
				common.VTagID("marketing"),
			},
		},
		{
			Content: `ROI isn't just Return on Investment. It's also "Really Overworked Individuals."`,
			TagIDs: []common.VTagID{
				common.VTagID("marketing"),
			},
		},
	}

	broadAreaPostsMap["Human Resources"] = []hub.AddPostRequest{
		{
			Content: `Onboarding new hires today! Step 1: Explain the coffee machine. Step 2: Everything else.`,
			TagIDs: []common.VTagID{
				common.VTagID("human-resources"),
			},
		},
		{
			Content: `"We're like a family here." Translation: Slightly dysfunctional, passive-aggressive emails are common.`,
			TagIDs: []common.VTagID{
				common.VTagID("human-resources"),
			},
		},
		{
			Content: `Updated the employee handbook. Added a clause about not microwaving fish in the breakroom. Progress.`,
			TagIDs: []common.VTagID{
				common.VTagID("human-resources"),
			},
		},
		{
			Content: `Performance review season: The annual festival of finding creative ways to say "meets expectations."`,
			TagIDs: []common.VTagID{
				common.VTagID("human-resources"),
				common.VTagID("management"),
			},
		},
		{
			Content: `Interviewing candidates. Asked one "Where do you see yourself in 5 years?" They said "Celebrating the 5-year anniversary of you asking me this question." Hired.`,
			TagIDs: []common.VTagID{
				common.VTagID("human-resources"),
				common.VTagID("job-search"),
			},
		},
		{
			Content: `Received an anonymous suggestion for "Mandatory Nap Time." Forwarding to leadership with my strongest recommendation.`,
			TagIDs: []common.VTagID{
				common.VTagID("human-resources"),
				common.VTagID("future-of-work"),
			},
		},
		{
			Content: `My job involves mediating disputes. Today's crisis: Who finished the good biscuits?`,
			TagIDs: []common.VTagID{
				common.VTagID("human-resources"),
			},
		},
		{
			Content: `Trying to foster a positive work environment while simultaneously enforcing 37 different compliance regulations. Easy peasy.`,
			TagIDs: []common.VTagID{
				common.VTagID("human-resources"),
			},
		},
		{
			Content: `Another team-building exercise successfully planned. Hope everyone likes lukewarm pizza and awkward icebreakers!`,
			TagIDs: []common.VTagID{
				common.VTagID("human-resources"),
			},
		},
		{
			Content: `"Please activate your synergy." - Just trying out phrases for the next all-hands meeting.`,
			TagIDs: []common.VTagID{
				common.VTagID("human-resources"),
				common.VTagID("communication"),
			},
		},
	}

	broadAreaPostsMap["Medical"] = []hub.AddPostRequest{
		{
			Content: `Survived a 12-hour shift fueled by caffeine, adrenaline, and the sheer terror of reading Dr. Smith's handwriting.`,
			TagIDs: []common.VTagID{
				common.VTagID("healthcare"),
			},
		},
		{
			Content: `Patient asked if WebMD was accurate. I prescribed two aspirin and a strong dose of skepticism.`,
			TagIDs: []common.VTagID{
				common.VTagID("healthcare"),
			},
		},
		{
			Content: `The hospital cafeteria coffee has medicinal properties. Specifically, the property of keeping you awake through sheer bitterness.`,
			TagIDs: []common.VTagID{
				common.VTagID("healthcare"),
			},
		},
		{
			Content: `"Stat" means "Drop everything and run," usually towards the sound of a beeping machine or an empty coffee pot.`,
			TagIDs: []common.VTagID{
				common.VTagID("healthcare"),
			},
		},
		{
			Content: `Charting: The art of describing a complex medical event in fewer characters than a tweet, but taking 10 times longer.`,
			TagIDs: []common.VTagID{
				common.VTagID("healthcare"),
			},
		},
		{
			Content: `My scrubs have seen things... unspeakable things. Mostly coffee spills, but still.`,
			TagIDs: []common.VTagID{
				common.VTagID("healthcare"),
			},
		},
		{
			Content: `Attended a medical conference. Learned about groundbreaking new treatments and confirmed the hotel biscuits are still mediocre.`,
			TagIDs: []common.VTagID{
				common.VTagID("healthcare"),
				common.VTagID("education"),
			},
		},
		{
			Content: `Trying to explain complex medical procedures using simple analogies. "Think of this artery like a really stubborn garden hose..."`,
			TagIDs: []common.VTagID{
				common.VTagID("healthcare"),
				common.VTagID("communication"),
			},
		},
		{
			Content: `The longest distance in the universe isn't between galaxies, it's between the nurses' station and a working pen.`,
			TagIDs: []common.VTagID{
				common.VTagID("healthcare"),
			},
		},
		{
			Content: `Night shift brain: Pretty sure I just tried to unlock my car with my stethoscope.`,
			TagIDs: []common.VTagID{
				common.VTagID("healthcare"),
				common.VTagID("mental-health-awareness"),
			},
		},
	}

	broadAreaPostsMap["Data Science"] = []hub.AddPostRequest{
		{
			Content: `80% of data science is cleaning data. The other 20% is complaining about cleaning data.`,
			TagIDs: []common.VTagID{
				common.VTagID("data-science"),
			},
		},
		{
			Content: `My model has an accuracy of 98%! (On the training data. Let's not talk about the test data.)`,
			TagIDs: []common.VTagID{
				common.VTagID("machine-learning"),
				common.VTagID("data-science"),
			},
		},
		{
			Content: `Built a dashboard so complex, even I don't know what it means anymore. But look at the pretty colors!`,
			TagIDs: []common.VTagID{
				common.VTagID("data-science"),
			},
		},
		{
			Content: `Feature engineering: Also known as "staring at data until it confesses."`,
			TagIDs: []common.VTagID{
				common.VTagID("data-science"),
			},
		},
		{
			Content: `Asked ChatGPT to explain my model. It gave a beautiful, confident, and completely wrong explanation. Like a tiny intern.`,
			TagIDs: []common.VTagID{
				common.VTagID("artificial-intelligence"),
				common.VTagID("data-science"),
			},
		},
		{
			Content: `The joy of finding a dataset that's already clean and perfectly formatted. ... Just kidding, that doesn't exist.`,
			TagIDs: []common.VTagID{
				common.VTagID("data-science"),
			},
		},
		{
			Content: `Stakeholder: "Can you just sprinkle some AI on this?" Me: *Adds random forest model* "Consider it sprinkled."`,
			TagIDs: []common.VTagID{
				common.VTagID("artificial-intelligence"),
				common.VTagID("data-science"),
				common.VTagID("management"),
			},
		},
		{
			Content: `My code is held together by Stack Overflow answers and sheer willpower.`,
			TagIDs: []common.VTagID{
				common.VTagID("data-science"),
				common.VTagID("software-engineering"),
			},
		},
		{
			Content: `Correlation does not imply causation, but it's really good at making graphs look convincing.`,
			TagIDs: []common.VTagID{
				common.VTagID("data-science"),
			},
		},
		{
			Content: `Successfully predicted customer churn with my model. It also predicted I'd eat pizza for dinner. 100% accuracy today!`,
			TagIDs: []common.VTagID{
				common.VTagID("data-science"),
				common.VTagID("machine-learning"),
			},
		},
	}

	broadAreaPostsMap["Consulting"] = []hub.AddPostRequest{
		{
			Content: `Just created a 100-slide deck explaining a 2-slide concept. Added value. Synergy. Impact.`,
			TagIDs:  []common.VTagID{},
		},
		{
			Content: `My travel schedule is optimized for maximum airport lounge access and minimal sleep.`,
			TagIDs: []common.VTagID{
				common.VTagID("work-life-balance"),
				common.VTagID("travel"),
			},
		},
		{
			Content: `Client asked for a "quick analysis." Three weeks and 5 frameworks later... "Here's the quick analysis."`,
			TagIDs:  []common.VTagID{},
		},
		{
			Content: `"Let's take this offline" = "I don't know the answer, but I'll find someone who does before the next meeting."`,
			TagIDs: []common.VTagID{
				common.VTagID("problem-solving"),
			},
		},
		{
			Content: `Survived another "blue sky thinking" session. Pretty sure my main contribution was suggesting more coffee.`,
			TagIDs: []common.VTagID{
				common.VTagID("creativity"),
			},
		},
		{
			Content: `Expense report submitted. Claimed "strategic alignment fuel" (it was coffee). Wish me luck.`,
			TagIDs:  []common.VTagID{},
		},
		{
			Content: `Just used a 2x2 matrix to decide what to have for lunch. Peak consulting achieved.`,
			TagIDs:  []common.VTagID{},
		},
		{
			Content: `Wearing a suit on a Zoom call from my living room. Peak professionalism or peak absurdity? Yes.`,
			TagIDs: []common.VTagID{
				common.VTagID("remote-work"),
			},
		},
		{
			Content: `The client loves the recommendations! (They were mostly the client's ideas repackaged with better graphics).`,
			TagIDs: []common.VTagID{
				common.VTagID("management"),
				common.VTagID("graphic-design"),
			},
		},
		{
			Content: `My primary skill is looking confident while presenting slides I finished 5 minutes ago.`,
			TagIDs: []common.VTagID{
				common.VTagID("public-speaking"),
			},
		},
	}

	broadAreaPostsMap["Design"] = []hub.AddPostRequest{
		{
			Content: `Client feedback: "Can you make it pop more?" *Increases saturation by 2%* Client: "Perfect!"`,
			TagIDs: []common.VTagID{
				common.VTagID("graphic-design"),
			},
		},
		{
			Content: `Spent 3 hours debating the perfect shade of grey. My life is a monochrome adventure.`,
			TagIDs: []common.VTagID{
				common.VTagID("graphic-design"),
			},
		},
		{
			Content: `My therapist told me to embrace whitespace. I told her I do, but stakeholders keep wanting to fill it with more content.`,
			TagIDs: []common.VTagID{
				common.VTagID("graphic-design"),
				common.VTagID("management"),
			},
		},
		{
			Content: `"Just a small tweak" - famous last words before a complete redesign.`,
			TagIDs: []common.VTagID{
				common.VTagID("graphic-design"),
			},
		},
		{
			Content: `Organized my layers in Photoshop. Feeling like I have my life together. It's an illusion, but a well-structured one.`,
			TagIDs: []common.VTagID{
				common.VTagID("graphic-design"),
				common.VTagID("productivity"),
			},
		},
		{
			Content: `Judging websites based solely on their font choices. It's not snobbery, it's *professional assessment*.`,
			TagIDs: []common.VTagID{
				common.VTagID("graphic-design"),
				common.VTagID("ui-design"),
			},
		},
		{
			Content: `User testing: Where you watch people completely miss the giant button you thought was obvious. Humbling.`,
			TagIDs: []common.VTagID{
				common.VTagID("user-experience"),
				common.VTagID("ui-design"),
			},
		},
		{
			Content: `My design process involves 10% inspiration, 40% perspiration, and 50% convincing people Comic Sans is not an option.`,
			TagIDs: []common.VTagID{
				common.VTagID("graphic-design"),
				common.VTagID("creativity"),
			},
		},
		{
			Content: `Exported final files. Named them Final_Final_ReallyFinal_ThisOne_v3.zip. Seems about right.`,
			TagIDs: []common.VTagID{
				common.VTagID("graphic-design"),
			},
		},
		{
			Content: `That feeling when the design just *clicks*. It's rarer than finding a unicorn riding a pixel-perfect skateboard, but it happens.`,
			TagIDs: []common.VTagID{
				common.VTagID("creativity"),
				common.VTagID("graphic-design"),
			},
		},
	}

	broadAreaPostsMap["Operations"] = []hub.AddPostRequest{
		{
			Content: `My day involves putting out fires. Sometimes metaphorical, sometimes literal (don't ask about the server room incident).`,
			TagIDs: []common.VTagID{
				common.VTagID("management"),
				common.VTagID("problem-solving"),
			},
		},
		{
			Content: `Optimized a process today. Saved 3 seconds per transaction. At scale, that's almost enough time for a coffee break!`,
			TagIDs: []common.VTagID{
				common.VTagID("productivity"),
				common.VTagID("management"),
			},
		},
		{
			Content: `Spreadsheets are my love language. Pivot tables are my sonnets.`,
			TagIDs:  []common.VTagID{},
		},
		{
			Content: `"Seamless integration" usually means "duct tape and hope."`,
			TagIDs: []common.VTagID{
				common.VTagID("technology"),
			},
		},
		{
			Content: `Supply chain issues again. Pretty sure our shipment is currently vacationing in Bermuda. Can't blame it.`,
			TagIDs: []common.VTagID{
				common.VTagID("supply-chain"),
				common.VTagID("logistics"),
			},
		},
		{
			Content: `Created a flowchart so detailed, it includes branches for existential crises during coffee breaks.`,
			TagIDs: []common.VTagID{
				common.VTagID("management"),
			},
		},
		{
			Content: `My job is to ensure things run smoothly. Which mostly means anticipating how they might spectacularly fail.`,
			TagIDs: []common.VTagID{
				common.VTagID("management"),
			},
		},
		{
			Content: `Attended a meeting about efficiency. It ran 30 minutes over schedule. The irony was noted. Silently.`,
			TagIDs: []common.VTagID{
				common.VTagID("management"),
				common.VTagID("productivity"),
			},
		},
		{
			Content: `The ops team: Unsung heroes keeping the lights on, metaphorically and sometimes literally.`,
			TagIDs: []common.VTagID{
				common.VTagID("management"),
			},
		},
		{
			Content: `Just found the bottleneck. It was me, trying to find the bottleneck. Moving on.`,
			TagIDs: []common.VTagID{
				common.VTagID("problem-solving"),
				common.VTagID("management"),
			},
		},
	}

	broadAreaPostsMap["Sales"] = []hub.AddPostRequest{
		{
			Content: `Hit quota! Time to celebrate by immediately worrying about next month's quota.`,
			TagIDs: []common.VTagID{
				common.VTagID("sales"),
			},
		},
		{
			Content: `My CRM is my best friend, my worst enemy, and the only thing that understands my obsession with pipeline velocity.`,
			TagIDs: []common.VTagID{
				common.VTagID("sales"),
				common.VTagID("technology"),
			},
		},
		{
			Content: `"Just checking in" - Polite sales code for "Please, for the love of commission, sign the contract."`,
			TagIDs: []common.VTagID{
				common.VTagID("sales"),
			},
		},
		{
			Content: `Closed a deal that's been in the works for months. Feeling like a superhero whose only power is persistent emailing.`,
			TagIDs: []common.VTagID{
				common.VTagID("sales"),
				common.VTagID("email-marketing"),
			},
		},
		{
			Content: `Survived a cold-calling session. My ears are ringing, my spirit is slightly bruised, but my resilience is +10.`,
			TagIDs: []common.VTagID{
				common.VTagID("sales"),
				common.VTagID("resilience"),
			},
		},
		{
			Content: `Prospect: "We'll review internally and get back to you." Translation: "Prepare for the follow-up abyss."`,
			TagIDs: []common.VTagID{
				common.VTagID("sales"),
			},
		},
		{
			Content: `My sales pitch is so smooth, I accidentally sold myself a pen this morning.`,
			TagIDs: []common.VTagID{
				common.VTagID("sales"),
				common.VTagID("public-speaking"),
			},
		},
		{
			Content: `Celebrating the end of the quarter like we just won the lottery. Except the prize is just... the start of the next quarter.`,
			TagIDs: []common.VTagID{
				common.VTagID("sales"),
			},
		},
		{
			Content: `Always Be Closing... the fridge door, the meeting tab, the deal. It's a lifestyle.`,
			TagIDs: []common.VTagID{
				common.VTagID("sales"),
			},
		},
		{
			Content: `They say rejection builds character. At this point, my character should be a skyscraper.`,
			TagIDs: []common.VTagID{
				common.VTagID("sales"),
				common.VTagID("personal-development"),
			},
		},
	}

	broadAreaPostsMap["Product Management"] = []hub.AddPostRequest{
		{
			Content: `Roadmap planning: The art of confidently predicting the future while knowing everything will change next week.`,
			TagIDs: []common.VTagID{
				common.VTagID("product-management"),
				common.VTagID("strategy"),
				common.VTagID("agile"),
			},
		},
		{
			Content: `My job is 50% saying "no," 40% explaining why, and 10% wondering if I should have said "yes."`,
			TagIDs: []common.VTagID{
				common.VTagID("product-management"),
				common.VTagID("decision-making"),
			},
		},
		{
			Content: `User story writing: Translating vague stakeholder wishes into something engineers can actually build, possibly with magic.`,
			TagIDs: []common.VTagID{
				common.VTagID("product-management"),
				common.VTagID("user-experience"),
			},
		},
		{
			Content: `"Let's put a pin in that" - Product Manager code for "Good idea, but it's going straight to the backlog abyss."`,
			TagIDs: []common.VTagID{
				common.VTagID("product-management"),
			},
		},
		{
			Content: `Just shipped a new feature! Now accepting bug reports and feature requests for version 2.0.`,
			TagIDs: []common.VTagID{
				common.VTagID("product-management"),
			},
		},
		{
			Content: `Attended 7 meetings today. Pretty sure I'm now qualified as a professional meeting attendee. Where's my certificate?`,
			TagIDs: []common.VTagID{
				common.VTagID("product-management"),
				common.VTagID("time-management"),
			},
		},
		{
			Content: `Trying to balance user needs, business goals, and engineering constraints. It's like juggling chainsaws, but with more spreadsheets.`,
			TagIDs: []common.VTagID{
				common.VTagID("product-management"),
				common.VTagID("management"),
			},
		},
		{
			Content: `That feeling when user feedback validates a feature you fought for. Briefly makes the endless meetings worth it. Briefly.`,
			TagIDs: []common.VTagID{
				common.VTagID("user-experience"),
				common.VTagID("product-management"),
			},
		},
		{
			Content: `My backlog is longer than a CVS receipt and possibly contains items from the Jurassic period.`,
			TagIDs: []common.VTagID{
				common.VTagID("product-management"),
				common.VTagID("agile"),
			},
		},
		{
			Content: `Explaining the product vision with passion, clarity, and a slight tremor of panic about the deadline.`,
			TagIDs: []common.VTagID{
				common.VTagID("product-management"),
				common.VTagID("communication"),
				common.VTagID("leadership"),
			},
		},
	}

	broadAreaPostsMap["Aerospace"] = []hub.AddPostRequest{
		{
			Content: `Calculated orbital mechanics before my first coffee. Just another Monday.`,
			TagIDs:  []common.VTagID{},
		},
		{
			Content: `My simulation crashed. Again. Either the physics is wrong, or the universe just enjoys messing with me.`,
			TagIDs:  []common.VTagID{},
		},
		{
			Content: `Building things that fly requires meticulous planning, precise engineering, and ignoring the little voice saying "What if it doesn't?"`,
			TagIDs: []common.VTagID{
				common.VTagID("software-engineering"),
			},
		},
		{
			Content: `"It's not rocket science." Oh, wait. Yes, it is. That's why it's taking so long.`,
			TagIDs:  []common.VTagID{},
		},
		{
			Content: `Reviewed launch readiness checklist. Item 347: Ensure coffee supply is adequate. Critical path item.`,
			TagIDs:  []common.VTagID{},
		},
		{
			Content: `Dealing with tolerances measured in microns. My patience is measured in nanometers today.`,
			TagIDs: []common.VTagID{
				common.VTagID("software-engineering"),
			},
		},
		{
			Content: `Attended a design review. Used the phrase "thrust-to-weight ratio" five times. Felt powerful.`,
			TagIDs:  []common.VTagID{},
		},
		{
			Content: `My code is designed to handle catastrophic failure scenarios. My brain, less so after a long week.`,
			TagIDs: []common.VTagID{
				common.VTagID("software-engineering"),
			},
		},
		{
			Content: `Explaining aerodynamics using hand gestures. Pretty sure I just invented a new form of interpretive dance.`,
			TagIDs: []common.VTagID{
				common.VTagID("communication-skills"),
			},
		},
		{
			Content: `That moment when the test results match the simulation. Pure, unadulterated, nerdy joy.`,
			TagIDs: []common.VTagID{
				common.VTagID("software-engineering"),
			},
		},
	}

	broadAreaPostsMap["Automotive"] = []hub.AddPostRequest{
		{
			Content: `Debugged CAN bus issues all morning. I now speak fluent hexadecimal and existential dread.`,
			TagIDs: []common.VTagID{
				common.VTagID("automotive"),
				common.VTagID("software-engineering"),
			},
		},
		{
			Content: `Designed a component to withstand extreme temperatures. Tested it by leaving my coffee mug on it. Passed.`,
			TagIDs: []common.VTagID{
				common.VTagID("automotive"),
			},
		},
		{
			Content: `Talking about torque and horsepower like it's gossip. "Did you hear about the new engine? Scandalous!"`,
			TagIDs: []common.VTagID{
				common.VTagID("automotive"),
			},
		},
		{
			Content: `"Let's just make this sensor 1mm smaller." - Famous last words before redesigning half the engine bay.`,
			TagIDs: []common.VTagID{
				common.VTagID("automotive"),
				common.VTagID("graphic-design"),
			},
		},
		{
			Content: `Attended a meeting on fuel efficiency. Drove my V8 gas guzzler home. Balance.`,
			TagIDs: []common.VTagID{
				common.VTagID("automotive"),
				common.VTagID("sustainability"),
			},
		},
		{
			Content: `My car has more lines of code than the lunar lander. And probably more bugs.`,
			TagIDs: []common.VTagID{
				common.VTagID("automotive"),
				common.VTagID("iot"),
			},
		},
		{
			Content: `Ran thermal simulations for the new EV battery pack. Conclusion: It gets hot. Groundbreaking.`,
			TagIDs: []common.VTagID{
				common.VTagID("automotive"),
				common.VTagID("renewable-energy"),
			},
		},
		{
			Content: `Trying to reduce vehicle weight. Considered replacing myself with a lighter engineer, but HR advised against it.`,
			TagIDs: []common.VTagID{
				common.VTagID("automotive"),
				common.VTagID("human-resources"),
			},
		},
		{
			Content: `Crash test ratings are important. My code's crash test rating after pulling an all-nighter? Less stellar.`,
			TagIDs: []common.VTagID{
				common.VTagID("automotive"),
				common.VTagID("software-engineering"),
			},
		},
		{
			Content: `That feeling when the prototype finally drives without catching fire. A good day in automotive.`,
			TagIDs: []common.VTagID{
				common.VTagID("automotive"),
			},
		},
	}

	broadAreaPostsMap["Hospitality"] = []hub.AddPostRequest{
		{
			Content: `"The customer is always right," except when they insist their room key opens the minibar for free.`,
			TagIDs: []common.VTagID{
				common.VTagID("customer-experience"),
			},
		},
		{
			Content: `Survived the check-in rush fueled by complimentary lobby coffee and the ability to smile while crying internally.`,
			TagIDs:  []common.VTagID{},
		},
		{
			Content: `Mastered the art of folding a fitted sheet. Next up: World peace.`,
			TagIDs:  []common.VTagID{},
		},
		{
			Content: `"Can I speak to the manager?" - Words that strike fear into the heart of every hospitality worker.`,
			TagIDs: []common.VTagID{
				common.VTagID("management"),
			},
		},
		{
			Content: `Dealing with bizarre guest requests. "Can you arrange for a unicorn to deliver my room service?" Let me check on that for you...`,
			TagIDs: []common.VTagID{
				common.VTagID("customer-experience"),
			},
		},
		{
			Content: `My steps counter goes crazy during a shift. Pretty sure I walk a marathon around this hotel daily.`,
			TagIDs: []common.VTagID{
				common.VTagID("health"),
			},
		},
		{
			Content: `That moment a guest leaves a genuinely nice review. Restores my faith in humanity (for about 5 minutes).`,
			TagIDs: []common.VTagID{
				common.VTagID("customer-experience"),
				common.VTagID("motivation"),
			},
		},
		{
			Content: `Explaining the difference between 'ocean view' and 'ocean glimpse' requires diplomatic skills worthy of the UN.`,
			TagIDs: []common.VTagID{
				common.VTagID("communication-skills"),
			},
		},
		{
			Content: `The night audit: Where time, logic, and basic math skills go on a little vacation.`,
			TagIDs:  []common.VTagID{},
		},
		{
			Content: `"Service with a smile" - even when the coffee machine exploded and there's a conga line forming at reception.`,
			TagIDs: []common.VTagID{
				common.VTagID("customer-experience"),
				common.VTagID("resilience"),
			},
		},
	}

	broadAreaPostsMap["Retail"] = []hub.AddPostRequest{
		{
			Content: `Just perfectly folded a mountain of sweaters. It will remain perfect for approximately 7 seconds.`,
			TagIDs:  []common.VTagID{},
		},
		{
			Content: `Customer: "Do you work here?" Me: *Wearing branded uniform, name tag, standing behind register* "Just visiting."`,
			TagIDs: []common.VTagID{
				common.VTagID("customer-experience"),
			},
		},
		{
			Content: `Survived the weekend rush. My feet hate me, but my ability to upsell socks is stronger than ever.`,
			TagIDs: []common.VTagID{
				common.VTagID("sales"),
			},
		},
		{
			Content: `"The item scanned at the wrong price? Let me just manually override reality for you."`,
			TagIDs:  []common.VTagID{},
		},
		{
			Content: `Working inventory day. Pretty sure I've seen boxes older than me in that stockroom.`,
			TagIDs: []common.VTagID{
				common.VTagID("logistics"),
			},
		},
		{
			Content: `Hearing holiday music in October triggers my retail PTSD. Fa-la-la-la-NO.`,
			TagIDs:  []common.VTagID{},
		},
		{
			Content: `That feeling when a customer genuinely thanks you for your help. It's like finding a $20 bill in an old coat.`,
			TagIDs: []common.VTagID{
				common.VTagID("customer-experience"),
			},
		},
		{
			Content: `Explaining the return policy for the 50th time today. My patience is also non-refundable after 30 days.`,
			TagIDs:  []common.VTagID{},
		},
		{
			Content: `Closing shift: The magical time when everything needs to be cleaned, restocked, and faced, usually by one person. Me.`,
			TagIDs:  []common.VTagID{},
		},
		{
			Content: `"Can I get a discount?" - The official soundtrack of my retail career.`,
			TagIDs: []common.VTagID{
				common.VTagID("sales"),
			},
		},
	}

	broadAreaPostsMap["Pharmaceuticals"] = []hub.AddPostRequest{
		{
			Content: `Spent the day pipetting tiny amounts of liquid. Felt like a giant playing with a very expensive, sterile dollhouse.`,
			TagIDs: []common.VTagID{
				common.VTagID("healthcare"),
			},
		},
		{
			Content: `Reading clinical trial data. Side effects may include drowsiness, dizziness, and questioning all your life choices.`,
			TagIDs: []common.VTagID{
				common.VTagID("healthcare"),
				common.VTagID("data-science"),
			},
		},
		{
			Content: `Drug naming convention meeting. Rejected "MiracleCureXtreme." Suggested "SlightlyHelpfulMaybe." We'll workshop it.`,
			TagIDs: []common.VTagID{
				common.VTagID("marketing"),
				common.VTagID("healthcare"),
			},
		},
		{
			Content: `"This formulation requires precise temperature control." *Looks nervously at the office thermostat 전쟁*`,
			TagIDs: []common.VTagID{
				common.VTagID("manufacturing"),
				common.VTagID("healthcare"),
			},
		},
		{
			Content: `Attended a regulatory affairs seminar. Learned 100 new ways to fill out forms incorrectly. Progress!`,
			TagIDs: []common.VTagID{
				common.VTagID("healthcare"),
				common.VTagID("education"),
			},
		},
		{
			Content: `My experiment failed. Again. Time to invoke the scientific method: Step 1, cry. Step 2, coffee. Step 3, try again.`,
			TagIDs: []common.VTagID{
				common.VTagID("healthcare"),
				common.VTagID("resilience"),
			},
		},
		{
			Content: `Explaining pharmacokinetics using breakfast analogies. "This drug is like slow-release oatmeal..."`,
			TagIDs: []common.VTagID{
				common.VTagID("science-communication"),
				common.VTagID("healthcare"),
			},
		},
		{
			Content: `That moment your Western Blot actually works. You feel like a wizard who just summoned a faint, blurry band.`,
			TagIDs: []common.VTagID{
				common.VTagID("biotechnology"),
				common.VTagID("healthcare"),
			},
		},
		{
			Content: `Trying to synthesize a new compound. Currently synthesizing new levels of frustration.`,
			TagIDs: []common.VTagID{
				common.VTagID("chemistry"),
				common.VTagID("healthcare"),
			},
		},
		{
			Content: `Celebrating a successful drug approval like we just cured everything. (Spoiler: We didn't, but let us have this moment).`,
			TagIDs: []common.VTagID{
				common.VTagID("healthcare"),
			},
		},
	}

	broadAreaPostsMap["Construction"] = []hub.AddPostRequest{
		{
			Content: `Project site visit today. Confirmed mud exists and hard hats mess up your hair. Crucial findings.`,
			TagIDs: []common.VTagID{
				common.VTagID("project-management"),
			},
		},
		{
			Content: `Reading blueprints that look like abstract spaghetti monsters. Pretty sure this line means "wall," maybe?`,
			TagIDs: []common.VTagID{
				common.VTagID("project-management"),
			},
		},
		{
			Content: `Meeting about budget overruns. Suggested replacing solid gold fixtures with slightly less solid gold. We'll see.`,
			TagIDs: []common.VTagID{
				common.VTagID("project-management"),
				common.VTagID("finance"),
			},
		},
		{
			Content: `"We need to value engineer this." Translation: "How can we make this cheaper without it collapsing immediately?"`,
			TagIDs: []common.VTagID{
				common.VTagID("project-management"),
				common.VTagID("software-engineering"),
			},
		},
		{
			Content: `Dealing with unexpected delays. Today's culprit: A flock of pigeons decided the scaffolding was prime real estate.`,
			TagIDs: []common.VTagID{
				common.VTagID("project-management"),
			},
		},
		{
			Content: `Safety briefing: "Don't stand under the thing being lifted." Profound stuff.`,
			TagIDs: []common.VTagID{
				common.VTagID("health"),
			},
		},
		{
			Content: `That satisfying feeling when the concrete pour goes smoothly. It's the little things... and the giant spinning truck.`,
			TagIDs:  []common.VTagID{},
		},
		{
			Content: `Coordinating subcontractors is like herding cats. Except the cats have power tools and opinions on rebar spacing.`,
			TagIDs: []common.VTagID{
				common.VTagID("project-management"),
				common.VTagID("management"),
			},
		},
		{
			Content: `Trying to explain project timelines to clients. "Yes, the building will magically appear on Tuesday, weather permitting."`,
			TagIDs: []common.VTagID{
				common.VTagID("project-management"),
				common.VTagID("communication-skills"),
			},
		},
		{
			Content: `End of the day. Covered in dust, slightly deafened, but the structure is still standing. Success!`,
			TagIDs:  []common.VTagID{},
		},
	}

	broadAreaPostsMap["Real Estate"] = []hub.AddPostRequest{
		{
			Content: `Just hosted an open house. Served artisanal cheese. Pretty sure people came for the cheese, not the house. Still counts?`,
			TagIDs: []common.VTagID{
				common.VTagID("real-estate"),
				common.VTagID("marketing"),
			},
		},
		{
			Content: `"Cozy" = Small. "Charming" = Old. "Needs TLC" = Bring a bulldozer. Mastering the art of real estate euphemisms.`,
			TagIDs: []common.VTagID{
				common.VTagID("real-estate"),
			},
		},
		{
			Content: `Negotiating offers like a high-stakes poker game, except with more paperwork and less cool sunglasses.`,
			TagIDs: []common.VTagID{
				common.VTagID("real-estate"),
				common.VTagID("sales"),
			},
		},
		{
			Content: `Showing houses all day. My car now permanently smells like air freshener and desperation.`,
			TagIDs: []common.VTagID{
				common.VTagID("real-estate"),
			},
		},
		{
			Content: `Client: "I want a 5-bedroom house with a pool, downtown, under $100k." Me: "Have you considered Mars?"`,
			TagIDs: []common.VTagID{
				common.VTagID("real-estate"),
				common.VTagID("customer-experience"),
			},
		},
		{
			Content: `The thrill of getting a signed contract! Almost makes up for the 50 unanswered calls that preceded it. Almost.`,
			TagIDs: []common.VTagID{
				common.VTagID("real-estate"),
				common.VTagID("sales"),
			},
		},
		{
			Content: `Market analysis: Prices are up! Prices are down! Prices are sideways! Basically, nobody knows, but buy now!`,
			TagIDs: []common.VTagID{
				common.VTagID("real-estate"),
				common.VTagID("finance"),
			},
		},
		{
			Content: `Staging a house: The art of making it look like nobody actually lives there, which is ironically what buyers want.`,
			TagIDs: []common.VTagID{
				common.VTagID("real-estate"),
				common.VTagID("marketing"),
			},
		},
		{
			Content: `My commission check is playing hide-and-seek. It's very good at hiding.`,
			TagIDs: []common.VTagID{
				common.VTagID("real-estate"),
				common.VTagID("finance"),
			},
		},
		{
			Content: `"Location, Location, Location!" Also important: "Paperwork, Paperwork, Paperwork!"`,
			TagIDs: []common.VTagID{
				common.VTagID("real-estate"),
			},
		},
	}

	broadAreaPostsMap["Entertainment"] = []hub.AddPostRequest{
		{
			Content: `On set today. 90% waiting, 10% frantic activity. Showbiz, baby!`,
			TagIDs: []common.VTagID{
				common.VTagID("film"),
				common.VTagID("project-management"),
			},
		},
		{
			Content: `Just read a script where the main character's motivation is "revenge... for his parking spot?" Bold choice.`,
			TagIDs: []common.VTagID{
				common.VTagID("writing"),
				common.VTagID("film"),
			},
		},
		{
			Content: `Budget meeting. Suggested cutting the craft services budget. Almost got fired. Lesson learned.`,
			TagIDs: []common.VTagID{
				common.VTagID("film-finance"),
				common.VTagID("project-management"),
			},
		},
		{
			Content: `"We'll fix it in post." - The magical phrase that solves all production problems (until post-production).`,
			TagIDs: []common.VTagID{
				common.VTagID("film"),
				common.VTagID("project-management"),
			},
		},
		{
			Content: `Casting call today. Saw 50 people convincingly pretend to be a talking squirrel. This industry is wild.`,
			TagIDs: []common.VTagID{
				common.VTagID("film"),
				common.VTagID("human-resources"),
			},
		},
		{
			Content: `Trying to secure distribution. It's easier to get a meeting with Bigfoot.`,
			TagIDs: []common.VTagID{
				common.VTagID("film-distribution"),
				common.VTagID("business-development"),
			},
		},
		{
			Content: `That feeling when the audience actually laughs at the joke you wrote. Pure, unadulterated validation.`,
			TagIDs: []common.VTagID{
				common.VTagID("writing"),
				common.VTagID("comedy"),
			},
		},
		{
			Content: `Dealing with talent agents. Requires the patience of a saint and the negotiating skills of a warlord.`,
			TagIDs: []common.VTagID{
				common.VTagID("talent-management"),
				common.VTagID("sales"),
			},
		},
		{
			Content: `Wrap party! Celebrating the end of sleep deprivation and the beginning of worrying about the reviews.`,
			TagIDs: []common.VTagID{
				common.VTagID("film"),
			},
		},
		{
			Content: `"It's got heart." Entertainment code for "The plot makes no sense, but maybe you'll cry?"`,
			TagIDs: []common.VTagID{
				common.VTagID("film"),
			},
		},
	}

	broadAreaPostsMap["Media"] = []hub.AddPostRequest{
		{
			Content: `Deadline looming. Fueled by coffee, adrenaline, and the fear of the editor's red pen.`,
			TagIDs: []common.VTagID{
				common.VTagID("journalism"),
				common.VTagID("writing"),
			},
		},
		{
			Content: `Chasing down a source who insists on speaking only in riddles. Investigative journalism or LARPing? Hard to tell.`,
			TagIDs: []common.VTagID{
				common.VTagID("journalism"),
			},
		},
		{
			Content: `Fact-checking an article. Discovered the "expert" quoted based their entire argument on a tweet they misread. Sigh.`,
			TagIDs: []common.VTagID{
				common.VTagID("journalism"),
			},
		},
		{
			Content: `"We need more clicks!" - The battle cry of modern media. Let's add a listicle about cats!`,
			TagIDs: []common.VTagID{
				common.VTagID("digital-media"),
				common.VTagID("content-creation"),
			},
		},
		{
			Content: `Attended a press conference. Got a free pen and vague non-answers. Success?`,
			TagIDs: []common.VTagID{
				common.VTagID("journalism"),
			},
		},
		{
			Content: `Trying to explain complex global events in 800 words. It's like summarizing War and Peace on a cocktail napkin.`,
			TagIDs: []common.VTagID{
				common.VTagID("writing"),
				common.VTagID("journalism"),
			},
		},
		{
			Content: `That feeling when your story gets picked up by major outlets. Briefly forget you're paid in exposure and coffee vouchers.`,
			TagIDs: []common.VTagID{
				common.VTagID("journalism"),
				common.VTagID("motivation"),
			},
		},
		{
			Content: `Dealing with angry commenters who clearly only read the headline. My block button is getting a workout.`,
			TagIDs: []common.VTagID{
				common.VTagID("social-media"),
				common.VTagID("journalism"),
			},
		},
		{
			Content: `The news never sleeps. Unfortunately, journalists do. Occasionally.`,
			TagIDs: []common.VTagID{
				common.VTagID("journalism"),
				common.VTagID("work-life-balance"),
			},
		},
		{
			Content: `"Off the record..." - Famous last words before someone tells you the juiciest story you can't publish.`,
			TagIDs: []common.VTagID{
				common.VTagID("journalism"),
			},
		},
	}

	broadAreaPostsMap["Telecommunications"] = []hub.AddPostRequest{
		{
			Content: `Traced a network outage to a squirrel chewing on a fiber optic cable. Never underestimate nature's chaos monkeys.`,
			TagIDs: []common.VTagID{
				common.VTagID("networking"),
			},
		},
		{
			Content: `Configuring routers all day. Pretty sure I can now communicate directly with machines via blinking lights.`,
			TagIDs: []common.VTagID{
				common.VTagID("networking"),
				common.VTagID("technology"),
			},
		},
		{
			Content: `Explaining bandwidth limitations to customers. "No sir, you can't download the entire internet in 5 seconds."`,
			TagIDs: []common.VTagID{
				common.VTagID("customer-experience"),
				common.VTagID("networking"),
			},
		},
		{
			Content: `"Five nines uptime" is the goal. Reality involves hoping the duct tape holds during peak hours.`,
			TagIDs: []common.VTagID{
				common.VTagID("networking"),
			},
		},
		{
			Content: `Attended a meeting about 5G deployment. Mostly understood the acronyms. Progress.`,
			TagIDs: []common.VTagID{
				common.VTagID("5g"),
				common.VTagID("technology"),
			},
		},
		{
			Content: `My troubleshooting process: 1. Reboot it. 2. Check the cables. 3. Blame the user. 4. Panic. 5. Coffee. 6. Actually fix it.`,
			TagIDs: []common.VTagID{
				common.VTagID("networking"),
				common.VTagID("problem-solving"),
			},
		},
		{
			Content: `That satisfying feeling when the signal bars go from one to full. It's like watching a tiny miracle unfold.`,
			TagIDs: []common.VTagID{
				common.VTagID("networking"),
			},
		},
		{
			Content: `Dealing with legacy systems held together by hope and undocumented Perl scripts.`,
			TagIDs: []common.VTagID{
				common.VTagID("legacy-tech"),
				common.VTagID("networking"),
			},
		},
		{
			Content: `Climbing a cell tower. Great views, slightly terrifying. Worth it for the 'gram? Debatable.`,
			TagIDs: []common.VTagID{
				common.VTagID("networking"),
				common.VTagID("travel"),
			},
		},
		{
			Content: `"The network is slow." - The universal complaint that could mean anything from sunspots to someone microwaving a burrito.`,
			TagIDs: []common.VTagID{
				common.VTagID("networking"),
				common.VTagID("problem-solving"),
			},
		},
	}

	broadAreaPostsMap["Renewable Energy"] = []hub.AddPostRequest{
		{
			Content: `Calculating solar panel efficiency. Mostly involves praying for sunny days and minimal bird poop.`,
			TagIDs: []common.VTagID{
				common.VTagID("renewable-energy"),
				common.VTagID("sustainability"),
			},
		},
		{
			Content: `Designing a wind turbine blade. It needs to be strong, efficient, and hopefully not scare the local cows.`,
			TagIDs: []common.VTagID{
				common.VTagID("renewable-energy"),
				common.VTagID("engineering-design"),
			},
		},
		{
			Content: `Site assessment for a new solar farm. Discovered the optimal location is currently occupied by very stubborn goats. Negotiations pending.`,
			TagIDs: []common.VTagID{
				common.VTagID("renewable-energy"),
			},
		},
		{
			Content: `"Grid integration" sounds simple. In reality, it's like teaching calculus to a toaster.`,
			TagIDs: []common.VTagID{
				common.VTagID("renewable-energy"),
				common.VTagID("technology"),
			},
		},
		{
			Content: `Attended a conference on battery storage. Learned that the future is bright, rechargeable, and slightly explosive if mishandled.`,
			TagIDs: []common.VTagID{
				common.VTagID("renewable-energy"),
				common.VTagID("energy-storage"),
			},
		},
		{
			Content: `Trying to explain renewable energy credits. Pretty sure I confused myself halfway through.`,
			TagIDs: []common.VTagID{
				common.VTagID("renewable-energy"),
				common.VTagID("finance"),
			},
		},
		{
			Content: `That feeling when the turbines start spinning and the power meter goes up. Saving the planet, one rotation at a time!`,
			TagIDs: []common.VTagID{
				common.VTagID("renewable-energy"),
				common.VTagID("sustainability"),
			},
		},
		{
			Content: `Dealing with intermittent power generation. The sun sets, the wind stops, my anxiety spikes. Normal Tuesday.`,
			TagIDs: []common.VTagID{
				common.VTagID("renewable-energy"),
				common.VTagID("grid-management"),
			},
		},
		{
			Content: `My job involves harnessing the power of nature. Which mostly means dealing with weather delays and unexpected wildlife encounters.`,
			TagIDs: []common.VTagID{
				common.VTagID("renewable-energy"),
				common.VTagID("sustainability"),
			},
		},
		{
			Content: `"Carbon neutral" is the goal. My coffee consumption? Less so. Baby steps.`,
			TagIDs: []common.VTagID{
				common.VTagID("sustainability"),
				common.VTagID("personal-goals"),
			},
		},
	}

}
