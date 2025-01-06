package controllers

import (
	"fmt"
	"os"
	"os/user"
)

func ReadGOFile(filename string) ([]byte, error){
	file, err := os.ReadFile(filename);
	if(err != nil){
		return nil, err;
	}
	return file, nil;
}

func GitGoCommit(commitmessage string){
	_, err := os.ReadDir(BASE_REL_PATH);
	if(err != nil){
		fmt.Println("-----GetGo Not Initialized----");
		return;
	}
	
	// load current file data 
	filename, err := ReadGOFile(BASE_REL_PATH+FILENAME_PATH_NAME)
	if(err != nil){
		fmt.Println("COMMIT_ERR: Fatal Error occured in committing")
		return
	}
	currFileData, err := ReadGOFile(string(filename));
	if(err != nil){
		fmt.Println("COMMIT_ERR: Fatal Error occured in committing")
		return
	}
	prev_hash , err := ReadGOFile(BASE_REL_PATH+"HEAD.gotem");
	if(err != nil){
		fmt.Println("COMMIT_ERR: Fatal Error occured in committing")
		return
	}
	// Check for any Change with Respect to recent Update
	author, err2 := user.Current();
	if(err2 != nil){
		fmt.Println("COMMIT_ERR: Fatal Error occured in committing")
		return
	}
	// get current Node (and branch) from HEAD's data and check if any other commits followed from here
	Tree := CommitTree{};
	Tree.ReadFromJSON(BASE_REL_PATH+COMMIT_TREE_REL_PATH+COMMIT_TREE_NAME);
	Tree.FixTree()
	currNode := Tree.SearchCommit(string(prev_hash));
	if(len(currNode.Children)!=0){
		fmt.Println("UP_TO_DATE_ERR: You are NOT Up-to-date wrt this branch, Please pull first, or create a new Branch.")
		return 
	}
	currHash := GenrateHash(currNode.Branch+string(currFileData))
	if(string(prev_hash) == currHash){
		fmt.Println("Already Up to date: No Change to Commit.")
		return
	}
	someNode := Tree.SearchCommit(currHash);
	if(someNode != nil){
		fmt.Println("IMPLEMENTATION_ERR: in my implementation, one cant commit some changes he/she has done already.\n The same changes were made in commit: "+someNode.Checksum+"\n at Time: "+someNode.Time.GoString()+" \n By Author: "+someNode.Author)
		return
	}
	cuurBranch := currNode.Branch;
	// make new Node form hash and add Node to Commitree
	newNode, err := Tree.AddCommit(currNode, currHash, author.Username, commitmessage);
	if(err != nil){
		fmt.Println("COMMIT_ERR: Fatal Error occured in committing")
		return
	}
	// Create Blob and save object at path
	blobData, err:= CompressData(currFileData);
	if(err != nil){
		fmt.Println("COMMIT_ERR: Fatal Error occured in committing")
		return
	}
	obj_dir := currHash[:2];
	createDirectory(BASE_REL_PATH+OBJECT_REL_PATH+obj_dir);
	file_name := currHash[2:];
	obj_path := BASE_REL_PATH+OBJECT_REL_PATH+obj_dir+"/"+file_name+".obj";
	err1 := os.WriteFile(obj_path, blobData, 0644);
	if err1 != nil{
		fmt.Println("COMMIT_ERR: Fatal Error occured in committing")
		return
	}
	// update current branch and HEAD Pointer
	branch_path := BASE_REL_PATH+REFS_REL_PATH+cuurBranch+".ref"
	err3 := os.WriteFile(branch_path, []byte(currHash), 0644);
	if(err3 != nil){
		fmt.Println("COMMIT_ERR: Fatal Error occured in committing")
		return
	}
	Head_path := BASE_REL_PATH+"HEAD.gotem"
	err4 := os.WriteFile(Head_path, []byte(currHash), 0644);
	if(err4 != nil){
		fmt.Println("COMMIT_ERR: Fatal Error occured in committing")
		return
	}
	// write tree to JSON
	err5 := Tree.WriteToJSON(BASE_REL_PATH+COMMIT_TREE_REL_PATH+COMMIT_TREE_NAME);
	if(err5 != nil){
		fmt.Println("COMMIT_ERR: Fatal Error occured in committing")
		fmt.Println(err);
		return
	}
	fmt.Printf("Commit Made! Here are the Details:\n Checksum : %s\n Time : %v\n Authored by: %s\n", newNode.Checksum, newNode.Time, newNode.Author);	
}