package main

import (
	"fmt"
)

// state struct
type state struct {
	symbol rune
	edge1  *state
	edge2  *state
} // state struct

// helper struct
type nfa struct {
	initial *state
	accept  *state
} // nfa struct

func poregtonfa(postfix string) *nfa {
	// an array of pointers to nfa's, the curley braces denote we want an empty one
	nfastack := []*nfa{}

	// loop throught the regular expression one rune at a time
	for _, r := range postfix {
		switch r {
		case '.':
			// get the last thing off the stack and store in frag2
			frag2 := nfastack[len(nfastack)-1]
			// get rid of the last thing on the stack, because it's already on frag2
			nfastack = nfastack[:len(nfastack)-1]
			// get the last thing off the stack and store in frag1
			frag1 := nfastack[len(nfastack)-1]
			// get rid of the last thing on the stack, because it's already on frag1
			nfastack = nfastack[:len(nfastack)-1]

			// join the fragments together by setting the accept state of frag1
			// to the initial state of frag2
			frag1.accept.edge1 = frag2.initial

			// then we append the new nfa we created above to the nfastack
			nfastack = append(nfastack, &nfa{initial: frag1.initial, accept: frag2.accept})
		case '|':
			// get the last thing off the stack and store in frag2
			frag2 := nfastack[len(nfastack)-1]
			// get rid of the last thing on the stack, because it's already on frag2
			nfastack = nfastack[:len(nfastack)-1]
			// get the last thing off the stack and store in frag1
			frag1 := nfastack[len(nfastack)-1]
			// get rid of the last thing on the stack, because it's already on frag1
			nfastack = nfastack[:len(nfastack)-1]

			// new accept state
			accept := state{}
			// the new initial state where it's two edges point to the two fragments initial states
			initial := state{edge1: frag1.initial, edge2: frag2.initial}
			// we then set the accept states of the fragments to the new accept state
			frag1.accept.edge1 = &accept
			frag2.accept.edge1 = &accept

			// then we append the new nfa accept state and initial state we created above to the nfastack
			nfastack = append(nfastack, &nfa{initial: &initial, accept: &accept})
		case '*':
			// get the last thing off the stack and store in frag1
			frag := nfastack[len(nfastack)-1]
			// get rid of the last thing on the stack, because it's already on frag1
			nfastack = nfastack[:len(nfastack)-1]

			accept := state{}
			// the new initial state that points to the initial of the fragment at edge1
			// and points to the new accept state at edge2
			initial := state{edge1: frag.initial, edge2: &accept}
			// join the fragment edge1 to it's initial state
			frag.accept.edge1 = frag.initial
			// join the fragment edge2 to the new accept state
			frag.accept.edge2 = &accept

			// then we append the new nfa accept state and initial state we created above to the nfastack
			nfastack = append(nfastack, &nfa{initial: &initial, accept: &accept})
		default:
			// new empty accept state
			accept := state{}
			// new initial state with the symbol value of r
			// and it's only edge points at the new accept state
			initial := state{symbol: r, edge1: &accept}

			// appending the new nfa to the stack
			nfastack = append(nfastack, &nfa{initial: &initial, accept: &accept})
		} // switch
	} // for

	if len(nfastack) != 1 {
		fmt.Println("Uh oh:", len(nfastack), nfastack)
	} // if

	// the only item will be the actual nfa that you want to return
	return nfastack[0]
} // poregtonfa

// addState helper function
func addState(l []*state, s *state, a *state) []*state {
	l = append(l, s)

	if s != a && s.symbol == 0 {
		l = addState(l, s.edge1, a)
		if s.edge2 != nil {
			l = addState(l, s.edge2, a)
		} // if
	} // if

	return l
} // addState

func pomatch(po string, s string) bool {
	ismatch := false

	// should convert from infix to postfix here using the shunting algorithm

	// creating a postfix nfa, based on the postfix string
	ponfa := poregtonfa(po)

	// an array of pointers that track the current positions in the nfa
	current := []*state{}
	//
	next := []*state{}

	// setting the current state for the first time
	current = addState(current[:], ponfa.initial, ponfa.accept)

	for _, r := range s {
		for _, c := range current {
			if c.symbol == r {
				next = addState(next[:], c.edge1, ponfa.accept)
			} // if
		} // inner for

		// when a character is read set current to the values of next
		// and reset the next array to a blank state pointer array
		current, next = next, []*state{}
	} // outer for

	// check to see if the postfix nfa's accept state is in current
	for _, c := range current {
		if c == ponfa.accept {
			ismatch = true
			break
		} // if
	} // for

	return ismatch
}

func main() {
	//nfa := poregtonfa("ab.c*|")
	//fmt.Println(nfa)

	fmt.Println(pomatch("ab.c*|", "cccc"))
} // main
