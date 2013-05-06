package adn

import (
	"time"
)

type Description struct {
	Text     string    `json:"text,omitempty"`
	HTML     string    `json:"html,omitempty"`
	Entities *Entities `json:"entities"`
}

type Image struct {
	Height    int    `json:"height,omitempty"`
	Width     int    `json:"width,omitempty"`
	URL       string `json:"url,omitempty"`
	IsDefault bool   `json:"is_default,omitempty"`
}

type Counts struct {
	Following int `json:"following,omitempty"`
	Followers int `json:"followers,omitempty"`
	Posts     int `json:"posts,omitempty"`
	Stars     int `json:"stars,omitempty"`
}

type User struct {
	ID              string        `json:"id,omitempty"`
	Username        string        `json:"username,omitempty"`
	Name            string        `json:"name,omitempty"`
	Description     Description   `json:"description,omitempty"`
	Timezone        string        `json:"timezone,omitempty"`
	Locale          string        `json:"locale,omitempty"`
	AvatarImage     *Image        `json:"avatar_image,omitempty"`
	CoverImage      *Image        `json:"cover_image,omitempty"`
	Type            string        `json:"type,omitempty"`
	CreatedAt       time.Time     `json:"created_at,omitempty"`
	Counts          *Counts       `json:"counts,omitempty"`
	FollowsYou      bool          `json:"follows_you,omitempty"`
	YouBlocked      bool          `json:"you_blocked,omitempty"`
	YouFollow       bool          `json:"you_follow,omitempty"`
	YouMuted        bool          `json:"you_muted,omitempty"`
	YouCanSubscribe bool          `json:"you_can_subscribe,omitempty"`
	VerifiedDomain  string        `json:"verified_domain,omitempty"`
	Annotations     []*Annotation `json:"annotations,omitempty"`
}
