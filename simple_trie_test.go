package tries

import (
	"strings"
	"testing"

	"github.com/smarty/assertions"
	"github.com/smarty/assertions/should"
	"github.com/smarty/benchy"
	"github.com/smarty/benchy/options"
	"github.com/smarty/benchy/providers"
)

func Test_SimpleTrie_Find_UInt8(t *testing.T) {
	trie, _ := NewTrie[uint8, int]()
	trie.Add(23, 1)
	trie.Add(100, 2)
	trie.Add(0, 3)
	trie.Add(64, 4)

	testTable := map[string]struct {
		Input    uint8
		Expected int
		OK       bool
	}{
		"23":          {Input: 23, Expected: 1, OK: true},
		"100":         {Input: 100, Expected: 2, OK: true},
		"0":           {Input: 0, Expected: 3, OK: true},
		"64":          {Input: 64, Expected: 4, OK: true},
		"not-in-data": {Input: 5, Expected: 0, OK: false},
	}

	for name, testCase := range testTable {
		t.Run(name, func(t *testing.T) {
			and := assertions.New(t)
			actual, ok := trie.Find(testCase.Input)
			and.So(actual, should.Equal, testCase.Expected)
			and.So(ok, should.Equal, testCase.OK)
		})
	}
}

func Test_SimpleTrie_Find_Int8(t *testing.T) {
	trie, _ := NewTrie[int8, int]()
	trie.Add(23, 1)
	trie.Add(100, 2)
	trie.Add(0, 3)
	trie.Add(64, 4)

	testTable := map[string]struct {
		Input    int8
		Expected int
		OK       bool
	}{
		"23":          {Input: 23, Expected: 1, OK: true},
		"100":         {Input: 100, Expected: 2, OK: true},
		"0":           {Input: 0, Expected: 3, OK: true},
		"64":          {Input: 64, Expected: 4, OK: true},
		"not-in-data": {Input: 5, Expected: 0, OK: false},
	}

	for name, testCase := range testTable {
		t.Run(name, func(t *testing.T) {
			and := assertions.New(t)
			actual, ok := trie.Find(testCase.Input)
			and.So(actual, should.Equal, testCase.Expected)
			and.So(ok, should.Equal, testCase.OK)
		})
	}
}

func Test_SimpleTrie_Find_Int64(t *testing.T) {
	trie, _ := NewTrie[int64, int]()
	trie.Add(0x01FF_ABAB_ABAB_ABAB, 1)
	trie.Add(0x01FF_ABAB_ABAB_ABBB, 2)
	trie.Add(100, 3)
	trie.Add(0, 4)
	trie.Add(64, 5)

	testTable := map[string]struct {
		Input    int64
		Expected int
		OK       bool
	}{
		"0x01FF_ABAB_ABAB_ABAB": {Input: 0x01FF_ABAB_ABAB_ABAB, Expected: 1, OK: true},
		"0x01FF_ABAB_ABAB_ABBB": {Input: 0x01FF_ABAB_ABAB_ABBB, Expected: 2, OK: true},
		"100":                   {Input: 100, Expected: 3, OK: true},
		"0":                     {Input: 0, Expected: 4, OK: true},
		"64":                    {Input: 64, Expected: 5, OK: true},
		"not-in-data":           {Input: 5, Expected: 0, OK: false},
	}

	for name, testCase := range testTable {
		t.Run(name, func(t *testing.T) {
			and := assertions.New(t)
			actual, ok := trie.Find(testCase.Input)
			and.So(actual, should.Equal, testCase.Expected)
			and.So(ok, should.Equal, testCase.OK)
		})
	}
}

func Test_SimpleTrie_Find_String(t *testing.T) {
	trie, _ := NewTrie[string, int]()
	trie.Add("Hello", 1)
	trie.Add("World", 2)
	trie.Add("Helicopter", 3)
	trie.Add("Fair", 4)
	trie.Add("Weather", 5)
	trie.Add("Whether", 6)
	trie.Add("Help", 7)
	trie.Add("", 8)

	testTable := map[string]struct {
		Input    string
		Expected int
		OK       bool
	}{
		"Hello":        {Input: "Hello", Expected: 1, OK: true},
		"Helloo":       {Input: "Helloo", Expected: 0, OK: false},
		"World":        {Input: "World", Expected: 2, OK: true},
		"W":            {Input: "W", Expected: 0, OK: false},
		"Helicopter":   {Input: "Helicopter", Expected: 3, OK: true},
		"Fair":         {Input: "Fair", Expected: 4, OK: true},
		"Weather":      {Input: "Weather", Expected: 5, OK: true},
		"Whether":      {Input: "Whether", Expected: 6, OK: true},
		"Help":         {Input: "Help", Expected: 7, OK: true},
		"empty-string": {Input: "", Expected: 8, OK: true},
		"not-in-data":  {Input: "North", Expected: 0, OK: false},
	}

	for name, testCase := range testTable {
		t.Run(name, func(t *testing.T) {
			and := assertions.New(t)
			actual, ok := trie.Find(testCase.Input)
			and.So(actual, should.Equal, testCase.Expected)
			and.So(ok, should.Equal, testCase.OK)
		})
	}
}

