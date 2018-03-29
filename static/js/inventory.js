let socket = io();

socket.on('inventoryUpdate', function (dataStr) {
    const data = JSON.parse(dataStr);
    if (data.skuId) {
        let $tr = $('*[data-skuId="' + data.skuId + '"]').closest('tr');
        setInventoryTr($tr, data);
    }
});

function updateBySku(dom, id, action) {
    BootstrapDialog.show({
        closable: false,
        size: BootstrapDialog.SIZE_SMALL,
        title: action === 0 ? $('#strWithdraw').text() : $('#strStock').text(),
        message: createUpdateForm(id, action),
        onshown: function (bd) {
            $("#inventoryForm").on('submit', function (ev) {
                ev.preventDefault();
                let data = getObjFromSerialize($('#inventoryForm').serializeArray(), ['amount', 'inventoryId', 'action']);
                $.ajax({
                    url: '/inventory',
                    method: 'PUT',
                    dataType: 'json',
                    data: JSON.stringify(data)
                }).done(function () {
                    bd.close();
                }).fail(alertErr);
            });
        },
        buttons: [
            {
                label: $('#strBtnConfirm').text(),
                action: function () {
                    $('#inventoryFormSubmit').click();
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

function createUpdateForm(id, action) {
    let html = `<form id="inventoryForm">`;
    html += `<label for="amount">${$('#strQuantity').text()}</label>`;
    html += `<input id="amount" class="form-control" type="number" name="amount" required>`;
    html += `<input type="hidden" name="inventoryId" value="${id}">`;
    html += `<input type="hidden" name="action" value="${action}">`;
    html += `<input id="inventoryFormSubmit" type="submit" style="display: none">`;
    html += '</form>';
    return html;
}

function setInventoryTr($tr, data) {
    setHtml($tr, 'level', data.level === 0 ? '0' : data.level);
    setHtml($tr, 'user', data.user);
    setHtml($tr, 'updatedAt', getDateFormatFromStr(data.updatedAt, 'YYYY-MM-DD HH:mm:ss'));
}

