package flowutils

//VersionSample stores the format of the function Version()
type VersionSample struct {
	Major   int
	Minor   int
	Comment string
}

//Version returns the current version of this utility
func Version() VersionSample {
	return VersionSample{
		Major:   1,
		Minor:   1,
		Comment: "Alpha",
	}
}
