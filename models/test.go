package Models
import "time"

type TestModel struct{
	Id 			uint 		`json:"id"`
	Name 		string  	`json:"name"`
	Passcode 	string 		`json:"passcode"`
	Created_at  time.Time 	`json:"created_at"`
	Updated_at  time.Time 	`json:"updated_at"`

}