package main

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"
)

var (
	staticDir    = "./static"
	userDataDir  = "./UserData/icons"
	chunksDir    = "./UserData/chunks"
	userDataRoot = "./UserData"
	validTokens  sync.Map
)

type PasswordData struct {
	PasswordHash string    `json:"passwordHash"`
	CreatedAt   time.Time `json:"createdAt"`
}

type LoginRequest struct {
	Password string `json:"password"`
}

type InitPasswordRequest struct {
	Password string `json:"password"`
}

func main() {
	os.MkdirAll(userDataDir, 0755)
	os.MkdirAll(chunksDir, 0755)
	os.MkdirAll(userDataRoot, 0755)

	absStatic, _ := filepath.Abs(staticDir)
	log.Printf("Static directory: %s", absStatic)

	if _, err := os.Stat(staticDir + "/index.html"); os.IsNotExist(err) {
		log.Printf("WARNING: index.html not found in %s", staticDir)
	} else {
		log.Printf("index.html found")
	}

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		log.Printf("%s %s", r.Method, r.URL.Path)
		if r.URL.Path == "/" {
			http.ServeFile(w, r, staticDir+"/index.html")
			return
		}
		http.FileServer(http.Dir(staticDir)).ServeHTTP(w, r)
	})

	http.HandleFunc("/api/check-init", checkInit)
	http.HandleFunc("/api/init-password", initPassword)
	http.HandleFunc("/api/login", login)
	http.HandleFunc("/api/logout", logout)
	http.HandleFunc("/api/list-upload-icons", listUploadIcons)
	http.HandleFunc("/api/upload/user_icon/init", withAuth(uploadInit))
	http.HandleFunc("/api/upload/user_icon/chunk", withAuth(uploadChunk))
	http.HandleFunc("/api/upload/user_icon/merge", withAuth(uploadMerge))
	http.HandleFunc("/api/rename/user_icon", withAuth(renameIcon))
	http.HandleFunc("/api/delete/user_icon", withAuth(deleteIcon))

	http.HandleFunc("/deskdata/user_icon/", func(w http.ResponseWriter, r *http.Request) {
		filePath := filepath.Join(userDataDir, filepath.Base(r.URL.Path))
		http.ServeFile(w, r, filePath)
	})

	log.Println("Server starting on :9168...")
	log.Fatal(http.ListenAndServe(":9168", nil))
}

func checkInit(w http.ResponseWriter, r *http.Request) {
	pwFile := filepath.Join(userDataRoot, "pw.json")
	_, err := os.Stat(pwFile)
	sendJSON(w, map[string]interface{}{
		"initialized": err == nil,
	})
}

func hashPassword(password string) string {
	hash := sha256.Sum256([]byte(password))
	return hex.EncodeToString(hash[:])
}

func initPassword(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	pwFile := filepath.Join(userDataRoot, "pw.json")
	if _, err := os.Stat(pwFile); err == nil {
		sendJSON(w, map[string]interface{}{"success": false, "message": "已设置过密码，请登录"})
		return
	}

	var req InitPasswordRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		sendJSON(w, map[string]interface{}{"success": false, "message": "Invalid request"})
		return
	}

	if len(req.Password) < 6 {
		sendJSON(w, map[string]interface{}{"success": false, "message": "密码长度至少6位"})
		return
	}

	pwData := PasswordData{
		PasswordHash: hashPassword(req.Password),
		CreatedAt:    time.Now(),
	}

	data, _ := json.Marshal(pwData)
	if err := os.WriteFile(pwFile, data, 0644); err != nil {
		sendJSON(w, map[string]interface{}{"success": false, "message": "保存密码失败"})
		return
	}

	token := generateToken()
	validTokens.Store(token, time.Now().Add(30*time.Minute))
	sendJSON(w, map[string]interface{}{"success": true, "token": token})
}

func login(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	pwFile := filepath.Join(userDataRoot, "pw.json")
	data, err := os.ReadFile(pwFile)
	if err != nil {
		sendJSON(w, map[string]interface{}{"success": false, "message": "未设置密码"})
		return
	}

	var pwData PasswordData
	if err := json.Unmarshal(data, &pwData); err != nil {
		sendJSON(w, map[string]interface{}{"success": false, "message": "密码数据损坏"})
		return
	}

	var req LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		sendJSON(w, map[string]interface{}{"success": false, "message": "Invalid request"})
		return
	}

	if hashPassword(req.Password) != pwData.PasswordHash {
		sendJSON(w, map[string]interface{}{"success": false, "message": "密码错误"})
		return
	}

	token := generateToken()
	validTokens.Store(token, time.Now().Add(30*time.Minute))
	sendJSON(w, map[string]interface{}{"success": true, "token": token})
}

func logout(w http.ResponseWriter, r *http.Request) {
	token := r.Header.Get("X-Auth-Token")
	if token != "" {
		validTokens.Delete(token)
	}
	sendJSON(w, map[string]interface{}{"success": true})
}

func generateToken() string {
	hash := sha256.Sum256([]byte(time.Now().String() + "iconstation-secret-key"))
	return hex.EncodeToString(hash[:])
}

func withAuth(handler http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		token := r.Header.Get("X-Auth-Token")
		if token == "" {
			token = r.URL.Query().Get("token")
		}

		if token == "" {
			http.Error(w, `{"error": "Unauthorized"}`, http.StatusUnauthorized)
			return
		}

		exp, ok := validTokens.Load(token)
		if !ok || time.Now().After(exp.(time.Time)) {
			validTokens.Delete(token)
			http.Error(w, `{"error": "Token无效或已过期"}`, http.StatusUnauthorized)
			return
		}

		handler(w, r)
	}
}

