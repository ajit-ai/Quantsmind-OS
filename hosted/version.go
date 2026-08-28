package main

// Version is stamped at build time via -ldflags "-X main.Version=<tag>".
// Falls back to "dev" for local builds.
var Version = "dev"

const Product = "quantis"
