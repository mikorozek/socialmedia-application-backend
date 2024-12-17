package models;

type User struct {
	ID				uint		`gorm:"primarykey"`
	Username		string		`gorm:"unique;not null"`
	Email			string 		`gorm:"unique;not null"`
	PasswordHash	string		`gorm:"not null"`
};
