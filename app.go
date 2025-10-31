package main
import (
	"fmt"
	"net/http"
	"io/ioutil"
	"bytes"
	"encoding/json"
)

type PayLoad struct{
	Status string `json:"status"`;
	Data string `json:"data"`;
}

func postRequest(url string, input PayLoad)(output PayLoad, err error){

	jsonData, err := json.Marshal(input); // convert struct to JSON
	if err !=nil{ 
		return PayLoad{}, err; 
	}
	req,err:=http.NewRequest("POST",url,bytes.NewBuffer(jsonData)); // create new POST request
	if err !=nil{
		return PayLoad{}, err;
	}
	req.Header.Set("Content-Type","application/json"); // set content type to JSON
	client := &http.Client{}; // create new HTTP client
	resp, err := client.Do(req); // send request
	if err !=nil{ // 
		return PayLoad{}, err; 
	}
	defer resp.Body.Close(); // close response body when done
	body,err:=ioutil.ReadAll(resp.Body); // read response body
	if err !=nil{ 
		return PayLoad{}, err;
	}
	
	err = json.Unmarshal(body,&output); // convert JSON response to struct
	if err !=nil{
		return PayLoad{}, err;
	}
	return output,nil;

func main(){
	url:="https://faas-ams3-2a2df116.doserverless.co/api/v1/web/fn-2f25c703-39e1-4645-abdb-ee0c9d620425/rst/app"; // define URL
	input,err:=postRequest(url,PayLoad{Status:"active",Data:"sample data"}); // call postRequest function:
	if err !=nil{
		fmt.Println("Error:",err);
	}
	fmt.Println("Data:",input.Data);
	fmt.Println("Status:",input.Status); 
