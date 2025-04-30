package main

type contextKey string

type pathKey string

const isAuthenticatedKey = contextKey("isAuthenticated")
