package task

import (
	"context"
	"errors"
	"seconda/internal/enum"
	"seconda/internal/model/team"
	"seconda/internal/model/user"
	"time"

	"gorm.io/gorm"
)

type Task struct {
	Id            int         `json:"id" gorm:"type:int;not null;primaryKey;autoIncrement"`
	AssigneeId    int         `json:"assignee_id" gorm:"type:int;not null;"`
	AssigneeUser  user.User   `json:"assignee_user" gorm:"foreignKey:AssigneeId;references:Id"`
	TeamId        int         `json:"team_id" gorm:"type:int;not null;"`
	Team          team.Team   `json:"team" gorm:"foreignKey:TeamId;references:Id"`
	CreatedBy     int         `json:"created_by" gorm:"type:int;not null;"`
	CreatedByUser user.User   `json:"created_by_user" gorm:"foreignKey:CreatedBy;references:Id"`
	Name          string      `json:"name" gorm:"type:varchar(256);not null;"`
	CreatedAt     time.Time   `json:"created_at" gorm:"type:timestamp;not null"`
	UpdatedAt     time.Time   `json:"updated_at" gorm:"type:timestamp;not null"`
	Status        enum.Status `json:"status" gorm:"type:int;not null;"`
}

var NotFoundErr = errors.New("task not found")

var UserNotTeamMemberErr = errors.New("assignee user is not a member of team")

func (t Task) TableName() string {
	return "tasks"
}

type Filter struct {
	TeamId     int
	Status     int
	AssigneeId int
	Limit      int
	Offset     int
}

func CreateTask(ctx context.Context, db *gorm.DB, t *Task) error {
	var isMember bool
	query := `
		SELECT EXISTS (
			SELECT 1 
			FROM team_members 
			WHERE user_id = ? AND team_id = ?
		)
	`

	err := db.WithContext(ctx).Raw(query, t.AssigneeId, t.TeamId).Scan(&isMember).Error
	if err != nil {
		return err
	}

	if !isMember {
		return UserNotTeamMemberErr
	}

	return db.WithContext(ctx).Create(t).Error
}

func GetTasksByFilter(ctx context.Context, db *gorm.DB, filter Filter) ([]Task, error) {
	var tasks []Task

	query := `
		SELECT 
			t.id, t.assignee_id, t.team_id, t.created_by, t.status, t.created_at, t.updated_at,
			au.id AS "AssigneeUser__id", au.name AS "AssigneeUser__name", au.phone AS "AssigneeUser__phone", 
				au.surname AS "AssigneeUser__surname",  au.login AS "AssigneeUser__login",
				au.password AS "AssigneeUser__password", au.created_at AS "AssigneeUser__createdAt",
				au.updated_at AS "AssigneeUser__updatedAt",
			cu.id AS "CreatedByUser__id", cu.name AS "CreatedByUser__name", cu.phone AS "CreatedByUser__phone", 
				cu.surname AS "CreatedByUser__surname",  cu.login AS "CreatedByUser__login",
				cu.password AS "CreatedByUser__password", cu.created_at AS "CreatedByUser__createdAt",
				cu.updated_at AS "CreatedByUser__updatedAt",
		    ut.id AS "Team__id", ut.name AS "Team__name", ut.created_by AS "Team__createdBy", 
		    	ut.created_at AS "Team__createdAt", ut.updated_at AS "Team__updatedAt"
		FROM tasks t
		LEFT JOIN users au ON t.assignee_id = au.id
		LEFT JOIN users cu ON t.created_by = cu.id
		LEFT JOIN teams ut ON t.team_id = ut.id
		WHERE 1=1
	`

	var args []any

	if filter.TeamId > 0 {
		query += " AND t.team_id = ?"
		args = append(args, filter.TeamId)
	}
	if filter.Status > 0 {
		query += " AND t.status = ?"
		args = append(args, filter.Status)
	}
	if filter.AssigneeId > 0 {
		query += " AND t.assignee_id = ?"
		args = append(args, filter.AssigneeId)
	}

	query += " ORDER BY t.created_at DESC"

	if filter.Limit > 0 {
		query += " LIMIT ?"
		args = append(args, filter.Limit)
	}
	if filter.Offset > 0 {
		query += " OFFSET ?"
		args = append(args, filter.Offset)
	}

	err := db.WithContext(ctx).Raw(query, args...).Scan(&tasks).Error
	if err != nil {
		return tasks, err
	}

	return tasks, nil
}

func UpdateTaskById(ctx context.Context, db *gorm.DB, nt *Task, id int) (Task, error) {
	var ot Task

	result := db.WithContext(ctx).First(&ot, id)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return ot, NotFoundErr
	}

	err := db.WithContext(ctx).Model(&ot).Updates(nt).Error
	return ot, err
}

func GetTaskById(ctx context.Context, db *gorm.DB, id int) (Task, error) {
	var t Task

	err := db.WithContext(ctx).First(&t, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return Task{}, NotFoundErr
		}
		return Task{}, err
	}

	return t, nil
}
