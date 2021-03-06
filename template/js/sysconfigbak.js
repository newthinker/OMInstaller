/**
 * @author Administrator
 */
//需要额外处理的项
var exceeptionParam = {
    "db_ip" : "数据库服务器IP",
    "db_h2" : "内存数据库IP",
    "gis_ip" : "GIS服务器IP",
    "gis_port" : "GIS服务器端口",
    "gis_ags_ip" : "ArcGISServer服务器IP",
    "wmts_ip" : "WMTS系统的IP",
    "wmts_port" : "WMTS系统的PORT",
    "sysrest_ip" : "SysRest服务IP",
    "sysrest_port" : "SysRest服务PORT",
    "main_ip" : "运维服务器IP",
    "main_port" : "运维服务器端口",
    "geocoding_ip" : "GeoCoding系统IP",
    "geocoding_port" : "GeoCoding系统PORT",
    "geoportal_ip" : "GeoPortal系统IP",
    "geoportal_port" : "GeoPortal系统PORT",
    "file_ip" : "FileServices系统IP",
    "file_port" : "FileServices系统的PORT",
    "tile_ip" : "TileServices系统IP",
    "tile_port" : "TileServices系统PORT",
    "aggregator_ip" : "Aggregator系统IP",
    "aggregator_port" : "Aggregator系统PORT",
    "agent_ip" : "监控代理服务器IP",
    "agent_port" : "监控代理服务器PORT",
    "jms_ip" : "JMS服务器IP",
    "token_ip" : "Token服务器IP",
    "token_port" : "Token服务器PORT",
    "web_ip" : "门户IP",
    "web_port" : "门户PORT",
    "flex_ip" : "Flex发布服务器IP",
    "flex_port" : "Flex发布服务器PORT",
    "sl_ip" : "Silverlight发布服务器IP",
    "sl_port" : "Silverlight发布服务器PORT",
    "ria_ip" : "快速搭建服务器IP",
    "ria_port" : "快速搭建服务器PORT",
    //传空值就行
    "iis_root" : "Silverlight文件生成路径"
};

var tabContainer;
var checkInfo;
window.onload = function() {

    getServerConfigDInfo();
    //注册提交事件
    document.getElementById("sys_submit_bt").onclick = postAllSeverConfigInfos;

}
//获取服务器配置信息
function getServerConfigDInfo() {
    jx.load("post_bak.json", getServerConfigInfoHandler, "json", "get");
}

function getServerConfigInfoHandler(data) {
    initServerConfigInfo(data);
    initTabContainer();
}

//动态创建服务器配置信息节点
function initServerConfigInfo(data) {
    //动态生成config node节点信息
    var dc = dom.createDom;
    var fragment = document.createDocumentFragment();
    var fieldset = dc("fieldset", {
        "class" : "pd10 mgt10"
    }, fragment);
    var legend = dc("legend", {
        "innerHTML" : "服务器配置",
        'class' : "twb"
    }, fieldset);

    fieldset.appendChild(initServerModules(data['server_modules']));
    fieldset.appendChild(initServerParams(data['server_params']));
    //将动态生成的config info节点添加到页面
    var configInfo = dom.getDom('.config-info',document.getElementById('customTabContent'),"div")[0];
    configInfo.appendChild(fragment);
}

//初始化tabContainer
function initTabContainer() {
    var tempNode = document.getElementById("customTabContent");
    var firstNodeIntempNode = dom.getFirstElementChild(tempNode);
    if (!firstNodeIntempNode) {
        alert("无法查找模板节点，请刷新页面重试!");
        return;
    }
    var num = 1;
    tabContainer = new TabContainer(document.getElementById('tabdemo'));
    tabContainer.addTab({
        title : "Machine" + num,
        content : firstNodeIntempNode.cloneNode(true),
        closable : false
    });

    tabContainer.addTabBtn.onclick = function() {
        num++;
        tabContainer.addTab({
            title : "Machine" + num,
            content : firstNodeIntempNode.cloneNode(true),
            closable : true
        });
    }
}

//创建服务器模块
function initServerModules(modules) {
    var dc = dom.createDom;
    var configCheckServer = dc("div", {
        'class' : "config-checks",
        "name" : "config-check-servers"
    });
    var fieldset = dc("fieldset", {
        'class' : 'pd10 mgt10'
    }, configCheckServer);
    var legend = dc("legend", {
        "innerHTML" : "请选择安装服务器类型",
        'class' : 'twb'
    }, fieldset);
    var checkControls = dc("div", {
        "class" : "check-controls"
    }, fieldset);
    for (var i = 0; i < modules.length; i++) {
        var m = modules[i];
        checkControls.appendChild(createCheckControl(m));
    }
    return configCheckServer;
}

