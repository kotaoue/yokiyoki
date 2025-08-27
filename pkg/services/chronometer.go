package services

import (
	"fmt"
	"time"
)

type Chronometer struct {
	Start time.Time
	End   time.Time
	jst   *time.Location
}

type ChronometerOption struct {
	Days      *int
	StartDate *string
	EndDate   *string
}

func NewChronometer(opt ChronometerOption) (*Chronometer, error) {
	return createChronometer(opt)
}

func createChronometer(opt ChronometerOption) (*Chronometer, error) {
	if opt.Days == nil && opt.StartDate == nil && opt.EndDate == nil {
		return createByDefault()
	}

	if opt.StartDate != nil && opt.EndDate != nil {
		return createByDateRange(*opt.StartDate, *opt.EndDate)
	}

	if opt.Days != nil {
		return createByDays(*opt.Days)
	}

	return nil, fmt.Errorf("invalid chronometer options")
}

func createByDefault() (*Chronometer, error) {
	jst, _ := time.LoadLocation("Asia/Tokyo")
	return &Chronometer{jst: jst}, nil
}

func createByDateRange(startDate, endDate string) (*Chronometer, error) {
	jst, _ := time.LoadLocation("Asia/Tokyo")

	start, err := time.Parse("2006-01-02", startDate)
	if err != nil {
		return nil, err
	}

	end, err := time.Parse("2006-01-02", endDate)
	if err != nil {
		return nil, err
	}

	// Set end time to end of day
	end = end.Add(23*time.Hour + 59*time.Minute + 59*time.Second)

	return &Chronometer{
		Start: start,
		End:   end,
		jst:   jst,
	}, nil
}

func createByDays(days int) (*Chronometer, error) {
	jst, _ := time.LoadLocation("Asia/Tokyo")
	end := time.Now().In(jst)
	start := end.AddDate(0, 0, -days)
	return &Chronometer{
		Start: start,
		End:   end,
		jst:   jst,
	}, nil
}

// Contains checks if a time is within the chronometer's period
func (c *Chronometer) Contains(t time.Time) bool {
	return !t.Before(c.Start) && !t.After(c.End)
}

func (c *Chronometer) StartTime() time.Time {
	return c.Start
}

func (c *Chronometer) EndTime() time.Time {
	return c.End
}

func (c *Chronometer) GetLast7Days() (time.Time, time.Time) {
	now := time.Now().In(c.jst)
	c.Start = now.AddDate(0, 0, -7)
	c.End = now
	return c.Start, c.End
}

func (c *Chronometer) GetLast30Days() (time.Time, time.Time) {
	now := time.Now().In(c.jst)
	c.Start = now.AddDate(0, 0, -30)
	c.End = now
	return c.Start, c.End
}

func (c *Chronometer) GetLastMonth() (time.Time, time.Time) {
	now := time.Now().In(c.jst)
	firstOfThisMonth := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, c.jst)
	firstOfLastMonth := firstOfThisMonth.AddDate(0, -1, 0)
	lastOfLastMonth := firstOfThisMonth.AddDate(0, 0, -1)
	c.Start = firstOfLastMonth
	c.End = lastOfLastMonth
	return c.Start, c.End
}

func (c *Chronometer) GetPreviousHalf() (time.Time, time.Time) {
	now := time.Now().In(c.jst)
	c.Start, c.End = c.calculatePreviousHalf(now)
	return c.Start, c.End
}

func (c *Chronometer) GetPreviousYear() (time.Time, time.Time) {
	now := time.Now().In(c.jst)
	firstOfLastYear := time.Date(now.Year()-1, 1, 1, 0, 0, 0, 0, c.jst)
	lastOfLastYear := time.Date(now.Year()-1, 12, 31, 0, 0, 0, 0, c.jst)
	c.Start = firstOfLastYear
	c.End = lastOfLastYear
	return c.Start, c.End
}

