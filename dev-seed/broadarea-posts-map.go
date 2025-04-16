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
	}
}
