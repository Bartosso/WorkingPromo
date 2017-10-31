package Utils

//Чекает есть ли ошибка, если есть то разводит панику
func CheckError(err error) {
	if err != nil {
		panic(err)
	}
}
