// Structure of .gitGo folder (in my implementation)
// .gitGo
//
//	|-->filename.gotem
//	|-->commits
//	|---->tree.gotem (Commit DA-Graph)
//	|-->objects
//	|---->(value-value store)
//	|-->refs
//	|---->Branch Pointers
//	|-->HEAD.gotem (Current Pointer)
package main


// TODO: 
	// implement status with branch info as well, like x commits ahead y commits behind
	// Implement better Branch Tree(now as a Graph)

import (
	"flag"
	"fmt"
	"os"
	"strconv"
	"github.com/Meh-Mehul/GetGo/controllers"
)
func main() {
	if len(os.Args) < 2 {
		fmt.Println("Error: Please specify a command (e.g., init, commit, stash).")
		os.Exit(1)
	}
	QueryType := os.Args[1]
	switch QueryType {
	case "init":
		if len(os.Args) < 3 {
			fmt.Println("INITIALIZATION_ERR: Please provide the name of the file you want to track.")
			return
		}
		path, err := os.Getwd()
		if err != nil {
			fmt.Println("Error: Unable to get the current working directory.")
			return
		}
		filename := os.Args[2]
		controllers.GitGoInit(path, filename)
	case "commit":
		commitCmd := flag.NewFlagSet("commit", flag.ExitOnError)
		commitMsg := commitCmd.String("m", "", "The message to be sent in the commit.")
		if len(os.Args) > 2 {
			commitCmd.Parse(os.Args[2:])
		}
		if *commitMsg == "" {
			fmt.Println("Error: -m flag is required to specify a commit message.")
			os.Exit(1)
		}
		controllers.GitGoCommit(*commitMsg)
	case "stash":
		if len(os.Args) == 2 {
			controllers.GitGoStash_noarg()
		} 
		if len(os.Args)>2 {
			stashArg := os.Args[2]
			switch stashArg {
			case "apply":
				controllers.GitGoStashApply()
			case "list":
				controllers.GitGoStashList()
			case "clear":
				controllers.GitGoStashClear()
			default:
				fmt.Printf("Error: Unknown stash command '%s'. Supported commands are: apply, list, clear.\n", stashArg)
				os.Exit(1)
			}
		}
	case "log":
		if len(os.Args) == 2 {
			controllers.GitGoLog_noarg()
		} 
		if len(os.Args)>3 {
			stashArg := os.Args[2]
			switch stashArg {
			case "since":
				since := os.Args[3]/// format is "01 Jan 15 10:00 UTC"
				controllers.GitGoLogSince(since)
			case "before":
				before := os.Args[3]/// format is "01 Jan 15 10:00 UTC"
				controllers.GitGoLogBefore(before)
			case "n":
				n, err := strconv.Atoi(os.Args[3]);
				if(err != nil){
					fmt.Println("Please Provide a Valid Integer")
				} else{
					controllers.GitGoLogN(n);
				}
			default:
				fmt.Printf("Error: Unknown stash command '%s'. Supported commands are: apply, list.\n", stashArg)
				os.Exit(1)
			}
		}
	case "status" :
		controllers.GitGoStatus()
	case "revert" :
		if len(os.Args) == 2 {
			controllers.GitGoRevert_Noarg()
		} 
		if len(os.Args)==3 {
			revertArgs := os.Args[2]
			
			controllers.GitGoRevert(revertArgs);
		}
	case "branch":
		if len(os.Args) == 2 {
			controllers.GitGoBranch_Noarg()
		} 
		if len(os.Args)==3 {
			bargs := os.Args[2]
			
			controllers.GitGoBranch(bargs);
		}
	case "get":
		if len(os.Args) == 3 {
			gargs := os.Args[2];
			controllers.GitGoGet(gargs);
		} else{
			fmt.Println("Error: Please give checksum of the commit you want to get")
		}
	case "diff":
		if(len(os.Args) != 4){
			fmt.Println("Usage of diff is gitgo diff <file1> <file2>")
			return
		}
		f1 := os.Args[2];
		f2 := os.Args[3];
		controllers.GitGoDiff(f1, f2);
	default:
		fmt.Printf("Error: Unknown command '%s'. Supported commands are: init, commit, stash, log, status, revert, branch.\n", QueryType)
		os.Exit(1)
	}
}






