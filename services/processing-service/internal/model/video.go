package model

type VideoQuality string

const(
	Quality360p VideoQuality = "360p"
	Quality480p VideoQuality = "480p"
	Quality720p VideoQuality = "720p"
	Quality1080p VideoQuality = "1080p"
)

type VideoVariant struct{
	ID uint `gorm:"primaryKey;autoIncrement"`
	VideoID string `gorm:"not null;index"`
	Quality VideoQuality `gorm:"not null"`
	Width int
	Height int
	Bitrate int
	Path string `gorm:"not null"`
	M3U8Path string
	Size int64
	Duration float64
}

type QualityProfile struct {
	Width   int
	Height  int
	Bitrate string
	Audio   string
}

// QualityProfiles defines transcoding targets
var QualityProfiles = map[VideoQuality]QualityProfile{
	Quality360p: {
		Width: 640,
		Height: 360,
		Bitrate: "800k",
		Audio: "96k",
	},
	Quality480p: {
		Width: 854,
		Height: 480,
		Bitrate: "1400k",
		Audio: "128k",
	},
	Quality720p: {
		Width: 1280,
		Height: 720,
		Bitrate: "2800k",
		Audio: "128k",
	},
	Quality1080p: {
		Width: 1920,
		Height: 1080,
		Bitrate: "5000k",
		Audio: "192k",
	},
}
