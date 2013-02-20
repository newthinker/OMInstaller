/**
 * @author Administrator
 */
(function(win) {

var doc=win.document;
    function trim(/*String*/str) {
        return str.replace(/^\s\s*/, '').replace(/\s\s*$/, '');
    }

    if (!Function.prototype.bind) {
        Function.prototype.bind = function(context) {
            var self = this, args = Array.prototype.slice.call(arguments);

            return function() {
                return self.apply(context, args.slice(1));
            }
        };
    }
    Array.prototype.remove = function(index) {
        if (index == this.length) {
            this.pop();
            return this;
        } else if (index == 0) {
            this.shift();
            return this;
        } else {
            return this.slice(0, index).concat(this.slice(index + 1, this.length));
        }
    }
    function stopPropagation(/*Event*/event) {
        event = event || window.event;
        if (event.stopPropagation) {
            event.stopPropagation();
        } else {
            event.cancelBubble = true;
        }

    }

    function srcElement(/*Event*/event) {
        event = event || window.event;
        return event.srcElement || event.target;
    }

    function isFunction(it) {
        return {}.toString.call(it) == "[object Function]";
    }

    function isString(it) {
        return {}.toString.call(it) == "[object String]";
    }

    function isIp(ip) {
        var pattern = /^(?:(?:25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)\.){3}(?:25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)$/;
        return pattern.test(ip);
    }

    function isNull(str) {
        return trim(str) === "";
    }

    function isLinuxPath(str) {
        var pattern = /^(\/\w+)*([\/])?$/;
        return pattern.test(str);
    }

    function isWinPath(str) {
        var pattern = /^[a-zA-Z]:(\\\w+)*([\\])?$/;
        return pattern.test(str);
    }

    function isNum(str) {
        var pattern = /^\d+$/;
        return pattern.test(str);
    }

    function hasChinese(str) {
        var pattern = /[\u4E00-\u9FA5\uf900-\ufa2d]/;
        return str.match(pattern) ? true : false;
    }

    //字符串首字母大写
    String.prototype.toFirstCharUpperCase = function() {
        var firstChar = this.charAt(0);
        return this.replace(firstChar, firstChar.toUpperCase());
    }

    window.lang = {
        trim : trim,
        stopPropagation : stopPropagation,
        srcElement : srcElement,
        isString : isString,
        isFunction : isFunction,
        isIp : isIp,
        isNull : isNull,
        isLinuxPath : isLinuxPath,
        isWinPath : isWinPath,
        isNum : isNum,
        hasChinese : hasChinese,
        docWidth:doc.documentElement.clientWidth || win.innerWidth,
        docHeight:doc.documentElement.clientHeight || win.innerHeight
    };
})(window)
