$(document).ready(init);

function init() {
    //$('body').on('click', 'div.appTable', fnAjaxLoadPage);
    $(".appTable tbody").on("click", fnAjaxLoadPage);
    $(".appMenuItemsLeft nav").on("click", fnAjaxLoadPage);
    $(".appMenuItemsRight li").on("click", fnAjaxLoadPage);
    fnLog("init!");
}

function fnAjaxLoadPage(e) {
    //fnLog(e.target.nodeName + ", " + $(e.target).text());
    cRID = 11 // Client Request ID
    var Rows = []
    indexData = 0
    switch (e.target.nodeName) {
        case 'I':
        case 'A':
            switch ($(e.target).text()) {
                case "person":
                    cRID = 11
                    break;
                case "edit":
                    cRID = 12
                    break;
                case "refresh":
                    cRID = 13
                    break;
                case "delete":
                    cRID = 14
                    break;
                case "account_circle":
                    cRID = 1
                    break;
                case "arrow_back":
                    cRID = 2
                    break;
                case "account_circleMy Account":
                    cRID = 1
                    break;
                case "arrow_backQuit":
                    cRID = 2
                    break;
            }
            break;
        case 'LI':
            switch ($(e.target).text()) {
                case "Create User":
                    cRID = 3
                    break;
                case "Upload Image":
                    cRID = 4
                    break;
                case "Create Album":
                    cRID = 5
                    break;
                case "Download Album":
                    cRID = 6
                    break;
            }
            break;
    }
    if (cRID > 9 && cRID < 20) {
        Rows[indexData++] = $(e.target).closest("tr").attr("id")
    }
    //fnLog("Node: " + e.target.nodeName + ", Text: " + e.target.text + ", ID: " + cRID);
    if (cRID > 0) {
        //fnLog("ID: " + cRID);
        jQuery.ajax({
            type: 'post',
            url: "/ajax",
            data: {ID: cRID, Data: Rows},
            dataType: 'json',
            success: function (result) {
                fnLog("Success:" + result.Data);
            },
            error: function (result) {
                fnLog("Failure:" + result);
            }
        });
    }
}

function fnFindStr(sourceStr, str) {
    if (sourceStr.indexOf(str) != -1) {
        return str;
    }
    return "";
}

function fnNum2ZPfxdStr(num, requiredLength) {
    var numStr = num.toString();
    var lenNumStr = numStr.length;
    var diffLen = requiredLength - lenNumStr;
    for (var i = 0; i < diffLen; i++) {
        numStr += '0';
    }
    return numStr;
}

function fnLog(logStr) {
    var dt = new Date();

    var dateTimeStamp = dt.getFullYear() + "/" + fnNum2ZPfxdStr(dt.getMonth() + 1, 2) + "/" + fnNum2ZPfxdStr(dt.getDate(), 2) + " " + fnNum2ZPfxdStr(dt.getHours(), 2) + ":" + fnNum2ZPfxdStr(dt.getMinutes(), 2) + ":" + fnNum2ZPfxdStr(dt.getSeconds(), 2) + "." + fnNum2ZPfxdStr(dt.getMilliseconds(), 3);
    //var date_str = Date($.now());
    console.log("@" + dateTimeStamp + "> " + logStr);
}