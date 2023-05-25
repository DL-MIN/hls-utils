package stats

// ErrLogIsDir indicates that there is a directory instead of a file in the given path.
// It contains the failed path.
type ErrLogIsDir string

func (e ErrLogIsDir) Error() string {
	return "fifo path is a directory: " + string(e)
}

// ErrPlaylistNotExist indicates that a playlist file is missing.
// It contains the failed path.
type ErrPlaylistNotExist string

func (e ErrPlaylistNotExist) Error() string {
	return "playlist does not exist: " + string(e)
}
