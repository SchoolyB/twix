package main

import (
	"fmt"

	// c "client" //uncomment when ready to use client functions
	"ostrichdb-go/src/lib"
	sdk "ostrichdb-go/src/sdk"
)

func main (){
    config:= sdk.NewConfigBuilder()
    client:= sdk.NewClientBuilder(config)

    project:= sdk.NewProjectBuilder(client, "twix")
    projectErr:= sdk.CreateProject(project)
    if projectErr != nil {

     	fmt.Println("Error creating project")
      	fmt.Println(projectErr)
    }

    //Note: Whenever implementing user auth on client please replace "Marshall" with username input
    collection:= sdk.NewCollectionBuilder(project, "Marshall")
    collectionErr:= sdk.CreateCollection(collection)

    if collectionErr != nil {
     	fmt.Println("Error creating collection")
    }


    //Any time the user creates a Post, a new cluster will be created
    newPostName:=create_post_name(collection)
    cluster:= sdk.NewClusterBuilder(project, collection, newPostName)
    sdk.CreateCluster(cluster)

}


//Helper functions for client will move later - Marshalla
func increment_post_count(collection *lib.Collection) int{
	currentCount:= sdk.GetClusterCount(collection)
	incremented := currentCount + 1
	return incremented
}


//Only used to create a name for the post within OstrichDB e.g Post1, Post2, etc
func create_post_name(collection *lib.Collection) string{
	postCount:= increment_post_count(collection)
	postName:= fmt.Sprintf("Post%d", postCount)
	return postName
}