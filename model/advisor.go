package model

type Advisor struct {
	Name           *string `structs:"name" json:"name" `
	Phone          *string `structs:"phone" json:"phone"`
	WorkExperience *int    `structs:"work_experience" json:"workExperience"`
	Bio            *string `structs:"bio" json:"bio" `
	About          *string `structs:"about" json:"about"`
}
