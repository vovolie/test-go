package data

import (
	"database/sql"
	"log"

	"strings"

	_ "github.com/go-sql-driver/mysql"
	"fmt"
)

var (
	db                  *sql.DB
	materialBatachStmts []*sql.Stmt
	materialBaseSQL     string

	MaxMultiQueryCount int
)

type Material struct {
	ID           int    `json:"id"`
	Category 	 int    `json:"category"`
	Cover        string `json:"cover"`
	Name         string `json:"name"`
	URL          string `json:"url"`
	Sha          string `json:"sha"`
	Version      string `json:"version"`
	MateInfo     string `json:"mate_info"`
	HiddenAt     int    `json:"hidden_at"`
	CreatedAt    string `json:"created_at"`
}

type Category struct {
	ID string `json:"id"`
	Name string `json:"name"`
}

var CategoryList []Category


func GetCategory(id string) Category {
	for _, category := range CategoryList {
		if category.ID == id {
			return category
		}
	}
	return Category{}
}

func GetMaterialById(materialId string) (*Material, error) {
	stmt, err := db.Prepare("SELECT id, category, cover, name, url, sha, version, mate_info, hidden_at, created_at FROM material_library WHERE id = ? LIMIT 1")
	if err != nil {
		log.Println("prepare sql error", err)
		return nil, err
	}
	defer stmt.Close()

	material := new(Material)
	var id, category, hiddenAt int
	var cover, name, url, sha, version, mateInfo, createdAt sql.NullString
	err = stmt.QueryRow(materialId).Scan(&id, &category, &cover, &name, &url, &sha, &version, &mateInfo, &hiddenAt, &createdAt)
	switch {
	case err == sql.ErrNoRows:
		log.Println("cannot get matrial data by id: ", materialId)
		return nil, err
	case err != nil:
		log.Println("get material data error: ", err)
		return nil, err
	default:
		material.ID = id
		material.Category = category
		material.Cover = cover.String
		material.Name = name.String
		material.URL = url.String
		material.Sha = sha.String
		material.Version = version.String
		material.MateInfo = mateInfo.String
		material.HiddenAt = hiddenAt
		material.CreatedAt = createdAt.String
	}
	return material, nil
}

func InsertMaterial(material *Material) error {
	instertSQL := "INSERT INTO material_library(cover, name, url, sha, version, mate_info, hidden_at, created_at, category)"
	instertSQL += "VALUES(?, ?, ?, ?, ?, ?, ?, ?, ?)"
	_, err := db.Exec(instertSQL, material.Cover, material.Name, material.URL,
		material.Sha, material.Version, material.MateInfo, material.HiddenAt, material.CreatedAt, material.Category)
	if err != nil {
		fmt.Println(err)
	}
	return err
}

func DelMaterialById(materialId int) bool {
	deleteSQL := "DELETE FROM material_library WHERE id = ?"
	_, err := db.Exec(deleteSQL, materialId)
	if err != nil {
		fmt.Println(err)
		return false
	}
	return true
}

