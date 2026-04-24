package ffmpeg

import "github.com/streamingplatform/processing-service/internal/model"

// TranscodeOptions holds all options for a transcode operation
type TranscodeOptions struct {
	InputPath string
	OutputPath string
	Quality model.VideoQuality
	Profile model.QualityProfile
	Threads int
}

// HLSOptions holds all options for HLS generation
type HLSOptions struct {
	InputPath string
	OutputDir string
	SegmentLength int
	PlaylistName string
	Qualities []model.VideoQuality
	Threads int
}

// ThumbnailOptions holds all options for thumbnail generation
type ThumbnailOptions struct {
	InputPath string
	OutputPath string
	TimeOffset string
	Width int
	Height int
}

// DefaultHLSOptions returns sensible defaults for HLS generation
func DefaultHLSOptions(inputPath string, outputDir string) HLSOptions {
	return HLSOptions{
		InputPath: inputPath,
		OutputDir: outputDir,
		SegmentLength: 6,
		PlaylistName: "master.m3u8",
		Qualities: []model.VideoQuality{
			model.Quality360p,
			model.Quality480p,
			model.Quality720p,
			model.Quality1080p,
		},
		Threads: 4,
	}
}
