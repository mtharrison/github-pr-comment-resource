package resource

// Output is written by in script
type Output struct {
	Version  Version    `json:"version"`
	Metadata []Metadata `json:"metadata"`
}

// Metadata is a key-value type
type Metadata struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

// Comment is the actual resource
type Comment struct {
	ID        string     `json:"id"`
	Owner     string     `json:"owner"`
	Repo      string     `json:"repo"`
	PR        string     `json:"pr_number"`
	Branch    string     `json:"branch"`
	Body      string     `json:"body"`
	CreatedAt string     `json:"created_at"`
	Matches   [][]string `json:"matches"`
	User      string     `json:"user"`
}

// Metadata returns the comment metadata
func (c *Comment) Metadata() []Metadata {

	result := []Metadata{
		Metadata{
			Name:  "Owner",
			Value: c.Owner,
		}, Metadata{
			Name:  "Repo",
			Value: c.Repo,
		}, Metadata{
			Name:  "PR Number",
			Value: c.PR,
		}, Metadata{
			Name:  "Branch",
			Value: c.Branch,
		}, Metadata{
			Name:  "Created At",
			Value: c.CreatedAt,
		}, Metadata{
			Name:  "User",
			Value: c.User,
		},
	}

	return result
}
