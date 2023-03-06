package bots

import (
	"context"
	"time"
)

type Bot interface {
	Run(ctx context.Context) error
	Close(ctx context.Context) error
}

type T struct {
	Id            string `json:"id"`
	ApplicationId string `json:"application_id"`
	Type          int    `json:"type"`
	Data          struct {
		CustomId      string `json:"custom_id"`
		ComponentType int    `json:"component_type"`
		Resolved      struct {
			Users    interface{} `json:"users"`
			Members  interface{} `json:"members"`
			Roles    interface{} `json:"roles"`
			Channels interface{} `json:"channels"`
		} `json:"resolved"`
		Values []string `json:"values"`
	} `json:"data"`
	GuildId   string `json:"guild_id"`
	ChannelId string `json:"channel_id"`
	Message   struct {
		Id              string        `json:"id"`
		ChannelId       string        `json:"channel_id"`
		Content         string        `json:"content"`
		Timestamp       time.Time     `json:"timestamp"`
		EditedTimestamp interface{}   `json:"edited_timestamp"`
		MentionRoles    []interface{} `json:"mention_roles"`
		Tts             bool          `json:"tts"`
		MentionEveryone bool          `json:"mention_everyone"`
		Author          struct {
			Id            string `json:"id"`
			Email         string `json:"email"`
			Username      string `json:"username"`
			Avatar        string `json:"avatar"`
			Locale        string `json:"locale"`
			Discriminator string `json:"discriminator"`
			Token         string `json:"token"`
			Verified      bool   `json:"verified"`
			MfaEnabled    bool   `json:"mfa_enabled"`
			Banner        string `json:"banner"`
			AccentColor   int    `json:"accent_color"`
			Bot           bool   `json:"bot"`
			PublicFlags   int    `json:"public_flags"`
			PremiumType   int    `json:"premium_type"`
			System        bool   `json:"system"`
			Flags         int    `json:"flags"`
		} `json:"author"`
		Attachments       []interface{} `json:"attachments"`
		Embeds            []interface{} `json:"embeds"`
		Mentions          []interface{} `json:"mentions"`
		Reactions         interface{}   `json:"reactions"`
		Pinned            bool          `json:"pinned"`
		Type              int           `json:"type"`
		WebhookId         string        `json:"webhook_id"`
		Member            interface{}   `json:"member"`
		MentionChannels   interface{}   `json:"mention_channels"`
		Activity          interface{}   `json:"activity"`
		Application       interface{}   `json:"application"`
		MessageReference  interface{}   `json:"message_reference"`
		ReferencedMessage interface{}   `json:"referenced_message"`
		Interaction       struct {
			Id   string `json:"id"`
			Type int    `json:"type"`
			Name string `json:"name"`
			User struct {
				Id            string `json:"id"`
				Email         string `json:"email"`
				Username      string `json:"username"`
				Avatar        string `json:"avatar"`
				Locale        string `json:"locale"`
				Discriminator string `json:"discriminator"`
				Token         string `json:"token"`
				Verified      bool   `json:"verified"`
				MfaEnabled    bool   `json:"mfa_enabled"`
				Banner        string `json:"banner"`
				AccentColor   int    `json:"accent_color"`
				Bot           bool   `json:"bot"`
				PublicFlags   int    `json:"public_flags"`
				PremiumType   int    `json:"premium_type"`
				System        bool   `json:"system"`
				Flags         int    `json:"flags"`
			} `json:"user"`
			Member interface{} `json:"member"`
		} `json:"interaction"`
		Flags        int         `json:"flags"`
		StickerItems interface{} `json:"sticker_items"`
	} `json:"message"`
	AppPermissions string      `json:"app_permissions"`
	Member         interface{} `json:"member"`
	User           struct {
		Id            string `json:"id"`
		Email         string `json:"email"`
		Username      string `json:"username"`
		Avatar        string `json:"avatar"`
		Locale        string `json:"locale"`
		Discriminator string `json:"discriminator"`
		Token         string `json:"token"`
		Verified      bool   `json:"verified"`
		MfaEnabled    bool   `json:"mfa_enabled"`
		Banner        string `json:"banner"`
		AccentColor   int    `json:"accent_color"`
		Bot           bool   `json:"bot"`
		PublicFlags   int    `json:"public_flags"`
		PremiumType   int    `json:"premium_type"`
		System        bool   `json:"system"`
		Flags         int    `json:"flags"`
	} `json:"user"`
	Locale      string      `json:"locale"`
	GuildLocale interface{} `json:"guild_locale"`
	Token       string      `json:"token"`
	Version     int         `json:"version"`
}
