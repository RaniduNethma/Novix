package model

// VideoManifest holds HLS manifest info for a video
type VideoManifest struct {
	VideoID       string        `json:"videoId"`
	MasterPath    string        `json:"masterPath"`
	Qualities     []QualityInfo `json:"qualities"`
	Duration      float64       `json:"duration"`
	ThumbnailPath string        `json:"thumbnailPath"`
}

type QualityInfo struct {
	Quality      string `json:"quality"`
	Width        int    `json:"width"`
	Height       int    `json:"height"`
	Bitrate      string `json:"bitrate"`
	PlaylistPath string `json:"playlistPath"`
}

// StreamRequest is what the client sends to start streaming
type StreamRequest struct {
	VideoID string `json:"videoId" binding:"required"`
	Quality string `json:"quality"`
}

// StreamResponse is what we return to the client
type StreamResponse struct {
	SessionID    string `json:"sessionId"`
	ManifestURL  string `json:"manifestUrl"`
	ThumbnailURL string `json:"thumbnailUrl"`
	ExpiresIn    int    `json:"expiresIn"`
}
