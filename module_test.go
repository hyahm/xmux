package xmux

import "testing"

func TestModule(t *testing.T) {
	m := module{}
	m.add(DefaultModuleTemplate)
}
