package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"strings"
	"syscall"

	"github.com/sidalikharef/goremote/config"
	"github.com/sidalikharef/goremote/internal/auth"
	"golang.org/x/term"
)

func main() {
	var (
		username = flag.String("username", "", "Username")
		conf     = flag.String("config", config.DefaultPath(), "Config path")
	)
	flag.Parse()
	if *username == "" {
		fmt.Fprintln(os.Stderr, "username required")
		os.Exit(1)
	}
	store := auth.New(*conf + ".auth.json")
	if store.HasUsers() {
		fmt.Fprintln(os.Stderr, "users already exist")
		os.Exit(1)
	}
	token, err := store.GenerateBootstrapToken()
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("Bootstrap token: %s\n", token)
	fmt.Printf("Run from loopback to enroll: POST /api/auth/enroll with this token\n")
	fmt.Printf("Or use this interactive enrollment (local only):\n")
	fmt.Printf("Password: ")
	bytepw, err := term.ReadPassword(int(syscall.Stdin))
	fmt.Println()
	if err != nil {
		// fallback
		reader := bufio.NewReader(os.Stdin)
		pw, _ := reader.ReadString('\n')
		pw = strings.TrimSpace(pw)
		if len(pw) < 8 {
			fmt.Fprintln(os.Stderr, "password too short")
			os.Exit(1)
		}
		if err := store.Enroll(strings.TrimSpace(*username), pw, token); err != nil {
			fmt.Fprintf(os.Stderr, "enroll error: %v\n", err)
			os.Exit(1)
		}
		fmt.Println("enrolled successfully")
		return
	}
	pw := strings.TrimSpace(string(bytepw))
	if len(pw) < 8 {
		fmt.Fprintln(os.Stderr, "password too short")
		os.Exit(1)
	}
	if err := store.Enroll(strings.TrimSpace(*username), pw, token); err != nil {
		fmt.Fprintf(os.Stderr, "enroll error: %v\n", err)
		os.Exit(1)
	}
	fmt.Println("enrolled successfully")
}
