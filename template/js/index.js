/**
 * @author Administrator
 */
function nextStep() {
    var flag = {};
    var omtoolType = getCheckedRadioValue();
    var subconfigChecked = getSubconfigChecked();
    if (omtoolType == 1 && subconfigChecked) {
        flag = 4;
    } else {
        flag = omtoolType;
    }
    document.getElementById('flag').value = flag;
    document.getElementById('onemaptool-form').submit();
}

function getCheckedRadioValue() {
    var radios = document.getElementsByName('onemaptool');
    for (var i = 0, len = radios.length; i < len; i++) {
        var rd = radios[i];
        if (rd.checked)
            return rd.value;
    }
}

function getSubconfigChecked() {
    var subconfig = document.getElementsByName('subconfig')[0];
    return subconfig.checked;
}
