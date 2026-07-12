package transaction

func toSearchItem(record SearchRecord) SearchItem {
	return SearchItem{
		ID:              record.ID,
		ReferenceCode:   record.ReferenceCode,
		AccountID:       record.AccountID,
		TransactionType: record.TransactionType,
		Side:            record.Side,
		SideName:        mapSideName(record.Side),
		Amount:          formatDecimal(record.Amount),
		Fee:             formatDecimal(record.Fee),
		Currency:        record.Currency,
		Status:          record.Status,
		Description:     record.Description,
		TransactedAt:    formatDate(record.TransactedAt),
	}
}

func mapSideName(side *string) *string {
	if side == nil {
		return nil
	}

	var name string
	switch *side {
	case "B":
		name = "BUY"
	case "S":
		name = "SELL"
	default:
		return nil
	}

	return &name
}
