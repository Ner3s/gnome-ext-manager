package models

type Extension struct {
	UUID       string `json:"uuid"`
	Enabled    bool   `json:"enabled"`
	URL        string `json:"url"`
	GnomeShell string `json:"gnome_shell_version"`
}
