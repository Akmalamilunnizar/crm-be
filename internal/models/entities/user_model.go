package entities

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Enums
type UserRole string

const (
	RoleAdmin      UserRole = "ADMIN"
	RoleTechnician UserRole = "TECHNICIAN"
	RoleFinance    UserRole = "FINANCE"
)

// Accounts model
type Accounts struct {
	ID        string    `json:"id" gorm:"column:id;primaryKey"`
	Name      string    `json:"name" gorm:"column:name"`
	Saldo     int64     `json:"saldo" gorm:"column:saldo;default:0"`
	CreatedAt time.Time `json:"createdAt" gorm:"column:createdAt;default:current_timestamp"`
	UpdatedAt time.Time `json:"updatedAt" gorm:"column:updatedAt;"`
}

func (c *Accounts) TableName() string {
	return "accounts"
}

func (u *Accounts) BeforeCreate(tx *gorm.DB) error {
	if u.ID == "" {
		u.ID = uuid.New().String()
	}
	return nil
}

type CustomerInstallation struct {
	ID           string    `gorm:"column:id;type:varchar;primaryKey" json:"id"`
	CustomerID   string    `gorm:"column:customer_id;type:varchar;index:archive_installation_ibfk_1" json:"customer_id,omitempty"`
	TechnicianID string    `gorm:"column:technician_id;type:varchar;index:technician_id" json:"technician_id"`
	Description  string    `gorm:"column:description;type:text" json:"description,omitempty"`
	Date         string    `gorm:"column:date;type:date" json:"date"`
	CreatedAt    time.Time `gorm:"column:createdAt;default:current_timestamp" json:"createdAt"`
	UpdatedAt    time.Time `gorm:"column:updatedAt" json:"updatedAt"`
	Customer     Customer  `gorm:"foreignKey:CustomerID;references:id;constraint:OnDelete:RESTRICT,OnUpdate:RESTRICT" json:"customer,omitempty"`
	Technician   User      `gorm:"foreignKey:TechnicianID;references:id;constraint:OnUpdate:RESTRICT" json:"technician,omitempty"`
	Images       []Image   `gorm:"foreignKey:archive_installation_id;constraint:OnUpdate:RESTRICT" json:"images,omitempty"`
}

func (c *CustomerInstallation) TableName() string {
	return "customer_installations"
}

func (u *CustomerInstallation) BeforeCreate(tx *gorm.DB) error {
	if u.ID == "" {
		u.ID = uuid.New().String()
	}
	return nil
}

// Assets model
type Asset struct {
	ID           string        `json:"id" gorm:"primaryKey"`
	Type         string        `json:"type" validate:"required"`
	Brand        string        `json:"brand" validate:"required"`
	Model        string        `json:"model" validate:"required"`
	SerialNumber string        `json:"serial_number" validate:"required"`
	MacAddress   string        `json:"mac_address" validate:"required"`
	Date         string        `json:"date" validate:"required"`
	Site         string        `json:"site"`
	Quantity     int64         `json:"quantity" validate:"required"`
	Status       string        `json:"status" `
	Price        int64         `json:"price" validate:"required"`
	Description  string        `json:"description"`
	StatusInOut  string        `json:"status_in_out" validate:"required"`
	CreatedAt    time.Time     `json:"createdAt" gorm:"column:createdAt;default:current_timestamp"`
	UpdatedAt    time.Time     `json:"updatedAt" gorm:"column:updatedAt;"`
	ReportAssets *ReportAssets `json:"report_assets" gorm:"foreignKey:ID"`
}

func (c *Asset) TableName() string {
	return "assets"
}

func (u *Asset) BeforeCreate(tx *gorm.DB) error {
	if u.ID == "" {
		u.ID = uuid.New().String()
	}
	return nil
}

// Company model
type Company struct {
	ID          string     `json:"id" gorm:"primaryKey"`
	Name        string     `json:"name"`
	URL         string     `json:"url"`
	Email       string     `json:"email"`
	Phone       string     `json:"phone"`
	LogoURL     string     `json:"logo_url"`
	Description string     `json:"description"`
	Npwp        string     `json:"npwp"`
	Address     string     `json:"address"`
	CreatedAt   time.Time  `json:"createdAt" gorm:"column:createdAt; default:current_timestamp"`
	UpdatedAt   time.Time  `json:"updatedAt"  gorm:"column:updatedAt" `
	Customers   []Customer `json:"customer" gorm:"foreignKey:company_id"`
}

func (c *Company) TableName() string {
	return "company"
}

func (u *Company) BeforeCreate(tx *gorm.DB) error {
	if u.ID == "" {
		u.ID = uuid.New().String()
	}
	return nil
}

