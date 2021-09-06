package flowutils

type VersionSample struct {
	Major   int
	Minor   int
	Comment string
}

func Version() VersionSample {
	return VersionSample{
		Major:   1,
		Minor:   1,
		Comment: "Alpha",
	}
}
