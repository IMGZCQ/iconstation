package main

import (
	"crypto/sha256"
	"embed"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"io/fs"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"
)

//go:embed static
var staticFS embed.FS

var (
	userDataDir  string
	chunksDir    string
	userDataRoot string
	validTokens  sync.Map
)

func init() {
	exePath, err := os.Executable()
	if err != nil {
		log.Fatal(err)
	}
	exeDir := filepath.Dir(exePath)
	userDataRoot = filepath.Join(exeDir, "UserData")
	userDataDir = filepath.Join(userDataRoot, "icons")
	chunksDir = filepath.Join(userDataRoot, "chunks")
}

type PasswordData struct {
	PasswordHash string    `json:"passwordHash"`
	CreatedAt    time.Time `json:"createdAt"`
	Open         *bool     `json:"open,omitempty"`
}

type LoginRequest struct {
	Password string `json:"password"`
	Open     bool   `json:"open"`
}

type InitPasswordRequest struct {
	Password string `json:"password"`
	Open     bool   `json:"open"`
}

func main() {
	_ = os.MkdirAll(userDataDir, 0755)
	_ = os.MkdirAll(chunksDir, 0755)
	_ = os.MkdirAll(userDataRoot, 0755)

	staticSubFS, err := fs.Sub(staticFS, "static")
	if err != nil {
		log.Fatal(err)
	}

	// 定时清理过期 token（10分钟一次）
	go func() {
		ticker := time.NewTicker(10 * time.Minute)
		defer ticker.Stop()
		for range ticker.C {
			now := time.Now()
			validTokens.Range(func(key, value interface{}) bool {
				if exp, ok := value.(time.Time); ok && now.After(exp) {
					validTokens.Delete(key)
				}
				return true
			})
		}
	}()

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Cache-Control", "max-age=86400")
		if r.URL.Path == "/" {
			http.ServeFileFS(w, r, staticSubFS, "index.html")
			return
		}
		http.FileServer(http.FS(staticSubFS)).ServeHTTP(w, r)
	})

	http.HandleFunc("/api/check-init", checkInit)
	http.HandleFunc("/api/init-password", initPassword)
	http.HandleFunc("/api/login", login)
	http.HandleFunc("/api/logout", logout)
	http.HandleFunc("/api/get-open", getOpen)
	http.HandleFunc("/api/list-upload-icons", listUploadIcons)
	http.HandleFunc("/api/create-category", withAuth(createCategory))
	http.HandleFunc("/api/list-categories", listCategories)
	http.HandleFunc("/api/upload/user_icon/init", withAuth(uploadInit))
	http.HandleFunc("/api/upload/user_icon/chunk", withAuth(uploadChunk))
	http.HandleFunc("/api/upload/user_icon/merge", withAuth(uploadMerge))
	http.HandleFunc("/api/rename/user_icon", withAuth(renameIcon))
	http.HandleFunc("/api/delete/user_icon", withAuth(deleteIcon))

	http.HandleFunc("/deskdata/user_icon/", func(w http.ResponseWriter, r *http.Request) {
		if !isUserIconAccessible(r) {
			http.Error(w, `{"error":"Forbidden"}`, http.StatusForbidden)
			return
		}

		path := strings.TrimPrefix(r.URL.Path, "/deskdata/user_icon/")
		path = filepath.Clean(path)
		target := filepath.Join(userDataDir, path)

		// 安全校验：禁止路径穿越
		if !strings.HasPrefix(target, userDataDir+string(filepath.Separator)) && target != userDataDir {
			http.Error(w, `{"error":"Forbidden"}`, http.StatusForbidden)
			return
		}

		http.ServeFile(w, r, target)
	})

	log.Println("Server started on :9168")
	log.Fatal(http.ListenAndServe(":9168", nil))
}

func checkInit(w http.ResponseWriter, r *http.Request) {
	pwFile := filepath.Join(userDataRoot, "pw.json")
	_, err := os.Stat(pwFile)
	sendJSON(w, map[string]interface{}{"initialized": err == nil})
}

func getOpen(w http.ResponseWriter, r *http.Request) {
	sendJSON(w, map[string]interface{}{"open": getOpenStatus()})
}

func isUserIconAccessible(r *http.Request) bool {
	token := r.Header.Get("X-Auth-Token")
	if token == "" {
		token = r.URL.Query().Get("token")
	}

	if token != "" {
		exp, ok := validTokens.Load(token)
		if ok && time.Now().Before(exp.(time.Time)) {
			return true
		}
	}
	return getOpenStatus()
}

