package main

import (
	"fmt"
)

type (
	collection interface {
		createIterator() iterator
	}

	iterator interface {
		hasNext() bool
		getNext() *user
	}

	userCollection struct {
		users []*user
	}

	userIterator struct {
		index int
		users []*user
	}

	user struct {
		name string
		age  int
	}
)

func (u *userCollection) createIterator() iterator {
	return &userIterator{
		users: u.users,
	}
}

func (u *userIterator) hasNext() bool {
	return u.index < len(u.users)
}

func (u *userIterator) getNext() *user {
	if !u.hasNext() {
		return nil
	}

	user := u.users[u.index]
	u.index++
	return user
}

func main() {
	user1 := &user{
		name: "Stacy",
		age:  28,
	}

	user2 := &user{
		name: "Andy",
		age:  25,
	}

	user3 := &user{
		name: "Mary",
		age:  30,
	}

	userCollection := &userCollection{
		users: []*user{user1, user2, user3},
	}

	iterator := userCollection.createIterator()

	for iterator.hasNext() {
		user := iterator.getNext()
		fmt.Printf("User is %+v\n", user)
	}
}
