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

Remember: It’s not about how many bugs you fix, it’s about how confidently you explain why they aren’t your fault.`,
			NewTags: []common.VTagName{
				common.VTagName("Leadership"),
				common.VTagName("Senior Engineer"),
				common.VTagName("Fake It Till You Deploy It"),
			},
		},
		{
			Content: `Dear Recruiters,
When you say “We’re like a startup, but with the stability of a big company,”
Do you mean: chaotic, underpaid, but with more meetings?

Asking for a friend who’s crying into his YAML files.`,
			NewTags: []common.VTagName{
				common.VTagName("Life Lessons"),
				common.VTagName("Points to Ponder"),
			},
		},
		{
			Content: `Jenkins failed. Docker misbehaved. Kubernetes restarted everything.
I asked for logs, it gave me 12,000 lines of vibes.

Anyway, I’m a backend engineer now, but spiritually, I’m frontend. I just pretend everything works.`,
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
}
