package check

type User struct {
	Id            int64   `db:"id"`
	CompanyName   string  `db:"company_name" json:"company_name"`
	City          string  `db:"city" json:"city"`
	State         string  `db:"state" json:"state"`
	Zip           string  `db:"zip" json:"zip"`
	Email         string  `db:"email" json:"email"`
	Phone         int64   `db:"phone" json:"phone"`
	AverageGrade  float64 `db: "average_grade" json:"average_grade"`
	TotalJobCount int     `db:"jobs_count" json:"jobs_count"`
}
