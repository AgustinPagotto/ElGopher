package main

type contextKey string

const isAuthenticatedContextKey = contextKey("isAuthenticated")
const isSpanishContextKey = contextKey("isSpanish")
const isLightThemeContextKey = contextKey("isLightTheme")