//创建每个服务器模块类型
function createCheckControl(module) {
    var dc = dom.createDom;
    var checkControl = dc('p', {
        'class' : "check-control"
    });
    var cbname = "config-check-" + module.srvname;
    var cb = dc("input", {
        "type" : "checkbox",
        "value" : module.srvname,
        "class" : "check-value",
        "name" : cbname,
        "onclick" : "checkControlClickHandler(this);"
    }, checkControl);
    var lb = dc("label", {
        "class" : "control-label",
        "for" : cbname,
        innerHTML : module.srvname
    }, checkControl);

    //agent默认选中且不能改变状态
    if (module.srvname === "agent") {
        cb.setAttribute("checked", "checked");
        cb.disabled = "true";
    }
    //cb.onclick = checkControlClickHandler;
    return checkControl;

}

//服务器配置中checkbox点击事件，判断是否选中
function checkControlClickHandler(obj) {
    var name = "server-" + obj.value;
    var node = dom.getDom(name,tabContainer.selectedTab,"fieldset")[0];
    //因为agent、token下面现在只包含ip、port属性节点，而且这两个属性节点不用显示，
    //所以判断node下面是否含有需要显示的p标签（即属性节点），如果有的话则切换显示状态，否则始终隐藏该node
    if (!lookupConfigServerHasBlockTag(node))
        return;
    obj.checked ? dom.removeClass("dn", node) : dom.addClass("dn", node);
}

//创建服务器配置信息
function initServerParams(params) {
    var configServers = dom.createDom("div", {
        "class" : "config-servers",
        name : "servers"
    });
    for (var i = 0; i < params.length; i++) {
        configServers.appendChild(createConfigServer(params[i]));
    }
    return configServers;
}

//创建每个服务器的配置信息
function createConfigServer(param) {
    var servername = param.srvname;
    var dc = dom.createDom;
    var configServer = dc("div", {
        'class' : "config-server"
    });
    var fieldset = dc("fieldset", {
        "class" : 'pd10 mgt10 dn',
        'name' : "server-" + servername
    }, configServer);
    var legend = dc("legend", {
        "class" : "twb",
        "innerHTML" : servername + "服务器信息配置"
    }, fieldset);
    var inputInfos = param.params;
    for (var i = 0; i < inputInfos.length; i++) {
        fieldset.appendChild(createInputControl(inputInfos[i]));
    }
    return configServer;
}

//创建配置信息中的单个配置
function createInputControl(obj) {
    var dc = dom.createDom, name = obj.paramname, desc = obj.paramdesc, encrypt = obj.encrypt, select = obj.select;
    var p = dc("p", {
        "class" : hasKeyInException(name) ? "dn" : ""
    });
    var label = dc("label", {
        "class" : "control-label",
        "for" : name,
        "innerHTML" : desc + "："
    }, p);
    //如果含有select属性，生成选择框
    if (select) {
        var slt = dc("select", {
            "class" : "select-value",
            "name" : name
        }, p);
        for (var i = 0; i < select.length; i++) {
            var s = select[i];
            dc("option", {
                "value" : s.option,
                "innerHTML" : s.selected
            }, slt);
        }
    }
    //默认生成input，
    else {
        var input = dc('input', {
            'class' : "input-value",
            "name" : name,
            "placeholder" : desc
        }, p);
        //input可能会含有是否加密的属性
        if (encrypt) {
            input.setAttribute("encrypt", encrypt);
        }
    }
    return p;
}

//检查key是否在需要例外处理的属性值里面
function hasKeyInException(key) {
    return exceeptionParam[key] == undefined ? false : true;
}

//检查节点里面是否有显示的标签
function lookupConfigServerHasBlockTag(serverNode) {
    var ps = serverNode.getElementsByTagName('p');
    for (var i = 0; i < ps.length; i++) {
        if (!dom.hasClass("dn", ps[i]))
            return true;
    }
    return false;
}

//////////////////
//获取所有服务器配置信息
//////////////////

function postAllSeverConfigInfos() {
    checkInfo = {
        flag : true,
        index : 1,
        info : ""
    };
    var infos = getAllServerConfigInfos();
    if (!checkInfo.flag) {
        window.open("error.html", "error", " location=no, directories=no, status=no, width=700,height=500").focus();
        return;
    }
    var data = JSON.stringify(infos);
    console.log(data);
    //jx.load("",function(){},"json","post");
}

//获取所有服务器配置信息
function getAllServerConfigInfos() {
    var i, infos = [], info, tabs = tabContainer.getTabs();
    for ( i = 0; i < tabs.length; i++) {
        info = getServerConfigInfo(tabs[i]);
        infos.push(info);
    }
    return infos;
}

//获取某一个tab页内的服务器配置信息
function getServerConfigInfo(tab) {
    var info = {};
    info.server_base = getServerBaseInfo(tab);
    info.server_params = getServerParamsInfo(tab);
    return info;
}

//获取服务器的基础配置信息
function getServerBaseInfo(tab) {
    var i, j, base, configNodes, baseInfo = {};
    base = dom.getDom("base",tab,"div")[0];
    configNodes = getConfigNodes(base);

    for ( i = 0; i < configNodes.length; i++) {
        var cn = configNodes[i];
        checkConfigNodeValue(tab, cn);
        baseInfo[cn.name.replace("base_", "")] = cn.value;
    }

    return baseInfo;

}

