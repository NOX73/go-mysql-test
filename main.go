package main

import "fmt"
import "os"
import dump "./mysqldump"

func main () {

  if len(os.Args) == 1 {
    fmt.Println("Set mysql credentials: [username[:password]@][protocol[(address)]]/dbname [username2[:password2]@][protocol2[(address2)]]/dbname2")
    return;
  }

  dump.DumpFromServers(os.Args[1:])
}
