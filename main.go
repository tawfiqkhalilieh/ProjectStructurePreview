package main

import "os"
import "fmt"
import "strings"
import "io/ioutil"
import "path/filepath"

type Stack struct {
	// https://medium.com/@crusty0gphr/tricky-golang-interview-questions-part-2-bigo-of-len-2bc73ebb3090
	elements []string
	length int
}

func NewStack() *Stack {
	s := Stack{ elements: []string{}, length: 0 }
	return &s;
}

func (stack *Stack) Pop() string {
	if stack.length == 0 {
		panic("Stack is empty")
		return "";
	}
	var value string = stack.elements[stack.length-1]
	stack.elements = stack.elements[:stack.length-1]

	stack.length--;
	return value
}


func (stack *Stack) Push(value string) {
	// The time complexity of using Append is surprisingly O(1)
	// https://glucn.com/posts/2019-07-04-golang-the-time-complexity-of-append
	stack.length++;
	stack.elements = append(stack.elements, value)
} 

func (stack *Stack) IsEmpty() bool {
	return stack.length == 0
}

func (stack *Stack) Head() string { 
       if stack.length == 0 {
	       return ""
       } else {
	       return stack.elements[stack.length-1]
       }
}

var path string = "."

func main() {
	if len(os.Args) > 1 {
           path = os.Args[1]
	}

	listDir(path, 0);
	
}


// documents: https://pkg.go.dev/os
func listDir(path string, spaces int) {
	entries, err := os.ReadDir(path);

	ignores := processGitIgnore();

	if err != nil {
		fmt.Println(err)
	}
	
	for _, e:= range entries {
	
		stk := processPathLine(path + "/" + e.Name());
		var flag bool = false;

		for _, ignore := range ignores {
			ignoreCopy := ignore 

			if isIgnored(stk, ignoreCopy) {
				flag = true
				break
			}

			defer func(ig Stack) {
			}(ignoreCopy)
		}

		

		if flag { 
			continue;
		}

		if strings.HasPrefix(e.Name(), ".git") {
			continue
		} else if e.IsDir() {
			fmt.Println(strings.Repeat("     ", spaces) +  e.Name() + "/");
			listDir(path + "/" + e.Name(), spaces + 1);
		} else {
			fmt.Println(strings.Repeat("     ", spaces) +  e.Name());
		}

	}
}


func processGitIgnore() []Stack {
	data, err := ioutil.ReadFile(path + "/.gitignore")
	if err != nil {
		// fmt.Println("No .gitignore found")
	}	
	
	lines := strings.Split(string(data), "\n")

        var ignores = []Stack {}
  	index := 0

	for _, line := range lines {
		if line == "" || strings.Contains(strings.Trim(line, " "), "#") {
			continue;
		} else {
			stk := processPathLine(line);
			ignores = append(ignores, stk) 
			index++;
		}


	}	


	return ignores
}


func processPathLine(line string) Stack {
	stack := NewStack();
	ln := strings.Split(line, "/");
		
	for _, v := range ln {
		stack.Push(v)
	}
	
	return *stack

}

func (s Stack) Join() string {
	return strings.Join(s.elements, "/")
}

func matchGitignorePattern(ignore Stack, file Stack) bool {
	pattern := ignore.Join()
	lpath := path + "/" + file.Join()
	

	lignore := ignore
	lfile := file

	// (*, ?, [])
	matched, _ := filepath.Match(pattern, lpath) // || strings.Contains(pattern, "*") // || strings.Contains(pattern, path)
	if matched {
		return true
	}

	if strings.HasSuffix(pattern, "/") {
		prefix := strings.TrimSuffix(pattern, "/")
		if strings.HasPrefix(lpath, prefix+"/") {
			return true
		}
	}

	if strings.HasPrefix(lpath, pattern+"/") {
		return true
	}

	if pattern == lpath {
		return true
	}

	for !lfile.IsEmpty() && !lignore.IsEmpty() {
		file_top := strings.Trim(lfile.Pop(), " ")
		ignore_top := strings.Trim(lignore.Pop(), " ")

		if file_top == "" || ignore_top == "" {
			continue
		}

		if strings.Contains(file_top, ignore_top) {
			return true
		}

		if strings.Contains(lpath, ".git") {
			return true
		}


	}

	return false
}

func isIgnored(file Stack, ignore Stack) bool {
	if matchGitignorePattern(ignore, file) {
		return true
	}
	return false
}


