package main

import (
	"database/sql"
	"log"
	"net/http"
	"text/template"

	_ "github.com/go-sql-driver/mysql"
)

type Employee struct {
	Id   int
	Name string
	City string
}

// KONEKSI DENGAN DATABASE
func db_con() (db *sql.DB) {
	driver := "mysql"
	user := "ikal"
	pwd := "Rakhid_16"
	db_name := "Go_blog"

	db, err := sql.Open(driver, user+":"+pwd+"@/"+db_name)
	if err != nil {
		panic(err.Error())
	}

	return db
}

// INISIALISASI VARIABEL UNTUK RENDER FILES .html
var tmpl = template.Must(template.ParseGlob("form/*"))

// READ SEMUA DATA DARI TABEL employee
func Index(w http.ResponseWriter, r *http.Request) {
	db := db_con()
	sel_db, err := db.Query("SELECT * FROM employee ORDER By id DESC")

	if err != nil {
		panic(err.Error())
	}

	emp := Employee{}
	res := []Employee{}

	for sel_db.Next() {
		var id int
		var name, city string

		if err := sel_db.Scan(&id, &name, &city); err != nil {
			panic(err.Error())
		}

		emp.Id = id
		emp.Name = name
		emp.City = city
		res = append(res, emp)
	}

	tmpl.ExecuteTemplate(w, "Index", res)
	defer db.Close()
}

// LIHAT DATA BERDASARKAN ID
func Show(w http.ResponseWriter, r *http.Request) {
	db := db_con()
	n_Id := r.URL.Query().Get("id")

	sel_db, err := db.Query("SELECT * FROM employee WHERE id=?", n_Id)

	if err != nil {
		panic(err.Error())
	}

	emp := Employee{}

	for sel_db.Next() {
		var id int
		var name, city string

		if err = sel_db.Scan(&id, &name, &city); err != nil {
			panic(err.Error())
		}

		emp.Id = id
		emp.Name = name
		emp.City = city
	}

	tmpl.ExecuteTemplate(w, "Show", emp)
	defer db.Close()
}

func New(w http.ResponseWriter, r *http.Request) {
	tmpl.ExecuteTemplate(w, "New", nil)
}

func Edit(w http.ResponseWriter, r *http.Request) {
	db := db_con()
	n_id := r.URL.Query().Get("id")

	sel_db, err := db.Query("SELECT * FROM employee WHERE id=?", n_id)

	if err != nil {
		panic(err.Error())
	}

	emp := Employee{}

	for sel_db.Next() {
		var id int
		var name, city string

		if err = sel_db.Scan(&id, &name, &city); err != nil {
			panic(err.Error())
		}

		emp.Id = id
		emp.Name = name
		emp.City = city
	}

	tmpl.ExecuteTemplate(w, "Edit", emp)
	defer db.Close()
}

func Insert(w http.ResponseWriter, r *http.Request) {
	db := db_con()

	if r.Method == "POST" {
		name := r.FormValue("name")
		city := r.FormValue("city")

		ins_form, err := db.Prepare("INSERT INTO employee(name, city) VALUES(?,?)")
		if err != nil {
			panic(err.Error())
		}

		ins_form.Exec(name, city)
		log.Println("INSERT: Name: " + name + " | City: " + city)
	}

	defer db.Close()
	http.Redirect(w, r, "/", 301)
}

func Update(w http.ResponseWriter, r *http.Request) {
	db := db_con()

	if r.Method == "POST" {
		name := r.FormValue("name")
		city := r.FormValue("city")
		id := r.FormValue("uid")

		ins_form, err := db.Prepare("UPDATE employee SET name=?, city=? WHERE id=?")

		if err != nil {
			panic(err.Error())
		}

		ins_form.Exec(name, city, id)
		log.Println("UPDATE: Name: " + name + " | City: " + city)
	}

	defer db.Close()
	http.Redirect(w, r, "/", 301)
}

func Delete(w http.ResponseWriter, r *http.Request) {
	db := db_con()
	emp := r.URL.Query().Get("id")

	del_form, err := db.Prepare("DELETE FROM employee WHERE id=?")

	if err != nil {
		panic(err.Error())
	}

	del_form.Exec(emp)
	log.Println("DELETE")

	defer db.Close()
	http.Redirect(w, r, "/", 301)
}

func main() {
	http.HandleFunc("/", Index)
	http.HandleFunc("/show", Show)
	http.HandleFunc("/new", New)
	http.HandleFunc("/edit", Edit)
	http.HandleFunc("/insert", Insert)
	http.HandleFunc("/update", Update)
	http.HandleFunc("/delete", Delete)

	if err := http.ListenAndServe(":7000", nil); err != nil {
		panic(err.Error())
	}
}
