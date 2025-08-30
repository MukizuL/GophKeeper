package main

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"errors"
	"fmt"
	"io"
	"strings"
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

// WrapNoSplitWords splits s into lines of length <= n without splitting words.
// It normalizes whitespace (multiple spaces/tabs/newlines -> single spaces between words).
func WrapNoSplitWords(s string, n int) ([]string, error) {
	if n <= 0 {
		return nil, errors.New("n must be > 0")
	}

	words := strings.Fields(s)
	if len(words) == 0 {
		return []string{}, nil
	}

	lines := make([]string, 0, len(words))
	var cur strings.Builder
	curLen := 0 // rune count in cur

	for _, w := range words {
		wLen := utf8.RuneCountInString(w)
		if wLen > n {
			return nil, fmt.Errorf("a word exceeds the limit of %d chars", n)
		}

		if curLen == 0 {
			// start a new line
			cur.WriteString(w)
			curLen = wLen
			continue
		}

		// try to add " " + w
		if curLen+1+wLen <= n {
			cur.WriteByte(' ')
			cur.WriteString(w)
			curLen += 1 + wLen
		} else {
			// finish current line, start a new one with w
			lines = append(lines, cur.String())
			cur.Reset()
			cur.WriteString(w)
			curLen = wLen
		}
	}

	if curLen > 0 {
		lines = append(lines, cur.String())
	}
	return lines, nil
}

func updateProgressBar(ch chan tea.Msg) tea.Cmd {
	return func() tea.Msg {
		for m := range ch {
			return m
		}
		return nil
	}
}
