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
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

const iterations = 1     // number of iterations
const memory = 64 * 1024 // 64 MB
const threads = 4        // parallelism
const keyLen = 32        // AES-256

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
	key := argon2.IDKey([]byte(password), salt, iterations, memory, uint8(threads), keyLen)
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

func handleGRPCError(err error) error {
	if e, ok := status.FromError(err); ok {
		switch e.Code() {
		case codes.DeadlineExceeded:
			return fmt.Errorf("server took to long to respond: %s", e.Message())
		case codes.Unauthenticated:
			return fmt.Errorf("%s", e.Message())
		case codes.Internal:
			return fmt.Errorf("server error: %s", e.Message())
		default:
			return fmt.Errorf("unknown error: %s", e.Message())
		}
	}
	return err
}
