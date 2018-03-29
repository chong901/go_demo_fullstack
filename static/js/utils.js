function getThisTr(dom) {
    return $(dom).closest('tr');
}

function alertErr(err) {
    BootstrapDialog.alert(err.responseJSON.message, () => {
        if (err.status === 410) {
            window.location = '/';
        }
    });
}

function save(to, redirect) {
    let data = $('#form').serialize();
    $.ajax({
        url: to,
        method: 'POST',
        data
    }).done(function () {
        window.location.href = redirect;
    }).fail(alertErr)
}

function remove(dom, url, key, value) {
    BootstrapDialog.show({
        title: $('#strPrompt').text(),
        message: $('#strDeleteMsg').text(),
        buttons: [
            {
                label: $('#strBtnConfirm').text(),
                action: function (dialogItself) {
                    let data = {};
                    data[key] = value;
                    $.ajax({
                        url: `${url}?${key}=${value}`,
                        dataType: 'json',
                        method: 'DELETE'
                    }).done(function () {
                        dialogItself.close();
                        BootstrapDialog.alert($('#strDeleted').text());
                        if (dom) {
                            getThisTr(dom).remove();
                            resetTableWhenEmpty('dataTable');
                        }
                    }).fail(alertErr).always(() => dialogItself.close());
                }
            },
            {
                label: $('#strBtnCancel').text(),
                action: function (dialogItself) {
                    dialogItself.close();
                }
            }]
    });
}

function removeFromEditPage(dom, url, id) {
    const $tr = $(dom).closest('div');
    if (!id) {
        $tr.remove();
        return
    }
    BootstrapDialog.show({
        title: $('#strPrompt').text(),
        message: $('#strDeleteMsgNotReversed').text(),
        buttons: [
            {
                label: $('#strBtnConfirm').text(),
                action: function (dialogItself) {
                    $.ajax({
                        url: `${url}?id=${id}`,
                        method: 'delete',
                        dataType: 'json'
                    }).done(() => {
                        BootstrapDialog.alert($('#strDeleted').text());
                        $tr.remove();
                    }).fail(alertErr).always(() => dialogItself.close());
                }
            },
            {
                label: $('#strBtnCancel').text(),
                action: function (dialogItself) {
                    dialogItself.close();
                }
            }]
    });
}

function removeTr(dom) {
    const resetName = (jqTable) => {
        jqTable.find('.dataRow').each(eachDataRow);
    };

    const eachDataRow = (index, val) => {
        $(val).find('.tdInput').each(function (_, val) {
            let tempName = $(val).attr('name');
            $(val).attr('name', tempName.replace(/\d+/, index));
        });
    };

    let table = $(dom).closest('.table');
    $(dom).closest('.tr').remove();
    resetName(table);
}

function redirect(uri) {
    window.location.assign(uri);
}

function resetTableWhenEmpty(id) {
    const table = $('#' + id);
    if (table.find('tbody').children().length === 0) {
        const thLength = table.find('tr').length;
        let html = '<thead>';
        html += `<tr id="emptyTr"><td class="text-center" colspan="${thLength}"><h2>${$('#strNoData').text()}</h2></td></tr>`;
        html += '</thead>';
        table.append(html);
    }
}

function setHtml(jqObj, cls, val) {
    if (!val) {
        return;
    }
    jqObj.find('.' + cls).html(val);
}

function setHtmlAcceptEmpty(jqObj, cls, val) {
    val = val || '';
    jqObj.find('.' + cls).html(val);
}

function getObjFromSerialize(serializeData, intKeys, dateKeys, boolKeys) {
    let obj = {};
    for (let k in serializeData) {
        const data = serializeData[k];
        if (intKeys && intKeys.indexOf(data.name) !== -1) {
            obj[data.name] = parseInt(data.value);
        } else if (dateKeys && dateKeys.indexOf(data.name) !== -1) {
            obj[data.name] = new Date(data.value);
        } else if (boolKeys && boolKeys.indexOf(data.name) !== -1) {
            obj[data.name] = Boolean(parseInt(data.value));
        } else {
            obj[data.name] = data.value;
        }
    }
    return obj;
}

function removeEmptyDataTr() {
    $('#emptyTr').remove();
}

function getDateFormatFromStr(dateStr, momentFormat) {
    if (!dateStr) {
        return '';
    }
    return moment(dateStr).format(momentFormat);
}

