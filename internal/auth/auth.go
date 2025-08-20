package auth

import (
	spotifyauth "github.com/zmb3/spotify/v2/auth"
)

var Scopes = []string{
	spotifyauth.ScopeUserReadPlaybackState,
	spotifyauth.ScopeUserReadCurrentlyPlaying,
	spotifyauth.ScopeStreaming,
	spotifyauth.ScopePlaylistReadPrivate,
	spotifyauth.ScopePlaylistReadCollaborative,
	spotifyauth.ScopeUserFollowRead,
	spotifyauth.ScopeUserTopRead,
	spotifyauth.ScopeUserReadRecentlyPlayed,
	spotifyauth.ScopeUserLibraryRead,
	spotifyauth.ScopeUserReadEmail,
	spotifyauth.ScopeUserReadPrivate,
}
