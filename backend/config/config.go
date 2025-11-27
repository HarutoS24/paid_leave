package config

var AuthSessionName = "auth_session"
var AuthSessionDuration = 1800

type ContextKey string

var AuthContextKey ContextKey = "auth_session_key"
var DBContextKey ContextKey = "db_key"
var LoginUserContextKey ContextKey = "login_user_key"
