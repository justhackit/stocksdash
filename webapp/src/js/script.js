(function(global){
    let loginPage = "snippets/login-snippet.html";
    var myP = {};
    // Convenience function for inserting innerHTML for 'select'
    myP.insertHtml = function(selector,html){
        let elem = document.querySelector(selector);
        elem.innerHTML = html
    };

    myP.showLoading = function(selector){
        let html = "<div class='text-center'>";
        html += "<img src='images/Vanilla-1s-280px.gif'></div>"
        myP.insertHtml(selector,html)
    }

    //on page load
    document.addEventListener("DOMContentLoaded",function(event){
        myP.showLoading("#main-content");
        $ajaxUtils.sendGETRequest(loginPage,function(getHttpReq){
            myP.insertHtml("#main-content",getHttpReq.responseText);
        })
    })

    global.$myP = myP

})(window);

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
            loadDashboardPage();
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

function logoutClicked(){
    console.log("User clicked logout")
    window.localStorage.removeItem("refrt")
    $myP.showLoading("#main-content");
    $ajaxUtils.sendGETRequest("snippets/login-snippet.html",function(getHttpReq){
        $myP.insertHtml("#main-content",getHttpReq.responseText);
    })
    document.accessToken=null
    stopRefreshes();
}

function stopRefreshes(){
    clearInterval(window.refreshIntervalId)
    clearInterval(window.tokenIntervalId)
}


function getNewAccessToken(){
    console.log(new Date()+" : Getting new access token using refresh token..")
    $ajaxUtils.sendGETRequest("https://ajaysquare.com/auth-api/refresh-token",function(getHttpReq){
        if(getHttpReq.status===200){
            let respObj = JSON.parse(getHttpReq.responseText)
            document.accessToken = respObj.data.access_token
        }else {
            console.log("Unable to get refresh token. Loading login page")
            stopRefreshes();
            $myP.showLoading("#main-content");
            $ajaxUtils.sendGETRequest("snippets/login-snippet.html",function(getHttpReq){
                $myP.insertHtml("#main-content",getHttpReq.responseText);
            })
        }
        
    },{"Authorization":"Bearer "+window.localStorage.getItem("refrt")})
    
}

function loadDashboardPage(){
    $ajaxUtils.sendGETRequest("snippets/dashboard-snippet.html",function(getHttpReq){
        $myP.insertHtml("#main-content",getHttpReq.responseText)
    })
}

function refreshDashboard(){
    console.log(new Date()+" : Refreshing the dashboard...")
    $ajaxUtils.sendGETRequest("https://ajaysquare.com/stocksdash/dashboard",function(getHttpReq){
        if(getHttpReq.status===200){
            document.getElementById("mainMsg").innerHTML=buildAndGetTableHTMLFromJSON(JSON.parse(getHttpReq.responseText))
        }else if(getHttpReq.status===400) {
            document.accessToken=getNewAccessToken()
            refreshDashboard();
            //document.getElementById("mainMsg").style.color = "red"
        }
        //document.getElementById("mainMsg").textContent=getHttpReq.responseText
        document.getElementById("ftrMsg").textContent="Last refreshed at : "+new Date();
    },{"Authorization":"Bearer "+document.accessToken})
}

function buildAndGetTableHTMLFromJSON(responseJson){
    var tableHtml = '<table> \
    <tr style="background-color:#808080"> \
      <th>Stock</th> \
      <th>Bought At$</th> \
      <th>Present Price$</th> \
      <th>Total shares</th> \
      <th>Profit or Loss</th> \
      <th>Percentage</th> \
    </tr>'
    
    for(var i = 0; i < responseJson.length; i++) {
        var aStock = responseJson[i];
        if(aStock["profitLoss"] <0){
            tableHtml += '<tr style="background-color:#FF0000">'
        }else{
            tableHtml += '<tr style="background-color:#008000">'
        }
        tableHtml += '<td>' + aStock["ticker"] + '</td>'
        tableHtml += '<td>' + aStock["avgCostPrice"] + '</td>'
        tableHtml += '<td>' + Math.round(aStock["currentPrice"] * 100) / 100 + '</td>'
        tableHtml += '<td>' + aStock["totalShares"] + '</td>'
        tableHtml += '<td>' + Math.round(aStock["profitLoss"] * 100) / 100 + '</td>'
        tableHtml += '<td>' + Math.round(aStock["profitLossPerc"] * 100) / 100 + '%</td>'
        tableHtml+='</tr>'
    }

    return tableHtml + '</table>'
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
}
