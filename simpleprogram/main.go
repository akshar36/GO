package main

import "fmt"

func main(){
	numbers := []int{}
	for i :=0; i< 10; i++ {
		numbers = append(numbers, i+1)
	}

	for _, number  := range numbers{
		if number % 2 != 0{
			fmt.Printf("The number %v is odd \n", number)
		}else {
			fmt.Printf("The number %v is even \n", number)
		}
	}
}