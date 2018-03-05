package main

import (
	"fmt"
)

type ConfigNotLoadedError struct{}

type KeyTooYoungError struct {
	key string
}

type MultipleKeysFoundError struct{}

type NoKeysFoundError struct{}

func (e *ConfigNotLoadedError) Error() string {
	return fmt.Sprintf("Unable to load SDK config")
}

func (e *KeyTooYoungError) Error() string {
	return fmt.Sprintf("Key %s is less than 30 days old and there is no need for age-related rotation.  Use --force to override this.", e.key)
}

func (e *MultipleKeysFoundError) Error() string {
	return fmt.Sprintf("You have multiple keys available.  Specify --key NUM to choose one")
}

func (e *NoKeysFoundError) Error() string {
	return fmt.Sprintf("You do not have any access keys, or none can be found.")
}
