package control

import (
	"encoding/json"
	"fmt"
	"html/template"
	"io"
	"log"
	"net/http"
	"net/url"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"csz.net/tgstate/assets"
	"csz.net/tgstate/conf"
	"csz.net/tgstate/utils"
)

// UploadAPI 上传图片api
func UploadAPI(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")

	if r.Method == http.MethodPost {
		// 获取上传的文件
		file, header, err := r.FormFile("file")

		if err != nil {
			errJsonMsg("Unable to get file", w)
			// http.Error(w, "Unable to get file", http.StatusBadRequest)
			return
		}
		defer file.Close()
		if conf.Mode != "p" && r.ContentLength > 20*1024*1024 {
			// 检查文件大小
			errJsonMsg("File size exceeds 20MB limit", w)
			return
		}
		// 检查文件类型
		allowedExts := []string{".jpg", ".jpeg", ".png"}

		// 如果设置了AllowedExts，则使用设置的文件类型
		if len(conf.AllowedExts) > 0 {
			allowedExts = append(allowedExts, strings.Split(conf.AllowedExts, ",")...)
		}

		var fileName = header.Filename
		ext := filepath.Ext(fileName)
		valid := false
		for _, allowedExt := range allowedExts {
			if ext == allowedExt {
				valid = true
				break
			}
		}
		if conf.Mode != "p" && !valid {
			errJsonMsg(fmt.Sprintf("Invalid file type. Only .jpg, .jpeg, and .png %s are allowed.", conf.AllowedExts), w)
			// http.Error(w, "Invalid file type. Only .jpg, .jpeg, and .png are allowed.", http.StatusBadRequest)
			return
		}
		res := conf.UploadResponse{
			Code:    0,
			Message: "error",
		}
		fileId := utils.UpDocument(utils.TgFileData(fileName, file))
		if "blob" != fileName {
			ip := r.RemoteAddr // 获取上传者IP
			// 插入数据到数据库
			err := SaveFileRecord(fileId, fileName, ip)
			if err != nil {
				errJsonMsg("Unable to save file record", w)
			}
		}

		downloadUrl := conf.FileRoute + fileId
		if downloadUrl != conf.FileRoute {
			imageUrl := strings.TrimSuffix(conf.BaseUrl, "/") + downloadUrl
			// url encode imageUrl
			proxyUrl := conf.ProxyUrl + "/" + url.QueryEscape(imageUrl)
			res = conf.UploadResponse{
				Code:     1,
				Message:  downloadUrl,
				ImgUrl:   imageUrl,
				ProxyUrl: proxyUrl,
			}
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(res)
		return
	}

	// 如果不是POST请求，返回错误响应
	http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
}

// D 下载文件
func D(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Path
	fileId := strings.TrimPrefix(path, conf.FileRoute)
	if fileId == "" {
		// 设置响应的状态码为 404
		w.WriteHeader(http.StatusNotFound)
		// 写入响应内容
		errJsonMsg("404 Not Found", w)
		return
	}
	record, err := GetFileNameByIDOrName(fileId)
	if err == nil && record.FileId != "" {
		fileId = record.FileId
	}
	// 发起HTTP GET请求来获取Telegram图片
	fileUrl, _ := utils.GetDownloadUrl(fileId)
	resp, err := http.Get(fileUrl)
	if err != nil {
		http.Error(w, "Failed to fetch content", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Disposition", "inline") // 设置为 "inline" 以支持在线播放
	// 检查Content-Type是否为图片类型
	if !strings.HasPrefix(resp.Header.Get("Content-Type"), "application/octet-stream") {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("404 Not Found"))
		return
	}
	contentLength, err := strconv.Atoi(resp.Header.Get("Content-Length"))
	if err != nil {
		log.Println("获取Content-Length出错:", err)
		return
	}
	buffer := make([]byte, contentLength)
	n, err := resp.Body.Read(buffer)
	defer resp.Body.Close()
	if err != nil && err != io.ErrUnexpectedEOF {
		log.Println("读取响应主体数据时发生错误:", err)
		return
	}
	// 输出文件内容到控制台
	if string(buffer[:12]) == "tgstate-blob" {
		content := string(buffer)
		lines := strings.Split(content, "\n")
		log.Println("分块文件:" + lines[1])
		var fileSize string
		var startLine = 2
		if strings.HasPrefix(lines[2], "size") {
			fileSize = lines[2][len("size"):]
			startLine = startLine + 1
		}
		w.Header().Set("Content-Type", "application/octet-stream")
		w.Header().Set("Content-Disposition", "attachment; filename=\""+lines[1]+"\"")
		w.Header().Set("Content-Length", fileSize)
		for i := startLine; i < len(lines); i++ {
			fileStatus := false
			var fileUrl string
			var reTry = 0
			for !fileStatus {
				if reTry > 0 {
					time.Sleep(5 * time.Second)
				}
				reTry = reTry + 1
				fileUrl, fileStatus = utils.GetDownloadUrl(strings.ReplaceAll(lines[i], " ", ""))
			}
			blobResp, err := http.Get(fileUrl)
			if err != nil {
				http.Error(w, "Failed to fetch content", http.StatusInternalServerError)
				return
			}
			_, err = io.Copy(w, blobResp.Body)
			blobResp.Body.Close()
			if err != nil {
				log.Println("写入响应主体数据时发生错误:", err)
				return
			}
		}
	} else {
		// 使用DetectContentType函数检测文件类型
		w.Header().Set("Content-Type", http.DetectContentType(buffer))

		if err == nil && record.Filename != "" {
			w.Header().Set("Content-Disposition", "attachment; filename=\""+record.Filename+"\"")
		}

		_, err = w.Write(buffer[:n])
		if err != nil {
			http.Error(w, "Failed to write content", http.StatusInternalServerError)
			log.Println(http.StatusInternalServerError)
			return
		}
		_, err = io.Copy(w, resp.Body)
		resp.Body.Close()
		if err != nil {
			log.Println(http.StatusInternalServerError)
			return
		}
	}
}

// Index 首页
func Index(w http.ResponseWriter, r *http.Request) {
	htmlPath := "templates/images.tmpl"
	if conf.Mode == "p" {
		htmlPath = "templates/files.tmpl"
	}
	file, err := assets.Templates.ReadFile(htmlPath)
	if err != nil {
		http.Error(w, "HTML file not found", http.StatusNotFound)
		return
	}
	// 读取头部模板
	headerFile, err := assets.Templates.ReadFile("templates/header.tmpl")
	if err != nil {
		http.Error(w, "Header template not found", http.StatusNotFound)
		return
	}

	// 读取页脚模板
	footerFile, err := assets.Templates.ReadFile("templates/footer.tmpl")
	if err != nil {
		http.Error(w, "Footer template not found", http.StatusNotFound)
		return
	}

	// 创建HTML模板并包括头部
	tmpl := template.New("html")
	tmpl, err = tmpl.Parse(string(headerFile))
	if err != nil {
		http.Error(w, "Error parsing header template", http.StatusInternalServerError)
		return
	}

	// 包括主HTML内容
	tmpl, err = tmpl.Parse(string(file))
	if err != nil {
		http.Error(w, "Error parsing HTML template", http.StatusInternalServerError)
		return
	}

	// 包括页脚
	tmpl, err = tmpl.Parse(string(footerFile))
	if err != nil {
		http.Error(w, "Error parsing footer template", http.StatusInternalServerError)
		return
	}

	// 直接将HTML内容发送给客户端
	w.Header().Set("Content-Type", "text/html")
	err = tmpl.Execute(w, nil)
	if err != nil {
		http.Error(w, "Error rendering HTML template", http.StatusInternalServerError)
	}
}

func Pwd(w http.ResponseWriter, r *http.Request) {
	// 输出 HTML 表单
	if r.Method != http.MethodPost {
		file, err := assets.Templates.ReadFile("templates/pwd.tmpl")
		if err != nil {
			http.Error(w, "HTML file not found", http.StatusNotFound)
			return
		}
		// 读取头部模板
		headerFile, err := assets.Templates.ReadFile("templates/header.tmpl")
		if err != nil {
			http.Error(w, "Header template not found", http.StatusNotFound)
			return
		}

		// 创建HTML模板并包括头部
		tmpl := template.New("html")
		if tmpl, err = tmpl.Parse(string(headerFile)); err != nil {
			http.Error(w, "Error parsing Header template", http.StatusInternalServerError)
			return
		}

		// 包括主HTML内容
		if tmpl, err = tmpl.Parse(string(file)); err != nil {
			http.Error(w, "Error parsing File template", http.StatusInternalServerError)
			return
		}

		// 直接将HTML内容发送给客户端
		w.Header().Set("Content-Type", "text/html")
		if err := tmpl.Execute(w, nil); err != nil {
			http.Error(w, "Error rendering HTML template", http.StatusInternalServerError)
		}
		return
	}
	// 设置cookie
	cookie := http.Cookie{
		Name:  "p",
		Value: r.FormValue("p"),
	}
	http.SetCookie(w, &cookie)
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func FilesAPI(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	password := r.URL.Query().Get("password")
	response := conf.ResponseResult{
		Code:    0,
		Message: "ok",
	}

	if conf.ApiPass != "" && password != conf.ApiPass {
		response.Message = "Unauthorized"
		response.Code = http.StatusUnauthorized
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(response)
		return
	}

	record, err := SelectAllRecord()
	response.Data = record
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

func errJsonMsg(msg string, w http.ResponseWriter) {
	// 这里示例直接返回JSON响应
	response := conf.UploadResponse{
		Code:    0,
		Message: msg,
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

func Middleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// 只有当密码设置并且不为"none"时，才进行检查
		if conf.Pass != "" && conf.Pass != "none" {
			if strings.HasPrefix(r.URL.Path, "/api") && r.URL.Query().Get("pass") == conf.Pass {
				return
			}
			if cookie, err := r.Cookie("p"); err != nil || cookie.Value != conf.Pass {
				http.Redirect(w, r, "/pwd", http.StatusSeeOther)
				return
			}
		}
		next(w, r)
	}
}
