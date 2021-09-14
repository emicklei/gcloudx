package im

import (
	"encoding/json"
	"log"
	"os"
)

const cacheFile = "latest-gcp-bindings.json"

func loadMemberWithRoleInProjects() (list []MemberWithRoleInProject) {
	if *verbose {
		log.Println("loading member,role memberships")
	}
	in, err := os.Open(cacheFile)
	// does not exist, create a fresh
	if err != nil {
		list = fetchMemberWithRoleInProjects()
		writeCache(list)
	} else {
		list = readCache(in)
	}
	return
}

func readCache(in *os.File) (list []MemberWithRoleInProject) {
	if *verbose {
		log.Println("reading cache file")
	}
	dec := json.NewDecoder(in)
	err := dec.Decode(&list)
	if err != nil {
		log.Fatalf("failed to read input file: %v", err)
	}
	return
}

func writeCache(list []MemberWithRoleInProject) {
	if *verbose {
		log.Println("writing cache file")
	}
	out, err := os.Create(cacheFile)
	if err != nil {
		log.Fatalf("failed to create output file: %v", err)
	}
	enc := json.NewEncoder(out)
	enc.SetIndent("", "\t")
	err = enc.Encode(list)
	if err != nil {
		log.Fatalf("failed to write output file: %v", err)
	}
}
