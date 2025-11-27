package service

import (
	"backend/model"
	"backend/repo"
	"database/sql"
	"fmt"
	"time"
)

type PaidLeaveInfo struct {
	TotalCount int
	Used       float64
	GivenAt    time.Time
}

type PaidLeaveInfoSum struct {
	TotalCount int
	Used       float64
}

var VacationCounts = []int{10, 11, 12, 14, 16, 18, 20}

type AddPaidLeaveParams struct {
	EmployeeID   string
	VacationDate time.Time
	StartAtHour  int
	Duration     int
}

// off年前の有給がいつ付与されたのか計算する
func CalculateVacationGivenDateByOffset(db *sql.DB, employeeID string, off int) (*time.Time, error) {
	if off < 0 {
		return nil, nil
	}
	employee, err := repo.GetEmployeeByID(db, employeeID)
	if err != nil {
		return nil, err
	}
	joiningDate := employee.JoiningDate
	baseDate := addMonthAccordingToCalendar(joiningDate, 6)
	fmt.Printf("basedate: %s\n", baseDate.String())
	today := time.Now()
	i := 0
	var vacationGivenPoints []time.Time
	for {
		vacationGivenPoint := addMonthAccordingToCalendar(baseDate, i*12)
		vacationGivenPoints = append(vacationGivenPoints, vacationGivenPoint)
		if vacationGivenPoint.After(today) {
			break
		}
		i++
	}
	offsetFromBack := len(vacationGivenPoints) - off - 1
	if offsetFromBack < 0 {
		return nil, nil
	}
	return &vacationGivenPoints[offsetFromBack], nil
}

// off年前の有給は総数（使用した分を含めて）何日分あるのかを計算する
func CalculateVacationCountByOffset(db *sql.DB, employeeID string, off int) (*int, error) {
	if off < 0 {
		return nil, nil
	}
	givenAtPtr, err := CalculateVacationGivenDateByOffset(db, employeeID, off)
	if err != nil {
		return nil, err
	}
	if givenAtPtr == nil {
		return nil, nil
	}
	givenAt := *givenAtPtr
	employee, err := repo.GetEmployeeByID(db, employeeID)
	if err != nil {
		return nil, err
	}
	joiningDate := employee.JoiningDate
	baseDate := addMonthAccordingToCalendar(joiningDate, 6)
	i := 0
	for {
		vacationGivenPoint := givenAt.AddDate(0, -12*i, 0)
		if vacationGivenPoint.Before(baseDate) || vacationGivenPoint.Equal(baseDate) {
			break
		}
		i++
	}
	lim := len(VacationCounts)
	if i >= lim {
		i = lim - 1
	}
	return &VacationCounts[i], nil
}

func GetRegisteredLeaveEmployeeListByOffset(db *sql.DB, employeeID string, off int) ([]model.LeaveEmployee, error) {
	if off < 0 {
		return []model.LeaveEmployee{}, nil
	}
	givenAt, err := CalculateVacationGivenDateByOffset(db, employeeID, off)
	if err != nil {
		return []model.LeaveEmployee{}, err
	}
	if givenAt == nil {
		return []model.LeaveEmployee{}, nil
	}
	leaveList, err := repo.GetRegisteredLeaveListByGivenAt(db, employeeID, *givenAt)
	if err != nil {
		return []model.LeaveEmployee{}, err
	}
	return leaveList, nil
}

func CalculateGeneralPaidLeaveInfo(db *sql.DB, employeeID string, off int) (*PaidLeaveInfo, error) {
	if off < 0 {
		return nil, nil
	}
	var info PaidLeaveInfo
	givenAtPtr, err := CalculateVacationGivenDateByOffset(db, employeeID, off)
	if err != nil {
		return nil, err
	}
	totalCountPtr, err := CalculateVacationCountByOffset(db, employeeID, off)
	if err != nil {
		return nil, err
	}
	leaveList, err := GetRegisteredLeaveEmployeeListByOffset(db, employeeID, off)
	if err != nil {
		return nil, err
	}
	if givenAtPtr == nil || totalCountPtr == nil {
		return nil, nil
	}
	info.GivenAt = *givenAtPtr
	today := time.Now()
	expireDate := givenAtPtr.AddDate(2, 0, 0)
	if today.After(expireDate) || today.Equal(expireDate) {
		return nil, nil
	}
	info.TotalCount = *totalCountPtr
	for _, v := range leaveList {
		info.Used += durationToDay(v.Duration)
	}
	fmt.Printf("info: %+v\n", info)
	return &info, nil
}

func CalculateSumPaidLeaveInfo(db *sql.DB, employeeID string) (*PaidLeaveInfoSum, error) {
	var infoSum = PaidLeaveInfoSum{}
	today := time.Now()
	for i := 0; i <= 3; i++ {
		info, err := CalculateGeneralPaidLeaveInfo(db, employeeID, i)
		if err != nil {
			return nil, err
		}
		if info == nil {
			break
		}
		deadLine := addMonthAccordingToCalendar(info.GivenAt, 24)
		if today.After(deadLine) || today.Equal(deadLine) {
			break
		}
		infoSum.TotalCount += info.TotalCount
		infoSum.Used += info.Used
	}
	return &infoSum, nil
}

func AddPaidLeaveByOffset(db *sql.DB, params AddPaidLeaveParams, off int) error {
	givenAt, err := CalculateVacationGivenDateByOffset(db, params.EmployeeID, off)
	if err != nil {
		return err
	}
	today := time.Now()
	var p = []repo.AddLeaveEmployeeParams{{
		EmployeeID:   params.EmployeeID,
		Duration:     params.Duration,
		StartAtHour:  params.StartAtHour,
		VacationDate: params.VacationDate,
		GivenAt:      *givenAt,
		RegisteredAt: today,
	},
	}
	tx, err := db.Begin()
	if err != nil {
		return err
	}

	defer func() {
		if p := recover(); p != nil {
			tx.Rollback()
			panic(p)
		} else if err != nil {
			tx.Rollback()
		} else {
			tx.Commit()
		}
	}()

	if err := repo.AddLeaveEmployee(tx, p); err != nil {
		return err
	}
	return err
}

func AddPaidLeave(db *sql.DB, params AddPaidLeaveParams) error {
	requiredDays := durationToDay(params.Duration)
	fmt.Printf("required: %f\n", requiredDays)
	added := 0
	for i := 3; i >= 0; i-- {
		fmt.Println(i)
		info, err := CalculateGeneralPaidLeaveInfo(db, params.EmployeeID, i)
		if err != nil {
			return err
		}
		if info == nil {
			continue
		}
		remainingDays := float64(info.TotalCount) - info.Used
		fmt.Printf("remaining: %f\n", remainingDays)
		if requiredDays > remainingDays {
			continue
		}
		err = AddPaidLeaveByOffset(db, params, i)
		if err != nil {
			return err
		}
		added += 1
		break
	}
	if added == 0 {
		return fmt.Errorf("有給の数が足りません")
	}
	return nil
}

func durationToDay(dur int) float64 {
	return float64(dur) / 8.0
}

func addMonthAccordingToCalendar(t time.Time, m int) time.Time {
	candDate := t.AddDate(0, m, 0)
	year, month, _ := t.Date()
	location := t.Location()
	firstOfNextMonth := time.Date(year, month+time.Month(m)+1, 1, 0, 0, 0, 0, location)
	lastOfDestMonth := firstOfNextMonth.AddDate(0, 0, -1)

	// キャリーオーバーが発生している場合は月末日を応答日とする
	if lastOfDestMonth.Month() != candDate.Month() {
		return lastOfDestMonth
	}
	return candDate
}
