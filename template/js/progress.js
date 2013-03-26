/**
 * @author Administrator
 */

(function(window) {

    var server = 'ws://'+window.location.host+'/json';
    var pregsSocket=null;

    var dlgWrapper, dlg, close, pregressinfos, progress, progressvalue;

    function connect() {
        //先注册配置进度提示的关闭事件
        registDialogCloseEvent();

        if (!window.WebSocket) {
            alert("该浏览器暂不支持websocket，进度条可能无法正常显示。\n推荐使用最新版的chrome或firefox浏览器！");
            return;
        }

        pregsSocket = new WebSocket(server);
        pregsSocket.onopen = socketOpenHandler;
        pregsSocket.onclose = socketCloseHandler;
        pregsSocket.onmessage = socketMessageHandler;
        pregsSocket.onerror = socketErrorHandler;

    }

    function socketOpenHandler(evt) {
        console.log('开启 websockect 连接...');
    }

    function socketCloseHandler(evt) {
        console.log('关闭 websockect 连接...');
        //关闭sockect连接时显示进度框关闭按钮
        //close.style.display = "inline-block";
    }

    function socketMessageHandler(evt) {
        var data =JSON.parse(evt.data);

        //返回错误消息提示
        if (data.Ret == -1) {
            socketErrorHandler(evt);
            return;
        }
		
		 showProgress(data, 'info');
		
        //处理完成
        if (data.Ret == 100) {
            pregsSocket.close();
			close.style.display = "inline-block";
         	alert("服务器配置成功!");
        }

    }

    function socketErrorHandler(evt) {
        var data =JSON.parse(evt.data);
        pregsSocket.close();
		 close.style.display = "inline-block";
        
        showProgress(data, 'error');
		alert("服务器配置出错，清查看日志说明!");
    }

    function showPregressDialog() {

        dom.removeClass('dn', dlgWrapper);
        dom.removeClass('dn', dlg);
        dlg.style.top = (lang.docHeight - 250) / 2 + "px";
        dlg.style.left = (lang.docWidth - dlg.clientWidth) / 2 + "px";
		
	//	testOK();
    }

    function registDialogCloseEvent() {
        dlgWrapper = dom.getDom('.infodialogwrapper',null,"div")[0];

        dlg = dom.getDom('.infodialog',null,"div")[0];
        close = dom.getDom('.close',dlg,"span")[0];
        pregressinfos = dom.getDom('.pregressinfos',dlg,"div")[0];
        progress = dom.getDom('.progress',dlg,"div")[0];
        progressvalue = dom.getDom('.progressvalue',dlg,"span")[0];

        close.onclick = function() {
            close.style.display = "none";
            dom.empty(pregressinfos);
            progress.style.width = "0";
            progressvalue.innerHTML = '0%';
            dom.addClass('dn', dlg);
            dom.addClass('dn', dlgWrapper);

        }
    }

    function mkInfo(info, type) {
        var color = type === "error" ? "red" : '';
		var p=document.createElement('p');
		p.style.color=color;
		p.style.fontWeight="bold";
		p.innerHTML=' >> '+info;
		return p;
        //return '<p style="color:' + color + ';"> >> ' + info + "</p>";
    }

    function showProgress(data, type) {
        //显示进度信息
       // pregressinfos.innerHTML += mkInfo(data.Reason, type);
		var firstInfo=pregressinfos.getElementsByTagName('p')[0];
		if(firstInfo){
			firstInfo.style.fontWeight="normal";
			pregressinfos.insertBefore( mkInfo(data.Reason, type),firstInfo);
		}else{
			pregressinfos.appendChild(mkInfo(data.Reason, type));
		}
        if (type === 'error')
            return;
        //显示进度
        progress.style.width = data.Ret + "%";
        progressvalue.innerHTML = data.Ret + "%";
    }

    function testOK() {
        var t = 1;
        setInterval(function() {
            socketMessageHandler({
                "data" : {
                    "Ret" : t * 10,
                    "Reason" : "服务器" + t + "安装成功..."
                }
            });
            t++;
        }, 1000);
    }

    function testError() {
        socketErrorHandler({
            "data" : {
                "Ret" : "-1",
                "Reason" : "服务器1安装失败..."
            }
        })
    }


    window.progress = {
        connect : connect,
        showPregressDialog : showPregressDialog,
        testOK : testOK,
        testError : testError
    };
})(window)
