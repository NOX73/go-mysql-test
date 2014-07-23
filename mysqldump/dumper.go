package mysqldump

import (
  "sync"
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
  var wg sync.WaitGroup

  var dbs = make([]*db, 0, len(addrs)) ;

  for i:= 0; i < len(addrs); i++ {
    db := newDB(addrs[i], d.storage)
    dbs = append(dbs, db);

    wg.Add(1);

    go func(){
      defer wg.Done()
      db.EstablishConnection();
    }()

  }

  wg.Wait();

  for i:= 0; i < len(dbs); i++ {
    wg.Add(1);
    db := dbs[i]

    go func(){
      defer wg.Done()
      db.DumpUsers();
    }()

  }

  wg.Wait();
  d.storage.CloseUsers();

  for i:= 0; i < len(dbs); i++ {
    wg.Add(1);
    db := dbs[i]

    go func(){
      defer wg.Done()
      db.DumpSales();
    }()

  }

  wg.Wait();

  d.storage.CloseSales();
}

