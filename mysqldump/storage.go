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

  usersCnt int
  salesCnt int
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

  archive := zip.NewWriter(zipFile)

  usersFile, err := archive.Create("users.csv")
  if err != nil {
    log.Fatal(err);
  }

  userCsv := csv.NewWriter(usersFile);
  for u := range s.userCh {
    s.usersCnt++;
    err := userCsv.Write([]string{strconv.Itoa(u.id), u.name});
    if err != nil { log.Fatal(err); }
  }

  userCsv.Flush();

  s.doneCh <- true;

  salesFile, err := archive.Create("sales.csv")
  if err != nil {
    log.Fatal(err);
  }

  salesCsv := csv.NewWriter(salesFile);
  for sl := range s.saleCh {
    s.salesCnt++;
    err := salesCsv.Write([]string{strconv.Itoa(sl.id), strconv.Itoa(sl.user_id), strconv.FormatFloat(sl.order_amount, 'f', 3, 64)});
    if err != nil { log.Fatal(err); }
  }

  salesCsv.Flush();

  archive.Close();
  if err != nil {
    log.Fatal(err);
  }

  err = zipFile.Close();
  if err != nil {
    log.Fatal(err);
  }

  log.Println("Done storage. Users:", s.usersCnt, "Sales:", s.salesCnt);

  s.doneCh <- true;
}

func (s *storage) StoreUser (u *user) {
  s.userCh <- u
}


func (s *storage) StoreSale (sl *sale) {
  s.saleCh <- sl
}

func (s *storage) CloseUsers() {
  close(s.userCh)
  <-s.doneCh
}

func (s *storage) CloseSales() {
  close(s.saleCh)
  <-s.doneCh
}
