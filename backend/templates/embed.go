package templates

import "embed"

//go:embed mail/*
var MailTemplateFiles embed.FS
