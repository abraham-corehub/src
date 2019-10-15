package main

import (
	"bytes"
	"crypto/rand"
	"crypto/sha1"
	"database/sql"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	_ "runtime/debug"
	"strconv"
	"strings"
	"text/template"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

//Response type is to send JSON data from Server to Client
type Response struct {
	Data []string
}

//MenuItems is a custom type to store Menu items loaded dynamicaly on the Web Page's Header Bar
type MenuItems struct {
	Item string
	Link string
 Icon string
}

//AppData is a custom type to store the Data related to the Application
type AppData struct {
	Title          string
	User           AppUser
	MenuItemsRight []MenuItems
 MenuItemsLeft  []MenuItems
	Page           PageData
	Table          DBTable
	State          string
	UI             string
	Error          string
}

//PageData is a custom type to store Title and 
//Content / Body of the Web Page to be displayed
type PageData struct {
	Name   string
	Title  string
	Body   string
	Author PageAuthor
}

//AppUser is a custom type to 
//store the User's Name and 
//access level (Role)
type AppUser struct {
	Name         string
	Role         int
	ID           int
	UN           string
	PW           string
	Status       string
	Token        string
	TC           time.Time
}

//PageAuthor is a custom type 
//to store the User's Name and 
//access level (Role)
type PageAuthor struct {
	Name string
	ID   int
}

//DBTable is custom
type DBTable struct {
	Header RowData
	Rows   []RowData
}

//RowData is custom
type RowData struct {
	Index int
	Row   []ColData
}

//ColData is custom
type ColData struct {
	Index int
	Value string
}

const dirRes = `res/`
const dirImg = dirRes + `img/`
const dirTmpl = `tmpl/mdl/`
const dirSttc = `static/`
const dirFS = dirTmpl + `static/`

const pathDB = `db/`
const fNDB = `pb.db`

var tmplts *template.Template
var aD AppData
var dB *sql.DB

func main() {
	//initDB()
	//testDB()
	fmt.Println(`Server Started!,Please access http://localhost:8080`)
	startPhotoBook()
}

func startPhotoBook() {
	parseTmplts()
	initApp()
	initDB()
	mux := http.NewServeMux()
	fS := http.FileServer(fSNeutered{http.Dir(dirFS)})
	mux.Handle(`/`+dirSttc, http.StripPrefix(`/`+dirSttc, fS))

	mux.HandleFunc(`/`, handlerRoot)
	mux.HandleFunc(`/login`, handlerLogin)
	mux.HandleFunc(`/logout`, handlerLogout)
	mux.HandleFunc(`/home`, handlerHome)
	mux.HandleFunc(`/user/view`, handlerViewUser)
	mux.HandleFunc(`/album/view`, handlerViewAlbum)
	//mux.HandleFunc("/image/view", handlerViewImage)
	//mux.HandleFunc("/user/edit", handlerUserEdit)
	//mux.HandleFunc("/user/reset", handlerUserReset)
	//mux.HandleFunc("/user/delete", handlerUserDelete)
	//mux.HandleFunc("/user", handlerUser)
	log.Fatal(http.ListenAndServe(`:8080`, mux))
}

func parseTmplts() {
	nUITs := []string{
		`head`,
		`login`,
		`home`,
		`users`,
		`albums`,
		`images`,
		`image`,
	}
	for i := 0; i < len(nUITs); i++ {
		nUITs[i] = dirTmpl + nUITs[i] + `.html`
	}
	tmplts = template.Must(template.ParseFiles(nUITs...))
}

func initApp() {
	aD = AppData{}
	aD.User = AppUser{}
	aD.Page = PageData{}
	aD.Table = DBTable{}
	aD.Page.Author = PageAuthor{}
	aD.Title = `PhotoBook`
}

func initDB() {
	dBT, err := sql.Open(`sqlite3`, pathDB+fNDB)
	if err != nil {
		log.Fatal(err)
	}
	dB = dBT
}

func handlerRoot(w http.ResponseWriter, r *http.Request) {
	if err := isAuthorized(w, r); err == nil {
		loadPage(w, r)
	} else {
		http.Redirect(w, r, `/login`, http.StatusSeeOther)
	}
}

func handlerLogin(w http.ResponseWriter, r *http.Request) {
	if err := isAuthorized(w, r); err != nil {
		initApp()
		aD.UI = `login`
		if r.Method == `POST` {
			if err := getFormData(w, r); err != nil {
				aD.User.Status = err.Error()
			} else {
				http.Redirect(w, r, `/home`, http.StatusSeeOther)
			}
		}
		loadPage(w, r)
		return
	}
	http.Redirect(w, r, `/home`, http.StatusSeeOther)
}

func handlerHome(w http.ResponseWriter, r *http.Request) {
	if err := isAuthorized(w, r); err != nil {
		showError(w, r, errors.New(`unauthorized access`))
	} else {
		aD.User.Status = `Login Successfull`
		aD.UI = `home`
		switch aD.User.Role {
		case -7:
			aD.Page.Name = `users`
			aD.Page.Title = `Dashboard`
			dBGetUsers(w, r)
		default:
			aD.Page.Name = `albums`
			aD.Page.Title = `My Albums`
			aD.Page.Author.ID = aD.User.ID
			aD.Page.Author.Name = aD.User.Name
			dBGetAlbums(w, r)
		}
		loadPage(w, r)
	}
}

func handlerLogout(w http.ResponseWriter, r *http.Request) {
	if err := isAuthorized(w, r); err == nil {
		dBDelSession(w, r)
		http.SetCookie(w, &http.Cookie{
			Name:   `sessionToken`,
			Value:  ``,
			MaxAge: 0,
		})
		if aD.User.Status == `Login Successfull` {
			aD.User.Status = `Logout Successfull`
		}
		aD.UI = `login`
		loadPage(w, r)
	} else {
		http.Redirect(w, r, `/`, http.StatusSeeOther)
	}
}

func handlerViewUser(w http.ResponseWriter, r *http.Request) {
	if err := isAuthorized(w, r); err == nil {
		if err := r.ParseForm(); err != nil {
			showError(w, r, err)
		}
		aD.Page.Author.ID, _ = strconv.Atoi(r.Form.Get(`id`))
		aD.Page.Author.Name = r.Form.Get(`name`)
		dBGetAlbums(w, r)
		aD.Page.Name = `albums`
		aD.Page.Title = aD.Page.Author.Name + `'s Albums`
		loadPage(w, r)
	} else {
		showError(w, r, err)
	}
}

func handlerViewAlbum(w http.ResponseWriter, r *http.Request) {
	if err := isAuthorized(w, r); err == nil {
		if err := r.ParseForm(); err != nil {
			showError(w, r, err)
		}
		idAlbum := r.Form.Get(`id`)
		nameAlbum := r.Form.Get(`name`)
		aD.Page.Name = `images`
		aD.Page.Title = aD.Page.Author.Name + `'s ` + nameAlbum
		dBGetImgs(w, r, idAlbum)
		loadPage(w, r)
	} else {
		showError(w, r, err)
	}
}

/*
func handlerViewImage(w http.ResponseWriter, r *http.Request) {
	if err := isAuthorized(w, r); err == nil {
		if err := r.ParseForm(); err != nil {
			showError(w, r, err)
		}
		idImage := r.Form.Get("id")
		nameImage := r.Form.Get("name")
		aD.Page.Name = "image"
		aD.Page.Title = aD.Page.Author.Name + "'s " + nameImage
		dbGetImage(w, r, idImage, nameImage)
		loadPage(w, r)
	} else {
		showError(w, r, err)
	}
}
*/

func isAuthorized(w http.ResponseWriter, r *http.Request) error {
	if r.Method == `GET` {
		c, err := r.Cookie(`sessionToken`)
		if err == nil {
			if c.Value == aD.User.Token && aD.User.Token != "" {
				return nil
			}
			return dbRestoreUser(w, r, c.Value)
		}
		return errors.New(`session expired`)
	}
	return errors.New(`unauthorized access`)
}

func getFormData(w http.ResponseWriter, r *http.Request) error {
	if err := r.ParseForm(); err != nil {
		return err
	}
	if err := auth(w, r); err != nil {
		return err
	}
	if err := dBDelSession(w, r); err != nil {
		return err
	}
	if err := setCookie(w, r); err != nil {
		return err
	}
	if err := dBStoreSession(w, r); err != nil {
		return err
	}
	return nil
}

func showError(w http.ResponseWriter, r *http.Request, err error) {
	//log.Panicf(err.Error())
	fmt.Println(err.Error())
	//debug.PrintStack()
	http.Redirect(w, r, `/logout`, http.StatusSeeOther)
}

func renderTmplt(w http.ResponseWriter, r *http.Request) {
	err := tmplts.ExecuteTemplate(w, aD.UI+`.html`, aD)
	if err != nil {
		showError(w, r, err)
	}
}

func loadMenuItems() {
  aD.MenuItemsLeft = []MenuItems{
  {Item: `Home`, Icon: `home`, Link: `/home`},
  {Item: `My Account`, Icon: `account_circle`, Link: `/myAccount`},
  {Item: `Logout`, Icon: `arrow_back`, Link: `/logout`},
 }
	switch aD.User.Role {
	case -7:
		aD.MenuItemsRight = []MenuItems{
			{Item: `Create User`, Link: `/createUser`},
			{Item: `Upload Image`, Link: `/uploadImage`},
			{Item: `Create Album`, Link: `/createAlbum`},
			{Item: `Download Album`, Link: `/downloadAlbum`},
  }
	default:
		aD.Page.Author.ID = aD.User.ID
		aD.MenuItemsRight = []MenuItems{
			{Item: `Upload Image`, Link: `/uploadImage`},
			{Item: `Create Album`, Link: `/createAlbum`},
			{Item: `Download Album`, Link: `/downloadAlbum`},
		}
	}
}

func auth(w http.ResponseWriter, r *http.Request) error {
	uN := r.Form.Get(`username`)
	pW := r.Form.Get(`password`)
	uN, err := cleanStr(uN)
	if err == nil {
		aD.User.UN = uN
		if dBIsUNValid(w, r, uN) {
			pW, err := cleanStr(pW)
			if err == nil {
				aD.User.PW = pW
				pWH := sha1.New()
				pWH.Write([]byte(pW))
				pWHS := hex.EncodeToString(pWH.Sum(nil))
				if err := dBAuth(w, r, uN, pWHS); err == nil {
					return nil
				}
				return errors.New(`unregistered password`)
			}
		}
	}
	return errors.New(`unregistered username`)
}

func loadPage(w http.ResponseWriter, r *http.Request) {
	switch aD.UI {
	case `home`:
		loadMenuItems()
		aD.Page.loadPageBody(w, r)
	case ``:
		showError(w, r, errors.New(`UI not set`))
	}
	renderTmplt(w, r)
}

func isTimeExpired(dTSExpr int64) bool {
	dTSNow := time.Now()
	if dTSNow.Unix()-dTSExpr > 120 {
		return true
	}
	return false
}

func setCookie(w http.ResponseWriter, r *http.Request) error {
	uuid, err := newUUID()
	if err != nil {
		return err
	}
	aD.User.Token = uuid
	aD.User.TC = time.Now()
	http.SetCookie(w, &http.Cookie{
		Name:    `sessionToken`,
		Value:   uuid,
		Expires: aD.User.TC.Add(120 * time.Second),
	})
	return nil
}

func (PageData) loadPageBody(w http.ResponseWriter, r *http.Request) {
	var tpl bytes.Buffer
	err := tmplts.ExecuteTemplate(&tpl, aD.Page.Name+`.html`, aD)
	if err != nil {
		showError(w, r, err)
	}
	aD.Page.Body = tpl.String()
}

func dBIsUNValid(w http.ResponseWriter, r *http.Request, uN string) bool {
	stmt, err := dB.Prepare(`select * from user where username = ?`)
	rows, err := stmt.Query(uN)
	if err != nil {
		showError(w, r, err)
	}
	defer rows.Close()
	if rows.Next() {
		return true
	}
	return false
}

func dbRestoreUser(w http.ResponseWriter, r *http.Request, sT string) error {
	stmt, err := dB.Prepare(
		`select id_user, 
  datetimestamp_lastlogin from session where id = ?`)
	rows, err := stmt.Query(sT)
	if err != nil {
		return err
	}
	defer rows.Close()
	if rows.Next() {
		var idUser int
		var dTS string
		err := rows.Scan(&idUser, &dTS)
		if err != nil {
			return err
		}
		dTSExpr, _ := strconv.ParseInt(dTS, 10, 64)
		if isTimeExpired(dTSExpr) {
			return errors.New(`time expired`)
		}
		aD.User.ID = idUser
		aD.User.Token = sT
	} else {
		return errors.New(`invalid session`)
	}

	stmt, err = dB.Prepare(`select name, role from user where id = ?`)
	rows, err = stmt.Query(aD.User.ID)
	if err != nil {
		return err
	}
	defer rows.Close()
	if rows.Next() {
		var name string
		var role int
		err = rows.Scan(&name, &role)
		if err != nil {
			return err
		}
		aD.User.Name = name
		aD.User.Role = role
		return nil
	}
	return errors.New(`session user not found in user list`)
}

func dBAuth(w http.ResponseWriter, r *http.Request, uN string, pW string) error {
	stmt, err := dB.Prepare(
		`select name,
		role, 
		id from user where 
		username = ? and 
		password = ?`,
	)
	rows, err := stmt.Query(uN, pW)

	if err != nil {
		return err
	}

	defer rows.Close()

	if rows.Next() {
		var name string
		var role int
		var id int
		err = rows.Scan(&name, &role, &id)
		if err != nil {
			return err
		}
		aD.User.Name = name
		aD.User.Role = role
		aD.User.ID = id
		return nil
	}
	return errors.New(`db empty`)
}

func dBGetUsers(w http.ResponseWriter, r *http.Request) error {
	aD.Table.Header = RowData{
		0,
		[]ColData{
			{Index: 0, Value: `id`},
			{Index: 1, Value: `name`},
			{Index: 2, Value: `username`},
		},
	}

	aD.Table.Rows = make([]RowData, 0)
	stmt, err := dB.Prepare(
		`select ` +
			aD.Table.Header.Row[0].Value + `, ` +
			aD.Table.Header.Row[1].Value +
			` from user where ` +
			aD.Table.Header.Row[2].Value +
			` != ?`)
	rows, err := stmt.Query(`admin`)
	if err != nil {
		return err
	}

	defer rows.Close()
	for rows.Next() {
		var name string
		var id int
		err = rows.Scan(&id, &name)
		if err != nil {
			return err
		}
		aD.Table.Rows = append(
			aD.Table.Rows, RowData{
				Index: id,
				Row: []ColData{
					{Value: name},
				},
			},
		)
	}

	if len(aD.Table.Rows) > 0 {
		return nil
	}
	return errors.New(`db empty`)
}

func dBGetAlbums(w http.ResponseWriter, r *http.Request) error {
	aD.Table.Header = RowData{
		0,
		[]ColData{
			{Index: 0, Value: `name`},
			{Index: 1, Value: `id_user`},
		},
	}

	aD.Table.Rows = make([]RowData, 0)

	stmt, err := dB.Prepare(
		`select ` +
			aD.Table.Header.Row[0].Value +
			` from album where ` +
			aD.Table.Header.Row[1].Value +
			`=?`)

	rows, err := stmt.Query(strconv.Itoa(aD.Page.Author.ID))
	if err != nil {
		return err
	}

	defer rows.Close()
	for rows.Next() {
		var name string
		err = rows.Scan(&name)
		if err != nil {
			return err
		}
		aD.Table.Rows = append(
			aD.Table.Rows,
			RowData{
				Index: len(aD.Table.Rows) + 1,
				Row: []ColData{
					{Value: name},
				},
			},
		)
	}

	if len(aD.Table.Rows) > 0 {
		return nil
	}
	return errors.New(`db empty`)
}

func dBGetImgs(w http.ResponseWriter, r *http.Request, idA string) error {
	aD.Table.Header = RowData{
		0,
		[]ColData{
			{Index: 0, Value: `name`},
			{Index: 1, Value: `id_user`},
			{Index: 2, Value: `id_album`},
		},
	}

	aD.Table.Rows = make([]RowData, 0)

	stmt, err := dB.Prepare(
		`select ` +
			aD.Table.Header.Row[0].Value +
			` from image where ` +
			aD.Table.Header.Row[1].Value + `=? and ` +
			aD.Table.Header.Row[2].Value + ` =?`)

	rows, err := stmt.Query(aD.Page.Author.ID, idA)
	if err != nil {
		return err
	}

	defer rows.Close()
	for rows.Next() {
		var name string
		err = rows.Scan(&name)
		if err != nil {
			return err
		}
		aD.Table.Rows = append(
			aD.Table.Rows,
			RowData{
				Index: len(aD.Table.Rows) + 1,
				Row: []ColData{
					{Value: name},
				},
			})
	}
	if len(aD.Table.Rows) > 0 {
		return nil
	}
	return errors.New("db empty")
}

func dBGetImg(w http.ResponseWriter, r *http.Request, idImg string) error {
	aD.Table.Header = RowData{
		0,
		[]ColData{
			{Index: 0, Value: "name"},
			{Index: 1, Value: "id_album"},
			{Index: 2, Value: "type"},
			{Index: 3, Value: "id"},
		},
	}
	aD.Table.Rows = make([]RowData, 0)

	stmt, err := dB.Prepare(
		`select ` +
			aD.Table.Header.Row[0].Value +
			`, ` +
			aD.Table.Header.Row[1].Value +
			`, ` +
			aD.Table.Header.Row[2].Value +
			` from image where ` +
			aD.Table.Header.Row[3].Value + `=?`)

	rows, err := stmt.Query(strconv.Atoi(idImg))
	if err != nil {
		return err
	}

	defer rows.Close()
	for rows.Next() {
		var name string
		var idAlbum int
		var imgType string
		err = rows.Scan(&name, &idAlbum, &imgType)
		if err != nil {
			return err
		}
		aD.Table.Rows = append(
			aD.Table.Rows,
			RowData{
				Index: len(aD.Table.Rows) + 1,
				Row: []ColData{
					{Value: dirImg +
						strconv.Itoa(idAlbum) +
						"/" +
						idImg +
						"." +
						imgType},
					{Value: name},
				},
			},
		)
	}
	if len(aD.Table.Rows) > 0 {
		return nil
	}
	return errors.New(`db empty`)
}

func dBDelSession(w http.ResponseWriter, r *http.Request) error {
	statement, err := dB.Prepare(`delete from session where id_user=?`)
	if err != nil {
		return err
	}
	_, err = statement.Exec(aD.User.ID)
	if err != nil {
		return err
	}
	return nil
}

func dBStoreSession(w http.ResponseWriter, r *http.Request) error {
	stmt, err := dB.Prepare(`PRAGMA foreign_keys = true;`)
	if err != nil {
		return err
	}
	_, err = stmt.Exec()
	if err != nil {
		return err
	}

	stmt, err = dB.Prepare(
		`INSERT INTO session (
 id, 
 id_user, 
 datetimestamp_lastlogin) 
 VALUES ( ?, ?, ?)`)

	if err != nil {
		return err
	}
	_, err = stmt.Exec(aD.User.Token, aD.User.ID, aD.User.TC)
	if err != nil {
		return err
	}
	return nil
}

func dBStoreTokenTx(w http.ResponseWriter, t string, dTSE time.Time) error {

	stmts := []string{
		`RAGMA foreign_keys = true;`,
		`INSERT INTO session (
  id, 
  id_user, 
  datetimestamp_lastlogin) 
  VALUES (?, ?, ?)`,
	}

	for i, stmt := range stmts {
		trashSQL, err := dB.Prepare(stmt)
		if err != nil {
			return err
		}
		tx, err := dB.Begin()
		if err != nil {
			return err
		}

		switch i {
		case 0:
			_, err = tx.Stmt(trashSQL).Exec()
		case 1:
			_, err = tx.Stmt(trashSQL).Exec(
				t,
				strconv.Itoa(aD.User.ID),
				strconv.FormatInt(dTSE.Unix(),
					10))
		}

		if err != nil {
			fmt.Println(`doing rollback`)
			tx.Rollback()
		} else {
			tx.Commit()
		}
	}
	return nil
}

func testDB() {
	/*
		queryString := "select name from user"
		rows, err := dB.Query(queryString)
		if err != nil {
			log.Fatal(err)
		}
		defer rows.Close()
		for rows.Next() {
			var name string
			rows.Scan(&name)
			fmt.Println(name)
		}
	*/

	id := 1
	trashSQL, err := dB.Prepare(`delete from session where id_user=?`)
	if err != nil {
		fmt.Println(err)
	}
	tx, err := dB.Begin()
	if err != nil {
		fmt.Println(err)
	}
	_, err = tx.Stmt(trashSQL).Exec(id)
	if err != nil {
		fmt.Println(`doing rollback`)
		tx.Rollback()
	} else {
		tx.Commit()
	}

	id = 2
	stmt, err := dB.Prepare(`select name from user where id=?`)
	if err != nil {
		log.Fatal(err)
	}
	rows, err := stmt.Query(id)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()
	for rows.Next() {
		var name string
		rows.Scan(&name)
		fmt.Println(name)
	}

}

func cleanStr(str string) (string, error) {
	strN := str
	/*
		charsTrim := []byte{
			' ','\\','"','&','<','>','(',')','|','/','=',';',':','`',
		}
	*/
	charsTrim := []byte(`\ ,"'<>{}[]|/=;:.?`)
	for _, cH := range charsTrim {
		strN = strings.ReplaceAll(strN, string(cH), "")
	}
	if len(str) != len(strN) {
		return str, errors.New(`unclean string`)
	}
	return str, nil
}

// newUUID generates a random UUID according to RFC 4122
// https://play.golang.org/p/w7qciopoosz
func newUUID() (string, error) {
	uuid := make([]byte, 16)
	n, err := io.ReadFull(rand.Reader, uuid)
	if n != len(uuid) || err != nil {
		return ``, err
	}
	// variant bits; see section 4.1.1
	uuid[8] = uuid[8]&^0xc0 | 0x80
	// version 4 (pseudo-random); see section 4.1.3
	uuid[6] = uuid[6]&^0xf0 | 0x40
	return fmt.Sprintf(`%x-%x-%x-%x-%x`, uuid[0:4], uuid[4:6], uuid[6:8], uuid[8:10], uuid[10:]), nil
}

//To disable Directory Listing
//https://www.alexedwards.net/blog/disable-http-fileserver-directory-listings
type fSNeutered struct {
	fS http.FileSystem
}

//To disable Directory Listing
//https://www.alexedwards.net/blog/disable-http-fileserver-directory-listings
func (fSN fSNeutered) Open(path string) (http.File, error) {
	f, err := fSN.fS.Open(path)
	if err != nil {
		return nil, err
	}

	s, err := f.Stat()
	if s.IsDir() {
		index := strings.TrimSuffix(path, `/`) + `/index.html`
		if _, err := fSN.fS.Open(index); err != nil {
			return nil, err
		}
	}

	return f, nil
}
