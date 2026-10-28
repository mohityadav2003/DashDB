package main

import "errors"

func decode(data []byte) (any, error) {
	if len(data) == 0 {
		return nil, errors.New("No Data")
	}
	value, _, err := decodeOne(data)
	return value, err
} 

func decodeOne(data []byte) (any,int, error){
	if(len(data) == 0){
		return nil, 0, errors.New("No data")
	}
	switch data[0]{
		case '+':
			return readSimpleString(data)
		case '-':
			return readSimpleError(data)
		case '$':
			return readBulkString(data)
		case ':':
			return readSimpleInteger(data)
		case '*':
			return readArray(data)
	}
	return nil, 0, errors.New("Unknown DataType in Command")
}

func readSimpleString(data []byte) (string, int, error){
	pos:=1

	interpretedString := ""

	for ; data[pos]!='\r' ; pos++	{
		interpretedString = interpretedString + string(data[pos])
	}
	return interpretedString, pos+2, nil
}

func readSimpleInteger(data []byte) (int, int, error){
	pos:= 1
	var interpretedInteger int = 0
	for ; data[pos] != '\r' ; pos++ {
		interpretedInteger = interpretedInteger*10 + int(data[pos]-'0')
	}
	return interpretedInteger, pos+2, nil
}

func readSimpleError(data []byte) (string, int, error){
	return readSimpleString(data)
}

func readBulkString(data []byte) (string, int, error) {
	length, delta, err := readLength(data)
	if err != nil {
		return "", 0, errors.New("something went wrong")
	}
	content := string(data[delta : delta+length])
	return content, delta + length + 2, nil
}

func readLength(data []byte) (int, int, error){
	var len int = 0

	var pos int = 1

	for ; data[pos] >= '0' && data[pos] <='9' ;pos++{
		len = (len*10) + int(data[pos] - '0')
	}
	return len, pos+2, nil
}

func readArray(data []byte) (any, int, error){
	var pos int = 1 

	length, delta , _ := readLength(data)
	pos = delta

	var array []any = make([]any, length)

	var i int = 0

	for ; i<length; i++{
		elem, delta, _ := decodeOne(data[pos:])
		array[i] = elem
		pos += delta
	}
	return array, pos, nil
}

func encodeToSimpleString(input string, err error)(string, error){
	if(err != nil){
		return ("-" + err.Error() +"\r\n"), nil
	}
	return ("+" + input + "\r\n"), nil
}
