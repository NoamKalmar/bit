package files

import "path/filepath"

var DEFAULT_BRANCH = "main"

var BIT_PATH = ".bit"

// data files
var HEAD_PATH = filepath.Join(BIT_PATH, "head")
var INDEX_PATH = filepath.Join(BIT_PATH, "index")
var IGNORE_PATH = ".bitignore"

// objects
var OBJECTS_PATH = filepath.Join(BIT_PATH, "objects")
var BLOBS_PATH = filepath.Join(OBJECTS_PATH, "blobs")
var COMMITS_PATH = filepath.Join(OBJECTS_PATH, "commits")
var TREES_PATH = filepath.Join(OBJECTS_PATH, "trees")
var BRANCHES_PATH = filepath.Join(BIT_PATH, "branches")
var DEFAULT_BRANCH_PATH = filepath.Join(BRANCHES_PATH, DEFAULT_BRANCH)
