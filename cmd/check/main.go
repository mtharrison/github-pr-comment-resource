package main

import (
	"encoding/json"
	"log"
	"os"

	"github.com/mtharrison/github-pr-comment-resource/internal/resource"
)

func main() {
	input, err := resource.GetInput(os.Stdin, false)
	if err != nil {
		log.Fatal("Could not read input: ", err)
	}

	versions, err := resource.GithubVersions(input)
	if err != nil {
		log.Fatal("Could not fetch Github versions: ", err)
	}

	// If no version provided, just return the very latest version
	if len(versions) != 0 && (input.Version == nil || input.Version.CommentID == "") {
		versions = versions[len(versions)-1:]
	}

	out, err := json.Marshal(versions)
	if err != nil {
		log.Fatal("JSON marshal error: ", err)
	}

	_, err = os.Stdout.Write(out)
	if err != nil {
		log.Fatal("STDOUT write error: ", err)
	}
}