func Test_SimpleTrie_Find_Int8Slice(t *testing.T) {
	trie, _ := NewTrie[[]byte, int]()
	trie.Add([]byte("Hello"), 1)
	trie.Add([]byte("World"), 2)
	trie.Add([]byte("Helicopter"), 3)
	trie.Add([]byte("Fair"), 4)
	trie.Add([]byte("Weather"), 5)
	trie.Add([]byte("Whether"), 6)
	trie.Add([]byte("Help"), 7)
	trie.Add([]byte(""), 8)

	testTable := map[string]struct {
		Input    []byte
		Expected int
		OK       bool
	}{
		"Hello":        {Input: []byte("Hello"), Expected: 1, OK: true},
		"Helloo":       {Input: []byte("Helloo"), Expected: 0, OK: false},
		"World":        {Input: []byte("World"), Expected: 2, OK: true},
		"W":            {Input: []byte("W"), Expected: 0, OK: false},
		"Helicopter":   {Input: []byte("Helicopter"), Expected: 3, OK: true},
		"Fair":         {Input: []byte("Fair"), Expected: 4, OK: true},
		"Weather":      {Input: []byte("Weather"), Expected: 5, OK: true},
		"Whether":      {Input: []byte("Whether"), Expected: 6, OK: true},
		"Help":         {Input: []byte("Help"), Expected: 7, OK: true},
		"empty-string": {Input: []byte(""), Expected: 8, OK: true},
		"not-in-data":  {Input: []byte("North"), Expected: 0, OK: false},
	}

	for name, testCase := range testTable {
		t.Run(name, func(t *testing.T) {
			and := assertions.New(t)
			actual, ok := trie.Find(testCase.Input)
			and.So(actual, should.Equal, testCase.Expected)
			and.So(ok, should.Equal, testCase.OK)
		})
	}
}

func Test_SimpleTrie_Find_Int64Slice(t *testing.T) {
	trie, _ := NewTrie[[]int64, int]()
	trie.Add([]int64{1, 2, 3, 4}, 1)
	trie.Add([]int64{1, 2, 3, 5}, 2)
	trie.Add([]int64{}, 3)
	trie.Add([]int64{6, 5, 4, 3, 2, 1}, 4)

	testTable := map[string]struct {
		Input    []int64
		Expected int
		OK       bool
	}{
		"1234":        {Input: []int64{1, 2, 3, 4}, Expected: 1, OK: true},
		"1235":        {Input: []int64{1, 2, 3, 5}, Expected: 2, OK: true},
		"empty":       {Input: []int64{}, Expected: 3, OK: true},
		"654321":      {Input: []int64{6, 5, 4, 3, 2, 1}, Expected: 4, OK: true},
		"6543211":     {Input: []int64{6, 5, 4, 3, 2, 1, 1}, Expected: 0, OK: false},
		"not-in-data": {Input: []int64{6, 5, 6, 3, 6}, Expected: 0, OK: false},
	}

	for name, testCase := range testTable {
		t.Run(name, func(t *testing.T) {
			and := assertions.New(t)
			actual, ok := trie.Find(testCase.Input)
			and.So(actual, should.Equal, testCase.Expected)
			and.So(ok, should.Equal, testCase.OK)
		})
	}
}

