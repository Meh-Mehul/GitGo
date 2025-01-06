package controllers

import (
	"fmt"
	"os"
	"time"
)

func GitGoLog_noarg(){
	// fetch current commit checksum from the HEAD.gotem file
	_, err := os.ReadDir(BASE_REL_PATH);
	if(err != nil){
		fmt.Println("-----GetGo Not Initialized----");
		return;
	}
	filename, err :=os.ReadFile(BASE_REL_PATH+FILENAME_PATH_NAME);
	if(err != nil){
		fmt.Println("Some Error Occured While Fetching commit Data.");
		return;
	}
	currHash, err := os.ReadFile(BASE_REL_PATH+"HEAD.gotem")
	if(err != nil){
		fmt.Println("Some Error Occured While Fetching commit Data.");
		return;
	}
	// Read the current commitTree
	Tree := CommitTree{}
	Tree.ReadFromJSON(BASE_REL_PATH+COMMIT_TREE_REL_PATH+COMMIT_TREE_NAME);
	// Search for current commmit
	currNode := Tree.SearchCommit(string(currHash));
	if(currNode == nil){
		fmt.Println("Some Error Occured While Fetching commit Data.");
		return;
	}
	// display current commit's data if found
	fmt.Printf("GetGo File Tracker:\n FileName: %s\n Commit checksum: %s\n Commit Message: %s\n Author: %s\n Last Commit Time: %s\n No. of Commits(All Branches): %d\n", filename, currHash, currNode.CommitMessage, currNode.Author, currNode.Time, Tree.lastID);
}


func GitGoLogSince(argTime string){
	_, err := os.ReadDir(BASE_REL_PATH);
	if(err != nil){
		fmt.Println("-----GetGo Not Initialized----");
		return;
	}
	currTime,_ := time.Parse(time.RFC822, argTime);
	Tree := CommitTree{}
	err1 := Tree.ReadFromJSON(BASE_REL_PATH+COMMIT_TREE_REL_PATH+COMMIT_TREE_NAME);
	if(err1 != nil){
		fmt.Println("Some Error Occured While Fetching commit Data.");
		return;
	}
	// fetch the commit root
	outputList := []*Node{}
	// Recursively Search for the Commit (dfs) until the current time is not met
	stack := []*Node{Tree.Root}
	for len(stack) > 0 {
		curr := stack[len(stack)-1]
		stack = stack[:len(stack)-1]
		if(currTime.Before(curr.Time)){
			outputList = append(outputList, curr);
		}
		for i := len(curr.Children) - 1; i >= 0; i-- {
			stack = append(stack, curr.Children[i])
		}
	}
	if(len(outputList)==0){
		fmt.Println("No Commits Were Made since then.");
		return;
	}
	// Store them in an array and return
	fmt.Println("Here are the Commit Made Before "+string(currTime.String()));
	for _,commit:=range outputList{
		fmt.Printf("Date: %v .\n Branch: %s\n Checksum: %s\n",commit.Time,commit.Branch, commit.Checksum);
	}
}


func GitGoLogBefore(argTime string){
	_, err := os.ReadDir(BASE_REL_PATH);
	if(err != nil){
		fmt.Println("-----GetGo Not Initialized----");
		return;
	}
	currTime,_ := time.Parse(time.RFC822, argTime);
	Tree := CommitTree{}
	err1 := Tree.ReadFromJSON(BASE_REL_PATH+COMMIT_TREE_REL_PATH+COMMIT_TREE_NAME);
	if(err1 != nil){
		fmt.Println("Some Error Occured While Fetching commit Data.");
		return;
	}
	// fetch the commit root
	outputList := []*Node{}
	// Recursively Search for the Commit (dfs) until the current time is not met
	stack := []*Node{Tree.Root}
	for len(stack) > 0 {
		curr := stack[len(stack)-1]
		stack = stack[:len(stack)-1]
		if(currTime.After(curr.Time)){
			outputList = append(outputList, curr);
		}
		for i := len(curr.Children) - 1; i >= 0; i-- {
			stack = append(stack, curr.Children[i])
		}
	}
	if(len(outputList)==0){
		fmt.Println("No Commits Were Made before that.");
		return;
	}
	// Store them in an array and return
	fmt.Println("Here are the Commit Made Before "+string(currTime.GoString()));
	for _,commit:=range outputList{
		fmt.Printf("Date: %v .\n Branch: %s\n Checksum: %s\n",commit.Time,commit.Branch, commit.Checksum);
	}
}

func GitGoLogN(n int){
	_, err := os.ReadDir(BASE_REL_PATH);
	if(err != nil){
		fmt.Println("-----GetGo Not Initialized----");
		return;
	}
	currHash, err := os.ReadFile(BASE_REL_PATH+"HEAD.gotem")
	if(err != nil){
		fmt.Println("Some Error Occured While Fetching commit Data.");
		return;
	}
	Tree := CommitTree{}
	err1 := Tree.ReadFromJSON(BASE_REL_PATH+COMMIT_TREE_REL_PATH+COMMIT_TREE_NAME);
	if(err1 != nil){
		fmt.Println("Some Error Occured While Fetching commit Data.");
		return;
	}
	// fmt.Println("Da")
	outputList := []*Node{};
	Tree.FixTree();
	currNode := Tree.SearchCommit(string(currHash));
	if(currNode == nil){
		fmt.Println("Some Error Occured While Fetching commit Data.");
		return;
	}	
	for currNode != nil {
		outputList = append(outputList, currNode)
		currNode = currNode.Parent
	}
	if(len(outputList)==0){
		fmt.Println("No Commits Were Made.");
		return;
	}
	// Store them in an array and return
	fmt.Printf("Here are the last %d commits:\n", min(len(outputList), n));
	for _,commit:=range outputList{
		fmt.Printf("Date: %v .\n Branch: %s\n Checksum: %s\n",commit.Time,commit.Branch, commit.Checksum);
	}

}

