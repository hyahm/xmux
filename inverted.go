package xmux

import (
	"fmt"
	"strings"
	"sync"
)

type DataSource struct {
}

type Condition struct {
	// 用在filter 里面，用来过滤， 比如按照id排序， 其他条件可以在
	ID int64 `json:"id"`
}

// 搜索核心算法
type InvertedIndex[T any, C any] struct {
	sync.RWMutex
	cache map[int64]*T
	// 某个信息排序值
	sortOrder map[string][]int64
	maxid     int64
	condition C
}

func NewInverted[T any, C any]() InvertedIndex[T, C] {
	var c C
	return InvertedIndex[T, C]{
		sync.RWMutex{},
		make(map[int64]*T),
		make(map[string][]int64),
		0,
		c,
	}
}

// 不动
func (i *InvertedIndex[T, C]) GetMaxid() int64 {
	i.RLock()
	defer i.RUnlock()
	return i.maxid
}

// 插入新值， id 排序， 从小到大
func (ii *InvertedIndex[T, C]) insertInerted(cn []*T) {

	ii.Lock()
	defer ii.Unlock()
	// 因为长度大于0才会进来， 所以直接取最后一个值就可以
	// ii.cache[v.ID] = v
	// ii.sortOrder[] = append([]int64{v.ID}, ii.sortOrder["000000"]...)
	fmt.Printf("自选总共缓存%d条数据\n", len(ii.sortOrder))
}

func (ii *InvertedIndex[T, C]) InitConsultNew() {
	// 第一次插入所有历史数据
	// 入库历史数据
	// 处理一批数据就释放锁，避免长时间持有

}

func (ii *InvertedIndex[T, C]) UpdateConsultNew() {
	// for 循环， 多长时间更新一次
	var list []*T
	// 更新数据的操作
	ii.insertInerted(list)
}

// 不用动
func (ii *InvertedIndex[T, C]) MergeSliceIds(stock []string, pageSize int) []int64 {
	// 检查初始化状态，如果未完成则返回空
	// if !IsInitialized() {
	// 	golog.Warn("索引初始化未完成，无法执行查询")
	// 	return []int64{}
	// }
	if len(stock) == 0 {
		return make([]int64, 0)
	}
	s := make([][]int64, 0, len(stock))

	for _, v := range stock {
		value := ii.Get(v)
		if len(value) > 0 {
			s = append(s, value)
		}
	}
	return ii.mergeUniqueDesc(s, pageSize)
}

// 不用动
func (ii *InvertedIndex[T, C]) Get(key string) []int64 {
	ii.RLock()
	defer ii.RUnlock()
	// 返回切片的副本，避免外部修改内部数据
	result := make([]int64, len(ii.sortOrder[key]))
	copy(result, ii.sortOrder[key])
	return result
}

// 不用动
func (ii *InvertedIndex[T, C]) GetConsult(key int64) *T {
	ii.RLock()
	defer ii.RUnlock()
	return ii.cache[key]
}

func (ii *InvertedIndex[T, C]) Search(searchWord string, pageSize int) []*T {

	// 检查初始化状态，如果未完成则返回空
	// if !IsInitialized() {
	// 	golog.Warn("索引初始化未完成，无法执行查询")
	// 	return []meilisearch.ConsultNews{}
	// }
	// 取出搜索词，多个就分成切片，
	words := strings.Split(searchWord, ",") // 只动这里
	//
	newlist := ii.MergeSliceIds(words, pageSize)
	cns := make([]*T, 0, len(newlist))
	for _, v := range newlist {
		cns = append(cns, ii.GetConsult(v))
	}
	return cns
}

type node struct {
	val      int64 // 当前值
	arrayIdx int   // 多个切片中的某切片的id
	idx      int   // 在该切片中的下标
}

func (ii *InvertedIndex[T, C]) filter(id int64, pageSize int) bool {
	// 过滤条件写在这里
	return true
}

