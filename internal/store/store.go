package store
import ("database/sql";"fmt";"os";"path/filepath";"time";_ "modernc.org/sqlite")
type DB struct{db *sql.DB}
type NewHire struct {
	ID string `json:"id"`
	Name string `json:"name"`
	Email string `json:"email"`
	Department string `json:"department"`
	StartDate string `json:"start_date"`
	Manager string `json:"manager"`
	Progress int `json:"progress"`
	Status string `json:"status"`
	Notes string `json:"notes"`
	CreatedAt string `json:"created_at"`
}
func Open(d string)(*DB,error){if err:=os.MkdirAll(d,0755);err!=nil{return nil,err};db,err:=sql.Open("sqlite",filepath.Join(d,"outfitter.db")+"?_journal_mode=WAL&_busy_timeout=5000");if err!=nil{return nil,err}
db.Exec(`CREATE TABLE IF NOT EXISTS new_hires(id TEXT PRIMARY KEY,name TEXT NOT NULL,email TEXT DEFAULT '',department TEXT DEFAULT '',start_date TEXT DEFAULT '',manager TEXT DEFAULT '',progress INTEGER DEFAULT 0,status TEXT DEFAULT 'pending',notes TEXT DEFAULT '',created_at TEXT DEFAULT(datetime('now')))`)
return &DB{db:db},nil}
func(d *DB)Close()error{return d.db.Close()}
func genID()string{return fmt.Sprintf("%d",time.Now().UnixNano())}
func now()string{return time.Now().UTC().Format(time.RFC3339)}
func(d *DB)Create(e *NewHire)error{e.ID=genID();e.CreatedAt=now();_,err:=d.db.Exec(`INSERT INTO new_hires(id,name,email,department,start_date,manager,progress,status,notes,created_at)VALUES(?,?,?,?,?,?,?,?,?,?)`,e.ID,e.Name,e.Email,e.Department,e.StartDate,e.Manager,e.Progress,e.Status,e.Notes,e.CreatedAt);return err}
func(d *DB)Get(id string)*NewHire{var e NewHire;if d.db.QueryRow(`SELECT id,name,email,department,start_date,manager,progress,status,notes,created_at FROM new_hires WHERE id=?`,id).Scan(&e.ID,&e.Name,&e.Email,&e.Department,&e.StartDate,&e.Manager,&e.Progress,&e.Status,&e.Notes,&e.CreatedAt)!=nil{return nil};return &e}
func(d *DB)List()[]NewHire{rows,_:=d.db.Query(`SELECT id,name,email,department,start_date,manager,progress,status,notes,created_at FROM new_hires ORDER BY created_at DESC`);if rows==nil{return nil};defer rows.Close();var o []NewHire;for rows.Next(){var e NewHire;rows.Scan(&e.ID,&e.Name,&e.Email,&e.Department,&e.StartDate,&e.Manager,&e.Progress,&e.Status,&e.Notes,&e.CreatedAt);o=append(o,e)};return o}
func(d *DB)Update(e *NewHire)error{_,err:=d.db.Exec(`UPDATE new_hires SET name=?,email=?,department=?,start_date=?,manager=?,progress=?,status=?,notes=? WHERE id=?`,e.Name,e.Email,e.Department,e.StartDate,e.Manager,e.Progress,e.Status,e.Notes,e.ID);return err}
func(d *DB)Delete(id string)error{_,err:=d.db.Exec(`DELETE FROM new_hires WHERE id=?`,id);return err}
func(d *DB)Count()int{var n int;d.db.QueryRow(`SELECT COUNT(*) FROM new_hires`).Scan(&n);return n}

func(d *DB)Search(q string, filters map[string]string)[]NewHire{
    where:="1=1"
    args:=[]any{}
    if q!=""{
        where+=" AND (name LIKE ? OR email LIKE ?)"
        args=append(args,"%"+q+"%");args=append(args,"%"+q+"%");
    }
    if v,ok:=filters["status"];ok&&v!=""{where+=" AND status=?";args=append(args,v)}
    rows,_:=d.db.Query(`SELECT id,name,email,department,start_date,manager,progress,status,notes,created_at FROM new_hires WHERE `+where+` ORDER BY created_at DESC`,args...)
    if rows==nil{return nil};defer rows.Close()
    var o []NewHire;for rows.Next(){var e NewHire;rows.Scan(&e.ID,&e.Name,&e.Email,&e.Department,&e.StartDate,&e.Manager,&e.Progress,&e.Status,&e.Notes,&e.CreatedAt);o=append(o,e)};return o
}

func(d *DB)Stats()map[string]any{
    m:=map[string]any{"total":d.Count()}
    rows,_:=d.db.Query(`SELECT status,COUNT(*) FROM new_hires GROUP BY status`)
    if rows!=nil{defer rows.Close();by:=map[string]int{};for rows.Next(){var s string;var c int;rows.Scan(&s,&c);by[s]=c};m["by_status"]=by}
    return m
}
