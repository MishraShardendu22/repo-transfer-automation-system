package main

import (
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/go-resty/resty/v2"
	"github.com/joho/godotenv"
)

// TransferRequest payload from frontend or API
type TransferRequest struct {
	TargetUser   string   `json:"target_user"`
	Repositories []string `json:"repositories"`
}

// TransferResult output for each repository
type TransferResult struct {
	Repo       string `json:"repo"`
	Success    bool   `json:"success"`
	StatusCode int    `json:"status_code"`
	Message    string `json:"message"`
}

// GitHubUser holds basic profile details
type GitHubUser struct {
	Login     string `json:"login"`
	ID        int64  `json:"id"`
	AvatarURL string `json:"avatar_url"`
	Name      string `json:"name"`
	Email     string `json:"email"`
}

// GitHubRepo representation
type GitHubRepo struct {
	Name        string `json:"name"`
	FullName    string `json:"full_name"`
	Private     bool   `json:"private"`
	Description string `json:"description"`
	UpdatedAt   string `json:"updated_at"`
	Owner       struct {
		Login string `json:"login"`
	} `json:"owner"`
}

// App configuration and in-memory session store
type AppServer struct {
	ClientID     string
	ClientSecret string
	AppURL       string
	Port         string
	Client       *resty.Client
	Sessions     map[string]string // sessionToken -> githubAccessToken
	SessionUsers map[string]GitHubUser
	Mu           sync.RWMutex
}

func newAppServer() *AppServer {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	appURL := os.Getenv("APP_URL")
	if appURL == "" {
		appURL = "http://localhost:" + port
	}

	client := resty.New().
		SetRetryCount(3).
		SetRetryWaitTime(2 * time.Second).
		SetRetryMaxWaitTime(10 * time.Second).
		SetCloseConnection(true)

	return &AppServer{
		ClientID:     os.Getenv("GITHUB_CLIENT_ID"),
		ClientSecret: os.Getenv("GITHUB_CLIENT_SECRET"),
		AppURL:       strings.TrimSuffix(appURL, "/"),
		Port:         port,
		Client:       client,
		Sessions:     make(map[string]string),
		SessionUsers: make(map[string]GitHubUser),
	}
}

func generateRandomState() string {
	b := make([]byte, 16)
	_, _ = rand.Read(b)
	return hex.EncodeToString(b)
}

// Resolve GitHub token from Authorization header, session cookie, or environment variable
func (s *AppServer) getToken(r *http.Request) string {
	authHeader := r.Header.Get("Authorization")
	if strings.HasPrefix(authHeader, "Bearer ") {
		token := strings.TrimPrefix(authHeader, "Bearer ")
		if token != "" {
			return token
		}
	}

	cookie, err := r.Cookie("gh_session")
	if err == nil && cookie.Value != "" {
		s.Mu.RLock()
		token, exists := s.Sessions[cookie.Value]
		s.Mu.RUnlock()
		if exists && token != "" {
			return token
		}
	}

	return os.Getenv("GITHUB_TOKEN_CLASSIC")
}

func (s *AppServer) handleHome(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	_, _ = w.Write([]byte(dashboardHTML))
}

func (s *AppServer) handleOAuthLogin(w http.ResponseWriter, r *http.Request) {
	if s.ClientID == "" {
		http.Error(w, "GITHUB_CLIENT_ID not configured in environment. Please set GITHUB_CLIENT_ID and GITHUB_CLIENT_SECRET.", http.StatusBadRequest)
		return
	}
	state := generateRandomState()
	redirectURI := fmt.Sprintf("%s/api/auth/callback", s.AppURL)
	scope := "repo,admin:repo_hook"

	authURL := fmt.Sprintf(
		"https://github.com/login/oauth/authorize?client_id=%s&redirect_uri=%s&scope=%s&state=%s",
		s.ClientID, redirectURI, scope, state,
	)
	http.Redirect(w, r, authURL, http.StatusTemporaryRedirect)
}

