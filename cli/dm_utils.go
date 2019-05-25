package main

import "hash/fnv"

func hash64(s string) int64 {
	hash := fnv.New32()
	hash.Write([]byte(s))
	return int64(hash.Sum32())
}
