package dag

import "testing"

func TestStart(t *testing.T) {

	b := &block{
		height: 0,
		prev:   "000",
		hash:   "111",
		addr:   "jie",
		refers: nil,
	}
	ta := newTask()
	err := ta.add(b)
	if err != nil {
		t.Fatal(err)
	}

	err = ta.add(&block{
		height: 1,
		prev:   "111",
		hash:   "222",
		addr:   "jie",
		refers: nil,
	})
	if err != nil {
		t.Fatal(err)
	}

	err = ta.add(&block{
		height: 1,
		prev:   "112",
		hash:   "223",
		addr:   "jie2",
		refers: nil,
	})
	if err != nil {
		t.Fatal(err)
	}

	err = ta.add(&block{
		height: 1,
		prev:   "113",
		hash:   "224",
		addr:   "jie3",
		refers: []string{"223"},
	})
	if err != nil {
		t.Fatal(err)
	}
	ta.print()
}
