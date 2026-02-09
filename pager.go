package xmux

// 分页用到的计算sql limit的数值
func GetPaginationLimit(total, page, pageSize int) (finalPage int, offset int, limit int) {
	// 1. 基础参数校验与修正
	// 处理总条数为负数的极端情况
	if total < 0 {
		total = 0
	}
	// 每页条数<=0时，默认设为10（可根据业务调整）
	if pageSize <= 0 {
		pageSize = 10
	}
	// 页码<1时，默认返回第一页
	if page < 1 {
		page = 1
	}

	// 2. 计算最大有效页码
	maxPage := 0
	if total > 0 {
		// 简化最大页码计算：向上取整 (total + pageSize - 1) / pageSize
		maxPage = (total + pageSize - 1) / pageSize
	}

	// 3. 修正请求页码（超出最大页码则返回最后一页）
	finalPage = page
	if maxPage > 0 && finalPage > maxPage {
		finalPage = maxPage
	}

	// 4. 计算OFFSET和LIMIT
	offset = (finalPage - 1) * pageSize
	// 最后一页的LIMIT需要修正（避免超出总条数）
	if offset+pageSize > total && total > 0 {
		limit = total - offset
	} else {
		limit = pageSize
	}

	// 总条数为0时，强制返回0
	if total == 0 {
		offset = 0
		limit = 0
	}

	return finalPage, offset, limit
}