func Test_SimpleTrie_Find_WithTransform(t *testing.T) {
	trie, _ := NewTrie[string, int](func(in byte) (out byte, use bool) {
		if in == '-' || in == '_' {
			return 0, false
		}

		if in >= 'A' && in <= 'Z' {
			return in - 'A' + 'a', true
		}

		return in, true
	})

	trie.Add("Hello", 1)
	trie.Add("World", 2)
	trie.Add("Helicopter", 3)
	trie.Add("Fair", 4)
	trie.Add("Weather", 5)
	trie.Add("Whether", 6)
	trie.Add("Help", 7)
	trie.Add("", 8)

	testTable := map[string]struct {
		Input    string
		Expected int
		OK       bool
	}{
		"Hello":          {Input: "Hello", Expected: 1, OK: true},
		"hellO":          {Input: "hellO", Expected: 1, OK: true},
		"-H-e-l-lo-":     {Input: "-H-e-l-lo-", Expected: 1, OK: true},
		"World":          {Input: "World", Expected: 2, OK: true},
		"worLd":          {Input: "worLd", Expected: 2, OK: true},
		"_World_":        {Input: "_World_", Expected: 2, OK: true},
		"Helicopter":     {Input: "Helicopter", Expected: 3, OK: true},
		"heliCopter":     {Input: "heliCopter", Expected: 3, OK: true},
		"He--licopte__r": {Input: "He--licopte__r", Expected: 3, OK: true},
		"Fair":           {Input: "Fair", Expected: 4, OK: true},
		"fAIR":           {Input: "fAIR", Expected: 4, OK: true},
		"F--a__i--r":     {Input: "F--a__i--r", Expected: 4, OK: true},
		"Weather":        {Input: "Weather", Expected: 5, OK: true},
		"weaTHer":        {Input: "weaTHer", Expected: 5, OK: true},
		"-----Weather":   {Input: "-----Weather", Expected: 5, OK: true},
		"Whether":        {Input: "Whether", Expected: 6, OK: true},
		"whether":        {Input: "whether", Expected: 6, OK: true},
		"Whether_______": {Input: "Whether_______", Expected: 6, OK: true},
		"Help":           {Input: "Help", Expected: 7, OK: true},
		"help":           {Input: "help", Expected: 7, OK: true},
		"H_-elp":         {Input: "H_-elp", Expected: 7, OK: true},
		"empty":          {Input: "", Expected: 8, OK: true},
		"_-_--_":         {Input: "_-_--_", Expected: 8, OK: true},
		"not-in-data":    {Input: "North", Expected: 0, OK: false},
	}

	for name, testCase := range testTable {
		t.Run(name, func(t *testing.T) {
			and := assertions.New(t)
			actual, ok := trie.Find(testCase.Input)
			and.So(actual, should.Equal, testCase.Expected)
			and.So(ok, should.Equal, testCase.OK)
		})
	}
}

