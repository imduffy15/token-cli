package help

func Context() string {
	return `A context represents a previously fetched access token and associated metadata
such as the scopes that token contains. The token CLI caches these results on a
local file so that they may be used when issuing requests that require an
Authorization header.
`
}
