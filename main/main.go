package main

import (
	"fmt"
	"ostrichdb-go/src/lib"
	c "client" //uncomment when ready to use client functions
	sdk "ostrichdb-go/src/sdk"
)

var config lib.Config = *sdk.NewConfigBuilder()
var client lib.Client = *sdk.NewClientBuilder(&config)
var project lib.Project =  *sdk.NewProjectBuilder(&client, "twix")


func main (){
	var err error

	// START | START | START | START | START | START | START | START | START
	// The following code is an example of the flow for creating a
	// Project(Dir of Users), Specific User Collection, and a first Post(Cluster)
	// Once 'END' is reached an exmaple of another flow begins

    err= sdk.CreateProject(&project)
    if err != nil {
     	fmt.Println("Error creating project")
      	fmt.Println("Error: ",err)
    }

    // //Note: Whenever implementing user auth on client please replace "Marshall" with username input
    collection:= sdk.NewCollectionBuilder(&project, "Marshall")
    err = sdk.CreateCollection(collection)

    if err != nil {
     	fmt.Println("Error creating Collection")
      	fmt.Println("Error: ",err)
    }

    //Testing things out
    currentUser:= c.NewUserBuilder("@Marshall")
    //TODO: If there is a space in the string value of Content is not stored. OstrichDB engine issue
    newPost:= c.NewPostBuilder(&currentUser, "Hello World")

    err = c.HandlePost(*collection, newPost)
    if err != nil {
   		fmt.Println("Error appending new post to a Collection")
    	fmt.Println("Error: ",err)
    }
    // END | END| END| END| END| END| END| END| END| END



    // Fetching a Post from OstrichDB
    specificPost:= sdk.NewClusterBuilder(collection, "Marshall")
    fmt.Println(c.FetchPostAuthor(specificPost, 1))
    // c.Run()

}

