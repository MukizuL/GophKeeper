package main

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"fmt"
	"io"
	"unicode/utf8"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"golang.org/x/crypto/argon2"
)

func validateLogin(login string) error {
	if utf8.RuneCountInString(login) < 3 {
		return fmt.Errorf("login must be at least 3 characters")
	}
	if utf8.RuneCountInString(login) > 255 {
		return fmt.Errorf("login must be at most 255 characters")
	}

	return nil
}

func validatePassword(password string) error {
	if utf8.RuneCountInString(password) < 8 {
		return fmt.Errorf("password must be at least 8 characters")
	}
	if utf8.RuneCountInString(password) > 36 {
		return fmt.Errorf("password must be at most 36 characters")
	}
	if len(password) > 72 {
		return fmt.Errorf("password is %d characters but is longer than 72 bytes", utf8.RuneCountInString(password))
	}

	return nil
}

func focusOrBlur(inputs []textinput.Model, focusIndex int) []tea.Cmd {
	cmds := make([]tea.Cmd, len(inputs))
	for i := range inputs {
		if i == focusIndex {
			cmds[i] = inputs[i].Focus()
			inputs[i].PromptStyle = formSelectedStyle
			inputs[i].TextStyle = formSelectedStyle
			continue
		}
		inputs[i].Blur()
		inputs[i].PromptStyle = formStyle
		inputs[i].TextStyle = formStyle
	}

	return cmds
}

// deriveKey generates a strong key from a password using Argon2id.
func deriveKey(password string, salt []byte) ([]byte, error) {
	const time = 1           // number of iterations
	const memory = 64 * 1024 // 64 MB
	const threads = 4        // parallelism
	const keyLen = 32        // AES-256

	key := argon2.IDKey([]byte(password), salt, time, memory, uint8(threads), keyLen)
	return key, nil
}

// Encrypt data using AES-GCM.
// key must be 16, 24, or 32 bytes (AES-128/192/256).
func encrypt(key []byte, plaintext []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	nonce := make([]byte, gcm.NonceSize())
	if _, err = io.ReadFull(rand.Reader, nonce); err != nil {
		return nil, err
	}

	ciphertext := gcm.Seal(nonce, nonce, plaintext, nil)
	return ciphertext, nil
}

// Decrypt data using AES-GCM.
func decrypt(key []byte, ciphertext []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	nonceSize := gcm.NonceSize()
	if len(ciphertext) < nonceSize {
		return nil, fmt.Errorf("ciphertext too short")
	}

	nonce, ciphertext := ciphertext[:nonceSize], ciphertext[nonceSize:]
	return gcm.Open(nil, nonce, ciphertext, nil)
}
