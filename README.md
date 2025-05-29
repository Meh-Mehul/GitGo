# GitGo
## Introduction
This is a simple Git-like File Tracker i made in Golang. It supports almost all of the fundamental features of a git-like system, but for a single file.

## Usage
1. ```gitGo.exe init <filename>``` to intialize a tracker for a file. Currently a single tracker per directory is supported.
2. ```gitGo.exe commit -m "Message"``` to commit change. (Note that since im tracking a single file, i only need to commit as there are no extra files to track).
3. ```gitGo.exe stash [args]``` to stash changes. Almost works like git.

4. ```gitGo.exe log [args]``` to log commits and get their hashes.

5. ```gitGo.exe status``` to get status of current HEAD pointer. (I might implement Branch Ancestor Info later on too).
6. ```gitGo.exe revert [args]``` to revert to previous commits, and rollback their changes completely.
7. ```gitGo.exe branch [name]``` to switch branches.
8. ```gitGo.exe get [checksum]``` to get the filedata at a specific commit.
9. ```gitGo.exe diff [file1] [file2]``` to get diff data between two files.

## Notes:
1. This implementation does NOT support merging two branches yet because i have implemented the Commits in form of a Tree and not a graph.
2. This is completely brute-force storing of compressed file-objects (and not deltas).
