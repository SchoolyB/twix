package client

import (
	"fmt"
	"ostrichdb-go/src/lib"
	sdk "ostrichdb-go/src/sdk"
)

// Helper functions for client will move later - Marshalla
func increment_post_count(collection *lib.Collection) int{
	currentCount:= sdk.GetClusterCount(collection)
	incremented := currentCount + 1
	return incremented
}


// Only used to create a name for the post within OstrichDB e.g Post1, Post2, etc
func create_post_name(collection *lib.Collection) string{
	postCount:= increment_post_count(collection)
	postName:= fmt.Sprintf("Post%d", postCount)
	return postName
}

