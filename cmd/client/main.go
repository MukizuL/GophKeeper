package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	pb "github.com/MukizuL/GophKeeper/internal/proto"
	tea "github.com/charmbracelet/bubbletea"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/insecure"
)

var (
	buildVersion string
	buildDate    string
	buildCommit  string

	conn pb.GophkeeperClient
)

func main() {
	cfg, err := newConfig()
	if err != nil {
		flag.Usage()
		log.Fatal(err)
	}

	cgrpc, err := newGRPConn(cfg)
	if err != nil {
		log.Fatal(err)
	}
	defer cgrpc.Close()

	conn = pb.NewGophkeeperClient(cgrpc)

	// TODO: Check if server is available

	s := newStart()
	a := newAbout()
	r := newRegister()
	l := newLogin()
	h := newHome()
	c := newCreate()
	cpass := newCreatePassword()
	cbank := newCreateBank()
	ctext := newCreateText()
	cdata := newCreateData()
	store := newStorage()

	m := model{
		start:          s,
		about:          a,
		register:       r,
		login:          l,
		home:           h,
		create:         c,
		createPassword: cpass,
		createBank:     cbank,
		createText:     ctext,
		createData:     cdata,
		storage:        store,
		window:         "start",
	}

	if _, err = tea.NewProgram(m).Run(); err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}
}

func newGRPConn(cfg *config) (*grpc.ClientConn, error) {
	if cfg.TLS {
		creds, err := credentials.NewClientTLSFromFile(cfg.Cert, "")
		if err != nil {
			return nil, fmt.Errorf("failed to load TLS credentials: %w", err)
		}

		c, err := grpc.NewClient(cfg.Addr, grpc.WithTransportCredentials(creds))
		if err != nil {
			return nil, fmt.Errorf("failed to connect: %w", err)
		}

		return c, nil
	} else {
		c, err := grpc.NewClient(cfg.Addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
		if err != nil {
			return nil, fmt.Errorf("failed to connect: %w", err)
		}

		return c, nil
	}
}
