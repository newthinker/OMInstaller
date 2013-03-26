/**
 * @author Administrator
 */
(function(window) {
    var document = window.document;

    function hasClass(/*String ('a')*/classStr, /*node*/node) {
        return ((" " + node.className + " ").indexOf(" " + classStr + " ")) >= 0;

    }

    function removeClass(/*String ('a' or 'a b')*/classStrs, /*node*/node) {
        var classNames = node.className ? " " + node.className + " " : " ";
        classStrs = classStrs.split(/\s+/);
        for (var i = 0; i < classStrs.length; i++) {
            var cls = classStrs[i];
            classNames = classNames.replace(" " + cls + " ", ' ');
        }
        classNames = lang.trim(classNames);
        if (node.className != classNames)
            node.className = classNames;
    }

    function addClass(/*String ('a' or 'a b')*/classStrs, /*node*/node) {
        var classNames = node.className ? " " + node.className + " " : " ";
        classStrs = classStrs.split(/\s+/);
        for (var i = 0; i < classStrs.length; i++) {
            var cls = classStrs[i];
            if (cls && classNames.indexOf(" " + cls + " ") < 0) {
                classNames += cls + " ";
            }
        }
        classNames = lang.trim(classNames);
        if (node.className != classNames)
            node.className = classNames;
    }

    function getDom(/*String (name:'name' or class:'.class')*/param, /*node*/root, /*String ('div')*/tag) {

        tag = tag || "*";
        root = root || document.body;
        return _getDom(param, root, tag);
    }

    function _getDom(/*String ((name:'name' or class:'.class')*/param, /*node*/root, /*String ('div')*/tag) {
        var nodes = [];
        var tagNodes = root.getElementsByTagName(tag);
        for (var i = 0; i < tagNodes.length; i++) {
            var n = tagNodes[i];
            if (_isParamMatchClassNameOrName(param, n))
                nodes.push(n);
        }
        return nodes;
    }

    function _isParamMatchClassNameOrName(param, node) {
        if (param.indexOf('.') >= 0) {
            return hasClass(param.substring(1, param.length), node);
        } else {
            return node.getAttribute("name") == param;
        }
    }

    function createDom(/*String ('div')*/tag, /*Object*/opts, /*domNode*/parentNode) {
        if (!tag)
            return null;
        var node = document.createElement(tag);
        if (opts && typeof opts == "object") {
            for (var att in opts) {
                if (opts.hasOwnProperty(att)) {
                    if (att == "innerHTML") {
                        node.innerHTML = opts.innerHTML;
                    } else {
                        node.setAttribute(att, opts[att]);
                    }
                }
            }
        }
        if (parentNode && parentNode.appendChild)
            parentNode.appendChild(node);

        return node;
    }

    function getFirstElementChild(/*domNode*/node) {
        var firstElementChild = node.firstElementChild;
        if (firstElementChild) {
            return firstElementChild;
        } else {
            var childNodes = node.childNodes;
            for (var i = 0; i < childNodes.length; i++) {
                var c = childNodes[i];
                if (c.nodeType === 1)
                    return c;
            }
            return null;
        }
    }

    function empty(node) {
        while (node.hasChildNodes())//当div下还存在子节点时 循环继续
        {
            node.removeChild(node.firstChild);
        }
    }


    window.dom = {
        hasClass : hasClass,
        removeClass : removeClass,
        addClass : addClass,
        getDom : getDom,
        createDom : createDom,
        getFirstElementChild : getFirstElementChild,
        empty:empty
    }
})(window)
