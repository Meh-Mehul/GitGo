package controllers

import (
	"fmt"
	"os"
)
func fileExists(path string) bool {
	_, err := os.Stat(path)
	return !os.IsNotExist(err)
}

func GitGoStash_noarg(){
	// get current file data, and hash it
	_, err := os.ReadDir(BASE_REL_PATH);
	if(err != nil){
		fmt.Println("-----GetGo Not Initialized----");
		return;
	}
	// load current file data 
	filename, err := ReadGOFile(BASE_REL_PATH+FILENAME_PATH_NAME)
	if(err != nil){
		fmt.Println("FATAL_ERR: Fatal Error occured in Stashing")
		return
	}
	currFileData, err := ReadGOFile(string(filename));
	if(err != nil){
		fmt.Println("FATAL_ERR: Fatal Error occured in Stashing")
		return
	}
	prev_hash , err := ReadGOFile(BASE_REL_PATH+"HEAD.gotem");
	if(err != nil){
		fmt.Println("FATAL_ERR: Fatal Error occured in stashing")
		return
	}
	Tree := CommitTree{};
	Tree.ReadFromJSON(BASE_REL_PATH+COMMIT_TREE_REL_PATH+COMMIT_TREE_NAME);
	currNode := Tree.SearchCommit(string(prev_hash));
	currHash := GenrateHash(currNode.Branch+string(currFileData));
	if(currHash == string(prev_hash)){
		fmt.Println("Nothing to Stash, Already at latest commit.");
		return
	}
	_, err1 := os.ReadDir(BASE_REL_PATH+STASH_REL_PATH);
	if(err1 != nil){
		createDirectory(BASE_REL_PATH+STASH_REL_PATH);
	}
	stackFilePath := BASE_REL_PATH+STASH_REL_PATH+STACK_FILE_NAME
	stack := Stack{}
	if(fileExists(stackFilePath)){
		stack.ReadFromJSON(stackFilePath);
	}

	stack.Push(currHash)
	blobData, err:= CompressData(currFileData);
	if(err != nil){
		fmt.Println("FATAL_ERR: Fatal Error occured in stashing")
		return
	}
	obj_dir := currHash[:2];
	createDirectory(BASE_REL_PATH+OBJECT_REL_PATH+obj_dir);
	file_name := currHash[2:];
	obj_path := BASE_REL_PATH+OBJECT_REL_PATH+obj_dir+"/"+file_name+".obj";
	err2 := os.WriteFile(obj_path, blobData, 0644);
	if err2 != nil{
		fmt.Println("FATAL_ERR: Fatal Error occured in stashing")
		return
	}
	tobeBlobData, err  := getBlobFromHash(string(prev_hash));
	if(err != nil){
		fmt.Println("FATAL_ERR: Fatal Error occured in stashing")
		return
	}
	tobeFileData, err := DecompressData(tobeBlobData);
	if(err != nil){
		fmt.Println("FATAL_ERR: Fatal Error occured in stashing")
		return
	}
	file, err := os.OpenFile(string(filename), os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0755)
	if(err != nil){
		fmt.Println("FATAL_ERR: Fatal Error occured in stashing")
		return
	}
	defer file.Close()
	_, err3 := file.Write(tobeFileData);
	if(err3 != nil){
		fmt.Println("FATAL_ERR: Fatal Error occured in stashing")
		return
	}
	stack.WriteToJSON(stackFilePath);
	fmt.Println("You are back to last Commit.")
	fmt.Printf("Stash Hash: %s\n", currHash);
}	

func GitGoStashList(){
	_, err := os.ReadDir(BASE_REL_PATH);
	if(err != nil){
		fmt.Println("-----GetGo Not Initialized----");
		return;
	}
	stackFilePath := BASE_REL_PATH+STASH_REL_PATH+STACK_FILE_NAME
	stack := Stack{}
	if(!fileExists(stackFilePath)){
		fmt.Println("Nothing to Stash, working tree clean.")
		return;
	}
	stack.ReadFromJSON(stackFilePath);
	fmt.Printf("Currenlty Stored %d Stashed Objects.\n", len(stack.Stack));
}

func GitGoStashClear(){
	_, err := os.ReadDir(BASE_REL_PATH);
	if(err != nil){
		fmt.Println("-----GetGo Not Initialized----");
		return;
	}
	stackFilePath := BASE_REL_PATH+STASH_REL_PATH+STACK_FILE_NAME
	stack := Stack{}
	if(!fileExists(stackFilePath)){
		fmt.Println("Working tree cleaned.")
		return;
	}
	stack.ReadFromJSON(stackFilePath);
	commit_arr := stack.Stack;
	for _, commitHash := range commit_arr{
		err := DeleteDirectory(commitHash);
		if(err != nil){
			fmt.Printf("Error Removing %s Stash\n", commitHash);
		}
	}
	emptyStack := Stack{}
	emptyStack.WriteToJSON(stackFilePath);
	fmt.Println("Working tree cleaned.")
}

// i have set it to apply the recent change in the stash
func GitGoStashApply(){
	_, err := os.ReadDir(BASE_REL_PATH);
	if(err != nil){
		fmt.Println("-----GetGo Not Initialized----");
		return;
	}
	filename, err := ReadGOFile(BASE_REL_PATH+FILENAME_PATH_NAME)
	if(err != nil){
		fmt.Println("FATAL_ERR: Fatal Error occured in Stashing")
		return
	}
	stackFilePath := BASE_REL_PATH+STASH_REL_PATH+STACK_FILE_NAME
	stack := Stack{}
	if(!fileExists(stackFilePath)){
		fmt.Println("Nothing was Stashed, so nothing to revert to.  :(")
		return
	}
	stack.ReadFromJSON(stackFilePath);
	if(stack.isEmpty()){
		fmt.Println("Nothing was Stashed, so nothing to revert to.  :(")
		return
	}
	stashHash, err := stack.Pop()
	if(err != nil){
		fmt.Println("FATAL_ERR: An error occured while stashing");
		return
	}
	stashBlobData, err := getBlobFromHash(stashHash);
	if(err != nil){
		fmt.Println("FATAL_ERR: An error occured while stashing");
		return
	}
	decomp_stashBlobData, err := DecompressData(stashBlobData);
	if(err != nil){
		fmt.Println("FATAL_ERR: An error occured while stashing");
		return
	}
	file, err := os.OpenFile(string(filename), os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0755)
	if(err != nil){
		fmt.Println("FATAL_ERR: Fatal Error occured in stashing")
		return
	}
	defer file.Close()
	_, err3 := file.Write(decomp_stashBlobData);
	if(err3 != nil){
		fmt.Println("FATAL_ERR: Fatal Error occured in stashing")
		return
	}
	err1 := DeleteDirectory(stashHash);
	if(err1 != nil){
		fmt.Printf("Error Removing %s Stash\n", stashHash);
	}
	fmt.Printf("Returned to the Stash: %s.\n", stashHash);
	stack.WriteToJSON(stackFilePath);
}