package mysqldump

import (
"log"
  "database/sql"
  _ "github.com/go-sql-driver/mysql"
    )

const SELECT_LIMIT = 100;

type db struct {
  addr string
  storage *storage
  db *sql.DB
}

func newDB (addr string, storage *storage) *db {
  return &db{ addr: addr, storage: storage }
}

func (d *db) EstablishConnection () {
  log.Println("Try establish connection with ", d.addr)

  db, err := sql.Open("mysql", d.addr);

  if err != nil {
    log.Fatal(err)
  }

  err = db.Ping()
  if err != nil {
    log.Fatal(err)
  }

  d.db = db
}

func (d *db) DumpUsers () {
  count := d.getCount("users");
  for i := 0; i < count; i += SELECT_LIMIT {
    d.processUserBatch(i);
  }
}

func (d *db) DumpSales () {
  count := d.getCount("sales");
  for i := 0; i < count; i += SELECT_LIMIT {
    d.processSalesBatch(i);
  }
}

func (d *db) getCount (table string) int {
  var count int;

  countQuery, err := d.db.Query("select COUNT(*) from " + table)

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

func (d *db) processUserBatch(start int) {
  rows, err := d.db.Query("select * from users limit ?, ?", start, SELECT_LIMIT);

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

func (d *db) processSalesBatch (start int) {
  rows, err := d.db.Query("select * from sales limit ?, ?", start, SELECT_LIMIT);

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
