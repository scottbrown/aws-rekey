package main

import (
	"fmt"
)

type UnsupportedCredentialsProviderError struct{}

func (e *UnsupportedCredentialsProviderError) Error() string {
	return fmt.Sprintf("Could not find credentials in the usual places.  You may be storing your AWS access keys in environment variables or in a custom credentials provider, both of which are unsupported at this time.")
}

type ConfigNotLoadedError struct{}

func (e *ConfigNotLoadedError) Error() string {
	return fmt.Sprintf("Unable to load SDK config")
}

type KeyTooYoungError struct {
	key string
}

func (e *KeyTooYoungError) Error() string {
	return fmt.Sprintf("Key %s is less than 30 days old and there is no need for age-related rotation.  Use --force to override this.", e.key)
}

type MultipleKeysFoundError struct{}

func (e *MultipleKeysFoundError) Error() string {
	return fmt.Sprintf("You have multiple keys available.  Specify --key NUM to choose one")
}

type NoKeysFoundError struct{}

func (e *NoKeysFoundError) Error() string {
	return fmt.Sprintf("You do not have any access keys, or none can be found.")
}
