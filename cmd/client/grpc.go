package main

import (
	"context"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"time"

	pb "github.com/MukizuL/GophKeeper/internal/proto"
	tea "github.com/charmbracelet/bubbletea"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

func Register(login, password string) error {
	err := validateLogin(login)
	if err != nil {
		return err
	}

	err = validatePassword(password)
	if err != nil {
		return err
	}

	req := pb.RegisterRequest{
		Login:    login,
		Password: password,
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	_, err = conn.Register(ctx, &req)
	err = handleGRPCError(err)
	if err != nil {
		return err
	}

	return nil
}

func Login(login, password string) (string, []byte, error) {
	err := validateLogin(login)
	if err != nil {
		return "", nil, err
	}

	err = validatePassword(password)
	if err != nil {
		return "", nil, err
	}

	req := pb.AuthRequest{
		Login:    login,
		Password: password,
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	var header metadata.MD
	_, err = conn.Authorize(ctx, &req, grpc.Header(&header))
	err = handleGRPCError(err)
	if err != nil {
		return "", nil, err
	}

	tokens := header.Get("access-token")
	if len(tokens) == 0 {
		return "", nil, fmt.Errorf("access token is missing")
	}

	dks := header.Get("dk-bin")
	if len(tokens) == 0 {
		return "", nil, fmt.Errorf("derived key is missing")
	}

	return tokens[0], []byte(dks[0]), nil
}

func CreatePassword(token string, dk []byte, name, login, password, description string) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	md := metadata.Pairs("access-token", token)
	ctxOut := metadata.NewOutgoingContext(ctx, md)

	data, err := json.Marshal(map[string]string{"name": name, "login": login, "password": password, "description": description})
	if err != nil {
		return err
	}

	encrypted, err := encrypt(dk, data)
	if err != nil {
		return err
	}

	req := pb.CreatePasswordRequest{
		Data: encrypted,
	}

	_, err = conn.CreatePassword(ctxOut, &req)
	err = handleGRPCError(err)
	if err != nil {
		return err
	}

	return nil
}

func CreateBank(token string, dk []byte, ccn, exp, cvv, name string) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	md := metadata.Pairs("access-token", token)
	ctxOut := metadata.NewOutgoingContext(ctx, md)

	data, err := json.Marshal(map[string]string{"ccn": ccn, "exp": exp, "cvv": cvv, "name": name})
	if err != nil {
		return err
	}

	encrypted, err := encrypt(dk, data)
	if err != nil {
		return err
	}

	req := pb.CreateBankRequest{
		Data: encrypted,
	}

	_, err = conn.CreateBank(ctxOut, &req)
	err = handleGRPCError(err)
	if err != nil {
		return err
	}

	return nil
}

func CreateTextual(token string, dk []byte, name, text string) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	md := metadata.Pairs("access-token", token)
	ctxOut := metadata.NewOutgoingContext(ctx, md)

	data, err := json.Marshal(map[string]string{"name": name, "text": text})
	if err != nil {
		return err
	}

	encrypted, err := encrypt(dk, data)
	if err != nil {
		return err
	}

	req := pb.CreateTextRequest{
		Data: encrypted,
	}

	_, err = conn.CreateText(ctxOut, &req)
	err = handleGRPCError(err)
	if err != nil {
		return err
	}

	return nil
}

func CreateData(token string, dk []byte, path string, ch chan tea.Msg) tea.Cmd {
	return func() tea.Msg {
		f, err := os.Open(path)
		if err != nil {
			return errMsg{fmt.Errorf("could not open file: %s", path)}
		}
		defer f.Close()

		_, filename := filepath.Split(path)

		stat, _ := f.Stat()
		filesize := stat.Size()

		md := metadata.Pairs("access-token", token)
		ctxOut := metadata.NewOutgoingContext(context.Background(), md)

		stream, err := conn.CreateData(ctxOut)
		if err != nil {
			return errMsg{errors.New("could not create stream")}
		}

		block, err := aes.NewCipher(dk)
		if err != nil {
			return errMsg{err}
		}

		gcm, err := cipher.NewGCM(block)
		if err != nil {
			return errMsg{err}
		}

		nonce := make([]byte, gcm.NonceSize())
		if _, err = io.ReadFull(rand.Reader, nonce); err != nil {
			return errMsg{err}
		}

		buf := make([]byte, 1024*32) // 32KB chunks
		var sent int64

		for {
			n, err := f.Read(buf)
			if err == io.EOF {
				break
			}
			if err != nil {
				return errMsg{errors.New("could not read file")}
			}

			encryptedFilename, err := encrypt(dk, []byte(filename))
			if err != nil {
				return errMsg{err}
			}

			encryptedChunk, err := encrypt(dk, buf[:n])
			if err != nil {
				return errMsg{err}
			}

			req := &pb.CreateDataRequest{
				Filename: encryptedFilename,
				Chunk:    encryptedChunk,
			}

			if err := stream.Send(req); err != nil {
				return errMsg{errors.New("could not send data")}
			}

			sent += int64(n)
			ch <- uploadProgressMsg(float64(sent) / float64(filesize))
		}

		_, err = stream.CloseAndRecv()
		if err != nil {
			return errMsg{errors.New("could not receive data")}
		}

		ch <- uploadProgressMsg(1.0)

		return nil
	}
}

