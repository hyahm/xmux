package xmux

// 分页用到的计算sql limit的数值

func GetLimit(count, page, limit int) (int, int, int) {
	// 如果limit是0， 那么返回 0 0
	if limit <= 0 || count == 0 {
		return 1, 0, 0
	}
	// 如果page小于1页， 默认返回第一页
	if page < 1 {
		page = 1
	}
	// 超出了最大页码，返回最大的页码

	if page*limit > count+limit {
		if count%limit == 0 {
			page = count / limit
		} else {
			page = count/limit + 1
		}
	}

	start := (page - 1) * limit
	// 计算最终返回的start, step
	if count-start < limit {
		return page, start, count - start
	} else {
		return page, start, limit
	}

}
