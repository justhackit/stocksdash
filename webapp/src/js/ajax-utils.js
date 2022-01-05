(function (global){
    //set up namespace for this utility
    var ajaxUtils={};

    //perform an AJAX GET call
    ajaxUtils.sendGETRequest = function (url,responseHandlerFunc,requestHeaders={}){
        let xmlHttp = new XMLHttpRequest();
        xmlHttp.onreadystatechange = function(){
            handleResponse(xmlHttp,responseHandlerFunc)
        }
        xmlHttp.open("GET",url,true)
        for(let key in requestHeaders){
            xmlHttp.setRequestHeader(key,requestHeaders[key])
        }
        xmlHttp.send(null)
    };

    //perform an AJAX POST call
    ajaxUtils.sendPOSTRequest = function (url,payload,responseHandlerFunc,requestHeaders={}){
        let xmlHttp = new XMLHttpRequest();
        xmlHttp.onreadystatechange = function(){
            handleResponse(xmlHttp,responseHandlerFunc)
        }
        xmlHttp.open("POST",url,true)
        for(let key in requestHeaders){
            xmlHttp.setRequestHeader(key,requestHeaders[key])
        }
        xmlHttp.send(payload)
    };

    function handleResponse(httpReqResp,responseHandlerFunc){
        if(httpReqResp.readyState===4){
            responseHandlerFunc(httpReqResp)
        }
    }

    // Expose utility to the global object
    global.$ajaxUtils = ajaxUtils;

})(window);

