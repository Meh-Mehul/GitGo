package controllers

// This file Stores diff Algorithms to track changes in a file

import (
	"fmt"
	"os"
	"time"
	"github.com/fatih/color"
	"github.com/sergi/go-diff/diffmatchpatch"
)

func GitGoStatus(){
	_, err := os.ReadDir(BASE_REL_PATH);
	if(err != nil){
		fmt.Println("-----GetGo Not Initialized----");
		return;
	}
	// fetch filename and currHash
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
	// make and fix tree, and get root and currNode
	Tree := CommitTree{}
	Tree.ReadFromJSON(BASE_REL_PATH+COMMIT_TREE_REL_PATH+COMMIT_TREE_NAME);
	currNode := Tree.SearchCommit(string(currHash));
	if(currNode == nil){
		fmt.Println("Some Error Occured While Fetching commit Data.");
		return;
	}
	currFileData, err := ReadGOFile(string(filename));
	if(err != nil){
		fmt.Println("COMMIT_ERR: Fatal Error occured in committing")
		return
	}
	// check for changes, and time since last commit, owner of og commit
	flag := true;
	if(GenrateHash(currNode.Branch+string(currFileData)) == string(currHash)){
		flag = false;
	}
	ogOwner := Tree.Root.Author;
	elapsed := time.Since(currNode.Time)
	// return all in a good manner(hard).
	fmt.Printf("Current Status of %s file\n Original Owner: %s\n Commits(all): %d\n Changed?: %v\n ",filename, ogOwner, Tree.lastID, flag)
	if(flag){
		fmt.Printf("Time since Last Change was Committed: %v\n", elapsed.String());
	}
}
func GitGoDiff(file1, file2 string) {
	data1, err := os.ReadFile(file1)
	if err != nil {
		fmt.Printf("Error reading file %s: %v\n", file1, err)
	}
	data2, err := os.ReadFile(file2)
	if err != nil {
		fmt.Printf("Error reading file %s: %v\n", file2, err)
	}
	dmp := diffmatchpatch.New()
	diffs := dmp.DiffMain(string(data1), string(data2), false)
	addColor := color.New(color.FgGreen).SprintFunc()
	delColor := color.New(color.FgRed).SprintFunc()
	neutralColor := color.New(color.FgWhite).SprintFunc()
	fmt.Println(neutralColor("diff --git a/" + file1 + " b/" + file2))
	fmt.Println(neutralColor("--- a/" + file1))
	fmt.Println(neutralColor("+++ b/" + file2))
	for _, diff := range diffs {
		switch diff.Type {
		case diffmatchpatch.DiffInsert:
			fmt.Println(addColor("+" + diff.Text))
		case diffmatchpatch.DiffDelete:
			fmt.Println(delColor("-" + diff.Text))
		case diffmatchpatch.DiffEqual:
			fmt.Println(neutralColor(" " + diff.Text))
		}
	}
}