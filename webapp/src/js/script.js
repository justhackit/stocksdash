function loginClicked(){
    let validationResult = validate()
    if(validationResult !==false){
        console.log("All validations passed")
        let postBody = JSON.stringify(validationResult)
        $ajaxUtils.sendPOSTRequest("https://ajaysquare.com/auth-api/login",postBody,function(loginHttpRespObj){
            if(loginHttpRespObj.status === 200)
        {
            console.log("Login is successfull")
            let respObj = JSON.parse(loginHttpRespObj.responseText)
            let accessToken = respObj.data.access_token
            document.accessToken = accessToken
            let refreshToken = respObj.data.refresh_token;
            //!!!Storing the refresh token in local storage
            window.localStorage.setItem("refrt",refreshToken);
            refreshDashboard();
            //Refresh the stocks dashboard for every ? secs
            window.refreshIntervalId=setInterval(refreshDashboard,1*60*1000);
            //Get new access token for every ? secs
            window.tokenIntervalId=setInterval(getNewAccessToken,14*60*1000);
        }else{
            console.warn("Unable to authorize the user!")
            let respObj = JSON.parse(loginHttpRespObj.responseText)
            let errMsg = respObj.message
            document.getElementById("mainMsg").style.color = "red"
            document.getElementById("mainMsg").textContent="Login Failed :"+errMsg
        }
        })
    }else {
        console.warn("Input validations failed")
    }
}


function getNewAccessToken(){
    console.log(new Date()+" : Getting new access token using refresh token..")
    $ajaxUtils.sendGETRequest("https://ajaysquare.com/auth-api/refresh-token",function(getHttpReq){
        let respObj = JSON.parse(getHttpReq.responseText)
        document.accessToken = respObj.data.access_token
    },{"Authorization":"Bearer "+window.localStorage.getItem("refrt")})
    
}

function refreshDashboard(){
    console.log(new Date()+" : Refreshing the dashboard...")
    $ajaxUtils.sendGETRequest("https://ajaysquare.com/stocksdash/dashboard",function(getHttpReq){
        if(getHttpReq.status===200){
            document.getElementById("mainMsg").style.color = "green"
        }else {
            document.getElementById("mainMsg").style.color = "red"
        }
        document.getElementById("mainMsg").textContent=getHttpReq.responseText +"\n"+"Last refreshed at : "+new Date();
    },{"Authorization":"Bearer "+document.accessToken})
}


function validate(){
    let unameElem=document.getElementById("uname")
    let pswdElem=document.getElementById("psw")
    let uname=unameElem.value
    let pswd=pswdElem.value
    let errMsg=""
    if(!uname.endsWith(".com")){
        errMsg += "Invalid Email Id!\n"
        unameElem.focus();
    }
    if(pswd.length <5){
        errMsg += "Password cannot be less than 5 chars!\n"
        pswdElem.focus();
    }
    if(errMsg.length > 0){
        alert(errMsg)
        return false
    }
    else {
        return {
            "email":uname,
            "password":pswd,
            "clientId": "SVEC"
        }
    }
}

function clearFields(){
    document.getElementById("uname").value='';
    document.getElementById("psw").value='';
    document.getElementById("mainMsg").textContent="Enter your credentials"
}
