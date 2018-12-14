package main

type Config struct {
	Processes  []string `json:"processes"`
	PIDFileDir string   `json:"pid_file_dir"`
}
