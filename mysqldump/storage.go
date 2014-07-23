package mysqldump

import (
    "os"
    "encoding/csv"
    "log"
    "strconv"
    "archive/zip"
    )

type storage struct {
  userCh chan *user;
  saleCh chan *sale;
  doneCh chan bool;
}

func newStorage () *storage {
  s := &storage{
    userCh: make(chan *user, 5),
    saleCh: make(chan *sale, 5),
    doneCh: make(chan bool),
  }

  go s.start();

  return s
}

func (s *storage) start () {
  os.Mkdir("./archive", 0755);

  zipFile, err := os.OpenFile("./archive/archive.zip", os.O_WRONLY | os.O_CREATE, 0644)
  if err != nil {
    log.Fatal(err);
  }

  defer func(){
    err := zipFile.Close();

    if err != nil {
      log.Fatal(err);
    }
  }();

  archive := zip.NewWriter(zipFile)

  usersFile, err := archive.Create("users.csv")
  if err != nil {
    log.Fatal(err);
  }

  salesFile, err := archive.Create("sales.csv")
  if err != nil {
    log.Fatal(err);
  }

  userCsv := csv.NewWriter(usersFile);
  salesCsv := csv.NewWriter(salesFile);

  for {
    select {

    case u, ok := <- s.userCh:
      if ok {
        userCsv.Write([]string{strconv.Itoa(u.id), u.name});
      } else {
        s.userCh = nil;
      }
    case sl, ok := <- s.saleCh:
      if ok {
        salesCsv.Write([]string{strconv.Itoa(sl.id), strconv.Itoa(sl.user_id), strconv.FormatFloat(sl.order_amount, 'f', 3, 64)});
      } else {
        s.saleCh = nil;
      }
    }

    if s.saleCh == nil && s.userCh == nil {
      break;
    }

  }

  s.doneCh <- true;
}

func (s *storage) StoreUser (u *user) {
  s.userCh <- u
}


func (s *storage) StoreSale (sl *sale) {
  s.saleCh <- sl
}

func (s *storage) Close() {
  close(s.userCh)
  close(s.saleCh)
  <-s.doneCh

  log.Println("Dump to CSV finished");
}
