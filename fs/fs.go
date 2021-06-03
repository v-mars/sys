package fs

import "embed"

////go:embed jobor.sql
//var SqlFs embed.FS

//go:embed dist/*
var DistFs embed.FS

////go:embed ./docs/*
//var SwaggerFs embed.FS