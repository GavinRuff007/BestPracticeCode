package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"reflect"
	"strings"

	_ "modernc.org/sqlite"
)

/*-------------------------------------------*/
/*Write your Entity Database: */
/*-------------------------------------------*/
type Product struct {
	ID          int
	ProductCode string
	Name        string
	Price       int
	Status      string
	Inventory   string
}

/*-------------------------------------------*/

/*-------------------------------------------*/
/*Write your Database Name: */
/*-------------------------------------------*/
var databaseName = "product"

/*-------------------------------------------*/

/*-------------------------------------------*/
/*Write your Table Name: */
/*-------------------------------------------*/
var tableName = "products"

/*-------------------------------------------*/

func createTableFromStruct(db *sql.DB, tableName string, model interface{}) error {
	v := reflect.TypeOf(model)
	if v.Kind() != reflect.Struct {
		return fmt.Errorf("ورودی باید struct باشد")
	}

	var columns []string
	for i := 0; i < v.NumField(); i++ {
		field := v.Field(i)
		columnName := strings.ToLower(field.Name)
		sqlType := goTypeToSQLType(field.Type.Kind())
		if columnName == "id" {
			columns = append(columns, fmt.Sprintf("%s INTEGER PRIMARY KEY", columnName))
		} else {
			columns = append(columns, fmt.Sprintf("%s %s", columnName, sqlType))
		}
	}

	query := fmt.Sprintf(`CREATE TABLE IF NOT EXISTS %s (%s);`, tableName, strings.Join(columns, ", "))
	_, err := db.Exec(query)
	return err
}

func goTypeToSQLType(kind reflect.Kind) string {
	switch kind {
	case reflect.Int, reflect.Int64:
		return "INTEGER"
	case reflect.Float32, reflect.Float64:
		return "REAL"
	case reflect.String:
		return "TEXT"
	case reflect.Bool:
		return "BOOLEAN"
	default:
		return "BLOB"
	}
}

func getDesktopPath() (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}

	desktopPath := filepath.Join(homeDir, "Desktop")
	return desktopPath, nil
}

func main() {
	desktop, err := getDesktopPath()
	if err != nil {
		log.Fatal("خطا در پیدا کردن دسکتاپ:", err)
	}

	dbPath := filepath.Join(desktop, databaseName+".db") // بهتره پسوند .db اضافه بشه

	db, err := sql.Open("sqlite", dbPath)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	err = createTableFromStruct(db, tableName, Product{})
	if err != nil {
		log.Fatal(err)
	}

	log.Println("فایل پایگاه داده در دسکتاپ ساخته شد:", dbPath)

	// ─────────────────────────────
	// نمایش جدول‌ها و ستون‌ها
	// ─────────────────────────────
	fmt.Println("\n--- لیست جداول در پایگاه داده ---")
	tables, err := db.Query(`SELECT name FROM sqlite_master WHERE type='table';`)
	if err != nil {
		log.Fatal("خطا در گرفتن لیست جداول:", err)
	}
	defer tables.Close()

	for tables.Next() {
		var tableName string
		if err := tables.Scan(&tableName); err != nil {
			log.Fatal(err)
		}
		fmt.Println("جدول:", tableName)

		// نمایش ستون‌های جدول
		fmt.Println("ستون‌ها:")
		columns, err := db.Query(fmt.Sprintf("PRAGMA table_info(%s);", tableName))
		if err != nil {
			log.Fatal(err)
		}
		defer columns.Close()

		for columns.Next() {
			var cid int
			var name, ctype string
			var notnull, pk int
			var dfltValue interface{}
			if err := columns.Scan(&cid, &name, &ctype, &notnull, &dfltValue, &pk); err != nil {
				log.Fatal(err)
			}
			fmt.Printf(" - %s (%s)\n", name, ctype)
		}
		fmt.Println()
	}
}
