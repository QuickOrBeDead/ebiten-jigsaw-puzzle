package assets

import "embed"

//go:embed fonts/*
//go:embed audio/*
var Assets embed.FS
