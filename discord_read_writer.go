package discordbot

// DiscordReadWriter implements io.ReadWriter by embedding structs that implement
// io.Reader and io.Writer
type DiscordReadWriter struct {
	*DiscordReader
	*DiscordWriter
}

// NewDiscordReadWriter returns struct that implements io.ReadWriter
func NewDiscordReadWriter(r *DiscordReader, w *DiscordWriter) *DiscordReadWriter {
	return &DiscordReadWriter{r, w}
}
