package client

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
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

//Used to fetch an entire Post (p) which is stored as a single cluster within OstrichDB
// Even though this function claims to fetch a post by id. its really looking for the post's name
// THis is the case becauser the ID is infact inside the name
func FetchPostByID(c *lib.Cluster, id uint64) Post{
	var post Post
	client:= c.Collection.Project.Client
	projName:= c.Collection.Project.Name
	colName:= c.Collection.Name
	cluName:= fmt.Sprintf("Post%d", id)

	path:= fmt.Sprintf("%s/projects/%s/collections/%s/clusters/%s", lib.OSTRICHDB_ADDRESS, projName, colName, cluName )

	response, err:= lib.Get(client, path)
	if err != nil {
		return post
	}

	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return post
	}

	// Read response body
	body, err := io.ReadAll(response.Body)
	if err != nil {
		return post
	}

	// Unmarshal the cluster JSON response
	var clusterResp ClusterResponse
	err = json.Unmarshal(body, &clusterResp)
	if err != nil {
		return post
	}

	// Convert ClusterResponse to Post
	fetchedPost, err := clusterResp.ToPost()
	if err != nil {
		return post
	}

	return *fetchedPost
}

// Fetch the number of likes for a specific post
func FetchNumOfLikes(c *lib.Cluster, id uint64) uint64 {
	post := FetchPostByID(c,id)
	return post.NumOfLikes
}

// Fetch the number of comments for a specific post
func FetchNumOfComments(c *lib.Cluster, id uint64) uint64 {
	post := FetchPostByID(c,id)
	return post.NumOfComments
}

// Fetch the list of users who liked a specific post
func FetchPostLikesList(c *lib.Cluster, id uint64) []string {
	post := FetchPostByID(c,id)
	return post.WhoLikedPost
}

// Fetch the list of users who commented on a specific post
func FetchPostCommentersList(c *lib.Cluster, id uint64) []string {
	post := FetchPostByID(c,id)
	return post.WhoCommented
}

// Fetch the author of a specific post
func FetchPostAuthor(c *lib.Cluster, id uint64) string {
	post := FetchPostByID(c,id)
	if post.Author != nil {
		return post.Author.Handle
	}
	return "Unknown Author"
}

// Fetch the actual comments (Post objects) for a specific post
func FetchPostCommnets(c *lib.Cluster, id uint64) []Post {
	post := FetchPostByID(c,id)
	return post.Comments
}