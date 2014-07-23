package mysqldump

type user struct {
  id int
  name string
}

type sale struct {
  id int
  user_id int
  order_amount float64
}
