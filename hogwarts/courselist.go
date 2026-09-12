//go:build !solution

package hogwarts

import "strings"

func GetPath(path map[string]struct{}) string {
	var err strings.Builder
	for c := range path {
		err.WriteString(c)
		err.WriteString(" -> ")
	}
	return err.String()
}

func GetCourseList(prereqs map[string][]string) []string {
	uniqueCourses := make(map[string]struct{})
	courses := make([]string, 0)

	var GetCourseListRec func(current string, path map[string]struct{})
	GetCourseListRec = func(
		current string,
		path map[string]struct{},
	) {
		path[current] = struct{}{}
		for _, c := range prereqs[current] {
			if _, ok := uniqueCourses[c]; !ok {
				if _, ok := path[c]; ok {
					panic(GetPath(path) + c)
				}
				path[c] = struct{}{}
				GetCourseListRec(c, path)
			}
		}
		if _, ok := uniqueCourses[current]; !ok {
			courses = append(courses, current)
		}
		uniqueCourses[current] = struct{}{}
	}

	for c := range prereqs {
		path := make(map[string]struct{})
		GetCourseListRec(c, path)
	}

	return courses
}
