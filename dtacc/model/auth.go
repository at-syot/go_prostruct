package model

import (
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/uptrace/bun"
)

// AuthSession represents an authentication session
type AuthSession struct {
	bun.BaseModel `bun:"table:auth_sessions,alias:as"`

	ID               uuid.UUID  `bun:"id,pk,type:uuid,default:uuid_generate_v4()" json:"id"`
	UserID           uuid.UUID  `bun:"user_id,notnull" json:"user_id"`
	RefreshTokenHash string     `bun:"refresh_token_hash,notnull" json:"-"` // Never expose in JSON
	UserAgent        *string    `bun:"user_agent" json:"user_agent,omitempty"`
	IPAddress        *string    `bun:"ip_address" json:"ip_address,omitempty"`
	DeviceName       *string    `bun:"device_name" json:"device_name,omitempty"`
	IsValid          bool       `bun:"is_valid,default:true" json:"is_valid"`
	CreatedAt        time.Time  `bun:"created_at,nullzero,notnull,default:now()" json:"created_at"`
	ExpiresAt        time.Time  `bun:"expires_at,notnull" json:"expires_at"`
	LastUsedAt       *time.Time `bun:"last_used_at" json:"last_used_at,omitempty"`

	// Relations
	User *User `bun:"rel:belongs-to,join:user_id=id" json:"user,omitempty"`
}

// AuthLog represents an authentication event log
type AuthLog struct {
	bun.BaseModel `bun:"table:auth_logs,alias:al"`

	ID        uuid.UUID  `bun:"id,pk,type:uuid,default:uuid_generate_v4()" json:"id"`
	UserID    uuid.UUID  `bun:"user_id,notnull" json:"user_id"`
	SessionID *uuid.UUID `bun:"session_id" json:"session_id,omitempty"`
	EventType string     `bun:"event_type,notnull" json:"event_type"` // Should be one of the EventType constants
	IPAddress *string    `bun:"ip_address" json:"ip_address,omitempty"`
	UserAgent *string    `bun:"user_agent" json:"user_agent,omitempty"`
	EventTime time.Time  `bun:"event_time,nullzero,notnull,default:now()" json:"event_time"`

	// Relations
	User    *User        `bun:"rel:belongs-to,join:user_id=id" json:"user,omitempty"`
	Session *AuthSession `bun:"rel:belongs-to,join:session_id=id" json:"session,omitempty"`
}

// AuthEventType constants for authentication events
const (
	EventTypeLogin   = "login"
	EventTypeLogout  = "logout"
	EventTypeRefresh = "refresh"
	EventTypeRevoked = "revoked"
)

// Session validation and utility methods

// IsExpired checks if the auth session has expired
func (as *AuthSession) IsExpired() bool {
	return time.Now().After(as.ExpiresAt)
}

// IsActive checks if the session is valid and not expired
func (as *AuthSession) IsActive() bool {
	return as.IsValid && !as.IsExpired()
}

// Invalidate marks the session as invalid (soft delete)
func (as *AuthSession) Invalidate() {
	as.IsValid = false
}

// UpdateLastUsed updates the last used timestamp to current time
func (as *AuthSession) UpdateLastUsed() {
	now := time.Now()
	as.LastUsedAt = &now
}

// TimeUntilExpiry returns the duration until the session expires
func (as *AuthSession) TimeUntilExpiry() time.Duration {
	return time.Until(as.ExpiresAt)
}

// AuthLog utility methods

// IsValidEventType checks if the event type is one of the allowed constants
func (al *AuthLog) IsValidEventType() bool {
	switch al.EventType {
	case EventTypeLogin, EventTypeLogout, EventTypeRefresh, EventTypeRevoked:
		return true
	default:
		return false
	}
}

// Validate performs basic validation on the AuthLog
func (al *AuthLog) Validate() error {
	if !al.IsValidEventType() {
		return fmt.Errorf("invalid event type: %s", al.EventType)
	}
	if al.UserID == uuid.Nil {
		return fmt.Errorf("user_id cannot be nil")
	}
	return nil
}

// Helper functions for creating auth logs

// NewAuthLog creates a new auth log entry
func NewAuthLog(userID uuid.UUID, eventType string, sessionID *uuid.UUID, ipAddress, userAgent *string) *AuthLog {
	return &AuthLog{
		UserID:    userID,
		SessionID: sessionID,
		EventType: eventType,
		IPAddress: ipAddress,
		UserAgent: userAgent,
		EventTime: time.Now(),
	}
}
