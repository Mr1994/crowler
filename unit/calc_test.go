package unit

import (
	"fmt"
	"testing"
)

func TestCalcTest1(t *testing.T) {
	testSlice := []int{0, 1, 2, 3, 4, 5, 6, 7, -1, -2, -3, -4, -1}
	sumSliceArray := [][]int{}
	for i := 0; i < len(testSlice); i++ {
		// 固定一个数
		num1 := testSlice[i]
		// 双指针查找另外两个数
		left, right := i+1, len(testSlice)-1
		for left < right {
			num2 := testSlice[left]
			num3 := testSlice[right]
			sum := num1 + num2 + num3
			if sum == 0 {
				sumSliceArray = append(sumSliceArray, []int{num1, num2, num3})
				// 如果找到了一个解，双指针同时收缩，寻找其他解
				left++
				right--
			} else if sum < 0 {
				left++
			} else {
				right--
			}
		}
	}
	fmt.Println(sumSliceArray)
}



func quickSort(arr []int, l, r int)[]int {
	if l >= r {
		return []int{}
	}

	pivot := arr[l] // 以第一个元素作为 pivot
	left, right := l, r
	for left < right {
		if left < right && arr[right] >= pivot {
			right--
		}
		arr[left] = arr[right]
		if left < right && arr[left] <= pivot {
			left++
		}
		arr[right] = arr[left]
	}
	arr[left] = pivot

	// 把左半部分进行排序
	quickSort(arr, l, left-1)
	// 把右半部门进行排序
	quickSort(arr, left+1, r)
	return arr
}
func TestQuickSort(t *testing.T) {
	arr := []int{3, 1, 4, 1, 5, 9, 2, 6}
	d := quickSort(arr, 0, len(arr)-1)
	fmt.Println(d)
}