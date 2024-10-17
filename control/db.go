package control

import (
	"database/sql"
	"errors"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
	"log"
	"sync"
)

var (
	db   *sql.DB
	once sync.Once
)

// InitDB 初始化数据库，创建表结构
func InitDB() (*sql.DB, error) {
	var err error
	// 使用 sync.Once 确保数据库只初始化一次
	once.Do(func() {
		db, err = sql.Open("sqlite3", "./files.db")
		if err != nil {
			log.Fatal("Failed to open database:", err)
		}

		query := `CREATE TABLE IF NOT EXISTS uploaded_files (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			fileId TEXT NOT NULL,
			filename TEXT NOT NULL,
			ip TEXT NOT NULL,
			time TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		);`
		_, err = db.Exec(query)
		if err != nil {
			log.Fatal("Failed to create table:", err)
		}
	})

	return db, err
}

// GetFileNameByID 查询文件名
func GetFileNameByID(id string) (string, error) {
	var fileName string
	// 执行查询，获取对应id的fileName
	query := "SELECT filename FROM uploaded_files WHERE fileId = ?"
	err := db.QueryRow(query, id).Scan(&fileName)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return "", fmt.Errorf("no file found with id %s", id)
		}
		return "", err
	}

	return fileName, nil
}

func SaveFileRecord(fileID string, fileName string, ip string) error {
	// 插入数据到数据库
	_, err := db.Exec("INSERT INTO uploaded_files (fileId, filename, ip) VALUES (?, ?, ?)", fileID, fileName, ip)
	return err
}