// 核心函数：去重排序， 不需要动
func (ii *InvertedIndex[T, C]) mergeUniqueDesc(arrays [][]int64, pageSize int) []int64 {
	if len(arrays) == 0 {
		return make([]int64, 0)
	}
	result := make([]int64, 0, pageSize)
	// ns 的key 是 arrays 的 index
	// ns := make(map[int]*node)
	// 长度为 0 会进来吗

	if len(arrays) == 1 {
		for _, v := range arrays[0] {
			if !ii.filter(v, pageSize) {
				continue
			}
			result = append(result, v)
			if len(result) >= pageSize {
				return result
			}

		}
		return result
	}
	// 只是为了去重
	m := make(map[int64]struct{})
	// 排序所有的array用的
	orderNode := make([]*node, 0, len(arrays))
	// if len(arrays) > pageSize {
	// 如果大于pagesize， 因为只需要返回 pagesize 条， 那么只需要取最大的 pagesize 条 即可， 其他的都可以丢弃
	// 直接取有效的第一个数值
	for i := range arrays {

		for j := range arrays[i] {
			// 有效值就直接返回当前的有效值， 进入下一个切片的计算
			// 1. 如果minid 大于0，  并且id的值比minid小的
			// 2, 如果查询的是资讯， 但是显示的资讯的就跳过，  因为这里只有资讯和公告
			// 3. 如果查询的是公告， 但是显示的公告
			// 4. 前面没有有这个id的值的
			if _, ok := m[arrays[i][j]]; ok {
				continue
			}

			if !ii.filter(arrays[i][j], pageSize) {
				continue
			}

			m[arrays[i][j]] = struct{}{}
			// if i == 2238 {
			// 	golog.Info(arrays[i][j], len(orderNode))
			// }
			ns := &node{
				arrayIdx: i,
				idx:      j,
				val:      arrays[i][j],
			}
			orderNode = append(orderNode, ns)
			break
		}

	}

	if len(orderNode) == 0 {
		return result
	}
	// 去掉多余的切片最多保留 pagesize 个切片
	newOrderNode := make([]*node, 0, pageSize)
	// 排序只需要前 pagesize 条
	for range pageSize {
		max := 0
		if len(orderNode) == 1 {
			newOrderNode = append(newOrderNode, orderNode[0])
			break
		}
		for j := range orderNode[1:] {
			if orderNode[j+1].val > orderNode[max].val {
				// 因为索引是从切片的1开始的， 所以要+1
				max = j + 1
			}
		}
		newOrderNode = append(newOrderNode, orderNode[max])
		orderNode = append(orderNode[:max], orderNode[max+1:]...)

	}

	for range pageSize {
		// newOrderNode 的下标
		max := 0
		// fmt.Println(len(newOrderNode))
		for i := range newOrderNode {
			if newOrderNode[i].val > newOrderNode[max].val {
				max = i
			}
		}
		result = append(result, newOrderNode[max].val)
		// 如果当前最大值取出来了，  那么要获取改切片的下一个有效值
		// 如果已经是最后一个元素了，那么把这条删掉
		if len(arrays[newOrderNode[max].arrayIdx])-1 == newOrderNode[max].idx {
			//
			if max == len(newOrderNode)-1 {
				newOrderNode = newOrderNode[:max]
			} else {
				newOrderNode = append(newOrderNode[:max], newOrderNode[max+1:]...)
			}
			continue
		}
		// 判断 如果到末尾了也没找到符合条件的值， 那么就也删除
		exsit := false
		// 遍历这个切片
		for j, value := range arrays[newOrderNode[max].arrayIdx][newOrderNode[max].idx:] {
			// 有效值就直接返回当前的有效值， 进入下一个切片的计算
			// 1. 如果minid 大于0，  并且id的值比minid小的
			// 2, 如果查询的是资讯， 但是显示的资讯的就跳过，  因为这里只有资讯和公告
			// 3. 如果查询的是公告， 但是显示的公告
			// 4. 前面没有有这个id的值的
			_, ok := m[value]
			if ok {
				continue
			}

			m[value] = struct{}{}
			newOrderNode[max].idx = j + newOrderNode[max].idx
			newOrderNode[max].val = value
			exsit = true
			break
		}
		if !exsit {
			newOrderNode = append(newOrderNode[:max], newOrderNode[max+1:]...)

		}

	}
	// 返回结果
	return result

}
