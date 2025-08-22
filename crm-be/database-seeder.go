package main

import (
	"log"
	"skripsi-be/internal/config/database"
	"skripsi-be/internal/models/entities"

	"github.com/google/uuid"
	"github.com/joho/godotenv"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

func main() {
	// Load environment variables
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	// Get database connection
	db := database.GetDB()

	log.Println("Starting database seeding...")

	// 1. Seed Roles
	log.Println("1. Seeding roles...")
	roles := seedRoles(db)

	// 2. Seed Users
	log.Println("2. Seeding users...")
	seedUsers(db, roles)

	// 3. Seed Areas
	log.Println("3. Seeding areas...")
	areas := seedAreas(db)

	// 4. Seed Products
	log.Println("4. Seeding products...")
	products := seedProducts(db)

	// 5. Seed Companies
	log.Println("5. Seeding companies...")
	companies := seedCompanies(db)

	// 6. Seed Customers
	log.Println("6. Seeding customers...")
	seedCustomers(db, companies, areas, products)

	// 5. Seed Trouble Types
	log.Println("5. Seeding trouble types...")
	seedTroubleTypes(db)

	log.Println("Database seeding completed successfully!")
}

func seedRoles(db *gorm.DB) map[string]entities.Role {
	roles := []entities.Role{
		{ID: uuid.New().String(), Name: "ADMIN"},
		{ID: uuid.New().String(), Name: "CUSTOMER_SERVICE"},
		{ID: uuid.New().String(), Name: "NOC"},
		{ID: uuid.New().String(), Name: "TECHNICIAN"},
		{ID: uuid.New().String(), Name: "FINANCE"},
	}

	roleMap := make(map[string]entities.Role)

	for _, role := range roles {
		var existingRole entities.Role
		if err := db.Where("name = ?", role.Name).First(&existingRole).Error; err != nil {
			// Role doesn't exist, create it
			if err := db.Create(&role).Error; err != nil {
				log.Printf("Error creating role %s: %v", role.Name, err)
			} else {
				log.Printf("Created role: %s (ID: %s)", role.Name, role.ID)
				roleMap[role.Name] = role
			}
		} else {
			log.Printf("Role %s already exists (ID: %s)", role.Name, existingRole.ID)
			roleMap[role.Name] = existingRole
		}
	}

	return roleMap
}

func seedAreas(db *gorm.DB) map[string]entities.Areas {
	areas := []entities.Areas{
		{
			NameCity:        "Jakarta Pusat",
			NameSubdistrict: "Menteng",
			NameVillage:     "Menteng",
		},
		{
			NameCity:        "Jakarta Selatan",
			NameSubdistrict: "Kebayoran Baru",
			NameVillage:     "Senayan",
		},
	}

	areaMap := make(map[string]entities.Areas)

	for _, area := range areas {
		var existingArea entities.Areas
		if err := db.Where("name_city = ? AND name_subdistrict = ?", area.NameCity, area.NameSubdistrict).First(&existingArea).Error; err != nil {
			// Area doesn't exist, create it
			if err := db.Create(&area).Error; err != nil {
				log.Printf("Error creating area %s: %v", area.NameCity, err)
			} else {
				log.Printf("Created area: %s", area.NameCity)
				areaMap[area.NameCity] = area
			}
		} else {
			log.Printf("Area %s already exists", area.NameCity)
			areaMap[area.NameCity] = existingArea
		}
	}

	return areaMap
}

func seedProducts(db *gorm.DB) map[string]entities.Products {
	products := []entities.Products{
		{
			Name:        "Internet 100 Mbps",
			Price:       500000,
			Description: "High-speed internet connection with 100 Mbps download speed",
		},
		{
			Name:        "Internet 50 Mbps",
			Price:       300000,
			Description: "Standard internet connection with 50 Mbps download speed",
		},
		{
			Name:        "Internet 25 Mbps",
			Price:       200000,
			Description: "Basic internet connection with 25 Mbps download speed",
		},
	}

	productMap := make(map[string]entities.Products)

	for _, product := range products {
		var existingProduct entities.Products
		if err := db.Where("name = ?", product.Name).First(&existingProduct).Error; err != nil {
			// Product doesn't exist, create it
			if err := db.Create(&product).Error; err != nil {
				log.Printf("Error creating product %s: %v", product.Name, err)
			} else {
				log.Printf("Created product: %s", product.Name)
				productMap[product.Name] = product
			}
		} else {
			log.Printf("Product %s already exists", product.Name)
			productMap[product.Name] = existingProduct
		}
	}

	return productMap
}

func seedCompanies(db *gorm.DB) map[string]entities.Company {
	companies := []entities.Company{
		{
			Name:        "PT Maju Bersama",
			URL:         "https://majubersama.com",
			Email:       "info@majubersama.com",
			Phone:       "021-5550123",
			Description: "Leading technology solutions provider",
			Address:     "Jl. Sudirman No. 123, Jakarta Pusat",
		},
		{
			Name:        "CV Sukses Mandiri",
			URL:         "https://suksesmandiri.co.id",
			Email:       "contact@suksesmandiri.co.id",
			Phone:       "021-5550456",
			Description: "Innovative business solutions",
			Address:     "Jl. Thamrin No. 45, Jakarta Pusat",
		},
		{
			Name:        "UD Makmur Jaya",
			URL:         "https://makmurjaya.com",
			Email:       "admin@makmurjaya.com",
			Phone:       "021-5550789",
			Description: "Reliable service provider",
			Address:     "Jl. Gatot Subroto No. 67, Jakarta Selatan",
		},
	}

	companyMap := make(map[string]entities.Company)

	for _, company := range companies {
		var existingCompany entities.Company
		if err := db.Where("name = ?", company.Name).First(&existingCompany).Error; err != nil {
			// Company doesn't exist, create it
			if err := db.Create(&company).Error; err != nil {
				log.Printf("Error creating company %s: %v", company.Name, err)
			} else {
				log.Printf("Created company: %s", company.Name)
				companyMap[company.Name] = company
			}
		} else {
			log.Printf("Company %s already exists", company.Name)
			companyMap[company.Name] = existingCompany
		}
	}

	return companyMap
}

func seedUsers(db *gorm.DB, roles map[string]entities.Role) {
	// Hash password for all users
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte("password"), bcrypt.DefaultCost)
	if err != nil {
		log.Fatal("Error hashing password:", err)
	}

	users := []entities.User{
		{
			ID:       uuid.New().String(),
			Email:    "admin@lillyisp.com",
			Name:     "System Administrator",
			Password: string(hashedPassword),
			RoleId:   roles["ADMIN"].ID,
			Phone:    "081234567890",
		},
		{
			ID:       uuid.New().String(),
			Email:    "cs@lillyisp.com",
			Name:     "Customer Service Agent",
			Password: string(hashedPassword),
			RoleId:   roles["CUSTOMER_SERVICE"].ID,
			Phone:    "081234567891",
		},
		{
			ID:       uuid.New().String(),
			Email:    "noc@lillyisp.com",
			Name:     "NOC Engineer",
			Password: string(hashedPassword),
			RoleId:   roles["NOC"].ID,
			Phone:    "081234567892",
		},
		{
			ID:       uuid.New().String(),
			Email:    "tech@lillyisp.com",
			Name:     "Field Technician",
			Password: string(hashedPassword),
			RoleId:   roles["TECHNICIAN"].ID,
			Phone:    "081234567893",
		},
		{
			ID:       uuid.New().String(),
			Email:    "finance@lillyisp.com",
			Name:     "Finance Manager",
			Password: string(hashedPassword),
			RoleId:   roles["FINANCE"].ID,
			Phone:    "081234567894",
		},
	}

	for _, user := range users {
		var existingUser entities.User
		if err := db.Where("email = ?", user.Email).First(&existingUser).Error; err != nil {
			// User doesn't exist, create it
			if err := db.Create(&user).Error; err != nil {
				log.Printf("Error creating user %s: %v", user.Email, err)
			} else {
				log.Printf("Created user: %s (Password: password)", user.Email)
			}
		} else {
			log.Printf("User %s already exists", user.Email)
		}
	}
}

