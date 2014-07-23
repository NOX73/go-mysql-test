package mysqldump

import (
  "sync"
  "database/sql"
  _ "github.com/go-sql-driver/mysql"
  "log"
)

const SELECT_LIMIT = 100;

type Dumper struct {
  storage *storage
}

type user struct {
  id int
  name string
}

type sale struct {
  id int
  user_id int
  order_amount float64
}

func NewDumper () *Dumper {
  d := &Dumper{
    storage: newStorage(),
  }
  return d;
}

func (d *Dumper) ProcessServers(addrs []string) {
  defer d.storage.Close()
  var wg sync.WaitGroup

  for i:= 0; i < len(addrs); i++ {
    var addr = addrs[i]

    wg.Add(1)

    go func() {
      defer wg.Done()
      d.processServer(addr, d.storage);
    } ()

  }

  wg.Wait();
}

func (d *Dumper) processServer (addr string, storage *storage) {

  log.Println("Try establish connection with ", addr)

  db, err := sql.Open("mysql", addr);

  if err != nil {
    log.Fatal(err)
  }

  err = db.Ping()
  if err != nil {
    log.Fatal(err)
  }

  defer db.Close()
  log.Println("Established connection with ", addr)

  count := d.getCount(db, "users");
  for i := 0; i < count; i += SELECT_LIMIT {
    d.processUserBatch(i, db);
  }

  count = d.getCount(db, "sales");
  for i := 0; i < count; i += SELECT_LIMIT {
    d.processSalesBatch(i, db);
  }

}

func (d *Dumper) getCount (db *sql.DB, table string) int {
  var count int;

  countQuery, err := db.Query("select COUNT(*) from " + table)

  if err != nil {
    log.Fatal(err)
  }

  defer countQuery.Close()

  countQuery.Next();
  err = countQuery.Scan(&count);

  if err != nil {
    log.Fatal(err)
  }

  return count;
}

func (d *Dumper) processUserBatch(start int, db *sql.DB) {
  rows, err := db.Query("select * from users limit ?, ?", start, start + SELECT_LIMIT);

  if err != nil {
    log.Fatal(err)
  }

  defer rows.Close()

  for rows.Next() {
    var u user;

    err := rows.Scan(&u.id, &u.name)

    if err != nil {
      log.Fatal(err)
    }

    d.storage.StoreUser(&u);
  }
}

func (d *Dumper) processSalesBatch (start int, db *sql.DB) {
  rows, err := db.Query("select * from sales limit ?, ?", start, start + SELECT_LIMIT);

  if err != nil {
    log.Fatal(err)
  }

  defer rows.Close()

  for rows.Next() {
    var s sale;

    err := rows.Scan(&s.id, &s.user_id, &s.order_amount)

    if err != nil {
      log.Fatal(err)
    }

    d.storage.StoreSale(&s);
  }

}
