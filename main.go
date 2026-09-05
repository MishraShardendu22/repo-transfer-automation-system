package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"

	"github.com/MishraShardendu22/pkg/server"
	"github.com/joho/godotenv"
)

func main() {
	_ = godotenv.Load()

	cliMode := flag.Bool("cli", false, "Run in headless CLI mode using hardcoded/configured parameters")
	sourceFlag := flag.String("source", "ShardenduMishra22", "Source GitHub account in CLI mode")
	targetFlag := flag.String("target", "MishraShardendu22", "Target GitHub account/org in CLI mode")
	flag.Parse()

	if *cliMode {
		repos := flag.Args()
		if len(repos) == 0 {
			repos = []string{
				"Tinder-LLD",
				"Go-TransferScript",
				"Zepto-LLD",
				"DiscountCoupon-LLD",
				"PaymentGateway-LLD",
			}
		}
		server.RunCLI(*sourceFlag, *targetFlag, repos)
		return
	}

	appServer := server.NewAppServer()
	handler := appServer.Handler()

	addr := ":" + appServer.Port
	fmt.Println("==========================================================")
	fmt.Printf("GitHub Repository Transfer Engine running on %s\n", appServer.AppURL)
	fmt.Printf("Web Dashboard: %s\n", appServer.AppURL)
	fmt.Printf("OAuth Callback: %s/api/auth/callback\n", appServer.AppURL)
	fmt.Println("==========================================================")

	if err := http.ListenAndServe(addr, handler); err != nil {
		log.Fatalf("Server exited: %v", err)
	}
}