func GetPasswords(token string, dk []byte) ([][]byte, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	md := metadata.Pairs("access-token", token)
	ctxOut := metadata.NewOutgoingContext(ctx, md)

	req := pb.GetPasswordsRequest{}

	data, err := conn.GetPasswords(ctxOut, &req)
	err = handleGRPCError(err)
	if err != nil {
		return nil, err
	}

	var passwords [][]byte
	for _, v := range data.Data {
		decrypted, err := decrypt(dk, v)
		if err != nil {
			return nil, err
		}

		passwords = append(passwords, decrypted)
	}

	return passwords, nil
}

func GetBank(token string, dk []byte) ([][]byte, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	md := metadata.Pairs("access-token", token)
	ctxOut := metadata.NewOutgoingContext(ctx, md)

	req := pb.GetBankRequest{}

	data, err := conn.GetBank(ctxOut, &req)
	err = handleGRPCError(err)
	if err != nil {
		return nil, err
	}

	var banks [][]byte
	for _, v := range data.Data {
		decrypted, err := decrypt(dk, v)
		if err != nil {
			return nil, err
		}

		banks = append(banks, decrypted)
	}

	return banks, nil
}

func GetText(token string, dk []byte) ([][]byte, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	md := metadata.Pairs("access-token", token)
	ctxOut := metadata.NewOutgoingContext(ctx, md)

	req := pb.GetTextRequest{}

	data, err := conn.GetText(ctxOut, &req)
	err = handleGRPCError(err)
	if err != nil {
		return nil, err
	}

	var texts [][]byte
	for _, v := range data.Data {
		decrypted, err := decrypt(dk, v)
		if err != nil {
			return nil, err
		}

		texts = append(texts, decrypted)
	}

	return texts, nil
}

func GetData(token string, dk []byte) ([]file, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	md := metadata.Pairs("access-token", token)
	ctxOut := metadata.NewOutgoingContext(ctx, md)

	req := pb.GetDataRequest{}

	data, err := conn.GetData(ctxOut, &req)
	err = handleGRPCError(err)
	if err != nil {
		return nil, err
	}

	var files []file
	for _, v := range data.File {
		decryptedFilename, err := decrypt(dk, v.Filename)
		if err != nil {
			return nil, err
		}

		temp := file{
			ID:       v.Id,
			Filename: string(decryptedFilename),
		}

		files = append(files, temp)
	}

	return files, nil
}

func DownloadFile(token string, dk []byte, id, filename string) tea.Cmd {
	return func() tea.Msg {
		md := metadata.Pairs("access-token", token)
		ctxOut := metadata.NewOutgoingContext(context.Background(), md)

		req := pb.DownloadRequest{Id: id}

		stream, err := conn.Download(ctxOut, &req)
		err = handleGRPCError(err)
		if err != nil {
			return errMsg{err}
		}

		f, err := os.Create(filename)
		if err != nil {
			return err
		}
		defer f.Close()

		for {
			chunk, err := stream.Recv()
			if err == io.EOF {
				break
			}
			if err != nil {
				return errMsg{errors.New("error receiving chunk")}
			}

			decrypted, err := decrypt(dk, chunk.Chunk)
			if err != nil {
				return errMsg{err}
			}

			if _, err := f.Write(decrypted); err != nil {
				return err
			}
		}

		return Done{}
	}
}