func listUploadIcons(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	files, err := os.ReadDir(userDataDir)
	if err != nil {
		sendJSON(w, map[string]interface{}{"files": []string{}, "error": err.Error()})
		return
	}

	fileList := []string{}
	for _, f := range files {
		if !f.IsDir() {
			ext := strings.ToLower(filepath.Ext(f.Name()))
			if ext == ".png" || ext == ".svg" || ext == ".jpg" || ext == ".jpeg" || ext == ".webp" {
				fileList = append(fileList, f.Name())
			}
		}
	}

	sendJSON(w, map[string]interface{}{"files": fileList})
}

type UploadInitRequest struct {
	FileName    string `json:"fileName"`
	FileSize    int64  `json:"fileSize"`
	TotalChunks int    `json:"totalChunks"`
	UploadID    string `json:"uploadId"`
}

func uploadInit(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req UploadInitRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		sendJSON(w, map[string]interface{}{"success": false, "message": "Invalid request"})
		return
	}

	uploadDir := filepath.Join(chunksDir, req.UploadID)
	if err := os.MkdirAll(uploadDir, 0755); err != nil {
		sendJSON(w, map[string]interface{}{"success": false, "message": "Failed to create upload directory"})
		return
	}

	sendJSON(w, map[string]interface{}{"success": true, "uploadId": req.UploadID})
}

func uploadChunk(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	r.ParseMultipartForm(32 << 20)

	uploadID := r.FormValue("uploadId")
	chunkIndex := r.FormValue("chunkIndex")

	if uploadID == "" || chunkIndex == "" {
		sendJSON(w, map[string]interface{}{"success": false, "message": "Missing parameters"})
		return
	}

	uploadDir := filepath.Join(chunksDir, uploadID)
	os.MkdirAll(uploadDir, 0755)

	file, _, err := r.FormFile("file")
	if err != nil {
		sendJSON(w, map[string]interface{}{"success": false, "message": "Failed to read chunk"})
		return
	}
	defer file.Close()

	chunkPath := filepath.Join(uploadDir, chunkIndex)
	outFile, err := os.Create(chunkPath)
	if err != nil {
		sendJSON(w, map[string]interface{}{"success": false, "message": "Failed to save chunk"})
		return
	}
	defer outFile.Close()

	io.Copy(outFile, file)
	sendJSON(w, map[string]interface{}{"success": true})
}

type UploadMergeRequest struct {
	UploadID string `json:"uploadId"`
	FileName string `json:"fileName"`
}

func uploadMerge(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req UploadMergeRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		sendJSON(w, map[string]interface{}{"success": false, "message": "Invalid request"})
		return
	}

	uploadDir := filepath.Join(chunksDir, req.UploadID)
	dstPath := filepath.Join(userDataDir, req.FileName)

	outFile, err := os.Create(dstPath)
	if err != nil {
		sendJSON(w, map[string]interface{}{"success": false, "message": "Failed to create file"})
		return
	}
	defer outFile.Close()

	files, err := os.ReadDir(uploadDir)
	if err != nil {
		sendJSON(w, map[string]interface{}{"success": false, "message": "Failed to read chunks"})
		return
	}

	for i := 0; i < len(files); i++ {
		chunkPath := filepath.Join(uploadDir, fmt.Sprintf("%d", i))
		chunkData, err := os.ReadFile(chunkPath)
		if err != nil {
			os.Remove(dstPath)
			sendJSON(w, map[string]interface{}{"success": false, "message": "Failed to read chunk"})
			return
		}
		outFile.Write(chunkData)
	}

	os.RemoveAll(uploadDir)
	sendJSON(w, map[string]interface{}{"success": true, "path": dstPath})
}

type RenameRequest struct {
	OldName string `json:"oldName"`
	NewName string `json:"newName"`
}

func renameIcon(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req RenameRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		sendJSON(w, map[string]interface{}{"success": false, "message": "Invalid request"})
		return
	}

	oldPath := filepath.Join(userDataDir, req.OldName)
	newPath := filepath.Join(userDataDir, req.NewName)

	if _, err := os.Stat(oldPath); os.IsNotExist(err) {
		sendJSON(w, map[string]interface{}{"success": false, "message": "File not found"})
		return
	}

	if err := os.Rename(oldPath, newPath); err != nil {
		sendJSON(w, map[string]interface{}{"success": false, "message": "Failed to rename"})
		return
	}

	sendJSON(w, map[string]interface{}{"success": true})
}

type DeleteRequest struct {
	FileName string `json:"fileName"`
}

func deleteIcon(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req DeleteRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		sendJSON(w, map[string]interface{}{"success": false, "message": "Invalid request"})
		return
	}

	filePath := filepath.Join(userDataDir, req.FileName)

	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		sendJSON(w, map[string]interface{}{"success": false, "message": "File not found"})
		return
	}

	if err := os.Remove(filePath); err != nil {
		sendJSON(w, map[string]interface{}{"success": false, "message": "Failed to delete"})
		return
	}

	sendJSON(w, map[string]interface{}{"success": true})
}

func sendJSON(w http.ResponseWriter, data map[string]interface{}) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(data)
}