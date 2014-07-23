package mysqldump

import (
  "sync"
  "database/sql"
  _ "github.com/go-sql-driver/mysql"
  "log"
)

type Dumper struct {
  storage *storage
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


}