func Benchmark_SimpleTrie(b *testing.B) {
	statesMap := map[string]int{
		"Alabama":                  0,
		"AL":                       0,
		"Kentucky":                 1,
		"KY":                       1,
		"Ohio":                     2,
		"OH":                       2,
		"Alaska":                   3,
		"AK":                       3,
		"Louisiana":                4,
		"LA":                       4,
		"Oklahoma":                 5,
		"OK":                       5,
		"Arizona":                  6,
		"AZ":                       6,
		"Maine":                    7,
		"ME":                       7,
		"Oregon":                   8,
		"OR":                       8,
		"Arkansas":                 9,
		"AR":                       9,
		"Maryland":                 10,
		"MD":                       10,
		"Pennsylvania":             11,
		"PA":                       11,
		"American Samoa":           12,
		"AS":                       12,
		"Massachusetts":            13,
		"MA":                       13,
		"Puerto Rico":              14,
		"PR":                       14,
		"California":               15,
		"CA":                       15,
		"Michigan":                 16,
		"MI":                       16,
		"Rhode Island":             17,
		"RI":                       17,
		"Colorado":                 18,
		"CO":                       18,
		"Minnesota":                19,
		"MN":                       19,
		"South Carolina":           20,
		"SC":                       20,
		"Connecticut":              21,
		"CT":                       21,
		"Mississippi":              22,
		"MS":                       22,
		"South Dakota":             23,
		"SD":                       23,
		"Delaware":                 24,
		"DE":                       24,
		"Missouri":                 25,
		"MO":                       25,
		"Tennessee":                26,
		"TN":                       26,
		"District of Columbia":     27,
		"DC":                       27,
		"Montana":                  28,
		"MT":                       28,
		"Texas":                    29,
		"TX":                       29,
		"Florida":                  30,
		"FL":                       30,
		"Nebraska":                 31,
		"NE":                       31,
		"Trust Territories":        32,
		"TT":                       32,
		"Georgia":                  33,
		"GA":                       33,
		"Nevada":                   34,
		"NV":                       34,
		"Utah":                     35,
		"UT":                       35,
		"Guam":                     36,
		"GU":                       36,
		"New Hampshire":            37,
		"NH":                       37,
		"Vermont":                  38,
		"VT":                       38,
		"Hawaii":                   39,
		"HI":                       39,
		"New Jersey":               40,
		"NJ":                       40,
		"Virginia":                 41,
		"VA":                       41,
		"Idaho":                    42,
		"ID":                       42,
		"New Mexico":               43,
		"NM":                       43,
		"Virgin Islands":           44,
		"VI":                       44,
		"Illinois":                 45,
		"IL":                       45,
		"New York":                 46,
		"NY":                       46,
		"Washington":               47,
		"WA":                       47,
		"Indiana":                  48,
		"IN":                       48,
		"North Carolina":           49,
		"NC":                       49,
		"West Virginia":            50,
		"WV":                       50,
		"Iowa":                     51,
		"IA":                       51,
		"North Dakota":             52,
		"ND":                       52,
		"Wisconsin":                53,
		"WI":                       53,
		"Kansas":                   54,
		"KS":                       54,
		"Northern Mariana Islands": 55,
		"MP":                       55,
		"Wyoming":                  56,
		"WY":                       56,
	}

	trie, _ := NewTrieFromMap(statesMap)
	lookupMap := make(map[string]int, 0)
	for name, value := range statesMap {
		lookupMap[strings.ToLower(name)] = value
	}

	provider := providers.New2(func(string, string) {})
	for name1 := range statesMap {
		for name2 := range statesMap {
			provider.Add(name1, name2)
		}
	}

	benchy.New(b, options.Medium).
		RegisterBenchmark("map", provider.WrapBenchmarkFunc(func(a, b string) {
			_ = lookupMap[strings.ToLower(string(a))] == statesMap[strings.ToLower(string(b))]
		})).
		RegisterBenchmark("simple_trie", provider.WrapBenchmarkFunc(func(a, b string) {
			v1, _ := trie.Find(a)
			v2, _ := trie.Find(b)
			_ = v1 == v2
		}) /*, options.PProfCPU*/).
		ShowMemoryStats().
		Run()
}

func Benchmark_SimpleTrie_WithTransform(b *testing.B) {
	mapped := map[string]int{
		"hello":      1,
		"world":      2,
		"helicopter": 3,
		"fair":       4,
		"weather":    5,
		"whether":    6,
		"help":       7,
		"":           8,
	}

	trie, _ := NewTrieFromMap(
		mapped,
		func(in byte) (out byte, use bool) {
			if in == '-' || in == '_' {
				return 0, false
			}

			if in >= 'A' && in <= 'Z' {
				return in - 'A' + 'a', true
			}

			return in, true
		},
	)

	provider := providers.New2(func(string, bool) {}).
		Add("Hello", true).
		Add("hellO", true).
		Add("-H-e-l-lo", true).
		Add("World", true).
		Add("worLd", true).
		Add("_World_", true).
		Add("Helicopter", true).
		Add("heliCopter", true).
		Add("He--licopte__r", true).
		Add("Fair", true).
		Add("fAIR", true).
		Add("F--a__i--r", true).
		Add("Weather", true).
		Add("weaTHer", true).
		Add("-----Weather", true).
		Add("Whether", true).
		Add("whether", true).
		Add("Whether_______", true).
		Add("Help", true).
		Add("help", true).
		Add("H_-elp", true).
		Add("", true).
		Add("_-_--_", true).
		Add("not-in-data", false)

	var ok bool
	benchy.New(b, options.Medium).
		RegisterBenchmark("trie", provider.WrapBenchmarkFunc(func(s string, expected bool) {
			_, ok = trie.Find(s)
			if ok != expected {
				b.Logf("failed on 'trie' with input '%s'", s)
				b.Fail()
			}
		})).
		RegisterBenchmark("map", provider.WrapBenchmarkFunc(func(s string, expected bool) {
			_, ok = mapped[strings.ReplaceAll(strings.ReplaceAll(strings.ToLower(s), "_", ""), "-", "")]
			if ok != expected {
				b.Logf("failed on 'map' with input '%s'", s)
				b.Fail()
			}
		})).
		Run()

	ok = !ok
}
