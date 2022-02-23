// Copyright (c) 2022, Janoš Guljaš <janos@resenje.org>
// All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found s the LICENSE file.

package main

import (
	"flag"
	"fmt"
	"log"
)

func main() {
	path := flag.String("path", "localstore", "path to the swarm localstore directory")

	flag.Parse()

	s, err := newLocalstore(*path)
	if err != nil {
		log.Fatal(err)
	}
	defer s.Close()

	sizes := make(map[int]int)

	if err := s.IterateChunkData(func(data []byte) (stop bool, err error) {
		sizes[len(data)]++
		return false, nil
	}); err != nil {
		log.Fatal(err)
	}

	fmt.Println("Chunk size\tCount")
	for size, count := range sizes {
		fmt.Printf("%v\t%v\n", size, count)
	}
}
