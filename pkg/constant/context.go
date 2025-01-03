package constant

type ContextKey string

const (
	ContextRequestIDKey          ContextKey = "requestID"
	ContextRequestEndpointKey    ContextKey = "requestEndpoint"
	ContextAccountIDKey          ContextKey = "accountID"
	ContextAccountTypeKey        ContextKey = "accountType"
	ContextAccountRoleKey        ContextKey = "accountRole"
	ContextRequestScope          ContextKey = "requestScope"
	ContextChainNextRequestIDKey ContextKey = "requestChainNextRequestID" // chaining from the previous request id
)
