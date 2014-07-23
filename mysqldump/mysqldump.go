package mysqldump

func DumpFromServers(addrs []string) {
  dumper := NewDumper();
  dumper.ProcessServers(addrs);
}
