package main

import ("fmt"
		"net"
		"net/rpc"
		"net/rpc/jsonrpc"
		"net/http"
		"io/ioutil"
		"encoding/json"
		"strings"
		"strconv")

type Server struct{}
var id int 
var store map[int]string
type Request struct{

	StockSymbolAndPercentage []NestedRequest `json:"stockSymbolAndPercentage"`
	Budget float32 `json:"budget"`
}

type SecReq struct{

	Tradeid int `json:"tradeid"`
}
type NestedRequest struct{

	Fields ActualFields `json:"fields"`
}

type ActualFields struct{

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


func (this *Server) PrintMessage(msg string,rep *string) error{
		var jsonInt interface{}

		var structResponse Response

		var jsonMsg Request

		var companyName string

		var RemainderVal float32=0.0

		json.Unmarshal([]byte(msg),&jsonMsg)

		for _, i:= range jsonMsg.StockSymbolAndPercentage{
			companyName += i.Fields.Name +","
		}
		companyName=strings.Trim(companyName,",")

		response,err:= http.Get("http://finance.yahoo.com/webservice/v1/symbols/"+companyName+"/quote?format=json")

		if(err!=nil) {
			fmt.Println(err)
		}else {
			defer response.Body.Close()
			contents,err:= ioutil.ReadAll(response.Body)
			json.Unmarshal(contents,&jsonInt)
			for i,index := range (jsonInt.(map[string]interface{})["list"]).(map[string]interface{})["resources"].([]interface{}){ 
				price := index.(map[string]interface{})["resource"].(map[string]interface{})["fields"].(map[string]interface{})["price"]
				priceconv,_ := strconv.ParseFloat(price.(string),64)
				RemainderVal1:=(float64(jsonMsg.StockSymbolAndPercentage[i].Fields.Percentage) * float64(jsonMsg.Budget))/100
				name := index.(map[string]interface{})["resource"].(map[string]interface{})["fields"].(map[string]interface{})["symbol"]
				number := int( RemainderVal1/priceconv)
		
				RemainderVal = RemainderVal +  (float32(priceconv)*float32(number))

				structResponseFieldsFinal:=ResponseFieldsFinal{Name:name.(string),Number:number,Price:strconv.FormatFloat(priceconv,'f',-1,64)}

				structNestedResponse := NestedResponse{ResponseFields:structResponseFieldsFinal}

				structResponse.Stocks = append(structResponse.Stocks,structNestedResponse)
			}
			RemainderVal=jsonMsg.Budget-RemainderVal

			result1 := &Response{

    		TradeId:id,

        	Stocks: structResponse.Stocks,

        	UnvestedAmount:RemainderVal} //Map the values to Request structure

    		result2, _ := json.Marshal(result1) //Convert the Request to JSON

    		*rep = string(result2)

			store[id]=string(result2)

			id++
			if(err!=nil){
				fmt.Println(err)
			}
				
		}
		
		return nil
}

func (this *Server) LossOrGain(msg string,rep *string) error{
	var jsonReq SecReq

	var jsonMsg Response

	var jsonInt interface{}

	var companyName string

	var price []float64

	var structSecondResponse SecondResponse

	json.Unmarshal([]byte(msg),&jsonReq)

	tradeid:= jsonReq.Tradeid

	data:= store[tradeid]

	json.Unmarshal([]byte(data),&jsonMsg)

	for _,index:= range jsonMsg.Stocks{

		companyName += index.ResponseFields.Name +","
	}
	companyName=strings.Trim(companyName,",")

	response,err:= http.Get("http://finance.yahoo.com/webservice/v1/symbols/"+companyName+"/quote?format=json")

	if(err!=nil){
		fmt.Println(err)
	}else{
		defer response.Body.Close()

		contents,_:= ioutil.ReadAll(response.Body)

		json.Unmarshal(contents,&jsonInt)

		for _,index := range (jsonInt.(map[string]interface{})["list"]).(map[string]interface{})["resources"].([]interface{}){ 
				priceconv,_ := strconv.ParseFloat((index.(map[string]interface{})["resource"].(map[string]interface{})["fields"].(map[string]interface{})["price"]).(string),64)
				price = append(price,priceconv)
			}
		var value float32=0.0

		var strprice string

		for i,index := range jsonMsg.Stocks{

				temp,_:= strconv.ParseFloat(index.ResponseFields.Price,64)
				fmt.Println(price[i],temp)
				if price[i] > temp{
					strprice = "$+"+strconv.FormatFloat(price[i],'f',-1,64)

				}
				if price[i] < temp {
					strprice = "$-"+strconv.FormatFloat(price[i],'f',-1,64)
				}else {
					strprice = "$"+strconv.FormatFloat(price[i],'f',-1,64)
				}
				structResponseFieldsFinal:=ResponseFieldsFinal{Name:index.ResponseFields.Name,Number:index.ResponseFields.Number,Price:strprice}

				structNestedResponse := NestedResponse{ResponseFields:structResponseFieldsFinal}

				structSecondResponse.Stocks = append(structSecondResponse.Stocks,structNestedResponse)

				value = value + (float32(index.ResponseFields.Number) * float32(price[i]))
		}
		result1 := &SecondResponse{
    	CurrentMarketValue:value,

        Stocks: structSecondResponse.Stocks,

        UnvestedAmount:jsonMsg.UnvestedAmount} //Map the values to Request structure

    	result2, _ := json.Marshal(result1) //Convert the Request to JSON
    	
    	*rep = string(result2)
	
	}		
	return nil
}

func main(){
	fmt.Println("Launching server on port 3000...")// Start to Launch Server

	id++
	store =make(map[int]string)

	rpc.Register(new(Server)) // Register for new server
	hear,err:= net.Listen("tcp",":3000") // Registered on port 1000
	// error handling
	if(err!=nil){
		fmt.Println(err)
		return
	}
	//listen and accept
	for{
		c,error:= hear.Accept()
		if(error!=nil){
			continue
		}
		go jsonrpc.ServeConn(c)
	}

}
