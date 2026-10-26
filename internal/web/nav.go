package web

// Nav carries the per-request bits the layout header needs:
// whether the user is signed in, and their email.
type Nav struct {
	LoggedIn bool
	Email    string
}
