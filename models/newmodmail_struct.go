package models

import (
	"encoding/json"
	"time"
)

// NewModmailConversation represents a conversation between a user & mods
// in the new modmail interface.
type NewModmailConversation struct {
	Conversation struct {
		IsAuto bool `json:"isAuto"`
		ObjIds []struct {
			ID  string `json:"id"`
			Key string `json:"key"`
		} `json:"objIds"`
		IsRepliable    bool       `json:"isRepliable"`
		LastUserUpdate *time.Time `json:"lastUserUpdate"`
		IsInternal     bool       `json:"isInternal"`
		LastModUpdate  *time.Time `json:"lastModUpdate"`
		LastUpdated    *time.Time `json:"lastUpdated"`
		Authors        []struct {
			IsMod         bool     `json:"isMod"`
			IsAdmin       bool     `json:"isAdmin"`
			Name          string   `json:"name"`
			IsOp          bool     `json:"isOp"`
			IsParticipant bool     `json:"isParticipant"`
			IsHidden      bool     `json:"isHidden"`
			ID            RedditID `json:"id"`
			IsDeleted     bool     `json:"isDeleted"`
		} `json:"authors"`
		Owner struct {
			DisplayName string   `json:"displayName"`
			Type        string   `json:"type"`
			ID          RedditID `json:"id"`
		} `json:"owner"`
		ID            string          `json:"id"`
		IsHighlighted bool            `json:"isHighlighted"`
		Subject       string          `json:"subject"`
		Participant   json.RawMessage `json:"participant"`
		State         int             `json:"state"`
		LastUnread    *time.Time      `json:"lastUnread"`
		NumMessages   int             `json:"numMessages"`
	} `json:"conversation"`
	Messages map[string]struct {
		Body   string `json:"body"`
		Author struct {
			IsMod         bool     `json:"isMod"`
			IsAdmin       bool     `json:"isAdmin"`
			Name          string   `json:"name"`
			IsOp          bool     `json:"isOp"`
			IsParticipant bool     `json:"isParticipant"`
			IsHidden      bool     `json:"isHidden"`
			ID            RedditID `json:"id"`
			IsDeleted     bool     `json:"isDeleted"`
		} `json:"author"`
		IsInternal   bool      `json:"isInternal"`
		Date         time.Time `json:"date"`
		BodyMarkdown string    `json:"bodyMarkdown"`
		ID           string    `json:"id"`
	} `json:"messages"`
	User struct {
		RecentComments map[RedditID]struct {
			Comment   string    `json:"comment"`
			Date      time.Time `json:"date"`
			Permalink string    `json:"permalink"`
			Title     string    `json:"title"`
		} `json:"recentComments"`
		MuteStatus struct {
			IsMuted bool       `json:"isMuted"`
			EndDate *time.Time `json:"endDate"`
			Reason  string     `json:"reason"`
		} `json:"muteStatus"`
		Name      string    `json:"name"`
		Created   time.Time `json:"created"`
		BanStatus struct {
			EndDate     *time.Time `json:"endDate"`
			Reason      string     `json:"reason"`
			IsBanned    bool       `json:"isBanned"`
			IsPermanent bool       `json:"isPermanent"`
		} `json:"banStatus"`
		IsSuspended    bool `json:"isSuspended"`
		IsShadowBanned bool `json:"isShadowBanned"`
		RecentPosts    map[RedditID]struct {
			Date      time.Time `json:"date"`
			Permalink string    `json:"permalink"`
			Title     string    `json:"title"`
		} `json:"recentPosts"`
		RecentConvos map[string]struct {
			Date      time.Time `json:"date"`
			Permalink string    `json:"permalink"`
			ID        string    `json:"id"`
			Subject   string    `json:"subject"`
		} `json:"recentConvos"`
		ID RedditID `json:"id"`
	} `json:"user"`
	ModActions map[string]struct {
		Date         time.Time `json:"date"`
		ActionTypeID int       `json:"actionTypeId"`
		ID           string    `json:"id"`
		Author       struct {
			Name      string   `json:"name"`
			IsMod     bool     `json:"isMod"`
			IsAdmin   bool     `json:"isAdmin"`
			IsHidden  bool     `json:"isHidden"`
			ID        RedditID `json:"id"`
			IsDeleted bool     `json:"isDeleted"`
		} `json:"author"`
	} `json:"modActions"`
}
