package main

import ("fmt"
        "net/rpc/jsonrpc"
        "strings"
        "bufio"
        "os"
        "strconv"
        "encoding/json")

type Request struct{
	StockSymbolAndPercentage []NestedRequest `json:"stockSymbolAndPercentage"`
	Budget float32 `json:"budget"`
}

type SecReq struct{
	Tradeid int `json:"tradeid"`
}
type NestedRequest struct{
	Fields Input `json:"fields"`
}

type Input struct{

	Name string `json:"name"`

	Percentage int `json:"perecentage"`
}
type Response struct{
	Stocks []NestedResponse `json:"stocks"`

	TradeId int `json:"tradeid"`

	UnvestedAmount float32 `json:"unvestedAmount"`
}

type NestedResponse struct{

	ResponseFields ResponseFieldsFinal `json:"fields"`
}

type ResponseFieldsFinal struct{

	Name string `json:"name"`

	Number int `json:"number"`

	Price string `json:"price"`
}


type SecondResponse struct{

	Stocks []NestedResponse `json:"stocks"`

	CurrentMarketValue float32 `json:"currentMarketValue"`
	
	UnvestedAmount float32 `json:"unvestedAmount"`
}


func SeePortfolio(){

	c,err:= jsonrpc.Dial("tcp","127.0.0.1:3000")
	if err!=nil{
		fmt.Println(err)
		return
	}
	structSRequest:=new(SecReq)

	fmt.Println("Enter the request")

	var sRequest string

	fmt.Scanf("%s",&sRequest)

	sRequest= strings.Replace(sRequest,"\"","",-1)

	newsRequest:=strings.SplitN(sRequest,":",-1)

	structSRequest.Tradeid,_= strconv.Atoi(newsRequest[1])
	result3 := &SecReq{
		Tradeid: structSRequest.Tradeid}

	result4,_:= json.Marshal(result3)

	var jsonMsg2 SecondResponse

	var rep string

	err = c.Call("Server.LossOrGain",string(result4),&rep)

	var output string

	output = "\"stocks\":"

	json.Unmarshal([]byte(rep),&jsonMsg2)

	for _, i:= range jsonMsg2.Stocks{

		output += i.ResponseFields.Name +":"+strconv.Itoa(i.ResponseFields.Number)+":"+i.ResponseFields.Price+","
	}
	output=strings.Trim(output,",")

	output+="\"\n\"currentMarketValue\":$"+strconv.FormatFloat(float64(jsonMsg2.CurrentMarketValue),'f',-1,32)

	output+="\n\"unvestedAmount\":$"+strconv.FormatFloat(float64(jsonMsg2.UnvestedAmount),'f',-1,32)

	if err!=nil {

		fmt.Println(err)

	}else{

		fmt.Println("\nResponse:\n")

		fmt.Println(output)
	}
}


func PurchaseStocks(){

	c,err:= jsonrpc.Dial("tcp","127.0.0.1:3000")

	if err!=nil{

		fmt.Println(err)

		return
	}
	var rep string

	var structRequest Request

	var msg,data,newData []string

	fmt.Println("Enter reuest in the correct form ")     //"stockAndPercent":YHOO:50%,GOOG:50% "budget":2400

	in := bufio.NewReader(os.Stdin)

	line, err := in.ReadString('\n') // Read Input

	msg = strings.SplitN(line," ",-1)

	data = strings.SplitN(msg[0],":",2)

	newData = strings.SplitN(msg[1],":",2)

	bValue,err:=strconv.ParseFloat(strings.TrimSpace(newData[1]),64)

	data[1]= strings.Replace(data[1],"\"","",-1)

	data[1]= strings.Replace(data[1],"%","",-1)

	fields := strings.SplitN(data[1],",",-1)

	for _,index:=range fields{

			c:= strings.SplitN(index,":",-1)

			a,_:=strconv.Atoi(c[1])

			structFields := Input{Name:c[0],Percentage:a} 

			structNestedRequest := NestedRequest {Fields:structFields}

			structRequest.StockSymbolAndPercentage =append(structRequest.StockSymbolAndPercentage,structNestedRequest)
	}
	result1 := &Request{

    	Budget:float32(bValue),

        StockSymbolAndPercentage: structRequest.StockSymbolAndPercentage} //Map the values to Request structure

    result2, _ := json.Marshal(result1) //Convert the Request to JSON

	err = c.Call("Server.PrintMessage",string(result2),&rep)

	var jsonMsg Response

	var output string

	output = "\"tradeid\":"

	json.Unmarshal([]byte(rep),&jsonMsg)

	output+=strconv.Itoa(jsonMsg.TradeId)+"\n"+"\"stocks\":\""

	for _, i:= range jsonMsg.Stocks{

		output += i.ResponseFields.Name +":"+strconv.Itoa(i.ResponseFields.Number)+":"+"$"+i.ResponseFields.Price+","
	}

	output=strings.Trim(output,",")

	output+="\"\n\"unvestedAmount\":$"+strconv.FormatFloat(float64(jsonMsg.UnvestedAmount),'f',-1,32)	

	if err!=nil {

		fmt.Println(err)

	}else{
		
		fmt.Println("\nResponse:\n")

		fmt.Println(output)
	}
}



func main(){
	fmt.Println("Enter your choice\nPress 1 to buy stocks\nPress 2 for checking your portfolio")
	var choice int64 
	fmt.Scanf("%d\n",&choice)
	switch choice{

		case 1:
			PurchaseStocks()
			break
		case 2:
			SeePortfolio()
			break
		default:
			fmt.Println("Invalid Choice please enter the correct choice press 1 for buying stock and Press 2 for Portfolio Info.")
			break
		}
}