func (s *AppServer) handleOAuthCallback(w http.ResponseWriter, r *http.Request) {
	code := r.URL.Query().Get("code")
	if code == "" {
		http.Error(w, "Missing authorization code from GitHub OAuth", http.StatusBadRequest)
		return
	}

	// Exchange code for access token
	resp, err := s.Client.R().
		SetHeader("Accept", "application/json").
		SetFormData(map[string]string{
			"client_id":     s.ClientID,
			"client_secret": s.ClientSecret,
			"code":          code,
		}).
		Post("https://github.com/login/oauth/access_token")

	if err != nil {
		http.Error(w, "Failed to exchange token with GitHub: "+err.Error(), http.StatusInternalServerError)
		return
	}

	var tokenResp struct {
		AccessToken string `json:"access_token"`
		TokenType   string `json:"token_type"`
		Scope       string `json:"scope"`
		Error       string `json:"error"`
	}

	if err := json.Unmarshal(resp.Body(), &tokenResp); err != nil || tokenResp.AccessToken == "" {
		http.Error(w, "Invalid token response from GitHub: "+string(resp.Body()), http.StatusBadRequest)
		return
	}

	// Fetch user details
	userResp, err := s.Client.R().
		SetHeader("Accept", "application/vnd.github+json").
		SetHeader("Authorization", "Bearer "+tokenResp.AccessToken).
		SetHeader("X-GitHub-Api-Version", "2022-11-28").
		Get("https://api.github.com/user")

	if err != nil || userResp.StatusCode() != 200 {
		http.Error(w, "Failed to fetch user profile: "+string(userResp.Body()), http.StatusInternalServerError)
		return
	}

	var user GitHubUser
	_ = json.Unmarshal(userResp.Body(), &user)

	sessionID := generateRandomState()
	s.Mu.Lock()
	s.Sessions[sessionID] = tokenResp.AccessToken
	s.SessionUsers[sessionID] = user
	s.Mu.Unlock()

	http.SetCookie(w, &http.Cookie{
		Name:     "gh_session",
		Value:    sessionID,
		Path:     "/",
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
		MaxAge:   86400 * 7, // 7 days
	})

	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func (s *AppServer) handleAuthMe(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("gh_session")
	if err != nil || cookie.Value == "" {
		http.Error(w, `{"error":"unauthenticated"}`, http.StatusUnauthorized)
		return
	}

	s.Mu.RLock()
	user, exists := s.SessionUsers[cookie.Value]
	s.Mu.RUnlock()

	if !exists {
		http.Error(w, `{"error":"session expired"}`, http.StatusUnauthorized)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(user)
}

func (s *AppServer) handleLogout(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("gh_session")
	if err == nil && cookie.Value != "" {
		s.Mu.Lock()
		delete(s.Sessions, cookie.Value)
		delete(s.SessionUsers, cookie.Value)
		s.Mu.Unlock()
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "gh_session",
		Value:    "",
		Path:     "/",
		MaxAge:   -1,
		HttpOnly: true,
	})

	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func (s *AppServer) handleListRepos(w http.ResponseWriter, r *http.Request) {
	token := s.getToken(r)
	if token == "" {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusUnauthorized)
		_ = json.NewEncoder(w).Encode(map[string]string{
			"error": "No GitHub authorization provided. Please sign in via GitHub OAuth or provide a token.",
		})
		return
	}

	var allRepos []GitHubRepo
	page := 1

	// Paginate through user repositories (affiliations: owner, collaborator, organization_member)
	for {
		url := fmt.Sprintf("https://api.github.com/user/repos?per_page=100&page=%d&sort=updated", page)
		resp, err := s.Client.R().
			SetHeader("Accept", "application/vnd.github+json").
			SetHeader("Authorization", "Bearer "+token).
			SetHeader("X-GitHub-Api-Version", "2022-11-28").
			Get(url)

		if err != nil {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusInternalServerError)
			_ = json.NewEncoder(w).Encode(map[string]string{"error": "Failed fetching repos: " + err.Error()})
			return
		}

		if resp.StatusCode() != 200 {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(resp.StatusCode())
			_ = json.NewEncoder(w).Encode(map[string]string{"error": fmt.Sprintf("GitHub API returned %d: %s", resp.StatusCode(), resp.String())})
			return
		}

		var repos []GitHubRepo
		if err := json.Unmarshal(resp.Body(), &repos); err != nil {
			break
		}

		if len(repos) == 0 {
			break
		}

		allRepos = append(allRepos, repos...)
		if len(repos) < 100 {
			break
		}
		page++
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(allRepos)
}

func (s *AppServer) handleTransfer(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	token := s.getToken(r)
	if token == "" {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusUnauthorized)
		_ = json.NewEncoder(w).Encode(map[string]string{"error": "Unauthorized: No GitHub token provided."})
		return
	}

	var req TransferRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid JSON payload: "+err.Error(), http.StatusBadRequest)
		return
	}

	if req.TargetUser == "" || len(req.Repositories) == 0 {
		http.Error(w, "target_user and at least one repository are required", http.StatusBadRequest)
		return
	}

	// Fetch current authenticated user to know original owner
	userResp, err := s.Client.R().
		SetHeader("Accept", "application/vnd.github+json").
		SetHeader("Authorization", "Bearer "+token).
		SetHeader("X-GitHub-Api-Version", "2022-11-28").
		Get("https://api.github.com/user")

	if err != nil || userResp.StatusCode() != 200 {
		http.Error(w, "Could not determine authenticated user from token", http.StatusUnauthorized)
		return
	}

	var currentUser GitHubUser
	_ = json.Unmarshal(userResp.Body(), &currentUser)
	originalUser := currentUser.Login

	results := make([]TransferResult, len(req.Repositories))
	var wg sync.WaitGroup

	for i, repoName := range req.Repositories {
		wg.Add(1)
		go func(idx int, repo string) {
			defer wg.Done()

			// Check if repo has owner prefix like owner/repo or just repo
			owner := originalUser
			actualRepo := repo
			if strings.Contains(repo, "/") {
				parts := strings.SplitN(repo, "/", 2)
				owner = parts[0]
				actualRepo = parts[1]
			}

			url := fmt.Sprintf("https://api.github.com/repos/%s/%s/transfer", owner, actualRepo)
			res, transferErr := s.Client.R().
				SetBody(map[string]string{
					"new_owner": req.TargetUser,
					"new_name":  actualRepo,
				}).
				SetHeader("Accept", "application/vnd.github+json").
				SetHeader("Authorization", "Bearer "+token).
				SetHeader("X-GitHub-Api-Version", "2022-11-28").
				Post(url)

			if transferErr != nil {
				results[idx] = TransferResult{
					Repo:       repo,
					Success:    false,
					StatusCode: 500,
					Message:    transferErr.Error(),
				}
				return
			}

			if res.StatusCode() == 202 {
				results[idx] = TransferResult{
					Repo:       repo,
					Success:    true,
					StatusCode: 202,
					Message:    "Transfer accepted by GitHub",
				}
			} else {
				results[idx] = TransferResult{
					Repo:       repo,
					Success:    false,
					StatusCode: res.StatusCode(),
					Message:    res.String(),
				}
			}
		}(i, repoName)
	}

	wg.Wait()

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(map[string]interface{}{
		"target_user": req.TargetUser,
		"results":     results,
	})
}