func (c *Chronometer) GetPreviousFiscalYear() (time.Time, time.Time) {
	now := time.Now().In(c.jst)
	var fiscalYear int
	if now.Month() >= 4 {
		// 4月以降なら現在の年度の前年度
		fiscalYear = now.Year() - 1
	} else {
		// 3月以前なら前々年度
		fiscalYear = now.Year() - 2
	}
	firstOfFiscalYear := time.Date(fiscalYear, 4, 1, 0, 0, 0, 0, c.jst)
	lastOfFiscalYear := time.Date(fiscalYear+1, 3, 31, 0, 0, 0, 0, c.jst)
	c.Start = firstOfFiscalYear
	c.End = lastOfFiscalYear
	return c.Start, c.End
}

func (c *Chronometer) calculatePreviousHalf(now time.Time) (time.Time, time.Time) {
	currentMonth := now.Month()
	if currentMonth >= 4 && currentMonth <= 9 {
		prevHalfStart := time.Date(now.Year()-1, 10, 1, 0, 0, 0, 0, c.jst)
		prevHalfEnd := time.Date(now.Year(), 3, 31, 0, 0, 0, 0, c.jst)
		return prevHalfStart, prevHalfEnd
	} else {
		if currentMonth >= 10 {
			prevHalfStart := time.Date(now.Year(), 4, 1, 0, 0, 0, 0, c.jst)
			prevHalfEnd := time.Date(now.Year(), 9, 30, 0, 0, 0, 0, c.jst)
			return prevHalfStart, prevHalfEnd
		} else {
			prevHalfStart := time.Date(now.Year()-1, 4, 1, 0, 0, 0, 0, c.jst)
			prevHalfEnd := time.Date(now.Year()-1, 9, 30, 0, 0, 0, 0, c.jst)
			return prevHalfStart, prevHalfEnd
		}
	}
}

func (c *Chronometer) GetLast7DaysDescription() string {
	start, end := c.GetLast7Days()
	return fmt.Sprintf("(%s to %s JST)", start.Format("2006-01-02"), end.Format("2006-01-02"))
}

func (c *Chronometer) GetLast30DaysDescription() string {
	start, end := c.GetLast30Days()
	return fmt.Sprintf("(%s to %s JST)", start.Format("2006-01-02"), end.Format("2006-01-02"))
}

func (c *Chronometer) GetLastMonthDescription() string {
	start, end := c.GetLastMonth()
	return fmt.Sprintf("(%s to %s JST)", start.Format("2006-01-02"), end.Format("2006-01-02"))
}

func (c *Chronometer) GetPreviousHalfDescription() string {
	start, end := c.GetPreviousHalf()
	return fmt.Sprintf("(%s to %s JST)", start.Format("2006-01-02"), end.Format("2006-01-02"))
}

func (c *Chronometer) GetPreviousYearDescription() string {
	start, end := c.GetPreviousYear()
	return fmt.Sprintf("(%s to %s JST)", start.Format("2006-01-02"), end.Format("2006-01-02"))
}

func (c *Chronometer) GetPreviousFiscalYearDescription() string {
	start, end := c.GetPreviousFiscalYear()
	return fmt.Sprintf("(%s to %s JST)", start.Format("2006-01-02"), end.Format("2006-01-02"))
}

func (c *Chronometer) GetLast7DaysResult() (int, string, string) {
	return 7, "", ""
}

func (c *Chronometer) GetLast30DaysResult() (int, string, string) {
	return 30, "", ""
}

func (c *Chronometer) GetLastMonthResult() (int, string, string) {
	start, end := c.GetLastMonth()
	return 0, start.Format("2006-01-02"), end.Format("2006-01-02")
}

func (c *Chronometer) GetPreviousHalfResult() (int, string, string) {
	start, end := c.GetPreviousHalf()
	return 0, start.Format("2006-01-02"), end.Format("2006-01-02")
}

func (c *Chronometer) GetPreviousYearResult() (int, string, string) {
	start, end := c.GetPreviousYear()
	return 0, start.Format("2006-01-02"), end.Format("2006-01-02")
}

func (c *Chronometer) GetPreviousFiscalYearResult() (int, string, string) {
	start, end := c.GetPreviousFiscalYear()
	return 0, start.Format("2006-01-02"), end.Format("2006-01-02")
}
