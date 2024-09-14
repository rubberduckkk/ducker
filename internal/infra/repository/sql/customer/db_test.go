package customer

import (
	"fmt"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"

	"github.com/rubberduckkk/ducker/internal/domain/customer"
	"github.com/rubberduckkk/ducker/internal/domain/customer/entity"
	"github.com/rubberduckkk/ducker/internal/domain/customer/valueobj"
)

func TestMain(m *testing.M) {
	code := m.Run()
	cleanupTestDB()
	os.Exit(code)
}

func prepTestDB(t *testing.T) *gorm.DB {
	if _, err := os.Stat("./sqlite.db"); err != nil {
		if os.IsNotExist(err) {
			if _, err = os.Create("sqlite.db"); err != nil {
				t.Fatalf("error creating sqlite.db: %s", err)
			}
		} else {
			t.Fatalf("get file info failed: %s", err)
		}
	}
	db, err := gorm.Open(sqlite.Open("sqlite.db"), &gorm.Config{})
	if err != nil {
		t.Fatalf("failed to prep test db: %s", err)
	}
	if db.Migrator().HasTable(Customer{}.TableName()) {
		if err := db.Exec(fmt.Sprintf("DELETE FROM %v", Customer{}.TableName())).Error; err != nil {
			t.Fatalf("failed to truncate table: %s", err)
		}
		if err := db.Exec("DELETE FROM SQLITE_SEQUENCE WHERE name=?", Customer{}.TableName()).Error; err != nil {
			t.Logf("failed to reset table metadata: %s", err)
		}
		return db
	}
	if err = db.AutoMigrate(&Customer{}); err != nil {
		t.Fatalf("failed to create test table: %s", err)
	}
	return db
}

func cleanupTestDB() {
	_ = os.Remove("./sqlite.db")
}

func Test_sqlRepository_Create(t *testing.T) {
	db := prepTestDB(t)
	type fields struct {
		db *gorm.DB
	}
	type args struct {
		customer *customer.Customer
	}
	tests := []struct {
		name      string
		fields    fields
		args      args
		wantErr   assert.ErrorAssertionFunc
		checkData func() bool
	}{
		{
			name: "case success",
			fields: fields{
				db: db,
			},
			args: args{
				customer: &customer.Customer{
					Person: &entity.Person{
						ID:   "123",
						Name: "Test Customer",
					},
					ContactInfo: valueobj.ContactInfo{AreaCode: "+86", PhoneNum: "1234567890"},
				},
			},
			wantErr: assert.NoError,
			checkData: func() bool {
				return db.First(&Customer{}, "id = ?", "123").Error == nil
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &sqlRepository{
				db: tt.fields.db,
			}
			tt.wantErr(t, s.Create(tt.args.customer), fmt.Sprintf("Create(%v)", tt.args.customer))
			if tt.checkData != nil {
				assert.True(t, tt.checkData())
			}
		})
	}
}

func Test_sqlRepository_Update(t *testing.T) {
	db := prepTestDB(t)
	type fields struct {
		db *gorm.DB
	}
	type args struct {
		customer *customer.Customer
	}
	tests := []struct {
		name      string
		fields    func() fields
		args      args
		wantErr   assert.ErrorAssertionFunc
		checkData func() bool
	}{
		{
			name: "case success",
			fields: func() fields {
				if err := db.Create(&Customer{
					ID:       "123",
					Name:     "Test Customer",
					AreaCode: "+86",
					PhoneNum: "1234567890",
				}).Error; err != nil {
					t.Fatalf("failed to prep test data: %v", err)
				}
				return fields{db}
			},
			args: args{
				customer: &customer.Customer{
					Person: &entity.Person{
						ID:   "123",
						Name: "Test Customer",
					},
					ContactInfo: valueobj.ContactInfo{AreaCode: "+86", PhoneNum: "9876543210"},
				},
			},
			wantErr: assert.NoError,
			checkData: func() bool {
				var model Customer
				if err := db.First(&model, "id = ?", "123").Error; err != nil {
					t.Logf("find data from db failed: %s", err)
					return false
				}
				return model.PhoneNum == "9876543210"
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fields := tt.fields()
			s := &sqlRepository{
				db: fields.db,
			}
			tt.wantErr(t, s.Update(tt.args.customer), fmt.Sprintf("Update(%v)", tt.args.customer))
			if tt.checkData != nil {
				assert.True(t, tt.checkData())
			}
		})
	}
}

func Test_sqlRepository_Get(t *testing.T) {
	db := prepTestDB(t)

	type fields struct {
		db *gorm.DB
	}
	type args struct {
		id string
	}
	tests := []struct {
		name    string
		fields  func() fields
		args    args
		want    *customer.Customer
		wantErr bool
	}{
		{
			name: "case object exist",
			fields: func() fields {
				if err := db.Create(&Customer{
					ID:       "123",
					Name:     "Test Customer",
					AreaCode: "+86",
					PhoneNum: "1234567890",
				}).Error; err != nil {
					t.Fatalf("failed to prep test data: %s", err)
				}
				return fields{db}
			},
			args: args{
				id: "123",
			},
			want: &customer.Customer{
				Person: &entity.Person{
					ID:   "123",
					Name: "Test Customer",
				},
				ContactInfo: valueobj.ContactInfo{AreaCode: "+86", PhoneNum: "1234567890"},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fields := tt.fields()
			s := &sqlRepository{
				db: fields.db,
			}
			got, err := s.Get(tt.args.id)
			assert.NoError(t, err)
			assert.Equal(t, tt.want, got)
		})
	}
}

func Test_sqlRepository_Remove(t *testing.T) {
	db := prepTestDB(t)
	type fields struct {
		db *gorm.DB
	}
	type args struct {
		id string
	}
	tests := []struct {
		name      string
		fields    func() fields
		args      args
		wantErr   assert.ErrorAssertionFunc
		checkData func() bool
	}{
		{
			name: "case success",
			fields: func() fields {
				if err := db.Create(&Customer{
					ID:       "123",
					Name:     "Test Customer",
					AreaCode: "+86",
					PhoneNum: "1234567890",
				}).Error; err != nil {
					t.Fatalf("failed to prep test data: %s", err)
				}
				return fields{db}
			},
			args: args{
				id: "123",
			},
			wantErr: assert.NoError,
			checkData: func() bool {
				result := db.Find(&Customer{}, "id = ?", "123")
				return result.RowsAffected == 0
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fields := tt.fields()
			s := &sqlRepository{
				db: fields.db,
			}
			tt.wantErr(t, s.Remove(tt.args.id), fmt.Sprintf("Remove(%v)", tt.args.id))
			if tt.checkData != nil {
				assert.True(t, tt.checkData())
			}
		})
	}
}
