function validate(){
    let unameElem=document.getElementById("uname")
    let pswdElem=document.getElementById("psw")
    let uname=unameElem.value
    uname='aj.nextstep@gmail.com'
    let pswd=pswdElem.value
    pswd='Paswd1234'
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

function sleep(ms) {
    return new Promise(resolve => setTimeout(resolve, ms));
}

//This is the AJAX call
function httpPost(url,jsonPayload){
    let xmlHttp = new XMLHttpRequest();
    xmlHttp.open("POST",url,true)
    xmlHttp.onreadystatechange = function (){
        loginResponseHandler(xmlHttp)
    };
    xmlHttp.send(jsonPayload)
}
//This is the AJAX call back reponse handler
async function loginResponseHandler(loginHttpRespObj){
    console.log("In loginResponseHandler()")
    if(loginHttpRespObj.readyState == 4 ){
        console.log("Waiting for 2 secs...")
        await sleep(2000)
        if(loginHttpRespObj.status === 200)
        {
            console.log("Login is successfull")
            let respObj = JSON.parse(loginHttpRespObj.responseText)
            let accessToken = respObj.data.access_token
            document.accessToken = accessToken
            let refreshToken = respObj.data.refreshToken
            window.localStorage.setItem("refrt",refreshToken)
            document.getElementById("mainMsg").style.color = "green"
            //show the dashboard data as plain text for now
            console.log("Getting data from dashboard API...(not Ajax)")
            let getHttpReq = new XMLHttpRequest();
            getHttpReq.open( "GET", "https://ajaysquare.com/stocksdash/dashboard", false ); // false for synchronous request
            getHttpReq.setRequestHeader("Authorization","Bearer "+accessToken)
            getHttpReq.send( null );
            document.getElementById("mainMsg").textContent=getHttpReq.responseText;
        }else{
            console.warn("Unable to authorize the user!")
            let respObj = JSON.parse(loginHttpRespObj.responseText)
            let errMsg = respObj.message
            document.getElementById("mainMsg").style.color = "red"
            document.getElementById("mainMsg").textContent="Login Failed :"+errMsg
    }
    }
        
}

function loginClicked(){
    let validationResult = validate()
    if(validationResult !==false){
        document.getElementById("mainMsg").textContent="Please wait...for 2 secs"
        console.log("All validations passed")
        httpPost("https://ajaysquare.com/auth-api/login",JSON.stringify(validationResult))
        console.log("returning from loginClicked()")
    }else {
        console.warn("Input validations failed")
    }
}