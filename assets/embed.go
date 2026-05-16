package assets

import "embed"

//go:embed fonts/*
//go:embed audio/*
//go:embed pictures/*
var Assets embed.FS
