package config

import "log"

// Application struct for dependency injection
type Application struct {
	ErrorLog *log.Logger
	InfoLog *log.Logger
}