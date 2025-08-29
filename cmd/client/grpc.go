package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"time"

	pb "github.com/MukizuL/GophKeeper/internal/proto"
	tea "github.com/charmbracelet/bubbletea"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
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
	if err != nil {
		if e, ok := status.FromError(err); ok {
			switch e.Code() {
			case codes.DeadlineExceeded:
				return fmt.Errorf("server took to long to respond: %s", e.Message())
			case codes.FailedPrecondition:
				return fmt.Errorf("user with same login already exists: %s", e.Message())
			case codes.Internal:
				return fmt.Errorf("server error: %s", e.Message())
			default:
				return fmt.Errorf("unknown error: %s", e.Message())
			}
		}
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
	if err != nil {
		if e, ok := status.FromError(err); ok {
			switch e.Code() {
			case codes.DeadlineExceeded:
				return "", nil, fmt.Errorf("server took to long to respond: %s", e.Message())
			case codes.Unauthenticated:
				return "", nil, fmt.Errorf("%s", e.Message())
			case codes.Internal:
				return "", nil, fmt.Errorf("server error: %s", e.Message())
			default:
				return "", nil, fmt.Errorf("unknown error: %s", e.Message())
			}
		}
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
	if err != nil {
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
	if err != nil {
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
	if err != nil {
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
	}

	return nil
}

func CreateData(token string, dk []byte, filename string, percent *float64) tea.Cmd {
	return func() tea.Msg {
		f, err := os.Open(filename)
		if err != nil {
			return errMsg{fmt.Errorf("could not open file: %s", filename)}
		}
		defer f.Close()

		stat, _ := f.Stat()
		filesize := stat.Size()

		md := metadata.Pairs("access-token", token)
		ctxOut := metadata.NewOutgoingContext(context.Background(), md)

		stream, err := conn.CreateData(ctxOut)
		if err != nil {
			return errMsg{errors.New("could not create stream")}
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

			req := &pb.CreateDataRequest{
				Filename: filename,
				Chunk:    buf[:n],
			}
			if err := stream.Send(req); err != nil {
				return errMsg{errors.New("could not send data")}
			}

			sent += int64(n)
			*percent = float64(sent) / float64(filesize)
		}

		_, err = stream.CloseAndRecv()
		if err != nil {
			return errMsg{errors.New("could not receive data")}
		}

		return uploadProgressMsg(1.0)
	}
}

func GetPasswords(token string, dk []byte) ([][]byte, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	md := metadata.Pairs("access-token", token)
	ctxOut := metadata.NewOutgoingContext(ctx, md)

	req := pb.GetPasswordsRequest{}

	data, err := conn.GetPasswords(ctxOut, &req)
	if err != nil {
		if e, ok := status.FromError(err); ok {
			switch e.Code() {
			case codes.DeadlineExceeded:
				return nil, fmt.Errorf("server took to long to respond: %s", e.Message())
			case codes.Unauthenticated:
				return nil, fmt.Errorf("%s", e.Message())
			case codes.Internal:
				return nil, fmt.Errorf("server error: %s", e.Message())
			default:
				return nil, fmt.Errorf("unknown error: %s", e.Message())
			}
		}
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
	if err != nil {
		if e, ok := status.FromError(err); ok {
			switch e.Code() {
			case codes.DeadlineExceeded:
				return nil, fmt.Errorf("server took to long to respond: %s", e.Message())
			case codes.Unauthenticated:
				return nil, fmt.Errorf("%s", e.Message())
			case codes.Internal:
				return nil, fmt.Errorf("server error: %s", e.Message())
			default:
				return nil, fmt.Errorf("unknown error: %s", e.Message())
			}
		}
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
	if err != nil {
		if e, ok := status.FromError(err); ok {
			switch e.Code() {
			case codes.DeadlineExceeded:
				return nil, fmt.Errorf("server took to long to respond: %s", e.Message())
			case codes.Unauthenticated:
				return nil, fmt.Errorf("%s", e.Message())
			case codes.Internal:
				return nil, fmt.Errorf("server error: %s", e.Message())
			default:
				return nil, fmt.Errorf("unknown error: %s", e.Message())
			}
		}
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
