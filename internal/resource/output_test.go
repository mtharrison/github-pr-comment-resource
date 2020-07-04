package resource

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMetadata(t *testing.T) {
	comment := Comment{
		ID:        "12345",
		Owner:     "golang",
		Repo:      "go",
		PR:        "99",
		Branch:    "master",
		Body:      "This is bigger than the both of us",
		CreatedAt: "Yesterday",
		Matches:   [][]string{},
		User:      "Jimmy",
	}

	metadata := comment.Metadata()

	expected := []Metadata{
		Metadata{
			Name:  "Owner",
			Value: "golang",
		}, Metadata{
			Name:  "Repo",
			Value: "go",
		}, Metadata{
			Name:  "PR Number",
			Value: "99",
		}, Metadata{
			Name:  "Branch",
			Value: "master",
		}, Metadata{
			Name:  "Created At",
			Value: "Yesterday",
		}, Metadata{
			Name:  "User",
			Value: "Jimmy",
		},
	}

	assert.Equal(t, expected, metadata)
}
