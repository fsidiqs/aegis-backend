package repository

import (
	"context"
	"fmt"

	"github.com/fsidiqs/aegis-backend/db"
	"github.com/fsidiqs/aegis-backend/helper"
	"github.com/fsidiqs/aegis-backend/model"
	"github.com/fsidiqs/aegis-backend/model/apperror"
	"github.com/google/uuid"

	"gorm.io/gorm"
)

type psqlOrganizationRepository struct {
	DB *gorm.DB
}

// NewOrganizationRepository is factory pattern to create OrganizationRepository with gorm as parameter
func NewOrganizationRepository(db *gorm.DB) IOrganizationRepository {
	return &psqlOrganizationRepository{
		DB: db,
	}
}

func (r *psqlOrganizationRepository) Create(ctx context.Context, org model.Organization, userID uuid.UUID) (*model.Organization, error) {
	var dbContext *gorm.DB
	// get db from context, if db doesnt exists then use the class db
	if userDBCtx := ctx.Value(db.DBTrxUserKey); userDBCtx != nil {
		dbContext = userDBCtx.(*gorm.DB)
	} else {
		dbContext = r.DB
	}

	err := dbContext.WithContext(ctx).Model(&model.Organization{}).Create(&org).Error
	if err != nil {
		errMsg := fmt.Sprintf("%s:err:%v", helper.TraceCurrentFuncArgs(org), err)
		return nil, apperror.NewRepoErrorMsg(err, errMsg)
	}
	return &org, nil
}

// List returns all users
func (r *psqlOrganizationRepository) List(ctx context.Context) ([]model.Organization, error) {
	var dbContext *gorm.DB
	// get db from context, if db doesnt exists then use the class db

	dbContext = r.DB

	var orgs []model.Organization
	err := dbContext.WithContext(ctx).
		Find(&orgs).Error
	if err != nil {
		errMsg := fmt.Sprintf("%s", helper.TraceCurrentFuncArgs(nil))
		return nil, apperror.NewRepoErrorMsg(err, errMsg)
	}
	return orgs, nil
}

func (r *psqlOrganizationRepository) FindByID(ctx context.Context, id uuid.UUID) (*model.Organization, error) {
	var org model.Organization
	var dbContext *gorm.DB
	// get db from context, if db doesnt exists then use the class db
	if userDBCtx := ctx.Value(db.DBTrxUserKey); userDBCtx != nil {
		dbContext = userDBCtx.(*gorm.DB)
	} else {
		dbContext = r.DB
	}

	err := dbContext.Debug().WithContext(ctx).
		Model(&model.Organization{}).First(&org, id).Error
	if err != nil {
		errMsg := fmt.Sprintf("%s:err:%v", helper.TraceCurrentFuncArgs(
			map[string]string{},
		), err)
		return nil, apperror.NewRepoErrorMsg(err, errMsg)
	}
	return &org, nil
}

// soft delete user
func (r *psqlOrganizationRepository) HardDelete(ctx context.Context, uid uuid.UUID) error {
	var dbContext *gorm.DB
	// get db from context, if db doesnt exists then use the class db
	u := model.UserRecordFlagUpdate{
		RecordFlag: model.RecDeleted,
	}

	dbContext = r.DB

	err := dbContext.WithContext(ctx).
		Where("id = ?", uid).
		Delete(&u).Error
	if err != nil {
		errMsg := fmt.Sprintf("%s", helper.TraceCurrentFuncArgs(
			map[string]string{
				"user_id": uid.String(),
			},
		))
		return apperror.NewRepoErrorMsg(err, errMsg)
	}
	return nil
}

func (r *psqlOrganizationRepository) Update(ctx context.Context, id uuid.UUID, org *model.Organization) error {
	var dbContext *gorm.DB

	if dbCtx := ctx.Value(db.DBTrxUserKey); dbCtx != nil {
		dbContext = dbCtx.(*gorm.DB)
	} else {
		dbContext = r.DB
	}

	rows := dbContext.WithContext(ctx).
		// Model(&model.User{}).
		Where("id = ?", id).
		Updates(org).
		RowsAffected

	if rows == 0 {
		errMsg := fmt.Sprintf("%s", helper.TraceCurrentFuncArgs(
			map[string]string{},
			*org,
		))
		return apperror.NewRepoErrorMsg(nil, errMsg)
	}
	return nil
}
