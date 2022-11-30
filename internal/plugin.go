package internal

type plugin struct {
	// Name is how plugins are referenced in feed configuration
	Name string
	// Path on disc to the plugin. This also determines how it's interpreted.
	// The endings .so, .dll, and .dylib will be treated as binary plugins.
	// Other suffixes will be treated as executable files. There are a
	// reserved value here: `regex` may be used here to specify that the
	// plugin is actually a regex search/replace pattern on the body of a
	// feed entry. This is
	Path string
}