// Customer model
type Customer struct {
	ID               string    `json:"id" gorm:"primaryKey"`
	Address          string    `gorm:"column:address" json:"address"`
	AreaID           string    `gorm:"column:area_id" json:"area_id"`
	Area             *Areas    `gorm:"foreignKey:AreaID" json:"area"`
	CardIdentition   string    `gorm:"column:card_identition" json:"card_identition"`
	CompanyID        string    `gorm:"column:company_id" json:"company_id"`
	Company          *Company  `gorm:"foreignKey:CompanyID" json:"company"`
	Email            string    `gorm:"column:email;unique" json:"email"`
	Gender           string    `gorm:"column:gender" json:"gender"`
	ProductID        string    `gorm:"column:product_id" json:"product_id"`
	Product          *Products `gorm:"foreignKey:ProductID" json:"product"`
	IPStatic         string    `gorm:"column:ip_static" json:"ip_static"`
	Job              string    `gorm:"column:job" json:"job"`
	Latitude         float64   `gorm:"column:latitude" json:"latitude"`
	Longitude        float64   `gorm:"column:longitude" json:"longitude"`
	MacAddress       string    `gorm:"column:mac_address" json:"mac_address"`
	Name             string    `gorm:"column:name" json:"name"`
	NoIdentition     int       `gorm:"column:no_identition" json:"no_identition"`
	Password         string    `gorm:"column:password" json:"password"`
	Phone            string    `gorm:"column:phone" json:"phone"`
	TypeOfService    string    `gorm:"column:type_of_service" json:"type_of_service"`
	CreatedAt        time.Time `gorm:"column:createdAt;autoCreateTime" json:"created_at"`
	UpdatedAt        time.Time `gorm:"column:updatedAt;autoUpdateTime" json:"updated_at"`
	InstallationDate time.Time `gorm:"column:installation_date;type:date" json:"installation_date"`
	NextPaymentDate  time.Time `gorm:"column:next_payment_date;type:date" json:"next_payment_date"`
}

func (u *Customer) TableName() string {
	return "customer"
}
func (u *Customer) BeforeCreate(tx *gorm.DB) error {
	if u.ID == "" {
		u.ID = uuid.New().String()
		u.InstallationDate = time.Now()
		u.NextPaymentDate = time.Now().AddDate(0, 1, 0)
	}
	return nil
}

// Device model
type Device struct {
	ID        string    `json:"id" gorm:"primaryKey"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"createdAt" gorm:"default:current_timestamp"`
	UpdatedAt time.Time `json:"updatedAt"`
}

// Groups model
type Areas struct {
	ID              string     `json:"id" gorm:"primaryKey"`
	NameCity        string     `json:"name_city"`
	NameSubdistrict string     `json:"name_subdistrict"`
	NameVillage     string     `json:"name_village"`
	CreatedAt       time.Time  `json:"createdAt" gorm:"column:createdAt;default:current_timestamp"`
	UpdatedAt       time.Time  `json:"updatedAt" gorm:"column:updatedAt;"`
	Customers       []Customer `json:"customer" gorm:"foreignKey:area_id"`
}

func (u *Areas) TableName() string {
	return "areas"
}
func (u *Areas) BeforeCreate(tx *gorm.DB) error {
	if u.ID == "" {
		u.ID = uuid.New().String()
	}
	return nil
}

// Log model
type Log struct {
	ID        string    `json:"id" gorm:"primaryKey"`
	UserID    string    `json:"user_id" gorm:"index:user_id"`
	Action    string    `json:"action"`
	CreatedAt time.Time `json:"createdAt" gorm:"default:current_timestamp"`
	UpdatedAt time.Time `json:"updatedAt"`
	User      User      `json:"user" gorm:"foreignKey:UserID;constraint:OnUpdate:RESTRICT"`
}

// Products model
type Products struct {
	ID          string    `json:"id" gorm:"primaryKey"`
	Name        string    `json:"name"`
	Price       int64     `json:"price"` // BigInt mapped to int64
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"createdAt" gorm:"default:current_timestamp; column:createdAt"`
	UpdatedAt   time.Time `json:"updatedAt" gorm:"column:updatedAt" omitempty"`
}

func (u *Products) TableName() string {
	return "products"
}
func (u *Products) BeforeCreate(tx *gorm.DB) error {
	if u.ID == "" {
		u.ID = uuid.New().String()
	}
	return nil
}

// ReportAssets model
type ReportAssets struct {
	ID          string    `json:"id" gorm:"primaryKey;index:report_assets_ibfk_1"`
	Description string    `json:"description"`
	Quantity    int64     `json:"quantity"` // BigInt mapped to int64
	CreatedAt   time.Time `json:"createdAt" gorm:"default:current_timestamp"`
	UpdatedAt   time.Time `json:"updatedAt"`
	Assets      []Asset   `json:"assets" gorm:"foreignKey:ID;constraint:OnUpdate:RESTRICT"`
}

