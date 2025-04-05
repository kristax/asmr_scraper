package guess_epsides

import (
	"fmt"
	"regexp"
	"sort"
	"strconv"
)

var numRegex = regexp.MustCompile(`[0-9０-９]+`) // 匹配ASCII数字和全角数字

type NoNumbersError struct{ s string }

func (e NoNumbersError) Error() string {
	return fmt.Sprintf("no numbers found in: %q", e.s)
}

func MapOrders(names []string) (map[string]*int, error) {
	result := make(map[string]*int, len(names))
	validEntries := make([]struct {
		name   string
		nums   []int
		hasNum bool
	}, 0, len(names))

	for _, name := range names {
		nums, err := extractNumbers(name)
		switch err.(type) {
		case NoNumbersError:
			result[name] = nil
		case nil:
			validEntries = append(validEntries, struct {
				name   string
				nums   []int
				hasNum bool
			}{name, nums, true})
			result[name] = new(int)
		default:
			return nil, err
		}
	}

	if len(names) == 1 {
		name := names[0]
		if result[name] == nil {
			val := 1
			result[name] = &val
		}
		return result, nil
	}

	var validNumbers [][]int
	var validNames []string
	for _, entry := range validEntries {
		if entry.hasNum {
			validNumbers = append(validNumbers, entry.nums)
			validNames = append(validNames, entry.name)
		}
	}

	switch len(validNumbers) {
	case 0:
		return result, fmt.Errorf("no valid numeric entries found")
	case 1:
		val := validNumbers[0][0]
		result[validNames[0]] = &val
		return result, nil
	default:
		//baseLen := len(validNumbers[0])
		//for _, nums := range validNumbers {
		//	if len(nums) != baseLen {
		//		return nil, fmt.Errorf("inconsistent number counts")
		//	}
		//}

		// 新增排序逻辑（基于所有提取到的数字序列）
		type entry struct {
			name string
			nums []int
		}
		entries := make([]entry, len(validNames))
		for i := range validNames {
			entries[i] = entry{validNames[i], validNumbers[i]}
		}

		sort.Slice(entries, func(i, j int) bool {
			a, b := entries[i].nums, entries[j].nums
			minLen := len(a)
			if len(b) < minLen {
				minLen = len(b)
			}
			for k := 0; k < minLen; k++ {
				if a[k] != b[k] {
					return a[k] < b[k]
				}
			}
			return len(a) < len(b)
		})

		for i := range entries {
			validNames[i] = entries[i].name
			validNumbers[i] = entries[i].nums
		}

		pos, err := findIncreasingPosition(validNumbers)
		if err != nil {
			return nil, err
		}

		for i, name := range validNames {
			val := validNumbers[i][pos]
			result[name] = &val
		}
		return result, nil
	}
}

func extractNumbers(s string) ([]int, error) {
	matches := numRegex.FindAllString(s, -1) // 恢复提取所有数字
	if len(matches) == 0 {
		return nil, NoNumbersError{s}
	}

	nums := make([]int, len(matches))
	for i, m := range matches {
		normalized := fullwidthToHalf(m)
		n, err := strconv.Atoi(normalized)
		if err != nil {
			return nil, fmt.Errorf("invalid number %q: %v", m, err)
		}
		nums[i] = n
	}
	return nums, nil
}

// 全角数字转半角函数
func fullwidthToHalf(s string) string {
	var result []rune
	for _, r := range s {
		switch {
		case r >= '０' && r <= '９': // 全角数字范围
			result = append(result, rune('0'+r-'０'))
		default:
			result = append(result, r)
		}
	}
	return string(result)
}

// 改进后的递增位置检测
func findIncreasingPosition(numbers [][]int) (int, error) {
	maxLen := 0
	for _, nums := range numbers {
		if len(nums) > maxLen {
			maxLen = len(nums)
		}
	}

	for pos := 0; pos < maxLen; pos++ {
		valid := true
		for i := 0; i < len(numbers)-1; i++ {
			// 检查所有条目在当前位是否都有数字
			if pos >= len(numbers[i]) || pos >= len(numbers[i+1]) {
				valid = false
				break
			}
			if numbers[i][pos] >= numbers[i+1][pos] {
				valid = false
				break
			}
		}
		if valid {
			return pos, nil
		}
	}
	return 0, fmt.Errorf("no increasing sequence found %v", numbers)
}