func SearchMaterialByName(keyword string, limit int64, offset int64) ([]*Material, error) {
	searchSQL := "SELECT id, cover, name, url, sha, version, mate_info, hidden_at, UNIX_TIMESTAMP(created_at) AS created_at,"
	searchSQL += "category FROM material_library WHERE name LIKE '%" + keyword + "%'  LIMIT ? OFFSET ?"
	rows, err := db.Query(searchSQL, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	result := make([]*Material, 0, limit)
	for rows.Next() {
		material := new(Material)
		var cover, name, url, sha, version, mateInfo, createdAt sql.NullString
		err := rows.Scan(&material.ID, &cover, &name, &url, &sha, &version, &mateInfo, &material.HiddenAt,
			&createdAt, &material.Category)
		if err != nil {
			return nil, err
		}
		material.Cover = cover.String
		material.Name = name.String
		material.URL = url.String
		material.Sha = sha.String
		material.Version = version.String
		material.MateInfo = mateInfo.String
		material.CreatedAt = createdAt.String
		result = append(result, material)
	}

	return result, nil
}

func GetAllMaterial(limit int64, offset int64) ([]*Material, error) {
	getSQL := "SELECT id, cover, name, url, sha, version, mate_info, hidden_at, UNIX_TIMESTAMP(created_at) AS created_at,"
	getSQL += "category FROM material_library LIMIT ? OFFSET ?"
	rows, err := db.Query(getSQL, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	result := make([]*Material, 0, limit)
	for rows.Next() {
		material := new(Material)
		var cover, name, url, sha, version, mateInfo, createdAt sql.NullString
		err := rows.Scan(&material.ID, &cover, &name, &url, &sha, &version, &mateInfo, &material.HiddenAt,
			&createdAt, &material.Category)
		if err != nil {
			return nil, err
		}
		material.Cover = cover.String
		material.Name = name.String
		material.URL = url.String
		material.Sha = sha.String
		material.Version = version.String
		material.MateInfo = mateInfo.String
		material.CreatedAt = createdAt.String
		result = append(result, material)
	}

	return result, nil
}

func GetMaterialByCategory(category string, limit, offset int64) ([]*Material, error) {
	catSQL := "SELECT id, cover, name, url, sha, version, mate_info, hidden_at, UNIX_TIMESTAMP(created_at) AS created_at,"
	catSQL += "category FROM material_library WHERE category=?"
	catStmt, err := db.Prepare(catSQL)
	if err != nil {
		log.Println("prepare sql error", err)
		return nil, err
	}
	defer catStmt.Close()
	result := make([]*Material, 0, limit)
	rows, err := catStmt.Query(category)
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		material := new(Material)
		var cover, name, url, sha, version, mateInfo, createdAt sql.NullString
		err := rows.Scan(&material.ID, &cover, &name, &url, &sha, &version, &mateInfo, &material.HiddenAt,
			&createdAt, &material.Category)
		if err != nil {
			return nil, err
		}
		material.Cover = cover.String
		material.Name = name.String
		material.URL = url.String
		material.Sha = sha.String
		material.Version = version.String
		material.MateInfo = mateInfo.String
		material.CreatedAt = createdAt.String
		result = append(result, material)
	}

	return result, nil

}


func GetMaterialByIds(materialIds []interface{}) ([]*Material, error) {
	length := len(materialIds)
	if length > len(materialBatachStmts) {
		log.Fatal("material ids count exceed: ", materialIds)
		err := fmt.Errorf("批量查询数量超出限制")
		return nil, err
	}
	rows, err := materialBatachStmts[length-1].Query(materialIds...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	result := make([]*Material, 0, length)
	for rows.Next() {
		material := new(Material)
		var cover, name, url, sha, version, mateInfo, createdAt sql.NullString
		err := rows.Scan(&material.ID, &cover, &name, &url, &sha, &version, &mateInfo, &material.HiddenAt,
			&createdAt, &material.Category)
		if err != nil {
			return nil, err
		}
		material.Cover = cover.String
		material.Name = name.String
		material.URL = url.String
		material.Sha = sha.String
		material.Version = version.String
		material.MateInfo = mateInfo.String
		material.CreatedAt = createdAt.String
		result = append(result, material)
	}

	return result, nil
}


func MaterialsToSliceInterface(materials []*Material) []interface{} {
	materialsIFace := []interface{}{}
	for _, material := range materials {
		materialsIFace = append(materialsIFace, material)
	}
	return materialsIFace
}


func UpdateMaterial(materialId int64, cover, name, url, sha, version, mateInfo string, hiddenAt int, createdAt string, category int) error {
	_, err := db.Exec("UPDATE matrial_library SET cover = ?, name = ?, url = ?, sha = ?, version = ?, mate_info = ?, hidden_at = ?, created_at = ?, category = ? WHERE id = ?", cover, name, url, sha, version, mateInfo, hiddenAt, createdAt, category, materialId)
	return err
}

func Init() (err error) {
	db, err = sql.Open("mysql", "zl:123@tcp(127.0.0.1:3306)/vue?timeout=3s")
	if err != nil {
		return err
	}

	categoryMusic := Category{ID: "1", Name: "音乐" }
	categoryDecorative := Category{ID: "2", Name: "贴纸"}

	CategoryList = append(CategoryList, categoryMusic, categoryDecorative)

	materialBaseSQL = "SELECT id, cover, name, url, sha, version, mate_info, hidden_at, UNIX_TIMESTAMP(created_at) AS created_at,"
	materialBaseSQL += "category FROM material_library WHERE "
	materialBatachStmts = make([]*sql.Stmt, 0, 300)
	for i := 0; i < MaxMultiQueryCount; i++ {
		clause := strings.TrimRight(strings.Repeat("?,", i+1), ",")
		stmt, err := db.Prepare(materialBaseSQL + "id IN (" + clause + ")")
		if err != nil {
			return err
		}
		materialBatachStmts = append(materialBatachStmts, stmt)
	}
	return
}