// ReportCash model
type ReportCash struct {
	ID          string    `json:"id" gorm:"primaryKey"`
	CreatedAt   time.Time `json:"createdAt" gorm:"default:current_timestamp"`
	UpdatedAt   time.Time `json:"updatedAt"`
	Credit      int64     `json:"credit"` // BigInt mapped to int64
	Debit       int64     `json:"debit"`  // BigInt mapped to int64
	Description string    `json:"description"`
}

// Transactions model
// TransactionsTypeCash represents the type of cash transaction
type TransactionsTypeCash string

// Constants for TransactionsTypeCash
const (
	TransactionsTypeCashInternet TransactionsTypeCash = "internet"
	TransactionsTypeCashCashFlow TransactionsTypeCash = "cash_flow"
)

// TransactionsTypeInOut represents the direction of the transaction
type TransactionsTypeInOut string

// Constants for TransactionsTypeInOut
const (
	TransactionsTypeInOutIn  TransactionsTypeInOut = "debit"
	TransactionsTypeInOutOut TransactionsTypeInOut = "credit"
)

// TransactionsType represents the general transaction type
type TransactionsType string

type Transaction struct {
	ID          string                `json:"id" gorm:"column:id;primaryKey;type:varchar"`
	AccountID   string                `json:"account_id" gorm:"column:account_id;index:transactions_account_id_fkey;type:varchar"`
	InvoiceID   string                `json:"invoice_id" gorm:"column:invoice_id;index:transactions_invoice_id_fkey;type:varchar"`
	TypeCash    TransactionsTypeCash  `json:"type_cash" gorm:"column:type_cash;type:varchar"`
	TypeInOut   TransactionsTypeInOut `json:"type_in_out" gorm:"column:type_in_out;type:varchar"`
	Date        string                `json:"date" gorm:"column:date;type:datetime"`
	Description string                `json:"description" gorm:"column:description;type:varchar"`
	Amount      int64                 `json:"amount" gorm:"column:amount;type:bigint"`
	Category    string                `json:"category" gorm:"column:category;type:varchar"`
	Method      string                `json:"method" gorm:"column:method;type:varchar"`
	CreatedAt   time.Time             `json:"createdAt" gorm:"column:createdAt;default:current_timestamp"`
	UpdatedAt   time.Time             `json:"updatedAt" gorm:"column:updatedAt;type:datetime"`
	Account     Accounts              `json:"account" gorm:"foreignKey:AccountID;constraint:OnUpdate:RESTRICT"`
}

func (u *Transaction) TableName() string {
	return "transactions"
}
func (u *Transaction) BeforeCreate(tx *gorm.DB) error {
	if u.ID == "" {
		u.ID = uuid.New().String()
		u.Date = time.Now().Format("2006-01-02 15:04:05")
	}
	return nil
}

// Transfers model
type Transfers struct {
	ID            string    `json:"id" gorm:"primaryKey"`
	FromAccountID string    `json:"from_account_id" gorm:"index:transfers_from_account_id_fkey"`
	ToAccountID   string    `json:"to_account_id" gorm:"index:transfers_to_account_id_fkey"`
	Date          time.Time `json:"date"`
	Description   string    `json:"description"`
	Amount        int64     `json:"amount"` // BigInt mapped to int64
	Tags          string    `json:"tags"`
	CreatedAt     time.Time `json:"createdAt" gorm:"default:current_timestamp"`
	UpdatedAt     time.Time `json:"updatedAt"`
	FromAccount   Accounts  `json:"accounts_transfers_from_account_idToaccounts" gorm:"foreignKey:FromAccountID;constraint:OnUpdate:RESTRICT"`
	ToAccount     Accounts  `json:"accounts_transfers_to_account_idToaccounts" gorm:"foreignKey:ToAccountID;constraint:OnUpdate:RESTRICT"`
}

// User model
type User struct {
	ID        string    `json:"id" gorm:"primaryKey"`
	Email     string    `json:"email" gorm:"unique"`
	Name      string    `json:"name" gorm:"default:null"`
	Password  string    `json:"password"`
	RoleId    string    `gorm:"column:role_id"`
	Role      Role      `gorm:"foreignKey:RoleId"`
	Token     string    `json:"token" gorm:"default:null"`
	Phone     string    `json:"phone"`
	CreatedAt time.Time `json:"createdAt" gorm:"default:current_timestamp; column:createdAt"`
	UpdatedAt time.Time `json:"updatedAt" gorm:"column:updatedAt" omitempty"`
	Log       []Log     `json:"log" gorm:"foreignKey:UserID"`
}

