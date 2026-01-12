package web

import "embed"

//go:embed static seo/* templates/*
var Files embed.FS

func GetRobots() ([]byte, error)  { return Files.ReadFile("seo/robots.txt") }
func GetSitemap() ([]byte, error) { return Files.ReadFile("seo/sitemap.xml") }
