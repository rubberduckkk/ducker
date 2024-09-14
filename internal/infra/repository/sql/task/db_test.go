package task

import (
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"

	"github.com/rubberduckkk/ducker/internal/domain/task"
	"github.com/rubberduckkk/ducker/internal/domain/task/entity"
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

	if db.Migrator().HasTable(Task{}.TableName()) {
		if err := db.Exec(fmt.Sprintf("DELETE FROM %v", Task{}.TableName())).Error; err != nil {
			t.Fatalf("failed to truncate table: %s", err)
		}
		if err := db.Exec("DELETE FROM SQLITE_SEQUENCE WHERE name=?", Task{}.TableName()).Error; err != nil {
			t.Logf("failed to reset table metadata: %s", err)
		}
		return db
	}
	if err = db.AutoMigrate(&Task{}); err != nil {
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
		task *task.Task
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
				return fields{db: db}
			},
			args: args{
				task: &task.Task{
					TaskInfo: &entity.TaskInfo{
						ID:        "456",
						Content:   "test content",
						CreatedAt: time.Now(),
						UpdatedAt: time.Now(),
					},
					CustomerID: "123",
				},
			},
			wantErr: assert.NoError,
			checkData: func() bool {
				return db.First(&Task{}, "id = ?", "456").Error == nil
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fields := tt.fields()
			s := &sqlRepository{
				db: fields.db,
			}
			tt.wantErr(t, s.Create(tt.args.task), fmt.Sprintf("Create(%v)", tt.args.task))
			if tt.checkData != nil {
				assert.True(t, tt.checkData())
			}
		})
	}
}

func Test_sqlRepository_Update(t *testing.T) {
	db := prepTestDB(t)
	now := time.Now()
	type fields struct {
		db *gorm.DB
	}
	type args struct {
		task *task.Task
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
				if err := db.Create(&Task{
					ID:         "123",
					CustomerID: "456",
					Content:    "test content",
					CreatedAt:  now.Unix(),
					UpdatedAt:  now.Unix(),
				}).Error; err != nil {
					t.Fatalf("failed to prep test data: %s", err)
				}
				return fields{db: db}
			},
			args: args{
				task: &task.Task{
					TaskInfo: &entity.TaskInfo{
						ID:        "123",
						Content:   "test new content",
						CreatedAt: now,
						UpdatedAt: time.Now(),
					},
					CustomerID: "456",
				},
			},
			wantErr: assert.NoError,
			checkData: func() bool {
				var model Task
				if err := db.First(&model, "id = ?", "123").Error; err != nil {
					t.Fatalf("failed to query test data: %s", err)
				}
				return model.Content == "test new content"
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fields := tt.fields()
			s := &sqlRepository{
				db: fields.db,
			}
			tt.wantErr(t, s.Update(tt.args.task), fmt.Sprintf("Update(%v)", tt.args.task))
			if tt.checkData != nil {
				assert.True(t, tt.checkData())
			}
		})
	}
}

func Test_sqlRepository_Get(t *testing.T) {
	db := prepTestDB(t)
	now := time.Now()
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
		want    *task.Task
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name: "case success",
			fields: func() fields {
				if err := db.Create(&Task{
					ID:         "123",
					CustomerID: "456",
					Content:    "test content",
					CreatedAt:  now.Unix(),
					UpdatedAt:  now.Unix(),
				}).Error; err != nil {
					t.Fatalf("failed to prep test data: %s", err)
				}
				return fields{db: db}
			},
			args: args{
				id: "123",
			},
			want: &task.Task{
				TaskInfo: &entity.TaskInfo{
					ID:        "123",
					Content:   "test content",
					CreatedAt: time.Unix(now.Unix(), 0),
					UpdatedAt: time.Unix(now.Unix(), 0),
				},
				CustomerID: "456",
			},
			wantErr: assert.NoError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fields := tt.fields()
			s := &sqlRepository{
				db: fields.db,
			}
			got, err := s.Get(tt.args.id)
			if !tt.wantErr(t, err, fmt.Sprintf("Get(%v)", tt.args.id)) {
				return
			}
			assert.Equalf(t, tt.want, got, "Get(%v)", tt.args.id)
		})
	}
}

func Test_sqlRepository_Remove(t *testing.T) {
	db := prepTestDB(t)
	now := time.Now()
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
				if err := db.Create(&Task{
					ID:         "123",
					CustomerID: "456",
					Content:    "test content",
					CreatedAt:  now.Unix(),
					UpdatedAt:  now.Unix(),
				}).Error; err != nil {
					t.Fatalf("failed to prep test data: %s", err)
				}
				return fields{db: db}
			},
			args: args{
				id: "123",
			},
			wantErr: assert.NoError,
			checkData: func() bool {
				result := db.Find(&Task{}, "id = ?", "123")
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
