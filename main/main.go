package main

import (
	"fmt"

	// c "client" //uncomment when ready to use client functions
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

}