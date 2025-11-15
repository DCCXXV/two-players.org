package realtime

import (
	"sync"
	"time"

	"github.com/TwiN/go-away"
)

const (
	MaxMessageLength     = 500
	MaxMessagesPerWindow = 5
	RateLimitWindow      = 10 * time.Second
)

type ChatModerator struct {
	rateLimiter *RateLimiter
}

func NewChatModerator() *ChatModerator {
	return &ChatModerator{
		rateLimiter: NewRateLimiter(),
	}
}

func (cm *ChatModerator) ValidateMessage(clientID, message string) (string, string) {
	if len(message) == 0 {
		return "", "Message cannot be empty"
	}
	if len(message) > MaxMessageLength {
		return "", "Message too long (max 500 characters)"
	}

	if !cm.rateLimiter.AllowMessage(clientID) {
		return "", "You're sending messages too quickly. Please slow down."
	}

	if goaway.IsProfane(message) {
		filtered := goaway.Censor(message)
		return filtered, ""
	}

	return message, ""
}

type RateLimiter struct {
	mu             sync.RWMutex
	clientMessages map[string][]time.Time
}

func NewRateLimiter() *RateLimiter {
	rl := &RateLimiter{
		clientMessages: make(map[string][]time.Time),
	}

	go rl.cleanup()

	return rl
}

func (rl *RateLimiter) AllowMessage(clientID string) bool {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	now := time.Now()
	cutoff := now.Add(-RateLimitWindow)

	timestamps, exists := rl.clientMessages[clientID]
	if !exists {
		timestamps = []time.Time{}
	}

	validTimestamps := []time.Time{}
	for _, ts := range timestamps {
		if ts.After(cutoff) {
			validTimestamps = append(validTimestamps, ts)
		}
	}

	if len(validTimestamps) >= MaxMessagesPerWindow {
		return false
	}

	validTimestamps = append(validTimestamps, now)
	rl.clientMessages[clientID] = validTimestamps

	return true
}

func (rl *RateLimiter) cleanup() {
	ticker := time.NewTicker(1 * time.Minute)
	defer ticker.Stop()

	for range ticker.C {
		rl.mu.Lock()
		now := time.Now()
		cutoff := now.Add(-RateLimitWindow * 2)

		for clientID, timestamps := range rl.clientMessages {
			validTimestamps := []time.Time{}
			for _, ts := range timestamps {
				if ts.After(cutoff) {
					validTimestamps = append(validTimestamps, ts)
				}
			}

			if len(validTimestamps) == 0 {
				delete(rl.clientMessages, clientID)
			} else {
				rl.clientMessages[clientID] = validTimestamps
			}
		}
		rl.mu.Unlock()
	}
}
