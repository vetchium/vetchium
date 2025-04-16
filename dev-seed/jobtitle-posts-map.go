package main

import (
	"github.com/vetchium/vetchium/typespec/common"
	"github.com/vetchium/vetchium/typespec/hub"
)

var jobTitlePostsMap map[string][]hub.AddPostRequest

func init() {
	jobTitlePostsMap = map[string][]hub.AddPostRequest{}

	jobTitlePostsMap["Software Engineer"] = []hub.AddPostRequest{
		{
			Content: "I'm a software engineer",
			NewTags: []common.VTagName{
				common.VTagName("Software Engineer"),
			},
		},
	}
}
