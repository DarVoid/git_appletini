package types

type Config struct {
	Contexts       ContextMap `json:"contexts"`
	DefaultContext string     `json:"default_context"`
	ItemCount      int        `json:"item_count"`
}
type ContextMap map[string]Context

type Context struct {
	Title         string       `json:"title"`
	ChromeProfile string       `json:"chrome_profile"`
	Github        GithubConfig `json:"github"`
	Poll          PollConfig   `json:"poll"`
}
type GithubConfig struct {
	Host     string `json:"host"`
	GraphQL  string `json:"gqlAPI"`
	Username string `json:"username"`
	Token    string `json:"token"`
	Name     string `json:"name"`
	Email    string `json:"email"`
}
type PollConfig struct {
	Enabled   bool `json:"enabled"`
	Frequency int  `json:"frequency_s"`
}
