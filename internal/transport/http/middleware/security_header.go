package middleware

import "github.com/gin-gonic/gin"

func SecurityHeader() gin.HandlerFunc {
	return func(c *gin.Context) {
		// HSTS - Force HTTPS
		c.Header("Strict-Transport-Security", "max-age=63072000; includeSubDomains; preload")
		
		// Prevent MIME type sniffing
		c.Header("X-Content-Type-Options", "nosniff")
		
		// Clickjacking protection
		c.Header("X-Frame-Options", "DENY")
		
		// XSS protection
		c.Header("X-XSS-Protection", "1; mode=block")
		
		// Referrer policy
		c.Header("Referrer-Policy", "strict-origin-when-cross-origin")
		
		// Content Security Policy
		c.Header("Content-Security-Policy", "default-src 'self'")
		
		c.Next()
	}
}