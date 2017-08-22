package cap

import (  
	"time"  
   )  

type CapError struct {  
	When time.Time  
	What string  
}

func NewCapError(message string) CapError{
	err := CapError{}
	err.What = message
	err.When = time.Now()
	return err
}

func (err CapError) Error() string{
	return err.What
}