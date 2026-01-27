// Copyright (c) 2026 WabiSaby
// All rights reserved.
//
// This source code is proprietary and confidential. Unauthorized copying,
// modification, distribution, or use of this software, via any medium is
// strictly prohibited without the express written permission of WabiSaby.
//
// This software contains confidential and proprietary information of
// WabiSaby and its licensors. Use, disclosure, or reproduction
// is prohibited without the prior express written permission of WabiSaby.

package sdk

// ContentDownloaderPlugin is a base plugin for content download plugins.
// Embed this in your plugin to get default implementations.
type ContentDownloaderPlugin struct {
	*BasePlugin
}

// NewContentDownloaderPlugin creates a new ContentDownloaderPlugin.
func NewContentDownloaderPlugin() *ContentDownloaderPlugin {
	return &ContentDownloaderPlugin{
		BasePlugin: NewBasePlugin(),
	}
}

// DownloadRequest represents a request to download content.
type DownloadRequest struct {
	URL         string // URL to download from
	Format      string // Desired format (e.g., "mp3", "mp4")
	Quality     string // Desired quality (e.g., "high", "medium", "low")
	MaxDuration *int   // Optional maximum duration in seconds
}

// DownloadResult represents the result of a download operation.
type DownloadResult struct {
	FilePath string        // Path to the downloaded file
	Metadata *SongMetadata // Metadata about the downloaded content
	Duration int           // Duration in seconds
	FileSize int64         // File size in bytes
}

// SongMetadata represents metadata about a song.
type SongMetadata struct {
	Title        string  // Song title
	Artist       *string // Artist name (optional)
	Album        *string // Album name (optional)
	Duration     *int    // Duration in seconds (optional)
	ThumbnailURL *string // Thumbnail image URL (optional)
}

// ContentDownloader is the interface that content download plugins must implement.
type ContentDownloader interface {
	Plugin

	// Download downloads content from the given URL.
	// Returns the download result or an error.
	Download(ctx *Context, req *DownloadRequest) (*DownloadResult, error)

	// CanHandle checks if this plugin can handle the given URL.
	CanHandle(url string) bool

	// SupportedDomains returns a list of domains this plugin can handle.
	SupportedDomains() []string
}

// MetadataResolverPlugin is a base plugin for metadata resolution plugins.
// Embed this in your plugin to get default implementations.
type MetadataResolverPlugin struct {
	*BasePlugin
}

// NewMetadataResolverPlugin creates a new MetadataResolverPlugin.
func NewMetadataResolverPlugin() *MetadataResolverPlugin {
	return &MetadataResolverPlugin{
		BasePlugin: NewBasePlugin(),
	}
}

// ResolveURLRequest represents a request to resolve metadata from a URL.
type ResolveURLRequest struct {
	URL string // URL to resolve
}

// ResolveResult represents the result of a metadata resolution operation.
type ResolveResult struct {
	Metadata    *SongMetadata // Resolved metadata
	DownloadURL *string       // Optional direct download URL
	StreamURL   *string       // Optional streaming URL
}

// SearchRequest represents a request to search for content.
type SearchRequest struct {
	Query string // Search query
	Limit *int   // Optional limit on number of results
}

// SearchResult represents a single search result.
type SearchResult struct {
	Metadata    *SongMetadata // Metadata about the content
	URL         string        // URL to the content
	DownloadURL *string       // Optional direct download URL
	StreamURL   *string       // Optional streaming URL
}

// MetadataResolver is the interface that metadata resolver plugins must implement.
type MetadataResolver interface {
	Plugin

	// ResolveURL resolves metadata for the given URL.
	// Returns the resolution result or an error.
	ResolveURL(ctx *Context, req *ResolveURLRequest) (*ResolveResult, error)

	// Search searches for content matching the given query.
	// Returns a list of search results or an error.
	Search(ctx *Context, req *SearchRequest) ([]*SearchResult, error)

	// CanHandle checks if this plugin can handle the given URL.
	CanHandle(url string) bool

	// SupportedDomains returns a list of domains this plugin can handle.
	SupportedDomains() []string
}

// StorageProviderPlugin is a base plugin for storage provider plugins.
type StorageProviderPlugin struct {
	*BasePlugin
}

// NewStorageProviderPlugin creates a new StorageProviderPlugin.
func NewStorageProviderPlugin() *StorageProviderPlugin {
	return &StorageProviderPlugin{
		BasePlugin: NewBasePlugin(),
	}
}

// UploadHLSRequest represents a request to upload HLS files.
type UploadHLSRequest struct {
	PlaylistPath string
	SegmentsDir  string
	BaseFilename string
}

// StorageProvider is the interface that storage provider plugins must implement.
type StorageProvider interface {
	Plugin

	// UploadHLSFiles uploads HLS files and returns the CDN URL.
	UploadHLSFiles(ctx *Context, req *UploadHLSRequest) (string, error)

	// GetFileSizeMB returns the total size of an audio file in MB.
	GetFileSizeMB(ctx *Context, cdnURL string) (float64, error)

	// DeleteAudio deletes an audio file from storage.
	DeleteAudio(ctx *Context, cdnURL string) error
}
