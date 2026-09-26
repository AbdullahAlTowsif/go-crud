package main

type User struct {
	Id    int    `json:"id"`
	Name  string `json:"name"`
	Age   int    `json:"age"`
	Email string `json:"email"`
}



var users = []User{
	{
		Id:    1,
		Name:  "Towsif",
		Age:   23,
		Email: "towsif@gmail.com",
	},
	{
		Id:    2,
		Name:  "Abdullah",
		Age:   24,
		Email: "abdullah@gmail.com",
	},
}
