package ffmpeg

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/novix/services/processing-service/internal/model"
	"github.com/novix/services/processing-service/pkg/logger"
	"go.uber.org/zap"
)

type FFmpeg struct {
	path string
	threads int
}

func New(path string, threads int) *FFmpeg {
	return &FFmpeg{
		path: path,
		threads: threads,
	}
}

// Transcode converts a video to a specific quality variant
func (f *FFmpeg) Transcode(ctx context.Context, opts TranscodeOptions) error {
	profile := opts.Profile
	args := []string{
		"-i", opts.InputPath,
		"-vf", fmt.Sprintf("scale=%d:%d", profile.Width, profile.Height),
		"-c:v", "libx264",
		"-b:v", profile.Bitrate,
		"-c:a", "aac",
		"-b:a", profile.Audio,
		"-threads", fmt.Sprintf("%d", f.threads),
		"-preset", "fast",
		"-movflags", "+faststart",
		"-y",
		opts.OutputPath,
	}
	logger.Info("Starting transcode",
		zap.String("quality", string(opts.Quality)),
		zap.String("input", opts.InputPath),
		zap.String("output", opts.OutputPath),
	)
	return f.run(ctx, args)
}

// GenerateHLS creates adaptive bitrate HLS output with master playlist
func (f *FFmpeg) GenerateHLS(ctx context.Context, opts HLSOptions) (string, error) {
	// Create output directory
	if err := os.MkdirAll(opts.OutputDir, 0755); err != nil {
		return "", fmt.Errorf("failed to create output dir: %w", err)
	}

	// Build multi-quality HLS command
	args := []string{"-i", opts.InputPath}

	var streamMaps []string
	variantPlaylists := []string{}
	for i, quality := range opts.Qualities {
		profile := model.QualityProfiles[quality]
		qualityDir := filepath.Join(opts.OutputDir, string(quality))
		if err := os.MkdirAll(qualityDir, 0755); err != nil {
			return "", fmt.Errorf("failed to create quality dir: %w", err)
		}
		segmentPath := filepath.Join(qualityDir, "segment%03d.ts")
		playlistPath := filepath.Join(qualityDir, "playlist.m3u8")

		// Video stream for this quality
		args = append(args,
			fmt.Sprintf("-map"), "0:v:0",
			fmt.Sprintf("-map"), "0:a:0",
			fmt.Sprintf("-c:v:%d", i), "libx264",
			fmt.Sprintf("-b:v:%d", i), profile.Bitrate,
			fmt.Sprintf("-s:v:%d", i),
			fmt.Sprintf("%dx%d", profile.Width, profile.Height),
			fmt.Sprintf("-c:a:%d", i), "aac",
			fmt.Sprintf("-b:a:%d", i), profile.Audio,
		)
		streamMaps = append(streamMaps, fmt.Sprintf("v:%d,a:%d", i, i))
		variantPlaylists = append(variantPlaylists,
			fmt.Sprintf("#EXT-X-STREAM-INF:BANDWIDTH=%s,RESOLUTION=%dx%d\n%s/playlist.m3u8",
				profile.Bitrate, profile.Width, profile.Height, string(quality)),
		)
		args = append(args,
			"-hls_time", fmt.Sprintf("%d", opts.SegmentLength),
			"-hls_playlist_type", "vod",
			"-hls_segment_filename", segmentPath,
			playlistPath,
		)
		_ = streamMaps
	}
	logger.Info("Generating HLS",
		zap.String("input", opts.InputPath),
		zap.String("outputDir", opts.OutputDir),
		zap.Int("qualities", len(opts.Qualities)),
	)
	if err := f.run(ctx, args); err != nil {
		return "", err
	}

	// Write master playlist manually
	masterPath := filepath.Join(opts.OutputDir, opts.PlaylistName)
	masterContent := "#EXTM3U\n#EXT-X-VERSION:3\n\n" + strings.Join(variantPlaylists, "\n\n")
	if err := os.WriteFile(masterPath, []byte(masterContent), 0644); err != nil {
		return "", fmt.Errorf("failed to write master playlist: %w", err)
	}
	logger.Info("HLS generation complete",
		zap.String("masterPlaylist", masterPath),
	)
	return masterPath, nil
}

// GenerateThumbnail extracts a single frame as a JPEG thumbnail
func (f *FFmpeg) GenerateThumbnail(ctx context.Context,
	opts ThumbnailOptions) error {
	args := []string{
		"-i", opts.InputPath,
		"-ss", opts.TimeOffset,
		"-vframes", "1",
		"-vf", fmt.Sprintf("scale=%d:%d", opts.Width, opts.Height),
		"-q:v", "2",
		"-y",
		opts.OutputPath,
	}
	logger.Info("Generating thumbnail",
		zap.String("input", opts.InputPath),
		zap.String("output", opts.OutputPath),
		zap.String("offset", opts.TimeOffset),
	)
	return f.run(ctx, args)
}

// GetVideoDuration returns video duration in seconds
func (f *FFmpeg) GetVideoDuration(ctx context.Context,
	inputPath string) (string, error) {
	args := []string{
		"-v", "error",
		"-show_entries", "format=duration",
		"-of", "default=noprint_wrappers=1:nokey=1",
		inputPath,
	}
	cmd := exec.CommandContext(ctx, "ffprobe", args...)
	out, err := cmd.Output()
	if err != nil {
		return "", fmt.Errorf("ffprobe failed: %w", err)
	}
	return strings.TrimSpace(string(out)), nil
}

// run executes an FFmpeg command
func (f *FFmpeg) run(ctx context.Context, args []string) error {
	cmd := exec.CommandContext(ctx, f.path, args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("ffmpeg command failed: %w", err)
	}
	return nil
}
