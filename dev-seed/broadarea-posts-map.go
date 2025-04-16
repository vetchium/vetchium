package main

import (
	"github.com/vetchium/vetchium/typespec/common"
	"github.com/vetchium/vetchium/typespec/hub"
)

var broadAreaPostsMap map[string][]hub.AddPostRequest

func init() {
	broadAreaPostsMap = map[string][]hub.AddPostRequest{}

	broadAreaPostsMap["Software Engineer"] = []hub.AddPostRequest{
		{
			Content: "I'm a software engineer",
			NewTags: []common.VTagName{
				common.VTagName("Software Engineer"),
			},
		},
	}
}
