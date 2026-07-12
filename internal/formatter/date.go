package formatter

import "time"

func FormatDate(value *time.Time) *string {
	if value == nil {
		return nil
	}

	formatted := value.Format("02/01/2006")
	return &formatted
}
