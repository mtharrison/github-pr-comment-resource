package resource

import (
	"context"
	"regexp"
	"sort"
	"strconv"

	"github.com/google/go-github/github"
	"golang.org/x/oauth2"
)

// GithubVersions retrieves all matching versions for the given input
func GithubVersions(input Input) ([]Version, error) {
	result := make([]Version, 0)

	client, err := githubClient(input.Source.AccessToken, input.Source.V3Endpoint)
	if err != nil {
		return result, err
	}

	// fetch prs
	prs, _, err := client.PullRequests.List(
		context.Background(),
		input.Source.Owner(),
		input.Source.Repo(),
		&github.PullRequestListOptions{
			State:     "open",
			Sort:      "created",
			Direction: "asc",
		},
	)
	if err != nil {
		return result, err
	}

	for _, pr := range prs {

		// fetch comments
		comments, _, err := client.Issues.ListComments(
			context.Background(),
			input.Source.Owner(),
			input.Source.Repo(),
			*pr.Number,
			&github.IssueListCommentsOptions{
				Sort:      "created",
				Direction: "asc",
			},
		)
		if err != nil {
			return result, err
		}

		// filter comments
		for _, comment := range comments {

			include, err := filterComment(input, comment)
			if err != nil {
				return result, err
			}

			if include {
				result = append(result, Version{
					PrNumber:  strconv.Itoa(*pr.Number),
					CommentID: strconv.FormatInt(*comment.ID, 10),
				})
			}

		}
	}

	// Sort versions

	sort.Slice(result, func(i, j int) bool {
		// err should not be possible because we just formatted it from i64
		l, _ := strconv.ParseInt(result[i].CommentID, 10, 64)
		r, _ := strconv.ParseInt(result[j].CommentID, 10, 64)
		return l < r
	})

	return result, nil
}

// GithubComment retrieves a specific comment
func GithubComment(input Input) (Comment, error) {
	var result Comment

	client, err := githubClient(input.Source.AccessToken, input.Source.V3Endpoint)
	if err != nil {
		return result, err
	}

	prNumber, _ := strconv.Atoi(input.Version.PrNumber)
	pr, _, err := client.PullRequests.Get(context.Background(), input.Source.Owner(), input.Source.Repo(), prNumber)
	if err != nil {
		return result, err
	}

	commentID, _ := strconv.ParseInt(input.Version.CommentID, 10, 64)
	comment, _, err := client.Issues.GetComment(context.Background(), input.Source.Owner(), input.Source.Repo(), commentID)
	if err != nil {
		return result, err
	}

	matches := make([][]string, 0)

	if input.Source.Regex != "" {
		re := regexp.MustCompile(input.Source.Regex)
		matches = re.FindAllStringSubmatch(*comment.Body, -1)
	}

	result = Comment{
		ID:        input.Version.CommentID,
		Owner:     input.Source.Owner(),
		Repo:      input.Source.Repo(),
		PR:        strconv.Itoa(*pr.Number),
		Branch:    *pr.Head.Ref,
		Body:      *comment.Body,
		CreatedAt: comment.CreatedAt.String(),
		Matches:   matches,
		User:      *comment.User.Login,
	}

	return result, nil
}

func githubClient(token, endpoint string) (*github.Client, error) {
	ctx := context.Background()
	ts := oauth2.StaticTokenSource(&oauth2.Token{AccessToken: token})
	tc := oauth2.NewClient(ctx, ts)

	if endpoint == "" {
		return github.NewClient(tc), nil
	}

	client, err := github.NewEnterpriseClient(endpoint, endpoint, tc)
	if err != nil {
		return nil, err
	}
	return client, nil
}

func filterComment(input Input, comment *github.IssueComment) (bool, error) {
	// filter older comments
	if input.Version != nil && input.Version.CommentID != "" {
		commentVersionInt64, err := strconv.ParseInt(input.Version.CommentID, 10, 64)
		if err != nil {
			return false, err
		}

		if *comment.ID < commentVersionInt64 {
			return false, nil
		}
	}

	// filter non-matching comments
	if input.Source.Regex != "" {
		match, err := regexp.MatchString(input.Source.Regex, *comment.Body)
		if err != nil {
			return false, err
		}
		if !match {
			return false, nil
		}
	}

	return true, nil
}
