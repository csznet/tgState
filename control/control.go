package control

import (
	"encoding/json"
	"io"
	"net/http"
	"path/filepath"
	"strings"

	"csz.net/tgstate/conf"
	"csz.net/tgstate/utils"
)

// 上传文件api
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

		// 检查文件大小
		fileSize := r.ContentLength
		if fileSize > 20*1024*1024 {
			errJsonMsg("File size exceeds 20MB limit", w)
			return
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

		if !valid {
			errJsonMsg("Invalid file type. Only .jpg, .jpeg, and .png are allowed.", w)
			// http.Error(w, "Invalid file type. Only .jpg, .jpeg, and .png are allowed.", http.StatusBadRequest)
			return
		}

		// 读取文件内容
		fileBytes, err := io.ReadAll(file)
		if err != nil {
			errJsonMsg("Failed to read file", w)
			// http.Error(w, "Failed to read file", http.StatusInternalServerError)
			return
		}
		res := conf.UploadResponse{
			Code:    0,
			Message: "error",
		}
		var img string
		img = "/d/" + utils.UpDocument(utils.TgFileData(header.Filename, fileBytes))
		// if conf.ImgOrigin || fileSize > 5*1024*1024 {
		// 	img = "/d/" + utils.UpDocument(utils.TgFileData(header.Filename, fileBytes))
		// } else {
		// 	img = "/img/" + utils.SendImageToTelegram(utils.TgFileData(header.Filename, fileBytes))
		// }
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
	// 发起HTTP GET请求来获取Telegram图片
	resp, err := http.Get("https://api.telegram.org/file/bot" + conf.BotToken + "/documents/file_" + id)
	if err != nil {
		http.Error(w, "Failed to fetch image", http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	// 检查Content-Type是否为图片类型
	contentType := resp.Header.Get("Content-Type")
	if !strings.HasPrefix(contentType, "application/octet-stream") {
		// 设置响应的状态码为 404
		w.WriteHeader(http.StatusNotFound)
		// 写入响应内容
		w.Write([]byte("404 Not Found"))
		return
	}
	lastDotIndex := strings.LastIndex(id, ".")
	// 检查是否找到点
	if lastDotIndex != -1 {
		// 从点的位置截取字符串的子串，即文件扩展名
		extension := id[lastDotIndex+1:]
		w.Header().Set("Content-Type", "image/"+extension)
	} else {
		http.Error(w, "Failed to show image", http.StatusInternalServerError)
		return
	}

	// 将图片内容写入响应正文
	_, err = io.Copy(w, resp.Body)
	if err != nil {
		http.Error(w, "Failed to show image", http.StatusInternalServerError)
		return
	}
}

func Img(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Path
	id := strings.TrimPrefix(path, "/img/")
	// 发起HTTP GET请求来获取Telegram图片
	resp, err := http.Get("https://api.telegram.org/file/bot" + conf.BotToken + "/photos/file_" + id)
	if err != nil {
		http.Error(w, "Failed to fetch image", http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	// 检查Content-Type是否为图片类型
	contentType := resp.Header.Get("Content-Type")
	if !strings.HasPrefix(contentType, "application/octet-stream") {
		// 设置响应的状态码为 404
		w.WriteHeader(http.StatusNotFound)
		// 写入响应内容
		w.Write([]byte("404 Not Found"))
		return
	}

	// 设置响应头
	w.Header().Set("Content-Type", "image/jpeg")

	// 将图片内容写入响应正文
	_, err = io.Copy(w, resp.Body)
	if err != nil {
		http.Error(w, "Failed to show image", http.StatusInternalServerError)
		return
	}
}

// 首页
func Index(w http.ResponseWriter, r *http.Request) {
	// 如果不是 POST 请求，显示上传图片的 HTML 表单
	htmlForm := `<!DOCTYPE html>
	<html>
	<head>
		<meta charset="UTF-8">
		<title>Telegram图床</title>
		<meta name="keywords" content="telegram图床,tg图床,免费图床,永久图床,图片外链,免费图片外链,纸飞机图床,电报图床" />
		<meta name="description" content="telegram图床,tg图床,免费图床,永久图床,图片外链,免费图片外链,纸飞机图床,电报图床" />
		<meta name="viewport" content="width=device-width, initial-scale=1, maximum-scale=1">
		<style>
			#uploadButton,
			#uploadFileLabel {
				display: block;
				max-width: 200px; /* 固定宽度 */
				margin: 0 auto; /* 居中 */
				margin-bottom: 10px; /* 添加底部间距 */
			}
	
			body {
				font-family: Arial, sans-serif;
				text-align: center;
			}
	
			h1 {
				color: #333;
			}
	
			.custom-file-input {
				display: none;
			}
	
			.custom-file-label {
				background-color: #007BFF;
				color: #fff;
				padding: 10px 20px;
				cursor: pointer;
			}
	
			.custom-file-label:hover {
				background-color: #0056b3;
			}
	
			#uploadButton {
				background-color: #007BFF;
				color: #fff;
				padding: 10px 20px;
				border: none;
				cursor: pointer;
			}
	
			#uploadButton[disabled] {
				background-color: #ccc;
				cursor: not-allowed;
			}
	
			#uploadButton:hover {
				background-color: #0056b3;
			}
	
			#response {
				margin-top: 20px;
				padding: 10px;
			}
	
			.response-item {
				margin-bottom: 10px;
				padding: 10px;
				border-radius: 5px;
			}
	
			.response-success {
				background-color: #d4edda;
				border-color: #c3e6cb;
				color: #155724;
			}
	
			.response-error {
				background-color: #f8d7da;
				border-color: #f5c6cb;
				color: #721c24;
			}
	
			#loading {
				display: none;
			}
			.copy-code{
				margin: 5px;
			}
			.copy-links{
				margin-top: 5px;
			}
			#uploadButton[disabled]:hover {
		background-color: #ccc;
		cursor: not-allowed;
	}
	
		</style>
		<script src="https://code.jquery.com/jquery-3.6.0.min.js"></script>
	</head>
	<body>
		<h1>上传图片到 Telegram</h1>
		<label for="uploadFile" id="uploadFileLabel" class="custom-file-label">选择文件</label>
		<input type="file" name="image" id="uploadFile" accept=".jpg, .jpeg, .png" class="custom-file-input">
		<button id="uploadButton">上传</button>
		<div id="loading">上传中...</div>
		<div id="response" class="ui-widget">
		</div>
	
		<script>
	// 监听粘贴事件
	document.addEventListener('paste', function (e) {
		var items = e.clipboardData.items;
		for (var i = 0; i < items.length; i++) {
			var item = items[i];
			if (item.type.indexOf('image') !== -1) {
				// 获取粘贴的图片文件
				var file = item.getAsFile();
				// 调用上传函数，将file传递给上传逻辑
				$('#uploadFileLabel').text("已选择剪贴板文件").css('background-color', '#0056b3');
				uploadImg(file);
			}
		}
	});
	$(document).ready(function () {
		$('#uploadFile').change(function () {
			var fileName = $(this).val().split('\\').pop();
			if (fileName) {
				$('#uploadFileLabel').text('已选择文件: ' + fileName).css('background-color', '#0056b3');
			} else {
				$('#uploadFileLabel').text('选择文件').css('background-color', '#007BFF');
			}
		});
		$('#uploadButton').click(function () {
			var fileInput = document.getElementById('uploadFile');
			var file = fileInput.files[0];
			if(file){
				uploadImg(file)
			}else{
				alert('请选择一个图片文件');
			}
		});
	});
	function uploadImg(file){
		var formData = new FormData();
				formData.append('image', file);
				// 禁用上传按钮并显示loading
				$('#uploadButton').prop('disabled', true);
				$('#uploadButton').text('上传中');
				$('#loading').show();
				var baseUrl = window.location.protocol + "//" + window.location.hostname;
				if(window.location.port !== "80" && window.location.port.length>0){
					baseUrl = baseUrl + ":" + window.location.port;
				}
				$.ajax({
					type: 'POST',
					url: baseUrl+'/api',
					data: formData,
					contentType: false,
					processData: false,
					success: function (response) {
						if (response.code === 1) {
							var imgUrl = baseUrl + response.message;
							var newItem = $('<div class="response-item response-success">上传成功，图片外链：<a target="_blank" href="' + imgUrl + '">' + imgUrl + '</a>' +
								'<div class="copy-links">' +
								'<span class="copy-code" data-clipboard-text="&lt;img src=&quot;' + imgUrl + '&quot; alt=&quot;Your Alt Text&quot;&gt;">HTML</span>' +
								'<span class="copy-code" data-clipboard-text="![Alt Text](' + imgUrl + ')">Markdown</span>' +
								'<span class="copy-code" data-clipboard-text="[img]' + imgUrl + '[/img]">BBCode</span>' +
								'</div></div>');
							$('#response').prepend(newItem); // 将新数据放在最前面
	
							// 清除文件输入框的值
							$('#uploadFile').val('');
							$('#uploadFileLabel').text('选择文件').css('background-color', '#007BFF');
	
							// 添加复制功能
							$('.copy-code').click(function () {
								var textToCopy = $(this).data('clipboard-text');
								var tempInput = $('<input>');
								$('body').append(tempInput);
								tempInput.val(textToCopy).select();
								document.execCommand('copy');
								tempInput.remove();
	
								// 显示复制成功文本
								var copyText = $(this);
								var originalText = copyText.text();
								copyText.text('复制成功');
	
								setTimeout(function () {
									copyText.text(originalText);
								}, 1000);
							});
						} else {
							var newItem = $('<div class="response-item response-error">上传失败,错误信息：' + response.message + '</div>');
							$('#response').prepend(newItem); // 将新数据放在最前面
						}
					},
					error: function () {
						var newItem = $('<div class="response-item response-error">上传失败</div>');
						$('#response').prepend(newItem); // 将新数据放在最前面
					},
					complete: function () {
						// 启用上传按钮并隐藏loading
						$('#uploadButton').prop('disabled', false);
						$('#uploadButton').text('上传');
						$('#loading').hide();
					}
				});
	}
		</script>
		<a target="_blank" href="https://github.com/csznet/tgState">
		<svg version="1.1" id="Layer_1" xmlns="http://www.w3.org/2000/svg" xmlns:xlink="http://www.w3.org/1999/xlink" x="0px" y="0px" width="44px" height="15px" viewBox="0 0 44 15" enable-background="new 0 0 44 15" xml:space="preserve">  <image id="image0" width="44" height="15" x="0" y="0"
    href="data:image/png;base64,iVBORw0KGgoAAAANSUhEUgAAACwAAAAPBAMAAABtkjCqAAAABGdBTUEAALGPC/xhBQAAACBjSFJN
AAB6JgAAgIQAAPoAAACA6AAAdTAAAOpgAAA6mAAAF3CculE8AAAAKlBMVEUAAAD/////ZgCZmWb/
aAT/+vj/agj//fz8/Pr9/f36+vj/+PT/ZwOammiamlTzAAAAAWJLR0QB/wIt3gAAAAlwSFlzAAAL
EwAACxMBAJqcGAAAAAd0SU1FB+cKAw8uKd154KgAAABtSURBVBjTY2DADhgFsQABBkYlOBA2hgHS
hR1FlRIFiwQFhYQtjC0sBScvFAQLB6kFBQHVaisJW062sACqbYYIKwFVK+kEKQlbXLxhKQgXVgNK
BKkCzbYw7kSoVhRVEhRUQpgNFaaGu7GHCXYAAPxWLJi8tpSVAAAAJXRFWHRkYXRlOmNyZWF0ZQAy
MDIzLTEwLTAzVDE1OjQ2OjQxKzAwOjAwZpEQ6AAAACV0RVh0ZGF0ZTptb2RpZnkAMjAyMy0xMC0w
M1QxNTo0Njo0MSswMDowMBfMqFQAAAAodEVYdGRhdGU6dGltZXN0YW1wADIwMjMtMTAtMDNUMTU6
NDY6NDErMDA6MDBA2YmLAAAAAElFTkSuQmCC" />
</svg>
		</a>
	</body>
	</html>
	`
	// 输出 HTML 表单
	io.WriteString(w, htmlForm)
}