//获取服务器的详细配置信息
function getServerParamsInfo(tab) {
    var i, configChecks, dg = dom.getDom, checkBoxs, paramsInfo = [];
    configChecks = dg("config-check-servers",tab,"div")[0];
    checkBoxs = configChecks.getElementsByTagName("input");
    for ( i = 0; i < checkBoxs.length; i++) {
        var cb = checkBoxs[i];
        if (cb.checked === true) {
            var name = cb.value;
            var params = getParamsOfSelectServer(name, tab);
            paramsInfo.push({
                srvname : name,
                params : params
            });
        }
    }
    return paramsInfo;
}

//获取每个服务器的配置信息
function getParamsOfSelectServer(serverName, tab) {
    var i, configNodes, name, cn, value, serverNode, params = [];
    var dg = dom.getDom;
    serverNode = dg("server-"+serverName,tab,"fieldset")[0];
    configNodes = getConfigNodes(serverNode);
    for ( i = 0; i < configNodes.length; i++) {

        cn = configNodes[i];
        name = cn.name;
        //如果是需要额外处理的项的话，值和公用的值一样
        console.log("name: " + name);
        value = getValueForConfigNode(cn, tab);
        console.log("value: " + value);
        var param = {
            paramname : name,
            paramvalue : value
        };
        if (cn.hasAttribute('encrypt'))
            param.encrypt = cn.getAttribute('encrypt');
        params.push(param);
    }
    return params;
}

//获取配置节点的值
function getValueForConfigNode(node, tab) {
    var name = node.name, defaultIpName = "base_ip", dg = dom.getDom;
    //不在隐藏数据里传空间本身的值
    if (!hasKeyInException(name)) {
        checkConfigNodeValue(tab, node);
        return node.value || node.getAttribute("placeholder");
    }
    //db_h2比较特殊，传通用ip
    else if ("db_h2" == name) {
        return dg(defaultIpName,tab,"input")[0].value;
    }
    //Silverlight文件生成路径暂时传空值
    else if ("iis_root" == name) {
        return "";
    }
    //其他隐藏数据传对应的公用数据
    else {
        var n = name.substring(0, name.lastIndexOf("_"));
        return dg(name.replace(n,"base"),tab,node.tagName.toLowerCase())[0].value
    }

}

//获取所有配置节点
function getConfigNodes(node) {
    var all = [];

    var inputs = nodeListToArray(node.getElementsByTagName("input"));
    var selects = nodeListToArray(node.getElementsByTagName("select"));
    return all.concat(inputs, selects);
}

//根据服务器的选择来设置默认端口号
function baseContainerValueChangeHandler(obj) {
    var value = obj.value;
    var port = dom.getDom('base_port',tabContainer.selectedTab,"input")[0];
    port.value = value == "Tomcat" ? "8080" : "7001";
}

//验证配置信息
function checkConfigNodeValue(tab, node) {
    if (node.tagName.toLowerCase() != "input")
        return;
    var name = node.name, value = node.value;
    if (lang.isNull(value)) {
        invalidValue(tab, node, "不能为空!");
        return;
    }
    if (lang.hasChinese(value)) {
        invalidValue(tab, node, "不能包含中文字符!");
        return;
    }
    if (name == "base_omhome" && (!lang.isLinuxPath(value) && !lang.isWinPath(value))) {
        invalidValue(tab, node, "请填入正确的路径!");
        return;
    }
    if (name.indexOf("_ip") != -1 && !lang.isIp(value)) {
        invalidValue(tab, node, "请输入正确的ip地址!");
        return;
    }
    if (name.indexOf("_port") != -1 && !lang.isNum(value)) {
        invalidValue(tab, node, "端口号只能为数字!");
        return;
    } else {
        changeInputBorderColor(node, "#eee");
        return;
    }
}

//生成错误信息
function makeErrorInfo(tabtitle, nodename, info) {
    var tr = '<tr class="' + (checkInfo.index % 2 == 0 ? "odd" : "" ) + '"><td>' + tabtitle + '</td><td>' + nodename + '</td><td>' + info + '</td></tr>';
    checkInfo.info += tr;
    checkInfo.index++;
}

//对不合格的值进行处理
function invalidValue(tab, node, info) {
    checkInfo.flag = false;
    changeInputBorderColor(node, "red");
    makeErrorInfo(tab.getAttribute('tabtitle'), node.placeholder || node.name, info);
}

function changeInputBorderColor(node, color) {
    node.style.borderColor = color;
}

// ie8不支持 Array.prototype.slice.call
//自定义方法将nodelist转换为array
function nodeListToArray(nodelist) {
    var arr = [];
    for (var i = 0; i < nodelist.length; i++) {
        arr.push(nodelist[i]);
    }
    return arr;
}