func (u *User) TableName() string {
	return "users"
}
func (u *User) BeforeCreate(tx *gorm.DB) error {
	if u.ID == "" {
		u.ID = uuid.New().String()
	}
	return nil
}

type Role struct {
	ID        string    `json:"id" gorm:"primaryKey"`
	Name      string    `json:"name" gorm:"default:null"`
	CreatedAt time.Time `json:"createdAt" gorm:"default:current_timestamp; column:createdAt"`
	UpdatedAt time.Time `json:"updatedAt" gorm:"column:updatedAt" omitempty"`
}

func (u *Role) TableName() string {
	return "roles"
}
func (u *Role) BeforeCreate(tx *gorm.DB) error {
	if u.ID == "" {
		u.ID = uuid.New().String()
	}
	return nil
}

type Image struct {
	ID                    string     `gorm:"column:id;type:varchar;primaryKey" json:"id"`
	File                  string     `gorm:"column:file;type:varchar;not null" json:"file"`
	FullPath              string     `gorm:"column:full_path;type:varchar;not null" json:"full_path"`
	ArchiveInstallationId string     `gorm:"column:archive_installation_id;type:varchar;not null" json:"archive_installation_id"`
	CreatedAt             time.Time  `gorm:"column:createdAt;not null;default:now()" json:"createdAt"`
	UpdatedAt             *time.Time `gorm:"column:updatedAt;default:null" json:"updatedAt"`
}

func (u *Image) TableName() string {
	return "images"
}

func (u *Image) BeforeCreate(tx *gorm.DB) error {
	if u.ID == "" {
		u.ID = uuid.New().String()
	}
	return nil
}

// TransactionsTypeInOut represents the direction of the transaction
type InvoiceStatus string

// Constants for TransactionsTypeInOut
const (
	InvoiceStatusPaid    InvoiceStatus = "paid"
	InvoiceStatusUnpaid  InvoiceStatus = "unpaid"
	InvoiceStatusPending InvoiceStatus = "pending"
)

type Invoice struct {
	ID           string         `gorm:"column:id;type:varchar;primaryKey" json:"id"`
	Amount       int64          `gorm:"column:amount;type:int;not null" json:"amount"`
	CustomerID   string         `gorm:"column:customer_id;type:varchar;not null" json:"customer_id"`
	Customer     Customer       `gorm:"foreignKey:CustomerID;references:id;constraint:OnUpdate:RESTRICT" json:"customer"`
	Link         string         `gorm:"column:link;type:varchar;not null" json:"link"`
	Status       InvoiceStatus  `gorm:"column:status;type:varchar;not null" json:"status"`
	CreatedAt    time.Time      `gorm:"column:createdAt;default:current_timestamp" json:"created_at"`
	UpdatedAt    time.Time      `gorm:"column:updatedAt;not null" json:"updated_at"`
	InvoiceItems []InvoiceItems `gorm:"foreignKey:invoices_id;constraint:OnUpdate:RESTRICT" json:"invoice_items"`
	Transaction  Transaction    `gorm:"foreignKey:invoice_id;constraint:OnUpdate:RESTRICT" json:"transaction""`
}

func (Invoice) TableName() string {
	return "invoices"
}
func (u *Invoice) BeforeCreate(tx *gorm.DB) error {
	if u.ID == "" {
		u.ID = uuid.New().String()
		u.Status = InvoiceStatusUnpaid
	}
	return nil
}

type InvoiceItems struct {
	ID        string    `gorm:"column:id;type:varchar;primaryKey" json:"id"`
	Name      string    `gorm:"column:name;type:int;not null" json:"name"`
	Qty       int64     `gorm:"column:qty;type:varchar;not null" json:"qty"`
	Price     int64     `gorm:"column:price;type:varchar;not null" json:"price"`
	Total     int64     `gorm:"column:total;type:varchar;not null" json:"total"`
	InvoiceID string    `gorm:"column:invoices_id;type:varchar;not null" json:"invoice_id"`
	Invoice   Invoice   `gorm:"foreignKey:InvoiceID;references:id" json:"invoice"`
	CreatedAt time.Time `gorm:"column:createdAt;default:current_timestamp" json:"created_at"`
	UpdatedAt time.Time `gorm:"column:updatedAt;not null" json:"updated_at"`
}

func (InvoiceItems) TableName() string {
	return "invoice_items"
}
func (u *InvoiceItems) BeforeCreate(tx *gorm.DB) error {
	if u.ID == "" {
		u.ID = uuid.New().String()
	}
	return nil
}
