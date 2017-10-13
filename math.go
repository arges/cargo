package main

import (
	"html/template"
	"math/rand"
	"net/http"
	"strconv"
	"time"
)

type Context struct {
	Car       *Cargo
	Current   Problem
	Correct   int
	Incorrect int
}

type Problem struct {
	Operand1 int
	Operand2 int
	Operator string
	Answer   int
}

// NewProblem returns an addition problem that sums to a number less than 10
func NewProblem() Problem {
	var p Problem
	p.Operand1 = rand.Intn(10)
	p.Operand2 = rand.Intn(10 - p.Operand1)
	p.Operator = "+"
	p.Answer = p.Operand1 + p.Operand2
	return p
}

// ServeHTTP will render the math page.
func (c *Context) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		c.Current = NewProblem()
		t, _ := template.ParseFiles("math.html")
		t.Execute(w, c)
	} else {
		r.ParseForm()
		result, _ := strconv.Atoi(r.Form["Result"][0])
		if result == c.Current.Answer {
			c.Correct++
			c.Car.Move(127, 127)
			time.Sleep(2 * time.Second)
			c.Car.Move(0, 0)
		} else {
			c.Incorrect++
			c.Car.Move(-127, 0)
			time.Sleep(1 * time.Second)
			c.Car.Move(0, 0)
		}
		http.Redirect(w, r, "/math", http.StatusSeeOther)
	}
}