func runCLI(sourceUser, targetUser string, repos []string) {
	fmt.Printf("Starting CLI repository transfer: %s -> %s\n", sourceUser, targetUser)
	token := os.Getenv("GITHUB_TOKEN_CLASSIC")
	if token == "" {
		token = os.Getenv("GITHUB_TOKEN")
	}
	if token == "" {
		log.Fatal("GITHUB_TOKEN_CLASSIC or GITHUB_TOKEN environment variable required for CLI mode")
	}

	client := resty.New().
		SetRetryCount(3).
		SetRetryWaitTime(2 * time.Second).
		SetRetryMaxWaitTime(10 * time.Second).
		SetCloseConnection(true)

	var wg sync.WaitGroup
	for _, repo := range repos {
		wg.Add(1)
		go func(r string) {
			defer wg.Done()
			url := fmt.Sprintf("https://api.github.com/repos/%s/%s/transfer", sourceUser, r)
			res, err := client.R().
				SetBody(map[string]string{
					"new_owner": targetUser,
					"new_name":  r,
				}).
				SetHeader("Accept", "application/vnd.github+json").
				SetHeader("Authorization", "Bearer "+token).
				SetHeader("X-GitHub-Api-Version", "2022-11-28").
				Post(url)

			if err != nil {
				fmt.Printf("[ERROR] %s: %v\n", r, err)
				return
			}

			if res.StatusCode() != 202 {
				fmt.Printf("[FAIL] %s: %d - %s\n", r, res.StatusCode(), res.String())
			} else {
				fmt.Printf("[SUCCESS] %s: %s\n", r, res.Status())
			}
		}(repo)
	}
	wg.Wait()
	fmt.Println("Transfer run finished.")
}

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
		runCLI(*sourceFlag, *targetFlag, repos)
		return
	}

	server := newAppServer()

	http.HandleFunc("/", server.handleHome)
	http.HandleFunc("/api/auth/login", server.handleOAuthLogin)
	http.HandleFunc("/api/auth/callback", server.handleOAuthCallback)
	http.HandleFunc("/api/auth/me", server.handleAuthMe)
	http.HandleFunc("/api/auth/logout", server.handleLogout)
	http.HandleFunc("/api/repos", server.handleListRepos)
	http.HandleFunc("/api/transfer", server.handleTransfer)

	addr := ":" + server.Port
	fmt.Println("==========================================================")
	fmt.Printf("GitHub Repository Transfer Engine running on %s\n", server.AppURL)
	fmt.Printf("Web Dashboard: %s\n", server.AppURL)
	fmt.Printf("OAuth Callback: %s/api/auth/callback\n", server.AppURL)
	fmt.Println("==========================================================")

	if err := http.ListenAndServe(addr, nil); err != nil {
		log.Fatalf("Server exited: %v", err)
	}
}
