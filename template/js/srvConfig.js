<!-- ���÷����� -->
function configServer() {
	var cur_form=dom.getDom('sysconfig',tabContainer.selectedTab,"form")[0];
	var chkboxes = dom.getDom('.check-value',cur_form,"input")
	var count = 0;		// ȷ����checkbox����
	var agent = 0;		// ��ش����������������
	// ����checkboxes����ȡ����ȷ��checkboxes
	for(var i=0,elm;elm=chkboxes[i];i++) {
		var type = elm.value;		// ��ȡȷ������������
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
	
	
	if(count<=1) {		// û������
		alert("�����÷�����!");
		return;
	} else if(agent!=1) {	// û�����ü�ش���ģ��
		alert("�����ü�ش��������!");
		return;
	} 
	
}