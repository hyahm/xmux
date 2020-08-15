package xmux

// func midware1(http.ResponseWriter, *http.Request) bool {
// 	return false
// }

// func midware2(http.ResponseWriter, *http.Request) bool {
// 	return false
// }

// func TestCompareMidware(t *testing.T) {
// 	m1 := midware1
// 	m2 := midware1
// 	tm1 := reflect.ValueOf(m1).Pointer()

// 	tm2 := reflect.ValueOf(m2).Pointer()

// 	name1 := runtime.FuncForPC(tm1).Name()
// 	name2 := runtime.FuncForPC(tm2).Name()

// 	if name1 != name2 {
// 		t.Error("func not compare")
// 	}

// }
