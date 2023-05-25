package stats

type ErrLogIsDir string

func (e ErrLogIsDir) Error() string {
	return "fifo path is a directory: " + string(e)
}

type ErrPlaylistNotExist string

func (e ErrPlaylistNotExist) Error() string {
	return "playlist does not exist: " + string(e)
}
