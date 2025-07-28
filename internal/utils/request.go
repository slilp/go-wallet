package utils

func GetPaginationParams(pageQuery, limitQuery *int) (int, int) {
	page := 1
	limit := 20

	if pageQuery != nil {
		page = *pageQuery
	}

	if limitQuery != nil {
		limit = *limitQuery
	}

	return page, limit
}
