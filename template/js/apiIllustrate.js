/**
 * @author Administrator
 */
//
// window.dom(dom)介绍
//      查看dom节点是否含有css 类         return boolean
//      dom.hasClass(/*String ('a')*/classStr, /*node*/node)
//
//      dom节点移出css类
//      dom.removeClass(/*String ('a' or 'a b')*/classStrs, /*node*/node)
//
//      dom节点添加css类
//      dom.addClass(/*String ('a' or 'a b')*/classStrs, /*node*/node)
//
//      根据css类('.a')或者节点name('b')、父节点、dom标签获取dom节点      return Array
//      dom.getDom(/*String (name:'name' or class:'.class')*/param, /*node*/root, /*String ('div')*/tag)
//
//      根据dom标签创建dom节点          return domNode
//      dom.createDom(/*String ('div')*/tag,/*Object*/ opts)
//
// window.lang(lang)介绍
//      去除字符串前后空白     return String
//      lang.trim(/*String*/str)
//
//      阻止事件冒泡
//      lang.stopPropagation(/*Event*/event)
//
//      获取发生鼠标事件的dom节点     return domNode
//      lang.srcElement(/*Event*/event)
//
// window.TabContainer(TabContainer) 介绍
//      tab控件对象
//      eg:
//      //根据dom来实例化tab控件，node为空的话默认为documen.body对象
//      var tabContainer = new TabContainer(/*domNode*/node);
//      tabContainer.addTab({
//          title: "title",                 //tab页标题
//          content: "<p>tab内容</p>"     //tab页内容
//          closable: false                 //tab页是否允许关闭
//      });
//      注册添加tab事件
//      tabContainer.addTabBtn.onclick = function(){
//          tabContainer.addTab({
//              title: "title",
//              content: "<p>tab内容</p>"
//              closable: false
//          });
//      };
//      获取当前选择tab
//      var selected=tabContainer.selectedTab;
//
//      获取所有tab
//      var allTabs=tabContainer.getTabs();
//
//
//