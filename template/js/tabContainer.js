/**
 * @author Administrator
 */
(function() {
    function TabContainer(/*domNode*/node) {
        this._defaultTabName = "tab_";
        this._defaultContentName = "content_";
        this._tabs = [];
        this.selectedTab = null;
        node = node || document.body;
        this._initTabContainer(node);
    }


    TabContainer.prototype.addTab = function(tab) {
        tab = tab || {
            title : 'tab' + this._tabs.length,
            content : 'content' + this._tabs.length,
            closable : true
        };
        this._addTab(tab);
    }

    TabContainer.prototype.getTabs = function() {
        var tabs = [];
        var t;
        for (var i = 0; i < this._tabs.length; i++) {
            t = this._tabs[i];
            if (t)
                tabs.push(t)
        }
        return tabs;
    }

    TabContainer.prototype._initTabContainer = function(node) {

        var temp = '<div class="tabContainer pr"><div class="tabMenus"><ul class="tabMenus-ul"><li class="addTab-li pr"><span class="tab dib  tac fl ml2 addTabSize  addTab pa">+</span></li></ul></div><div id="tabContents"class="tabContents pr"></div></div>';
        node.innerHTML += temp;

        //避免单页面多个tabContainer，使用父节点id做前缀
        if (node.id) {
            this._defaultTabName = node.id + "_tab_";
            this._defaultContentName = node.id + "_content_";
        }
        this._initMenuAndContent(node);
    }
    TabContainer.prototype._initMenuAndContent = function(node) {
        this._tabContainer = dom.getDom(".tabContainer", node, "div")[0];
        this._tabMenus = dom.getDom(".tabMenus-ul", this._tabContainer, "ul")[0];
        this._tabContents = dom.getDom(".tabContents", this._tabContainer, "div")[0];
        this.addTabBtn = dom.getDom(".addTab", this._tabMenus, "span")[0];
        //默认按照tab总数来添加tab
        //如果需添加自定义tab标签，重新监听this.addTabBtn.onclick事件
        this.addTabBtn.onclick = this.addTab.bind(this, null);
    }
    TabContainer.prototype._addTab = function(tab) {
        //改变当前选择的tab标签为没被选择状态
        this._changeSelectTabStatus();

        this._addTabMenu(tab.title, tab.closable);
        this._addTabContent(tab.content, tab.title);
    }
    TabContainer.prototype._addTabMenu = function(title, closable) {
        var menu_li = dom.createDom('li', {
            'class' : 'tabMenus-li'
        });
        var menu = dom.createDom('span', {
            'class' : 'tab dib  tac fl normalTabSize pr tabCur',
            id : this._defaultTabName + this._tabs.length,
            'for' : this._defaultContentName + this._tabs.length,
            innerHTML : title,
            tabtitle : title
        });
        var close = dom.createDom('span', {
            'class' : closable ? 'pa tab-closeBtn' : 'pa tab-closeBtn  dn',
            innerHTML : "x"
        });
        close.onclick = this._closeTabHandler.bind(this);
        menu.onclick = this._tabClickHandler.bind(this);
        menu.appendChild(close);
        menu_li.appendChild(menu);
        this._tabMenus.insertBefore(menu_li, this.addTabBtn.parentNode);
    }
    TabContainer.prototype._addTabContent = function(content, customTitle) {
        var div = dom.createDom('div', {
            'class' : 'content contentCur',
            id : this._defaultContentName + this._tabs.length,
            tabtitle : customTitle
        })
        if (lang.isString(content)) {
            div.innerHTML = content;
        } else {
            div.appendChild(content);
        }
        this._tabContents.appendChild(div);
        this._tabs.push(div);
        this.selectedTab = div;
    }
    TabContainer.prototype._tabClickHandler = function(event) {
        event = event || window.event;
        var target = lang.srcElement(event);
        this._changeSelectTabStatus();
        dom.addClass('tabCur', target);
        this.selectedTab = document.getElementById(target.getAttribute('for'));
        dom.addClass("contentCur", this.selectedTab);
    }
    TabContainer.prototype._changeSelectTabStatus = function() {
        if (!this.selectedTab)
            return;
        var smenu = dom.getDom(".tabCur", this._tabMenus, "span")[0];
        dom.removeClass('tabCur', smenu);
        dom.removeClass('contentCur', this.selectedTab);
    }
    TabContainer.prototype._closeTabHandler = function(event) {
        lang.stopPropagation(event);
        event = event || window.event;
        var srcElement = lang.srcElement(event);
        var selectmenu = srcElement.parentNode;

        var result = window.confirm("是否删除服务器: " + selectmenu.getAttribute("tabtitle") + " ？");
        if (!result) {
            return;
        }

        this._removeTab(selectmenu);

    }

    TabContainer.prototype._removeTab = function(selectmenu) {
        var forId = selectmenu.getAttribute('for');
        var content = document.getElementById(forId);

        this._tabMenus.removeChild(selectmenu.parentNode);
        this._tabContents.removeChild(content);
        var index = forId.substring(forId.lastIndexOf("_") + 1, forId.length);
        this._tabs[index] = null;
        //删除当前tab，将前一个tab置为当前tab
        if (forId === this.selectedTab.id) {
            do {
                index--;
            } while (!this._tabs[index])
            this.selectedTab = this._tabs[index];
            dom.addClass("contentCur", this.selectedTab);
            dom.addClass("tabCur", document.getElementById(this._defaultTabName + (index)));
        }
    }
    window.TabContainer = TabContainer;
})()