func seedCustomers(db *gorm.DB, companies map[string]entities.Company, areas map[string]entities.Areas, products map[string]entities.Products) {
	customers := []entities.Customer{
		{
			Name:      "PT Maju Bersama",
			Email:     "contact@majubersama.com",
			Phone:     "021-5550123",
			Address:   "Jl. Sudirman No. 123, Jakarta Pusat",
			CompanyID: companies["PT Maju Bersama"].ID,
			AreaID:    areas["Jakarta Pusat"].ID,
			ProductID: products["Internet 100 Mbps"].ID,
			Latitude:  -6.2088,
			Longitude: 106.8456,
		},
		{
			Name:      "CV Sukses Mandiri",
			Email:     "info@suksesmandiri.co.id",
			Phone:     "021-5550456",
			Address:   "Jl. Thamrin No. 45, Jakarta Pusat",
			CompanyID: companies["CV Sukses Mandiri"].ID,
			AreaID:    areas["Jakarta Pusat"].ID,
			ProductID: products["Internet 50 Mbps"].ID,
			Latitude:  -6.1865,
			Longitude: 106.8222,
		},
		{
			Name:      "UD Makmur Jaya",
			Email:     "admin@makmurjaya.com",
			Phone:     "021-5550789",
			Address:   "Jl. Gatot Subroto No. 67, Jakarta Selatan",
			CompanyID: companies["UD Makmur Jaya"].ID,
			AreaID:    areas["Jakarta Selatan"].ID,
			ProductID: products["Internet 25 Mbps"].ID,
			Latitude:  -6.2088,
			Longitude: 106.8233,
		},
	}

	for _, customer := range customers {
		var existingCustomer entities.Customer
		if err := db.Where("email = ?", customer.Email).First(&existingCustomer).Error; err != nil {
			// Customer doesn't exist, create it
			if err := db.Create(&customer).Error; err != nil {
				log.Printf("Error creating customer %s: %v", customer.Name, err)
			} else {
				log.Printf("Created customer: %s", customer.Name)
			}
		} else {
			log.Printf("Customer %s already exists", customer.Name)
		}
	}
}

func seedTroubleTypes(db *gorm.DB) {
	troubleTypes := []entities.TroubleTypeRow{
		{ID: uuid.New().String(), Name: stringPtr("Internet Down")},
		{ID: uuid.New().String(), Name: stringPtr("Slow Connection")},
		{ID: uuid.New().String(), Name: stringPtr("Modem Issue")},
		{ID: uuid.New().String(), Name: stringPtr("Cable Damage")},
		{ID: uuid.New().String(), Name: stringPtr("Power Outage")},
		{ID: uuid.New().String(), Name: stringPtr("Configuration Error")},
	}

	for _, tt := range troubleTypes {
		var existingTT entities.TroubleTypeRow
		if err := db.Table("trouble_type").Where("name = ?", *tt.Name).First(&existingTT).Error; err != nil {
			// Trouble type doesn't exist, create it
			if err := db.Table("trouble_type").Create(&tt).Error; err != nil {
				log.Printf("Error creating trouble type %s: %v", *tt.Name, err)
			} else {
				log.Printf("Created trouble type: %s", *tt.Name)
			}
		} else {
			log.Printf("Trouble type %s already exists", *tt.Name)
		}
	}
}

// Helper function to create string pointers
func stringPtr(s string) *string {
	return &s
}
