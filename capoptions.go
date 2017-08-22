package cap

type CapOptions struct{
	options []capOption
}

type capOption struct{
	Name string
	Value interface{}
}

func (capOptions *CapOptions) Add(name string, value interface{}){
	option := capOption{}
	option.Name = name
	option.Value = value
	capOptions.options = append(capOptions.options, option)
}

func (capOptions *CapOptions) Get(name string)(interface{},error){
	var value interface{}
	for i := 0; i < len(capOptions.options); i++ {
		_option := capOptions.options[i]
		if(_option.Name == name){
			value = _option.Value
			break
		}else{
			continue
		}
	}
	if value != nil{
		return value, nil
	}else{
		return nil, NewCapError("Could not find key [" + name + "] in configured options.");
	}
}

func (capOptions *CapOptions) UseMySql(connectionString string){
	capOptions.Add("CONNECTION_STRING", connectionString)
}

func (capOptions *CapOptions) GetConnectionString() (string, error){
	value, err := capOptions.Get("CONNECTION_STRING")
	if err != nil{
		return "", err
	}else{
		return value.(string), nil
	}
}