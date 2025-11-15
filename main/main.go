package main

import (
	c "client" //uncomment when ready to use client functions
	// "fmt"
	"ostrichdb-go/src/lib"
	sdk "ostrichdb-go/src/sdk"
)
var config lib.Config = *sdk.NewConfigBuilder()
var client lib.Client = *sdk.NewClientBuilder(&config)
var project lib.Project =  *sdk.NewProjectBuilder(&client, "twix")


func main (){

    // projectErr:= sdk.CreateProject(&project)
    // if projectErr != nil {
    //  	fmt.Println("Error creating project")
    //   	fmt.Println(projectErr)
    // }

    //Note: Whenever implementing user auth on client please replace "Marshall" with username input
    collection:= sdk.NewCollectionBuilder(&project, "Marshall")
    // collectionErr:= sdk.CreateCollection(collection)

    // if collectionErr != nil {
    //  	fmt.Println("Error creating collection")
    // }

    //Testing things out
    currentUser:= c.NewUserBuilder("@Marshall")
    newPost:= c.NewPostBuilder(&currentUser, "Hello World")
    c.HandlePost(project, *collection, newPost)


    // c.Run()

}

