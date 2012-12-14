<!-- 添加服务模块 -->
function configServer() {
	var cur_form=dom.getDom('sysconfig',tabContainer.selectedTab,"form")[0];
	var chkboxes = dom.getDom('.check-value',cur_form,"input")
	var count = 0;		// 添加的服务模块个数，也即是选择的checkbox个数
	var agent = 0;		// 代理模块是否安装标识(必安装模块)
	// 遍历所有checkboxes，将选择的进行显示
	for(var i=0,elm;elm=chkboxes[i];i++) {
		var type = elm.value;		// checkbox所对应的服务器名称
		switch (type)
		{
		case "agent":
			if(elm.checked) {
				agent = 1;
				count++;
			}
			break;
		case "db":
			if(elm.checked) {
				dom.getDom("server-db",cur_form,"fieldset")[0].style.display='block';
				count++;
			} else {
				dom.getDom("server-db",cur_form,"fieldset")[0].style.display='none';
			}
			break;
		case "gis":
			if(elm.checked) {
				dom.getDom("server-gis",cur_form,"fieldset")[0].style.display='block';
				count++;
			} else {
				dom.getDom("server-gis",cur_form,"fieldset")[0].style.display='none';
			}
			break;
		case "main":
			if(elm.checked) {
				dom.getDom("server-main",cur_form,"fieldset")[0].style.display='block';
				count++;
			} else {
				dom.getDom("server-main",cur_form,"fieldset")[0].style.display='none';
			}
			break;
		case "msg":
			if(elm.checked) {
				dom.getDom("server-msg",cur_form,"fieldset")[0].style.display='block';
				count++;
			} else {
				dom.getDom("server-msg",cur_form,"fieldset")[0].style.display='none';
			}
			break;
		case "token":
			if(elm.checked) {
				//dom.getDom("server-token",cur_form,"fieldset")[0].style.display='block';
				count++;
			} else {
				dom.getDom("server-token",cur_form,"fieldset")[0].style.display='none';
			}
			break;
		case "web":
			if(elm.checked) {
				dom.getDom("server-web",cur_form,"fieldset")[0].style.display='block';
				count++;
			} else {
				dom.getDom("server-web",cur_form,"fieldset")[0].style.display='none';
			}
			break;
		}		
	}
	
	
	if(count<=1) {		// 没有选择安装任何服务器
		alert("请选择安装的服务器!");
		return;
	} else if(agent!=1) {	// 没有安装代理服务器
		alert("请选择必须安装的代理服务器模块!");
		return;
	} 
	
	// 将提交按钮设置为可提交状态
	document.getElementById("sys_submit_bt").disabled="";
}

<!-- 系统配置页面参数整理提交 -->
function sysSubmitHandler() {
	var sysSrvConfig = {};			// 所有服务器参数
	var num = 0;					// 服务器计数
	
	for (var cur_tab in tabContainer.getTabs()) {
		num++;
		
		var cur_form=dom.getDom('sysconfig', cur_tab, "form")[0];
		
		var server = {};
	
		// 服务器基本信息
		var base = {
			os: 
		};
		
		// 服务器安装模块信息
		
	}
	
	
	
	
	
	// 输出成json
	
}