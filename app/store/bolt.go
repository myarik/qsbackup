package store

import "github.com/coreos/bbolt"


// BoltDB implements store.
// there are 3 types of top-level buckets:
//  - dir_hashe a dir hash value
//  - history of all backups.
//  - user to comment references in "users" bucket. It used to get comments for user. Key is userID and value
//    is a nested bucket named userID with kv as ts:reference
//  - blocking info sits in "block" bucket. Key is userID, value - ts
//  - counts per post to keep number of comments. Key is post url, value - count

type BoltDB struct {
	db *bolt.DB
}