func getOpenStatus() bool {
	pwFile := filepath.Join(userDataRoot, "pw.json")
	data, err := os.ReadFile(pwFile)
	if err != nil {
		return false
	}
	var pwData PasswordData
	if err := json.Unmarshal(data, &pwData); err != nil {
		return false
	}
	return pwData.Open != nil && *pwData.Open
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
		sendJSON(w, map[string]interface{}{"success": false, "message": "已设置密码"})
		return
	}
	var req InitPasswordRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		sendJSON(w, map[string]interface{}{"success": false, "message": "参数错误"})
		return
	}
	if len(req.Password) < 6 {
		sendJSON(w, map[string]interface{}{"success": false, "message": "密码至少6位"})
		return
	}
	pwData := PasswordData{
		PasswordHash: hashPassword(req.Password),
		CreatedAt:    time.Now(),
		Open:         &req.Open,
	}
	data, _ := json.Marshal(pwData)
	_ = os.WriteFile(pwFile, data, 0644)

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
		sendJSON(w, map[string]interface{}{"success": false, "message": "未初始化"})
		return
	}
	var pwData PasswordData
	_ = json.Unmarshal(data, &pwData)

	var req LoginRequest
	_ = json.NewDecoder(r.Body).Decode(&req)

	if hashPassword(req.Password) != pwData.PasswordHash {
		sendJSON(w, map[string]interface{}{"success": false, "message": "密码错误"})
		return
	}

	pwData.Open = &req.Open
	newData, _ := json.Marshal(pwData)
	_ = os.WriteFile(pwFile, newData, 0644)

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
	hash := sha256.Sum256([]byte(time.Now().String() + "iconstation-secure-v2"))
	return hex.EncodeToString(hash[:])
}

func withAuth(handler http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		token := r.Header.Get("X-Auth-Token")
		if token == "" {
			token = r.URL.Query().Get("token")
		}
		if token == "" {
			http.Error(w, `{"error":"Unauthorized"}`, http.StatusUnauthorized)
			return
		}
		exp, ok := validTokens.Load(token)
		if !ok || time.Now().After(exp.(time.Time)) {
			validTokens.Delete(token)
			http.Error(w, `{"error":"Token无效"}`, http.StatusUnauthorized)
			return
		}
		handler(w, r)
	}
}

type IconInfo struct {
	Name     string `json:"name"`
	Category string `json:"category"`
}

func listUploadIcons(w http.ResponseWriter, r *http.Request) {
	if !isUserIconAccessible(r) {
		sendJSON(w, map[string]interface{}{"files": []IconInfo{}})
		return
	}
	var list []IconInfo
	_ = filepath.Walk(userDataDir, func(p string, info fs.FileInfo, err error) error {
		if err != nil || info.IsDir() {
			return nil
		}
		rel, _ := filepath.Rel(userDataDir, p)
		ext := strings.ToLower(filepath.Ext(rel))
		allow := map[string]bool{".png": true, ".svg": true, ".jpg": true, ".jpeg": true, ".webp": true}
		if allow[ext] {
			cat := ""
			dir := filepath.Dir(rel)
			if dir != "." {
				cat = dir
			}
			list = append(list, IconInfo{Name: rel, Category: cat})
		}
		return nil
	})
	sendJSON(w, map[string]interface{}{"files": list})
}

type UploadInitRequest struct {
	FileName    string `json:"fileName"`
	FileSize    int64  `json:"fileSize"`
	TotalChunks int    `json:"totalChunks"`
	UploadID    string `json:"uploadId"`
}

func uploadInit(w http.ResponseWriter, r *http.Request) {
	var req UploadInitRequest
	_ = json.NewDecoder(r.Body).Decode(&req)
	_ = os.MkdirAll(filepath.Join(chunksDir, req.UploadID), 0755)
	sendJSON(w, map[string]interface{}{"success": true, "uploadId": req.UploadID})
}

func uploadChunk(w http.ResponseWriter, r *http.Request) {
	_ = r.ParseMultipartForm(32 << 20)
	uploadID := r.FormValue("uploadId")
	chunkIdx := r.FormValue("chunkIndex")
	dir := filepath.Join(chunksDir, uploadID)
	_ = os.MkdirAll(dir, 0755)
	file, _, _ := r.FormFile("file")
	defer file.Close()
	dst, _ := os.Create(filepath.Join(dir, chunkIdx))
	defer dst.Close()
	_, _ = io.Copy(dst, file)
	sendJSON(w, map[string]interface{}{"success": true})
}

type UploadMergeRequest struct {
	UploadID string `json:"uploadId"`
	FileName string `json:"fileName"`
}

func getUniqueFileName(fileName, dir string) (string, bool) {
	ext := filepath.Ext(fileName)
	base := strings.TrimSuffix(fileName, ext)
	for i := 0; i < 1000; i++ {
		var candidate string
		if i == 0 {
			candidate = fileName
		} else {
			candidate = fmt.Sprintf("%s_%d%s", base, i, ext)
		}
		if _, err := os.Stat(filepath.Join(dir, candidate)); os.IsNotExist(err) {
			return candidate, i > 0
		}
	}
	return fileName, false
}

