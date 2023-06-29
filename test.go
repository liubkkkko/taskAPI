package main

import (
	"time"

	"github.com/jinzhu/gorm"
)

type Organization struct {
	ID       uint       `gorm:"column:id;primaryKey"`
	Name     string     `gorm:"column:name"`
	Code     string     `gorm:"column:code"`
	Users    []*User    `gorm:"many2many:user_organizations;"`
	Licenses []*License `gorm:"foreignKey:organization_id"`
	gorm.Model
}

type User struct {
	ID            uint            `gorm:"column:id;primaryKey"`
	Name          string          `gorm:"column:name"`
	Token         string          `gorm:"column:token"`
	Admin         bool            `gorm:"column:admin"`
	Organizations []*Organization `gorm:"many2many:user_organizations;"`
	Email         string          `gorm:"column:email;unique"`
	Password      []byte          `gorm:"column:password" json:"-"`
	Licenses      []*License      `gorm:"many2many:user_licenses;"`
	gorm.Model
}

type License struct {
	ID             uint      `gorm:"column:id;primaryKey"`
	Endpoint       string    `gorm:"column:endpoint"`
	ServiceName    string    `gorm:"column:service_name"`
	Expiry         string    `gorm:"column:expiry"`
	CurrentUser    int       `gorm:"column:current_user"`
	OrganizationID uint      `gorm:"column:organization_id;foreignKey:ID;default:null"`
	PaymentDue     int       `gorm:"column:payment_due"`
	LastActive     time.Time `gorm:"column:last_active"`
	Users          []*User   `gorm:"many2many:user_licenses;"`
	gorm.Model
}

func test() {
	// dsn := "host=localhost user=postgres password=postgres dbname=auth port=6789 sslmode=disable"
	// AuthDB, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	// if err != nil {
	// 	fmt.Println("Failed to connect to Auth db")
	// }

	// AuthDB.AutoMigrate(&Organization{})
	// AuthDB.AutoMigrate(&User{})
	// AuthDB.AutoMigrate(&License{})

	// password, _ := bcrypt.GenerateFromPassword([]byte("test"), 12)
	// org := Organization{Name: "TestOrg"}

	// user := User{Name: "Admin", Email: "test", Password: password, Admin: true}

	// licenses := []*License{{Endpoint: "/python-test/", ServiceName: "generic-webservice1"},
	// 	{Endpoint: "/generic-react1/", ServiceName: "generic-react1"}}

	// AuthDB.Create(&org)
	// AuthDB.Create(&user)
	// AuthDB.Create(&licenses)

	// // ONE TO MANY tests
	// // This works for updating licenses
	// for _, license := range licenses {
	// 	license.OrganizationID = org.ID
	// 	AuthDB.Save(&license)
	// }
	// // This does not work. Don't understand why.
	// org.Licenses = licenses
	// AuthDB.Session(&gorm.Session{FullSaveAssociations: true}).Updates(&org)
	// // This does not work. Don't understand why.
	// AuthDB.Model(&org).Association("Licenses").Append(licenses)

	// // MANY TO MANY tests
	// // This does not work.
	// user.Licenses = licenses
	// AuthDB.Session(&gorm.Session{FullSaveAssociations: true}).Updates(&user)
	// // This does not work.
	// AuthDB.Model(&user).Association("Licenses").Append(licenses)

	// // Lets try many to many without array
	// // Does not work
	// user.Organizations = append(user.Organizations, &org)
	// AuthDB.Session(&gorm.Session{FullSaveAssociations: true}).Updates(&user)
	// // Does not work.
	// AuthDB.Model(user).Association("Licenses").Append(&org)
}
