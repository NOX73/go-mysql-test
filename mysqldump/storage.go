package mysqldump

type storage struct {
}

func newStorage () *storage {
  return &storage{}
}

func (s *storage) Close() {
}