func safeJoinPath(root, sub string) (string, error) {
	clean := filepath.Clean(sub)
	path := filepath.Join(root, clean)
	if !strings.HasPrefix(path, root+string(filepath.Separator)) && path != root {
		return "", fmt.Errorf("invalid path")
	}
	return path, nil
}

func uploadMerge(w http.ResponseWriter, r *http.Request) {
	var req UploadMergeRequest
	_ = json.NewDecoder(r.Body).Decode(&req)

	// 安全路径校验
	dstRel := filepath.Clean(req.FileName)
	dst, err := safeJoinPath(userDataDir, dstRel)
	if err != nil {
		sendJSON(w, map[string]interface{}{"success": false, "message": "非法路径"})
		return
	}

	uploadDir := filepath.Join(chunksDir, req.UploadID)
	finalName, renamed := getUniqueFileName(filepath.Base(dst), filepath.Dir(dst))
	finalPath := filepath.Join(filepath.Dir(dst), finalName)

	_ = os.MkdirAll(filepath.Dir(finalPath), 0755)
	out, _ := os.Create(finalPath)
	defer out.Close()

	entries, _ := os.ReadDir(uploadDir)
	for i := 0; i < len(entries); i++ {
		chunk, _ := os.ReadFile(filepath.Join(uploadDir, fmt.Sprintf("%d", i)))
		_, _ = out.Write(chunk)
	}
	_ = os.RemoveAll(uploadDir)

	msg := ""
	if renamed {
		msg = "已自动重命名"
	}
	sendJSON(w, map[string]interface{}{"success": true, "path": finalName, "renamed": renamed, "message": msg})
}

type RenameRequest struct {
	OldName string `json:"oldName"`
	NewName string `json:"newName"`
}

func renameIcon(w http.ResponseWriter, r *http.Request) {
	var req RenameRequest
	_ = json.NewDecoder(r.Body).Decode(&req)

	old, err1 := safeJoinPath(userDataDir, req.OldName)
	new, err2 := safeJoinPath(userDataDir, req.NewName)
	if err1 != nil || err2 != nil {
		sendJSON(w, map[string]interface{}{"success": false, "message": "非法路径"})
		return
	}

	if _, err := os.Stat(old); os.IsNotExist(err) {
		sendJSON(w, map[string]interface{}{"success": false, "message": "文件不存在"})
		return
	}

	newDir := filepath.Dir(new)
	newBase := filepath.Base(new)
	finalName, renamed := getUniqueFileName(newBase, newDir)
	finalPath := filepath.Join(newDir, finalName)

	_ = os.Rename(old, finalPath)
	sendJSON(w, map[string]interface{}{"success": true, "renamed": renamed})
}

type DeleteRequest struct {
	FileName string `json:"fileName"`
}

func deleteIcon(w http.ResponseWriter, r *http.Request) {
	var req DeleteRequest
	_ = json.NewDecoder(r.Body).Decode(&req)

	target, err := safeJoinPath(userDataDir, req.FileName)
	if err != nil {
		sendJSON(w, map[string]interface{}{"success": false, "message": "非法路径"})
		return
	}

	if _, err := os.Stat(target); os.IsNotExist(err) {
		sendJSON(w, map[string]interface{}{"success": false, "message": "不存在"})
		return
	}
	_ = os.Remove(target)
	sendJSON(w, map[string]interface{}{"success": true})
}

type CreateCategoryRequest struct {
	Name string `json:"name"`
}

func listCategories(w http.ResponseWriter, r *http.Request) {
	if !isUserIconAccessible(r) {
		sendJSON(w, map[string]interface{}{"categories": []string{}})
		return
	}
	entries, _ := os.ReadDir(userDataDir)
	var cats []string
	for _, e := range entries {
		if e.IsDir() {
			cats = append(cats, e.Name())
		}
	}
	sendJSON(w, map[string]interface{}{"categories": cats})
}

func createCategory(w http.ResponseWriter, r *http.Request) {
	var req CreateCategoryRequest
	_ = json.NewDecoder(r.Body).Decode(&req)
	name := strings.TrimSpace(req.Name)
	if name == "" || strings.ContainsAny(name, `/\..`) {
		sendJSON(w, map[string]interface{}{"success": false, "message": "名称非法"})
		return
	}
	dir := filepath.Join(userDataDir, name)
	_ = os.MkdirAll(dir, 0755)
	sendJSON(w, map[string]interface{}{"success": true, "name": name})
}

func sendJSON(w http.ResponseWriter, data map[string]interface{}) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	_ = json.NewEncoder(w).Encode(data)
}
