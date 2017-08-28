package cap

type Looper struct{

}

func NewLooper() *Looper{
	looper := &Looper{}
	return looper
}

func (this *Looper) While(predicate func()bool,body func()error)error{
	if predicate() == true {
		for{
			if predicate() == true {
				err := body()
				if err != nil {
					return err
				}
			}else{
				break
			}
		}
	}
}