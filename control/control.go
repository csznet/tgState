package control

import (
	"encoding/json"
	"html/template"
	"io"
	"log"
	"net/http"
	"path/filepath"
	"regexp"
	"strings"

	"csz.net/tgstate/assets"
	"csz.net/tgstate/conf"
	"csz.net/tgstate/utils"
)

// UploadImageAPI 上传文件api
func UploadImageAPI(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	if r.Method == http.MethodPost {
		// 解析上传的文件
		err := r.ParseMultipartForm(5 * 1024 * 1024) // 限制上传文件大小为 5MB
		if err != nil {
			errJsonMsg("Unable to parse form", w)
			// http.Error(w, "Unable to parse form", http.StatusBadRequest)
			return
		}

		// 获取上传的文件
		file, header, err := r.FormFile("image")
		if err != nil {
			errJsonMsg("Unable to get file", w)
			// http.Error(w, "Unable to get file", http.StatusBadRequest)
			return
		}
		defer file.Close()
		if conf.Mode != "pan" {
			// 检查文件大小
			fileSize := r.ContentLength
			if fileSize > 20*1024*1024 {
				errJsonMsg("File size exceeds 20MB limit", w)
				return
			}
		}
		// 检查文件类型
		allowedExts := []string{".jpg", ".jpeg", ".png"}
		ext := filepath.Ext(header.Filename)
		valid := false
		for _, allowedExt := range allowedExts {
			if ext == allowedExt {
				valid = true
				break
			}
		}
		if conf.Mode != "pan" {
			if !valid {
				errJsonMsg("Invalid file type. Only .jpg, .jpeg, and .png are allowed.", w)
				// http.Error(w, "Invalid file type. Only .jpg, .jpeg, and .png are allowed.", http.StatusBadRequest)
				return
			}
		}
		res := conf.UploadResponse{
			Code:    0,
			Message: "error",
		}
		var img string
		img = "/d/" + utils.UpDocument(utils.TgFileData(header.Filename, file))
		res = conf.UploadResponse{
			Code:    1,
			Message: img,
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(res)
		return
	}

	// 如果不是POST请求，返回错误响应
	http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
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
func D(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Path
	id := strings.TrimPrefix(path, "/d/")
	if id == "" {
		// 设置响应的状态码为 404
		w.WriteHeader(http.StatusNotFound)
		// 写入响应内容
		w.Write([]byte("404 Not Found"))
		return
	}

	// 发起HTTP GET请求来获取Telegram图片
	resp, err := http.Get(utils.GetDownloadUrl(id))
	if err != nil {
		http.Error(w, "Failed to fetch content", http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()
	rType := resp.Header.Get("Content-Type")
	w.Header().Set("Content-Disposition", "inline") // 设置为 "inline" 以支持在线播放
	// 检查Content-Type是否为图片类型
	if !strings.HasPrefix(resp.Header.Get("Content-Type"), "application/octet-stream") {
		// 设置响应的状态码为 404
		w.WriteHeader(http.StatusNotFound)
		// 写入响应内容
		w.Write([]byte("404 Not Found"))
		return
	}
	// 读取前512个字节以用于文件类型检测
	buffer := make([]byte, 512)
	n, err := resp.Body.Read(buffer)
	if err != nil {
		log.Println("读取响应主体数据时发生错误:", err)
		return
	}
	// 输出文件内容到控制台
	if string(buffer[:12]) == "tgstate-blob" {
		content := string(buffer)
		lines := strings.Fields(content)
		log.Println("这是一个分块文件,文件名:" + lines[1])
		w.Header().Set("Content-Type", "application/octet-stream")
		w.Header().Set("Content-Disposition", "attachment; filename=\""+lines[1]+"\"")
		for i := 2; i < len(lines); i++ {
			blobResp, err := http.Get(utils.GetDownloadUrl(regexp.MustCompile("[^a-zA-Z0-9_-]").ReplaceAllString(lines[i], "")))
			if err != nil {
				http.Error(w, "Failed to fetch content", http.StatusInternalServerError)
				return
			}

			// 将文件名设置到Content-Disposition标头
			blobResp.Header.Set("Content-Disposition", "attachment; filename=\""+lines[1]+"\"")

			defer blobResp.Body.Close()

			_, err = io.Copy(w, blobResp.Body)
			if err != nil {
				log.Println("写入响应主体数据时发生错误:", err)
				return
			}
		}

	} else {
		// 使用DetectContentType函数检测文件类型
		rType = http.DetectContentType(buffer)
		w.Header().Set("Content-Type", rType)
		// 写入前512个字节到响应w
		_, err = w.Write(buffer[:n])
		if err != nil {
			http.Error(w, "Failed to write content", http.StatusInternalServerError)
			log.Println(http.StatusInternalServerError)
			return
		}
		_, err = io.Copy(w, resp.Body)
		if err != nil {
			//http.Error(w, "Failed to show content", http.StatusInternalServerError)
			log.Println(http.StatusInternalServerError)
			return
		}
	}
}

// Index 首页
func Index(w http.ResponseWriter, r *http.Request) {
	htmlPath := "templates/images.tmpl"
	if conf.Mode == "pan" {
		htmlPath = "templates/files.tmpl"
	}
	file, err := assets.Templates.ReadFile(htmlPath)
	if err != nil {
		http.Error(w, "HTML file not found", http.StatusNotFound)
		return
	}
	// 读取头部模板
	headerPath := "templates/header.tmpl"
	headerFile, err := assets.Templates.ReadFile(headerPath)
	if err != nil {
		http.Error(w, "Header template not found", http.StatusNotFound)
		return
	}

	// 读取页脚模板
	footerPath := "templates/footer.tmpl"
	footerFile, err := assets.Templates.ReadFile(footerPath)
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
		headerPath := "templates/header.tmpl"
		headerFile, err := assets.Templates.ReadFile(headerPath)
		if err != nil {
			http.Error(w, "Header template not found", http.StatusNotFound)
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

		// 直接将HTML内容发送给客户端
		w.Header().Set("Content-Type", "text/html")
		err = tmpl.Execute(w, nil)
		if err != nil {
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

func Middleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if conf.Pass != "" && conf.Pass != "none" {
			// 在这里检查cookie
			cookie, err := r.Cookie("p")
			if err != nil || cookie.Value != conf.Pass {
				// 如果cookie不存在或值不为110，则重定向到/pwd
				http.Redirect(w, r, "/pwd", http.StatusSeeOther)
				return
			}
		}
		// 如果cookie值为110，调用下一个处理程序
		next(w, r)
	}
}
